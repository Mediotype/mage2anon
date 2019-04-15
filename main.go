package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/schollz/progressbar"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/jkenneydaniel/mage2anon/src"
	"golang.org/x/crypto/ssh/terminal"
)

func writePlainFile(filePath string, data string) (int) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatal(os.Stderr, "Failed to open plain file located at: " + filePath)
	}

	defer file.Close()
	writer := bufio.NewWriter(file)

	bytesWritten, err := writer.WriteString(data)

	if err != nil {
		log.Fatal(os.Stderr, "Failed to write SQL to plain file located at: " + filePath)
	}

	writer.Flush()

	return bytesWritten
}

func writeGzippedFile(filePath string, data string) (int) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatal(os.Stderr, "Failed to open gzip file located at: " + filePath)
	}

	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	bytesWritten, err := gzipWriter.Write([]byte(data))

	if err != nil {
		log.Fatal(os.Stderr, "Failed to write SQL to gzip file located at: " + filePath)
	}

	gzipWriter.Flush()

	return bytesWritten
}

func removeFile(filePath string) {
	fileRemoveErr := os.Remove(filePath)

	if fileRemoveErr != nil {
		log.Fatal(os.Stderr, fileRemoveErr)
	}
}

func countFileLines(filePath string) (int) {
	file, fileError := os.OpenFile(filePath, os.O_RDONLY, 0644)

	if fileError != nil {
		log.Fatal(os.Stderr, fileError)
	}

	fileLineCounter := bufio.NewScanner(file)
	totalLines := 0

	for fileLineCounter.Scan() {
		totalLines++
	}

	return totalLines
}

func main() {

	// Define parameters for command
    requestedConfig := flag.String("config", "base", "Configuration to use. A \"base\" configuration is included out-of-box. Alternately, supply path to file")
	mysqlHost := flag.String("mysql-host", "127.0.0.1", "MySQL Host - defaults to 127.0.0.1")
	mysqlUser := flag.String("mysql-user", "root", "MySQL User - defaults to root")
	mysqlPass := flag.String("mysql-pass", "", "MySQL Password - defaults to root")
	mysqlPort := flag.String("mysql-port", "3306", "MySQL Port - defaults to 3306")
	mysqlTables := flag.String("mysql-tables", "", "MySQL tables - defaults to nil, useful for small exports")
	mysqlDb := flag.String("mysql-db", "", "MySQL Database - *Required*")

	// Parse the parameters
	flag.Parse()

	// Load the table configuration
	config, configErr := mage2anon.NewConfig(*requestedConfig)
	if configErr != nil {
		log.Println("Unable to load configuration!")
		log.Fatal(os.Stderr, configErr)
	}

	if len(*mysqlDb) <= 0 {
		log.Fatal(os.Stderr, "You must provide a database for us to access!")
	}

	currentTime := time.Now().Local()
	dumpFilenameFormat := fmt.Sprintf("%s-" + currentTime.Format("2006-01-01") + ".sql", *mysqlDb)
	tmpDumpLocation := "/tmp"
	tmpFilePath := path.Join(tmpDumpLocation, dumpFilenameFormat)

	// Define our MySQL config into the config variable (so it is not stored on FS)
	config.MysqlHost = *mysqlHost
	config.MysqlUser = *mysqlUser
	config.MysqlPass = *mysqlPass
	config.MysqlPort = *mysqlPort
	config.MysqlTables = *mysqlTables
	config.MysqlDb = *mysqlDb

	if len(config.MysqlPass) <= 0 {
		fmt.Println("You did not provide a password, please enter a password (or leave it empty):")
		bytePassword, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			log.Fatal(os.Stderr, "Encountered error when trying to process your password!")
		}

		config.MysqlPass = string(bytePassword)

		fmt.Println()
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dumpErr := mage2anon.PerformDump(config, tmpFilePath)

	if dumpErr != nil {
		removeFile(tmpFilePath)
		log.Println(os.Stderr, "Failed to dump temporary SQL data for processing!")
		log.Fatal(os.Stderr, dumpErr)
	}

	provider := mage2anon.NewProvider()
	eavProcessor := mage2anon.CreateEavGenerator(config, provider)
	sqlProcessor := mage2anon.CreateSQLProcessor(config, provider)
	sqlStringBuffer := bytes.Buffer{}

	newFilePath := path.Join(cwd, dumpFilenameFormat)
	newFileGzipPath := path.Join(cwd, dumpFilenameFormat) + ".gz"

	tmpFile, tmpFileError := os.OpenFile(tmpFilePath, os.O_RDONLY, 0644)

	if tmpFileError != nil {
		log.Fatal(os.Stderr, tmpFileError)
	}

	attributeParseScanner := bufio.NewScanner(tmpFile)
	attributeMapping := mage2anon.NewAttributeMapping()

	// sqlparser can be noisy
	// https://github.com/xwb1989/sqlparser/blob/120387863bf27d04bc07db8015110a6e96d0146c/ast.go#L52
	// We don't want to hear about it
	log.SetOutput(ioutil.Discard)

	totalLines := countFileLines(tmpFilePath)

	eavGeneratorProgressBar := progressbar.New(totalLines)
	eavGeneratorProgressBar.RenderBlank()

	/**
		Generate the EAV mapping for us to update EAV values on the fly
	 */
	for attributeParseScanner.Scan() {
		attrMapping := eavProcessor.GenerateEavMapping(attributeParseScanner.Text())
		attributeMapping.Attributes = append(attributeMapping.Attributes, attrMapping.Attributes...)
		eavGeneratorProgressBar.Add(1)
	}

	/**
		Reset the position of the file pointer back to the beginning
	 */
	_, seekErr := tmpFile.Seek(io.SeekStart, 0)

	if seekErr != nil {
		log.Println(os.Stderr, "Failed to seek back to beginning of temporary dump file!\n")
		log.Fatal(os.Stderr, seekErr)
	}

	/**
		Configure a new scanner
	 */
	dataScanner := bufio.NewScanner(tmpFile)

	processorProgressBar := progressbar.New(totalLines)
	processorProgressBar.RenderBlank()

	/**
		Process the SQL
	 */
	for dataScanner.Scan() {
		sqlStringBuffer.WriteString(sqlProcessor.ProcessSQL(dataScanner.Text(), attributeMapping) + "\n")
		processorProgressBar.Add(1)
	}

	removeFile(tmpFilePath)

	writePlainFile(newFilePath, sqlStringBuffer.String())
	writeGzippedFile(newFileGzipPath, sqlStringBuffer.String())
}