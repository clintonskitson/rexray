#ScaleIO

Running scale-out block storage as software

---

## Overview
The ScaleIO registers a storage driver named `scaleio` with the `REX-Ray`
driver manager and is used to connect and manage ScaleIO storage.

## Configuration
The following are all of the parameters for the `ScaleIO` driver in YAML.

```yaml
scaleio:
    endpoint:             https://domain.com/scaleio
    insecure:             false
    useCerts:             true
    userName:             admin
    password:             mypassword
    systemID:             0
    systemName:           sysv
    protectionDomainID:   0
    protectionDomainName: corp
    storagePoolID:        0
    storagePoolName:      gold
```

For information on the equivalent environment variable and CLI flag names
please see the section on how non top-level configuration properties are
[transformed](./config/#all-other-properties).

## Activating the Driver
To activate the ScaleIO driver please follow the instructions for
[activating storage drivers](/user-guide/config#activating-storage-drivers),
using `scaleio` as the driver name.

## Examples
The following is an example of a working configuration for the `ScaleIO` driver.

```yaml
storageDrivers:
    - scaleio
scaleio:
    authUrl: https://<your_rest_host>/api
    insecure: true
    username: <your_username>
    password: <your_password>
    systemName: <your_cluster_system_name>
    protectionDomainName: <your_protection_domain_name>
    storagePoolName: <your_storage_pool_name
```
