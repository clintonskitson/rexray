#Rackspace

They manage your services, we manage your storage

---

## Overview
The Rackspace driver registers a storage driver named `rackspace` with the
`REX-Ray` driver manager and is used to connect and manage storage on Rackspace
instances.

## Configuration
The following are all of the parameters for the `Rackspace` driver in YAML.

```yaml
rackspace:
    authURL:    https://domain.com/rackspace
    userID:     0
    userName:   admin
    password:   mypassword
    tenantID:   0
    tenantName: customer
    domainID:   0
    domainName: corp
```

## Activating the Driver
To activate the Rackspace driver please follow the instructions for
[activating storage drivers](/user-guide/config#activating-storage-drivers),
using `rackspace` as the driver name.

## Examples
The following is an example of a working configuration for the `Rackspace` driver.

```yaml
storageDrivers:
    - rackspace
rackspace:
    authURL:    https://identity.api.rackspacecloud.com/v2.0
    userName:   <your_username>
    password:   <your_password>
    tenantName: <your_tenant_name>
```
