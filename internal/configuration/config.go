package configuration

type Configurations struct {
	Version      string
	OutputDir    string
	C2lPtr       bool
	L2cPtr       bool
	HtmlPtr      bool
	TextPtr      bool
	InputPathPtr string
}

// SomeConfigurations exported
type SomeConfigurations struct {
	SomeName string
}
