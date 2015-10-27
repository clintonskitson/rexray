package cli

import (
	"testing"

	flag "github.com/spf13/pflag"

	"github.com/emccode/rexray/core/config"
)

func TestInitUsageTemplates(t *testing.T) {
	config.Register(testRegistration())
	c := NewWithArgs("-v")
	c.Execute()
}

func testRegistration() *config.Registration {
	r := &config.Registration{}
	r.Yaml = `mockProvider:
    userName: admin
    useCerts: true
    docker:
        MinVolSize: 16
`
	r.FlagSetName = "Mock Provider Flags"
	r.FlagSet = &flag.FlagSet{}
	r.FlagSet.String("mockProviderUserName", "admin", "")
	r.FlagSet.String("mockProviderPassword", "", "")
	r.FlagSet.Bool("mockProviderUseCerts", true, "")
	r.FlagSet.Int32("mockProviderDockerMinVolSize", 16, "")

	r.FlagBindings = map[string]*flag.Flag{
		"mockProvider.userName":          r.FlagSet.Lookup("mockProviderUserName"),
		"mockProvider.password":          r.FlagSet.Lookup("mockProviderPassword"),
		"mockProvider.useCerts":          r.FlagSet.Lookup("mockProviderUseCerts"),
		"mockProvider.docker.minVolSize": r.FlagSet.Lookup("mockProviderDockerMinVolSize"),
	}

	r.EnvBindings = map[string]string{
		"mockProvider.userName":          "MOCKPROVIDER_USERNAME",
		"mockProvider.password":          "MOCKPROVIDER_PASSWORD",
		"mockProvider.useCerts":          "MOCKPROVIDER_USECERTS",
		"mockProvider.docker.minVolSize": "MOCKPROVIDER_DOCKER_MINVOLSIZE",
	}

	return r
}
