package httptest

import (
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"testing"
)

type (
	// OptionSetter sets a configuration field to the configuration
	OptionSetter interface {
		// Set receives a pointer to the Configuration type and does the job of filling it
		Set(c *Configuration)
	}
	// OptionSet implements the OptionSetter
	OptionSet func(c *Configuration)
)

// Set is the func which makes the OptionSet an OptionSetter, this is used mostly
func (o OptionSet) Set(c *Configuration) {
	o(c)
}


// Configuration httptest configuration
type Configuration struct {
	// Debug if true then debug messages from the httpexpect will be shown when a test runs
	// Default is false
	Debug bool

	BaseURL string
}

// DefaultConfiguration returns the default configuration for the httptest
// all values are defaulted to false for clarity
func DefaultConfiguration() *Configuration {
	return &Configuration{BaseURL: "", Debug: false}
}


func New(h http.Handler, t *testing.T, setters ...OptionSetter) *httpexpect.Expect {
	conf := DefaultConfiguration()
	for _, setter := range setters {
		setter.Set(conf)
	}
	baseURL := ""

	if conf.BaseURL != "" {
		baseURL = conf.BaseURL
	}
	testConfiguration := httpexpect.Config{
		BaseURL: baseURL,
		Reporter: httpexpect.NewAssertReporter(t),
		Client: &http.Client{
			Transport: httpexpect.NewBinder(h),
			Jar: httpexpect.NewJar(),
		},
	}

	if conf.Debug {
		testConfiguration.Printers = []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		}
	}

	return httpexpect.WithConfig(testConfiguration)
}