package mage2anon

import (
	"syreclabs.com/go/faker"
)

type Provider struct{}

func NewProvider() *Provider {
	p := &Provider{}

	return p
}

type ProviderInterface interface {
	Get(s string) string
}

func (p Provider) Get(s string) string {
	switch s {
	case "firstname":
		return faker.Name().FirstName()
	case "lastname":
		return faker.Name().LastName()
	case "fullname":
		return faker.Name().FirstName() + " " + faker.Name().LastName()
	case "email":
		return faker.RandomString(16) + "-" + faker.Internet().Email()
	case "username":
		return faker.Internet().UserName()
	case "password":
		return faker.Internet().Password(8, 14)
	case "datetime":
		return faker.Date().Birthday(0, 40).Format("2006-01-02 15:04:05")
	case "customer_suffix":
		return faker.Name().Suffix()
	case "website":
		return faker.Internet().Url()
	case "ipv4":
		return faker.Internet().IpV4Address()
	case "state":
		return faker.Address().State()
	case "city":
		return faker.Address().City()
	case "postcode":
		return faker.Address().Postcode()
	case "street":
		return faker.Address().StreetAddress()
	case "telephone":
		return faker.PhoneNumber().PhoneNumber()
	case "title":
		return faker.Name().Prefix()
	case "company":
		return faker.Company().Name()
	case "md5":
		return faker.Lorem().Characters(32)
	case "note255":
		return faker.Lorem().Characters(50)
	case "region_id":
		// https://github.com/meanbee/magedbm2/blob/fc8bbf9a97db2c27d0cd8a1153dda8c95b6f5996/src/Anonymizer/Formatter/Address/RegionId.php#L24
		return faker.Number().Between(1, 550)
	case "gender":
		// https://github.com/meanbee/magedbm2/blob/fc8bbf9a97db2c27d0cd8a1153dda8c95b6f5996/src/Anonymizer/Formatter/Person/Gender.php#L20
		return faker.Number().Between(1, 3)
	case "country_code":
		return faker.Address().CountryCode()
	case "vat_number":
		// https://github.com/meanbee/magedbm2/blob/fc8bbf9a97db2c27d0cd8a1153dda8c95b6f5996/src/Anonymizer/Formatter/Company/VatNumber.php#L21
		return "GB" + faker.Number().Between(100000000, 999999999)
	}

	return ""
}