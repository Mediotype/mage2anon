package mage2anon

import (
	"os"
	"os/exec"
	"path"
)

func (config *Config) CreateDump(dumpDir string, dumpName string) error {

	DumpCmd := exec.Command("mysqldump")
	DumpCmd.Args = append(DumpCmd.Args, "--complete-insert")
	DumpCmd.Args = append(DumpCmd.Args, "-P"+config.MysqlPort)
	DumpCmd.Args = append(DumpCmd.Args, "-h"+config.MysqlHost)
	DumpCmd.Args = append(DumpCmd.Args, "-u"+config.MysqlUser)
	DumpCmd.Args = append(DumpCmd.Args, "-p\""+config.MysqlPass + "\"")
	DumpCmd.Args = append(DumpCmd.Args, config.MysqlDb)

	if len(config.MysqlTables) > 0 {
		DumpCmd.Args = append(DumpCmd.Args, config.MysqlTables)
	}

	out, err := os.OpenFile(path.Join(dumpDir, dumpName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	DumpCmd.Stdout = out

	return DumpCmd.Run()
}