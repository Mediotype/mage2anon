package mage2anon

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

func NewAttributeMapping() (*AttributeMapping) {

	am := &AttributeMapping{}

	return am
}