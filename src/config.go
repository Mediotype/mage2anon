package mage2anon

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type TableConfig struct {
	Name string `yaml:"name"`
	Columns map[string]string `yaml:"columns"`
}

type EavConfig struct {
	Name string `yaml:"name"`
	Attributes map[string]string `yaml:"attributes"`
}

type Config struct {
	Tables		[]TableConfig
	Eav 		[]EavConfig
	MysqlHost	string
	MysqlUser	string
	MysqlPass	string
	MysqlPort	string
	MysqlDb		string
	MysqlTables	string
}

func NewConfig(requested string) (*Config, error) {

	c := &Config{}

	if requested == "" {
		return c, errors.New("You must specify a config")
	}

	var source []byte
	var err error

	switch requested {
	case "base":
		source, _ = Asset("etc/magento2.yml")
	default:
		source, err = ioutil.ReadFile(requested)
		if err != nil {
			return c, err
		}
	}

	yaml.Unmarshal(source, &c)

	return c, nil
}


func (c Config) IdentifyTable(t string) string {
	for _, table := range c.Tables {
		if table.Name == t {
			return "table"
		}
	}

	for _, e := range c.Eav {
		if e.Name == t {
			return "eav"
		}
	}

	return ""
}

func (c Config) ProcessColumn(t string, col string) (bool, string) {
	for _, table := range c.Tables {
		if table.Name != t {
			continue
		}

		for k, v := range table.Columns {
			if k == col {
				return true, v
			}
		}
	}

	return false, ""
}

func (c Config) ProcessAttribute(tableName string, attributeCode string) (bool, string) {
	for _, attribute := range c.Eav {
		if attribute.Name != tableName {
			continue
		}

		for k, v := range attribute.Attributes {
			if k == attributeCode {
				return true, v
			}
		}
	}

	return false, ""
}