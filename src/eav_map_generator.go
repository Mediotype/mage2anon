package mage2anon

import (
	"github.com/xwb1989/sqlparser"
	"strings"
)

func CreateEavGenerator(c *Config, p ProviderInterface) *LineProcessor {
	return &LineProcessor{Config: c, Provider: p}
}

func (p LineProcessor) GenerateEavMapping(data string) *AttributeMapping {
	attributeMapping := NewAttributeMapping()
	i := strings.Index(data, "INSERT")
	if i != 0 {
		// We are only processing lines that begin with INSERT
		return attributeMapping
	}

	stmt, _ := sqlparser.Parse(data)
	insert := stmt.(*sqlparser.Insert)

	table := insert.Table.Name.String()

	switch table {
	case "eav_attribute":
		// Pull out EAV Attribute ID's by Attribute Code
		var attributeCodeColInd = insert.Columns.FindColumn(sqlparser.NewColIdent("attribute_code"))
		var attributeIdColInd = insert.Columns.FindColumn(sqlparser.NewColIdent("attribute_id"))

		rows := insert.Rows.(sqlparser.Values)

		for _, vt := range rows {
			var attributeId string
			var attributeCode string
			for i, e := range vt {
				if i == attributeCodeColInd {
					switch v := e.(type) {
					case *sqlparser.SQLVal:
						attributeCode = string(v.Val)
					}

				}

				if i == attributeIdColInd {
					switch v := e.(type) {
					case *sqlparser.SQLVal:
						attributeId = string(v.Val)
					}

				}
			}

			attributeMapping.Attributes = append(attributeMapping.Attributes, AttributeDefinition{
				AttributeId: attributeId,
				AttributeCode: attributeCode,
			})
		}

		return attributeMapping
	default:
		return attributeMapping
	}
}