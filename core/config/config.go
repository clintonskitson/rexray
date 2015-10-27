package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/emccode/rexray/core/errors"
	"github.com/emccode/rexray/util"
)

var (
	envVarRx      *regexp.Regexp
	registrations []*Registration
)

func init() {
	envVarRx = regexp.MustCompile(`^\s*([^#=]+)=(.+)$`)
	loadEtcEnvironment()
	Register(globalRegistration())
	Register(driverRegistration())
}

func globalRegistration() *Registration {
	r := &Registration{}
	r.Yaml = `host: tcp://:7979
logLevel: warn
`
	r.FlagSetName = "Global Flags"
	r.FlagSet = &flag.FlagSet{}
	r.FlagSet.StringP(
		"host", "h", "tcp://:7979", "The REX-Ray host")
	r.FlagSet.StringP(
		"logLevel", "l", "warn", "The log level (error, warn, info, debug)")

	r.BindFlagSet = true

	return r
}

func driverRegistration() *Registration {
	r := &Registration{}
	r.Yaml = `osDrivers:
- linux
volumeDrivers:
- docker
`

	r.FlagSetName = "Driver Flags"
	r.FlagSet = &flag.FlagSet{}
	r.FlagSet.StringSlice(
		"osDrivers", []string{"linux"}, "The OS drivers to consider")
	r.FlagSet.StringSlice(
		"storageDrivers", []string{}, "The storage drivers to consider")
	r.FlagSet.StringSlice(
		"volumeDrivers", []string{"docker"}, "The volume drivers to consider")

	r.BindFlagSet = true

	return r
}

// Config contains the configuration information
type Config struct {
	FlagSets map[string]*flag.FlagSet `json:"-"`
	v        *viper.Viper
}

// Registration is used to register configuration information with the config
// package.
type Registration struct {
	Yaml         string
	EnvBindings  map[string]string
	FlagSet      *flag.FlagSet
	FlagSetName  string
	BindFlagSet  bool
	FlagBindings map[string]*flag.Flag
}

// Register registers a new configuration with the config package.
func Register(r *Registration) {
	registrations = append(registrations, r)
}

// New initializes a new instance of a Config struct
func New() *Config {
	return NewConfig(true, true, "config", "yml")
}

// NewConfig initialies a new instance of a Config object with the specified
// options.
func NewConfig(
	loadGlobalConfig, loadUserConfig bool,
	configName, configType string) *Config {

	log.Debug("initializing configuration")

	c := &Config{
		v:        viper.New(),
		FlagSets: map[string]*flag.FlagSet{},
	}
	c.v.SetTypeByDefaultValue(true)
	c.v.SetConfigName(configName)
	c.v.SetConfigType(configType)

	c.processRegistrations()

	cfgFile := fmt.Sprintf("%s.%s", configName, configType)
	etcRexRayFile := util.EtcFilePath(cfgFile)
	usrRexRayFile := fmt.Sprintf("%s/.rexray/%s", util.HomeDir(), cfgFile)

	if loadGlobalConfig && util.FileExists(etcRexRayFile) {
		log.WithField("path", etcRexRayFile).Debug("loading global config file")
		if err := c.ReadConfigFile(etcRexRayFile); err != nil {
			log.WithFields(log.Fields{
				"path":  etcRexRayFile,
				"error": err}).Error(
				"error reading global config file")
		}
	}

	if loadUserConfig && util.FileExists(usrRexRayFile) {
		log.WithField("path", usrRexRayFile).Debug("loading user config file")
		if err := c.ReadConfigFile(usrRexRayFile); err != nil {
			log.WithFields(log.Fields{
				"path":  usrRexRayFile,
				"error": err}).Error(
				"error reading user config file")
		}
	}

	return c
}

func (c *Config) processRegistrations() {
	for _, r := range registrations {

		// bind the env vars
		if r.EnvBindings != nil {
			for k, v := range r.EnvBindings {
				c.v.BindEnv(k, v)
			}
		}

		// add the flag set
		if r.FlagSet != nil {
			c.FlagSets[r.FlagSetName] = r.FlagSet
		}

		// read the config
		if r.Yaml != "" {
			c.ReadConfig(bytes.NewReader([]byte(r.Yaml)))
		}
	}
}

