package mage2anon

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func PerformDump(c *Config, dumpDir string, dumpName string) error {
	DumpCmd := exec.Command("mysqldump")
	DumpCmd.Args = append(DumpCmd.Args, "--complete-insert")
	DumpCmd.Args = append(DumpCmd.Args, "--single-transaction")
	DumpCmd.Args = append(DumpCmd.Args, "-P"+c.MysqlPort)
	DumpCmd.Args = append(DumpCmd.Args, "-h"+c.MysqlHost)
	DumpCmd.Args = append(DumpCmd.Args, "-u"+c.MysqlUser)
	DumpCmd.Args = append(DumpCmd.Args, "-p"+c.MysqlPass + "")
	DumpCmd.Args = append(DumpCmd.Args, c.MysqlDb)

	if len(c.MysqlTables) > 0 {
		additionalTables := strings.Fields(c.MysqlTables)
		DumpCmd.Args = append(DumpCmd.Args, additionalTables...)
	}

	outFile, err := os.OpenFile(path.Join(dumpDir, dumpName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	DumpCmd.Stdout = outFile

	return DumpCmd.Run()
}