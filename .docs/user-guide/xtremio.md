#XtremIO

Not just a flash in the pan

---

## Overview
The XtremIO registers a storage driver named `xtremio` with the `REX-Ray`
driver manager and is used to connect and manage XtremIO storage.

## Configuration
The following are all of the parameters for the `XtremIO` driver in YAML.

```yaml
scaleio:
    endpoint:             https://domain.com/scaleio
    insecure:             false
    userName:             admin
    password:             mypassword
    deviceMapper:         false
    multiPath:            true
    remoteManagement:     false
```

For information on the equivalent environment variable and CLI flag names
please see the section on how non top-level configuration properties are
[transformed](./config/#all-other-properties).

## Activating the Driver
To activate the ScaleIO driver please follow the instructions for
[activating storage drivers](/user-guide/config#activating-storage-drivers),
using `xtremio` as the driver name.

## Examples
The following is an example of a working configuration for the `XtremIO` driver.

```yaml
storageDrivers:
    - xtremio
xtremio:
    endpoint: https://<your_api_host>/api/json
    insecure: true
    username: <your_username>
    password: <your_password>
    multipath: true
```
