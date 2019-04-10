package mage2anon

import (
	"github.com/xwb1989/sqlparser"
	"strings"
	"fmt"
)

func (p LineProcessor) ProcessTable(s string) string {
	i := strings.Index(s, "INSERT")
	if i != 0 {
		// We are only processing lines that begin with INSERT
		return s
	}

	stmt, _ := sqlparser.Parse(s)
	attributeMapping := &AttributeMapping{}

	switch stmt := stmt.(type) {
	case *sqlparser.Insert:

		table := stmt.Table.Name.String()

		processor := p.Config.ProcessTable(table)

		switch processor {
		case "":
			// This table doesn't need to be processed
			return s
		case "table":
			// "Classic" processing
			rows := stmt.Rows.(sqlparser.Values)
			for _, vt := range rows {
				for i, e := range vt {
					column := stmt.Columns[i].String()

					result, dataType := p.Config.ProcessColumn(table, column)

					if !result {
						continue
					}

					switch v := e.(type) {
					case *sqlparser.SQLVal:
						switch v.Type {
						default:
							v.Val = []byte(p.Provider.Get(dataType))
						}
					}
				}
			}
			return sqlparser.String(stmt) + ";\n"
		case "eav":
			// EAV processing
			var attributeId string
			rows := stmt.Rows.(sqlparser.Values)
			for _, vt := range rows {
				for i, e := range vt {
					column := stmt.Columns[i].String()
					if column == "attribute_id" {
						switch v := e.(type) {
						case *sqlparser.SQLVal:
							switch v.Type {
							default:
								fmt.Println(attributeMapping.GetAttributeCodeById(attributeId))
								attributeId = string(v.Val)
							}
						}
					}
				}
			}

			return "FOOO" + attributeId + "\n"
		default:
			return s
		}
	default:
		return s
	}
}