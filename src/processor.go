package mage2anon

type LineProcessor struct {
	Config *Config
	Provider ProviderInterface
}

func ProcessTable(c *Config, p ProviderInterface) *LineProcessor {
	return &LineProcessor{Config: c, Provider: p}
}

func ProcessEav(c *Config, p ProviderInterface) *LineProcessor {
	return &LineProcessor{Config: c, Provider: p}
}