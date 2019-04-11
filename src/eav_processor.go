package mage2anon

import (
	"github.com/xwb1989/sqlparser"
	"strings"
)

type AttributeDefinition struct {
	AttributeId		string
	AttributeCode	string
}

type AttributeMapping struct {
	Attributes	[]AttributeDefinition
}

func (am AttributeMapping) GetAttributeCodeById(attributeId string) string {
	var attributeCode string
	for _, attribute := range am.Attributes {
		if attribute.AttributeId == attributeId {
			attributeCode = attribute.AttributeCode
		}
	}

	return attributeCode
}

func (p LineProcessor) ProcessEav(s string) string {
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

		if table == "eav_attribute" {
			// Pull out EAV Attribute ID's by Attribute Code
			var attributeCodeColInd = stmt.Columns.FindColumn(sqlparser.NewColIdent("attribute_code"))
			var attributeIdColInd = stmt.Columns.FindColumn(sqlparser.NewColIdent("attribute_id"))

			rows := stmt.Rows.(sqlparser.Values)

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

			return s
		}
	default:
		return s
	}

	return s
}