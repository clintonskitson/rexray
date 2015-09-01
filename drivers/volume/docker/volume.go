package docker

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/emccode/rexray/drivers/volume"
	rros "github.com/emccode/rexray/os"
	"github.com/emccode/rexray/storage"
)

var (
	providerName string
)

const (
	defaultVolumeSize int64 = 16
)

type Driver struct {
	sdm *storage.StorageDriverManager
}

func init() {
	providerName = "docker"
	volumedriver.Register("docker", Init)
}

func Init(
	storageDriverManager *storage.StorageDriverManager) (volumedriver.Driver, error) {
	driver := &Driver{
		sdm: storageDriverManager,
	}

	if os.Getenv("REXRAY_DEBUG") == "true" {
		log.Println("Volume Manager Driver Initialized: " + providerName)
	}

	return driver, nil
}

func getVolumeMountPath(name string) (string, error) {
	if name == "" {
		return "", errors.New("Missing volume name")
	}

	return fmt.Sprintf("/var/lib/docker/volumes/%s", name), nil
}

// Mount will perform the steps to get an existing Volume with or without a fileystem mounted to a guest
func (driver *Driver) Mount(volumeName, volumeID string, overwriteFs bool, newFsType string) (string, error) {
	if volumeName == "" && volumeID == "" {
		return "", errors.New("Missing volume name or ID")
	}

	instances, err := driver.sdm.GetInstance()
	if err != nil {
		return "", err
	}

	switch {
	case len(instances) == 0:
		return "", errors.New("No instances")
	case len(instances) > 1:
		return "", errors.New("Too many instances returned, limit the storagedrivers")
	}

	volumes, err := driver.sdm.GetVolume(volumeID, volumeName)
	if err != nil {
		return "", err
	}

	switch {
	case len(volumes) == 0:
		return "", errors.New("No volumes returned by name")
	case len(volumes) > 1:
		return "", errors.New("Multiple volumes returned by name")
	}

	volumeAttachment, err := driver.sdm.GetVolumeAttach(
		volumes[0].VolumeID, instances[0].InstanceID)
	if err != nil {
		return "", err
	}

	if len(volumeAttachment) == 0 {
		volumeAttachment, err = driver.sdm.AttachVolume(
			false, volumes[0].VolumeID, instances[0].InstanceID)
		if err != nil {
			return "", err
		}
	}

	if len(volumeAttachment) == 0 {
		return "", errors.New("Volume did not attach")
	}

	mounts, err := rros.GetMounts(volumeAttachment[0].DeviceName, "")
	if err != nil {
		return "", err
	}

	if len(mounts) > 0 {
		return mounts[0].Mountpoint, nil
	}

	switch {
	case os.Getenv("REXRAY_DOCKER_VOLUMETYPE") != "":
		newFsType = os.Getenv("REXRAY_DOCKER_VOLUMETYPE")
	case newFsType == "":
		newFsType = "ext4"
	}

	if err := rros.Format(volumeAttachment[0].DeviceName, newFsType, overwriteFs); err != nil {
		return "", err
	}

	mountPath, err := getVolumeMountPath(volumes[0].Name)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(mountPath, 0755); err != nil {
		return "", err
	}

	if err := rros.Mount(volumeAttachment[0].DeviceName, mountPath, "", ""); err != nil {
		return "", err
	}

	return mountPath, nil
}

// Unmount will perform the steps to unmount and existing volume and detach
func (driver *Driver) Unmount(volumeName, volumeID string) error {
	if volumeName == "" && volumeID == "" {
		return errors.New("Missing volume name or ID")
	}

	instances, err := driver.sdm.GetInstance()
	if err != nil {
		return err
	}

	switch {
	case len(instances) == 0:
		return errors.New("No instances")
	case len(instances) > 1:
		return errors.New("Too many instances returned, limit the storagedrivers")
	}

	volumes, err := driver.sdm.GetVolume(volumeID, volumeName)
	if err != nil {
		return err
	}

	switch {
	case len(volumes) == 0:
		return errors.New("No volumes returned by name")
	case len(volumes) > 1:
		return errors.New("Multiple volumes returned by name")
	}

	volumeAttachment, err := driver.sdm.GetVolumeAttach(volumes[0].VolumeID, instances[0].InstanceID)
	if err != nil {
		return err
	}

	if len(volumeAttachment) == 0 {
		return nil
	}

	mounts, err := rros.GetMounts(volumeAttachment[0].DeviceName, "")
	if err != nil {
		return err
	}

	if len(mounts) == 0 {
		return nil
	}

	err = rros.Unmount(mounts[0].Mountpoint)
	if err != nil {
		return err
	}

	err = driver.sdm.DetachVolume(false, volumes[0].VolumeID, "")
	if err != nil {
		return err
	}
	return nil

}

