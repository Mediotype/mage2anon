package main

import (
	"flag"
	"fmt"
	"github.com/jkenneydaniel/mage2anon/src"
	"io/ioutil"
	"log"
	"os/exec"
	"os"
)

func main() {

	// Define parameters for command
    requestedConfig := flag.String("config", "base", "Configuration to use. A \"base\" configuration is included out-of-box. Alternately, supply path to file")
	mysqlHost := flag.String("mysql-host", "127.0.0.1", "MySQL Host - defaults to 127.0.0.1")
	mysqlUser := flag.String("mysql-user", "root", "MySQL User - defaults to root")
	mysqlPass := flag.String("mysql-pass", "root", "MySQL Password - defaults to root")
	mysqlPort := flag.String("mysql-port", "3306", "MySQL Port - defaults to 3306")
	mysqlDb := flag.String("mysql-db", "", "MySQL Database - *Required*")

	// Parse the parameters
	flag.Parse()

	// Load the table configuration
	config, err := mage2anon.NewConfig(*requestedConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Define our MySQL config into the config variable (so it is not stored on FS)
	config.MysqlHost = *mysqlHost
	config.MysqlUser = *mysqlUser
	config.MysqlPass = *mysqlPass
	config.MysqlPort = *mysqlPort
	config.MysqlDb = *mysqlDb

	DumpCmd := exec.Command(
		"mysqldump",
		"--complete-insert",
		"-P"+config.MysqlPort,
		"-h"+config.MysqlHost,
		"-u"+config.MysqlUser,
		"-p"+config.MysqlPass,
		config.MysqlDb,
	)

	stdout, err := DumpCmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := DumpCmd.Start(); err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}

	config.MysqlData = string(bytes)

	fmt.Println(config.MysqlData)
}