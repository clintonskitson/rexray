package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	flag "github.com/spf13/pflag"

	"github.com/emccode/rexray/util"
)

var (
	tmpPrefixDirs []string
	usrRexRayFile string
)

func TestMain(m *testing.M) {
	usrRexRayDir := fmt.Sprintf("%s/.rexray", util.HomeDir())
	os.MkdirAll(usrRexRayDir, 0755)
	usrRexRayFile = fmt.Sprintf("%s/%s.%s", usrRexRayDir, "config", "yml")
	usrRexRayFileBak := fmt.Sprintf("%s.bak", usrRexRayFile)

	os.Remove(usrRexRayFileBak)
	os.Rename(usrRexRayFile, usrRexRayFileBak)

	exitCode := m.Run()
	for _, d := range tmpPrefixDirs {
		os.RemoveAll(d)
	}

	os.Remove(usrRexRayFile)
	os.Rename(usrRexRayFileBak, usrRexRayFile)
	os.Exit(exitCode)
}

func newPrefixDir(testName string, t *testing.T) string {
	tmpDir, err := ioutil.TempDir(
		"", fmt.Sprintf("rexray-core-config_test-%s", testName))
	if err != nil {
		t.Fatal(err)
	}

	util.Prefix(tmpDir)
	os.MkdirAll(tmpDir, 0755)
	tmpPrefixDirs = append(tmpPrefixDirs, tmpDir)
	return tmpDir
}

func TestAssertConfigDefaults(t *testing.T) {
	newPrefixDir("TestAssertConfigDefaults", t)
	wipeEnv()
	c := New()

	osDrivers := c.GetStringSlice("osDrivers")
	volDrivers := c.GetStringSlice("volumeDrivers")

	assertString(t, c, "host", "tcp://:7979")
	assertString(t, c, "logLevel", "warn")

	if len(osDrivers) != 1 || osDrivers[0] != "linux" {
		t.Fatalf("osDrivers != []string{\"linux\"}, == %v", osDrivers)
	}

	if len(volDrivers) != 1 || volDrivers[0] != "docker" {
		t.Fatalf("volumeDrivers != []string{\"docker\"}, == %v", volDrivers)
	}
}

func TestAssertTestRegistration(t *testing.T) {
	newPrefixDir("TestAssertTestRegistration", t)
	wipeEnv()
	Register(testRegistration())
	c := New()

	userName := c.GetString("mockProvider.username")
	password := c.GetString("mockProvider.password")
	useCerts := c.GetBool("mockProvider.useCerts")
	minVolSize := c.GetInt("mockProvider.Docker.minVolSize")

	if userName != "admin" {
		t.Fatalf("mockProvider.userName != admin, == %s", userName)
	}

	if password != "" {
		t.Fatalf("mockProvider.password != '', == %s", password)
	}

	if !useCerts {
		t.Fatalf("mockProvider.useCerts != true, == %v", useCerts)
	}

	if minVolSize != 16 {
		t.Fatalf("minVolSize != 16, == %d", minVolSize)
	}
}

func TestToJSON(t *testing.T) {
	newPrefixDir("TestToJSON", t)
	wipeEnv()
	Register(testRegistration())
	c := New()

	var err error
	var strJSON string
	if strJSON, err = c.ToJSON(); err != nil {
		t.Fatal(err)
	}
	t.Log(strJSON)
}

func TestToEnvVars(t *testing.T) {
	newPrefixDir("TestToEnv", t)
	wipeEnv()
	Register(testRegistration())
	c := New()

	fev := c.EnvVars()

	for _, v := range fev {
		t.Log(v)
	}
}

