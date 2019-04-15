package mage2anon

import (
	"github.com/xwb1989/sqlparser"
	"reflect"
	"strings"
)

type LineProcessor struct {
	Config *Config
	Provider ProviderInterface
}

func CreateSQLProcessor(c *Config, p ProviderInterface) *LineProcessor {
	return &LineProcessor{Config: c, Provider: p}
}

func (p LineProcessor) ProcessSQL(data string, mapping *AttributeMapping) string {
	i := strings.Index(data, "INSERT")
	if i != 0 {
		// We are only processing lines that begin with INSERT
		return data
	}

	stmt, _ := sqlparser.Parse(data)
	insert := stmt.(*sqlparser.Insert)

	table := insert.Table.Name.String()

	processor := p.Config.IdentifyTable(table)

	switch processor {
	case "":
		return data
	case "table":
		// "Classic" processing
		rows := insert.Rows.(sqlparser.Values)
		for _, vt := range rows {
			for i, e := range vt {
				column := insert.Columns[i].String()

				result, dataType := p.Config.ProcessColumn(table, column)

				if !result {
					continue
				}

				sqlValue := e.(*sqlparser.SQLVal)
				sqlValue.Val = []byte(p.Provider.Get(dataType))
			}
		}
		return sqlparser.String(insert) + ";"
	case "eav":
		// EAV processing
		var attributeId string
		rows := insert.Rows.(sqlparser.Values)
		for _, vt := range rows {
			for i, e := range vt {
				column := insert.Columns[i].String()

				switch column {
				// Fetch the attribute ID so we can later get the code when processing the mapping
				case "attribute_id":
					sqlValue := e.(*sqlparser.SQLVal)
					attributeId = string(sqlValue.Val)
				case "value":
					attributeCode := mapping.GetAttributeCodeById(attributeId)
					result, dataType := p.Config.ProcessAttribute(table, attributeCode)

					if !result {
						continue
					}

					if reflect.TypeOf(e) != reflect.TypeOf((*sqlparser.SQLVal)(nil)) {
						continue
					}

					sqlValue := e.(*sqlparser.SQLVal)
					sqlValue.Val = []byte(p.Provider.Get(dataType))
				}
			}
		}

		return sqlparser.String(insert) + ";"
	default:
		return data
	}
}