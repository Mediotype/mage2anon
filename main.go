package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jkenneydaniel/mage2anon/src"
	"io/ioutil"
	"log"
	"os"
	"time"
	"compress/gzip"
)

func main() {

	// Define parameters for command
    requestedConfig := flag.String("config", "base", "Configuration to use. A \"base\" configuration is included out-of-box. Alternately, supply path to file")
	mysqlHost := flag.String("mysql-host", "127.0.0.1", "MySQL Host - defaults to 127.0.0.1")
	mysqlUser := flag.String("mysql-user", "root", "MySQL User - defaults to root")
	mysqlPass := flag.String("mysql-pass", "root", "MySQL Password - defaults to root")
	mysqlPort := flag.String("mysql-port", "3306", "MySQL Port - defaults to 3306")
	mysqlTables := flag.String("mysql-tables", "", "MySQL tables - defaults to nil, useful for small exports")
	mysqlDb := flag.String("mysql-db", "", "MySQL Database - *Required*")

	// Parse the parameters
	flag.Parse()

	// Load the table configuration
	config, err := mage2anon.NewConfig(*requestedConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(*mysqlDb) <= 0 {
		fmt.Fprintln(os.Stderr, "You must provide a database for us to access")
		os.Exit(1)
	}

	currentTime := time.Now().Local()
	dumpFilenameFormat := fmt.Sprintf("%s-" + currentTime.Format("2006-01-01") + ".sql", *mysqlDb)
	tmpDumpLocation := "/tmp/" + dumpFilenameFormat

	// Define our MySQL config into the config variable (so it is not stored on FS)
	config.MysqlHost = *mysqlHost
	config.MysqlUser = *mysqlUser
	config.MysqlPass = *mysqlPass
	config.MysqlPort = *mysqlPort
	config.MysqlTables = *mysqlTables
	config.MysqlDb = *mysqlDb

	CreateDump(tmpDumpLocation, dumpFilenameFormat)

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	provider := mage2anon.NewProvider()
	eavProcessor := mage2anon.ProcessEav(config, provider)
	tableProcessor := mage2anon.ProcessTable(config, provider)
	tmpFile, err := os.Open(tmpDumpLocation)
	newFile, err := os.Create(cwd + "/" + dumpFilenameFormat)
	newCompressedFile, err := os.Create(cwd + "/" + dumpFilenameFormat + ".gz")
	reader := bufio.NewReader(tmpFile)
	writer := bufio.NewWriter(newFile)
	gzipWriter := gzip.NewWriter(newCompressedFile)

	// sqlparser can be noisy
	// https://github.com/xwb1989/sqlparser/blob/120387863bf27d04bc07db8015110a6e96d0146c/ast.go#L52
	// We don't want to hear about it
	log.SetOutput(ioutil.Discard)

	for {
		text, err := reader.ReadString('\n')

		eavProcessor.ProcessEav(text)

		if err != nil {
			break
		}
	}

	for {
		text, err := reader.ReadString('\n')

		writer.WriteString(tableProcessor.ProcessTable(text))
		gzipWriter.Write([]byte(tableProcessor.ProcessTable(text)))

		if err != nil {
			break
		}
	}

	writer.Flush()
	gzipWriter.Flush()
}