// Path returns the mounted path of the volume
func (driver *Driver) Path(volumeName, volumeID string) (string, error) {
	if volumeName == "" && volumeID == "" {
		return "", errors.New("Missing volume name or ID")
	}

	instances, err := driver.sdm.GetInstance()
	if err != nil {
		return "", err
	}

	switch {
	case len(instances) == 0:
		return "", errors.New("No instances")
	case len(instances) > 1:
		return "", errors.New("Too many instances returned, limit the storagedrivers")
	}

	volumes, err := driver.sdm.GetVolume(volumeID, volumeName)
	if err != nil {
		return "", err
	}

	switch {
	case len(volumes) == 0:
		return "", errors.New("No volumes returned by name")
	case len(volumes) > 1:
		return "", errors.New("Multiple volumes returned by name")
	}

	volumeAttachment, err := driver.sdm.GetVolumeAttach(volumes[0].VolumeID, instances[0].InstanceID)
	if err != nil {
		return "", err
	}

	if len(volumeAttachment) == 0 {
		return "", nil
	}

	mounts, err := rros.GetMounts(volumeAttachment[0].DeviceName, "")
	if err != nil {
		return "", err
	}

	if len(mounts) == 0 {
		return "", nil
	}

	return mounts[0].Mountpoint, nil
}

// Create will create a remote volume
func (driver *Driver) Create(volumeName string) error {
	if volumeName == "" {
		return errors.New("Missing volume name")
	}

	instances, err := driver.sdm.GetInstance()
	if err != nil {
		return err
	}

	switch {
	case len(instances) == 0:
		return errors.New("No instances")
	case len(instances) > 1:
		return errors.New("Too many instances returned, limit the storagedrivers")
	}

	volumes, err := driver.sdm.GetVolume("", volumeName)
	if err != nil {
		return err
	}

	switch {
	case len(volumes) == 1:
		return nil
	case len(volumes) > 1:
		return errors.New(fmt.Sprintf("Too many volumes returned by name of %s", volumeName))
	}

	volumeType := os.Getenv("REXRAY_DOCKER_VOLUMETYPE")
	IOPSi, _ := strconv.Atoi(os.Getenv("REXRAY_DOCKER_IOPS"))
	IOPS := int64(IOPSi)
	sizei, _ := strconv.Atoi(os.Getenv("REXRAY_DOCKER_SIZE"))
	size := int64(sizei)
	if size == 0 {
		size = defaultVolumeSize
	}
	availabilityZone := os.Getenv("REXRAY_DOCKER_AVAILABILITYZONE")

	_, err = driver.sdm.CreateVolume(
		false, volumeName, "", "", volumeType, IOPS, size, availabilityZone)
	if err != nil {
		return err
	}
	return nil
}

// Remove will remove a remote volume
func (driver *Driver) Remove(volumeName string) error {
	if volumeName == "" {
		return errors.New("Missing volume name")
	}

	instances, err := driver.sdm.GetInstance()
	if err != nil {
		return err
	}

	switch {
	case len(instances) == 0:
		return errors.New("No instances")
	case len(instances) > 1:
		return errors.New("Too many instances returned, limit the storagedrivers")
	}

	volumes, err := driver.sdm.GetVolume("", volumeName)
	if err != nil {
		return err
	}

	switch {
	case len(volumes) == 0:
		return errors.New("No volumes returned by name")
	case len(volumes) > 1:
		return errors.New("Multiple volumes returned by name")
	}

	err = driver.Unmount("", volumes[0].VolumeID)
	if err != nil {
		return err
	}

	err = driver.sdm.RemoveVolume(volumes[0].VolumeID)
	if err != nil {
		return err
	}

	return nil
}

// Attach will attach a volume to an instance
func (driver *Driver) Attach(volumeName, instanceID string) (string, error) {
	volumes, err := driver.sdm.GetVolume("", volumeName)
	if err != nil {
		return "", err
	}

	switch {
	case len(volumes) == 0:
		return "", errors.New("No volumes returned by name")
	case len(volumes) > 1:
		return "", errors.New("Multiple volumes returned by name")
	}

	_, err = driver.sdm.AttachVolume(true, volumes[0].VolumeID, instanceID)
	if err != nil {
		return "", err
	}

	volumes, err = driver.sdm.GetVolume("", volumeName)
	if err != nil {
		return "", err
	}

	return volumes[0].NetworkName, nil
}

// Remove will remove a remote volume
func (driver *Driver) Detach(volumeName, instanceID string) error {
	volume, err := driver.sdm.GetVolume("", volumeName)
	if err != nil {
		return err
	}

	return driver.sdm.DetachVolume(true, volume[0].VolumeID, instanceID)
}

// NetworkName will return relevant information about how a volume can be discovered on an OS
func (driver *Driver) NetworkName(volumeName, instanceID string) (string, error) {
	volumes, err := driver.sdm.GetVolume("", volumeName)
	if err != nil {
		return "", err
	}

	switch {
	case len(volumes) == 0:
		return "", errors.New("No volumes returned by name")
	case len(volumes) > 1:
		return "", errors.New("Multiple volumes returned by name")
	}

	volumeAttachment, err := driver.sdm.GetVolumeAttach(
		volumes[0].VolumeID, instanceID)
	if err != nil {
		return "", err
	}

	if len(volumeAttachment) == 0 {
		return "", errors.New("Volume not attached")
	}

	volumes, err = driver.sdm.GetVolume("", volumeName)
	if err != nil {
		return "", err
	}

	return volumes[0].NetworkName, nil
}
