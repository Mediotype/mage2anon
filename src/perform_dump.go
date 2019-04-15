package mage2anon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func PerformDump(c *Config, tmpFilePath string) error {
	DumpCmdArgs := []string{
		"--complete-insert",
		"--single-transaction",
		"--net_buffer_length=4096",
		"--protocol=tcp",
		"-P"+c.MysqlPort,
		"-h"+c.MysqlHost,
		"-p"+c.MysqlPass,
		c.MysqlDb,
	}

	if len(c.MysqlTables) > 0 {
		additionalTables := strings.Fields(c.MysqlTables)
		DumpCmdArgs = append(DumpCmdArgs, additionalTables...)
	}
	DumpCmd := exec.Command("mysqldump", DumpCmdArgs...)
	fmt.Println(DumpCmd.Args)

	outFile, err := os.OpenFile(tmpFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(os.Stderr, err)
	}


	DumpCmd.Stdout = outFile

	return DumpCmd.Run()
}