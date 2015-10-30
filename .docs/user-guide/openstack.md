#OpenStack

Making storage management as transparent as the stack

---

## Overview
The OpenStack driver registers a storage driver named `openstack` with the
`REX-Ray` driver manager and is used to connect and manage storage on OpenStack
instances.

## Configuration
The following are all of the parameters for the `OpenStack` driver in YAML.

```yaml
openstack:
    authURL:              https://domain.com/openstack
    userID:               0
    userName:             admin
    password:             mypassword
    tenantID:             0
    tenantName:           customer
    domainID:             0
    domainName:           corp
    regionName:           USNW
    availabilityZoneName: Gold
```

For information on the equivalent environment variable and CLI flag names
please see the section on how non top-level configuration properties are
[transformed](./config/#all-other-properties).

## Activating the Driver
To activate the OpenStack driver please follow the instructions for
[activating storage drivers](/user-guide/config#activating-storage-drivers),
using `openstack` as the driver name.

## Examples
The following is an example of a working configuration for the `OpenStack` driver.

```yaml
storageDrivers:
    - openstack
openstack:
    authUrl: https://<your_keystone_host>:35357/v2.0/
    tenantName: <your_tenant>
    regionName: <your_region>
    username: <your_username>
    password: <your_password>
```
