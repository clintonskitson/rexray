package commands

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/emccode/rexray/config"
	osm "github.com/emccode/rexray/os"
	"github.com/emccode/rexray/storage"
	"github.com/emccode/rexray/volume"
)

const REXHOME = "/opt/rexray"
const EXEFILE = "/opt/rexray/rexray"
const ENVFILE = "/etc/rexray/rexray.env"
const CFGFILE = "/etc/rexray/rexray.conf"
const UNTFILE = "/etc/systemd/system/rexray.service"
const INTFILE = "/etc/init.d/rexray"

// init system types
const (
	UNKNOWN = iota
	SYSTEMD
	UPDATERCD
	CHKCONFIG
)

var (
	c *config.Config

	sdm  *storage.StorageDriverManager
	vdm  *volume.VolumeDriverManager
	osdm *osm.OSDriverManager

	client                  string
	fg                      bool
	cfgFile                 string
	snapshotID              string
	volumeID                string
	runAsync                bool
	description             string
	volumeType              string
	IOPS                    int64
	size                    int64
	instanceID              string
	volumeName              string
	snapshotName            string
	availabilityZone        string
	destinationSnapshotName string
	destinationRegion       string
	deviceName              string
	mountPoint              string
	mountOptions            string
	mountLabel              string
	fsType                  string
	overwriteFs             bool
	moduleTypeId            int32
	moduleInstanceId        int32
	moduleInstanceAddress   string
	moduleInstanceStart     bool
	moduleConfig            []string
)

type VerboseFlagPanic struct{}

//Exec function
func Exec() {
	defer func() {
		r := recover()
		if r != nil {
			switch r.(type) {
			case VerboseFlagPanic:
				// Do nothing
			default:
				panic(r)
			}
		}
	}()

	RexrayCmd.Execute()
}

func init() {
	c = config.New()
	updateLogLevel()
	initCommands()
	initFlags()
	initUsageTemplates()
}

func updateLogLevel() {
	switch c.LogLevel {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}

	log.WithField("logLevel", c.LogLevel).Debug("updated log level")
}

func preRun(cmd *cobra.Command, args []string) {
	updateLogLevel()

	cflags := cmd.Flags()
	if v, _ := cflags.GetBool("verbose"); v {
		cmd.Help()
		panic(&VerboseFlagPanic{})
	}

	initDriverManagers()
}

func initDriverManagers() {

	var osdmErr error
	osdm, osdmErr = osm.NewOSDriverManager(c)
	if osdmErr != nil {
		panic(osdmErr)
	}

	var sdmErr error
	sdm, sdmErr = storage.NewStorageDriverManager(c)
	if sdmErr != nil {
		panic(sdmErr)
	}

	var vdmErr error
	vdm, vdmErr = volume.NewVolumeDriverManager(c, osdm, sdm)
	if vdmErr != nil {
		panic(vdmErr)
	}
}