// BindFlags binds the parsed flags to the configuration.
func (c *Config) BindFlags() {
	for _, r := range registrations {
		if r.BindFlagSet {
			c.v.BindPFlags(r.FlagSet)
		} else if r.FlagBindings != nil {
			for k, v := range r.FlagBindings {
				c.v.BindPFlag(k, v)
			}
		}
	}
}

// Copy creates a copy of this Config instance
func (c *Config) Copy() (*Config, error) {
	newC := New()
	c.v.Unmarshal(newC.v.AllSettings())
	return newC, nil
}

// FromJSON initializes a new Config instance from a JSON string
func FromJSON(from string) (*Config, error) {
	c := New()
	if err := json.Unmarshal([]byte(from), c.v.AllSettings()); err != nil {
		return nil, err
	}
	return c, nil
}

// ToJSON exports this Config instance to a JSON string
func (c *Config) ToJSON() (string, error) {
	var err error
	var buf []byte
	if buf, err = json.MarshalIndent(c, "", "  "); err != nil {
		return "", err
	}
	return string(buf), nil
}

// MarshalJSON implements the encoding/json.Marshaller interface. It allows
// this type to provide its own marshalling routine.
func (c *Config) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(c.v.AllSettings(), "", "  ")
}

// ReadConfig reads a configuration stream into the current config instance
func (c *Config) ReadConfig(in io.Reader) error {

	if in == nil {
		return errors.New("config reader is nil")
	}

	c.v.ReadConfigNoNil(in)

	return nil
}

// ReadConfigFile reads a configuration files into the current config instance
func (c *Config) ReadConfigFile(filePath string) error {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return c.ReadConfig(bytes.NewBuffer(buf))
}

func (c *Config) updateFlag(name string, flags *flag.FlagSet) {
	if f := flags.Lookup(name); f != nil {
		val := c.v.Get(name)
		strVal := fmt.Sprintf("%v", val)
		f.DefValue = strVal
	}
}

// EnvVars returns an array of the initialized configuration keys as key=value
// strings where the key is configuration key's environment variable key and
// the value is the current value for that key.
func (c *Config) EnvVars() []string {
	as := c.v.AllSettings()
	ev := make(map[string]string)
	c.flattenEnvVars("", as, ev)

	var evArr []string
	for k, v := range ev {
		evArr = append(evArr, fmt.Sprintf("%s=%s", k, v))
	}

	return evArr
}

//flattenEnvVars returns a map of configuration keys coming from a config
//which may have been nested.
func (c *Config) flattenEnvVars(pk string, as map[string]interface{}, ev map[string]string) {
	if pk == "" {
		pk = "REXRAY"
	}

	for k, v := range as {
		sk := fmt.Sprintf("%s_%s", pk, k)
		switch v.(type) {
		case string:
			{
				ev[strings.ToUpper(sk)] = v.(string)
			}
		case []interface{}:
			{
				var vArr []string
				for _, iv := range v.([]interface{}) {
					vArr = append(vArr, iv.(string))
				}
				ev[strings.ToUpper(sk)] = strings.Join(vArr, " ")
			}
		case map[string]interface{}:
			{
				c.flattenEnvVars(sk, v.(map[string]interface{}), ev)
			}
		case bool:
			{
				ev[strings.ToUpper(sk)] = fmt.Sprintf("%v", v.(bool))
			}
		case int:
			{
				ev[strings.ToUpper(sk)] = fmt.Sprintf("%v", v.(int))
			}
		default:
			{
				panic(fmt.Sprintf("invalid type: %s", reflect.TypeOf(v)))
			}
		}

	}

	return
}

// GetString returns the value associated with the key as a string
func (c *Config) GetString(k string) string {
	return c.v.GetString(k)
}

// GetBool returns the value associated with the key as a bool
func (c *Config) GetBool(k string) bool {
	return c.v.GetBool(k)
}

// GetStringSlice returns the value associated with the key as a string slice
func (c *Config) GetStringSlice(k string) []string {
	return c.v.GetStringSlice(k)
}

// GetInt returns the value associated with the key as an int
func (c *Config) GetInt(k string) int {
	return c.v.GetInt(k)
}

func loadEtcEnvironment() {
	lr := util.LineReader("/etc/environment")
	if lr == nil {
		return
	}
	for l := range lr {
		m := envVarRx.FindStringSubmatch(l)
		if m == nil || len(m) < 3 || os.Getenv(m[1]) != "" {
			continue
		}
		os.Setenv(m[1], m[2])
	}
}
