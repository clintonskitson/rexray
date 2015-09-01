package volume

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	"github.com/emccode/rexray/config"
	volumedriver "github.com/emccode/rexray/drivers/volume"
	"github.com/emccode/rexray/storage"
)

type VolumeDriverManager struct {
	Drivers map[string]volumedriver.Driver
}

func Init(
	cfg *config.Config,
	storageDriverManager *storage.StorageDriverManager) (*VolumeDriverManager, error) {

	vd, vdErr := volumedriver.GetDrivers(cfg.VolumeDrivers, storageDriverManager)
	if vdErr != nil {
		return nil, vdErr
	}

	if len(vd) == 0 {
		log.Debug("no volume manager adapters initialized")
	}

	return &VolumeDriverManager{
		Drivers: vd,
	}, nil
}

func (vdm *VolumeDriverManager) Mount(volumeName, volumeID string, overwriteFs bool, newFsType string) (string, error) {
	for _, driver := range vdm.Drivers {
		return driver.Mount(volumeName, volumeID, overwriteFs, newFsType)
	}
	return "", errors.New("no volume manager specified")
}

func (vdm *VolumeDriverManager) Unmount(volumeName, volumeID string) error {
	for _, driver := range vdm.Drivers {
		return driver.Unmount(volumeName, volumeID)
	}
	return errors.New("no volume manager specified")
}

func (vdm *VolumeDriverManager) Path(volumeName, volumeID string) (string, error) {
	for _, driver := range vdm.Drivers {
		return driver.Path(volumeName, volumeID)
	}
	return "", errors.New("no volume manager specified")
}

func (vdm *VolumeDriverManager) Create(volumeName string) error {
	for _, driver := range vdm.Drivers {
		return driver.Create(volumeName)
	}
	return errors.New("no volume manager specified")
}

func (vdm *VolumeDriverManager) Remove(volumeName string) error {
	for _, driver := range vdm.Drivers {
		return driver.Remove(volumeName)
	}
	return errors.New("no volume manager specified")
}

func (vdm *VolumeDriverManager) Attach(volumeName, instanceID string) (string, error) {
	for _, driver := range volumedriver.Adapters {
		return driver.Attach(volumeName, instanceID)
	}
	return "", errors.New("no volume manager specified")
}

func (vdm *VolumeDriverManager) Detach(volumeName, instanceID string) error {
	for _, driver := range vdm.Drivers {
		return driver.Detach(volumeName, instanceID)
	}
	return errors.New("no volume manager specified")
}

func (vdm *VolumeDriverManager) NetworkName(volumeName, instanceID string) (string, error) {
	for _, driver := range vdm.Drivers {
		return driver.NetworkName(volumeName, instanceID)
	}
	return "", errors.New("no volume manager specified")
}