func TestCopy(t *testing.T) {
	newPrefixDir("TestCopy", t)
	wipeEnv()
	Register(testRegistration())

	etcRexRayCfg := util.EtcFilePath("config.yml")
	t.Logf("etcRexRayCfg=%s", etcRexRayCfg)
	util.WriteStringToFile(string(yamlConfig1), etcRexRayCfg)

	c := New()

	assertString(t, c, "logLevel", "error")
	assertStorageDrivers(t, c)
	assertOsDrivers1(t, c)

	cc, _ := c.Copy()

	assertString(t, c, "logLevel", "error")
	assertStorageDrivers(t, cc)
	assertOsDrivers1(t, cc)

	cJSON, _ := c.ToJSON()
	ccJSON, _ := cc.ToJSON()

	cMap := map[string]interface{}{}
	ccMap := map[string]interface{}{}
	json.Unmarshal([]byte(cJSON), cMap)
	json.Unmarshal([]byte(ccJSON), ccJSON)

	if !reflect.DeepEqual(cMap, ccMap) {
		t.Fail()
	}
}

func wipeEnv() {
	evs := os.Environ()
	for _, v := range evs {
		k := strings.Split(v, "=")[0]
		os.Setenv(k, "")
	}
}

func printConfig(c *Config, t *testing.T) {
	for k, v := range c.v.AllSettings() {
		t.Logf("%s=%v\n", k, v)
	}
}

func testRegistration() *Registration {
	r := &Registration{}
	r.Yaml = `mockProvider:
    userName: admin
    useCerts: true
    docker:
        MinVolSize: 16
`
	r.FlagSetName = "Mock Provider"
	r.FlagSet = &flag.FlagSet{}
	r.FlagSet.String("mockProviderUserName", "admin", "")
	r.FlagSet.String("mockProviderPassword", "", "")
	r.FlagSet.Bool("mockProviderUseCerts", true, "")
	r.FlagSet.Int32("mockProviderDockerMinVolSize", 16, "")

	r.EnvBindings = map[string]string{
		"mockProvider.userName":          "MOCKPROVIDER_USERNAME",
		"mockProvider.password":          "MOCKPROVIDER_PASSWORD",
		"mockProvider.useCerts":          "MOCKPROVIDER_USECERTS",
		"mockProvider.docker.minVolSize": "MOCKPROVIDER_DOCKER_MINVOLSIZE",
	}

	return r
}

func assertString(t *testing.T, c *Config, key, expected string) {
	v := c.GetString(key)
	if v != expected {
		t.Fatalf("%s != %s; == %v", key, expected, v)
	}
}

func assertStorageDrivers(t *testing.T, c *Config) {
	sd := c.GetStringSlice("storageDrivers")
	if sd == nil {
		t.Fatalf("storageDrivers == nil")
	}

	if len(sd) != 2 {
		t.Fatalf("len(storageDrivers) != 2; == %d", len(sd))
	}

	if sd[0] != "ec2" {
		t.Fatalf("sd[0] != ec2; == %v", sd[0])
	}

	if sd[1] != "xtremio" {
		t.Fatalf("sd[1] != xtremio; == %v", sd[1])
	}
}

func assertOsDrivers1(t *testing.T, c *Config) {
	od := c.GetStringSlice("osDrivers")
	if od == nil {
		t.Fatalf("osDrivers == nil")
	}
	if len(od) != 1 {
		t.Fatalf("len(osDrivers) != 1; == %d", len(od))
	}
	if od[0] != "linux" {
		t.Fatalf("od[0] != linux; == %v", od[0])
	}
}

func assertOsDrivers2(t *testing.T, c *Config) {
	od := c.GetStringSlice("osDrivers")
	if od == nil {
		t.Fatalf("osDrivers == nil")
	}
	if len(od) != 2 {
		t.Fatalf("len(osDrivers) != 2; == %d", len(od))
	}
	if od[0] != "darwin" {
		t.Fatalf("od[0] != darwin; == %v", od[0])
	}
	if od[1] != "linux" {
		t.Fatalf("od[1] != linux; == %v", od[1])
	}
}

var yamlConfig1 = []byte(`logLevel: error
storageDrivers:
- ec2
- xtremio
osDrivers:
- linux`)

var yamlConfig2 = []byte(`logLevel: debug
osDrivers:
- darwin
- linux`)
