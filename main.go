package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jkenneydaniel/mage2anon/src"
	"io/ioutil"
	"log"
	"os"
)

func main() {

    requestedConfig := flag.String("config", "", "Configuration to use. A \"base\" configuration is included out-of-box. Alternately, supply path to file")
	mysqlHost := flag.String("mysql-host", "127.0.0.1", "MySQL Host - defaults to 127.0.0.1")
	mysqlUser := flag.String("mysql-user", "root", "MySQL User - defaults to root")
	mysqlPass := flag.String("mysql-pass", "root", "MySQL Password - defaults to root")
	mysqlPort := flag.String("mysql-port", "3306", "MySQL Port - defaults to 3306")
	mysqlDb := flag.String("mysql-db", "", "MySQL Database - *Required*")

	flag.Parse()

	config, err := mage2anon.NewConfig(*requested)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}



	provider := dbanon.NewProvider()
	processor := dbanon.NewLineProcessor(config, provider)
	reader := bufio.NewReader(os.Stdin)

	// sqlparser can be noisy
	// https://github.com/xwb1989/sqlparser/blob/120387863bf27d04bc07db8015110a6e96d0146c/ast.go#L52
	// We don't want to hear about it
	log.SetOutput(ioutil.Discard)

	for {
		text, err := reader.ReadString('\n')
		fmt.Print(processor.ProcessLine(text))

		if err != nil {
			break
		}
	}
}