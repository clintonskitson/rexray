# Amazon EC2

Consume any storage you want any time

---

## Overview
The Amazon EC2 driver registers a storage driver named `ec2` with the `REX-Ray`
driver manager and is used to connect and manage storage on EC2 instances. The
EC2 driver is made possible by the
[goamz project](https://github.com/mitchellh/goamz).

## Configuration
The following are all of the parameters for the `AWS EC2` driver in YAML.

```yaml
aws:
    accessKey: MyAccessKey
    secretKey: MySecretKey
    region:    USNW
```

For information on the equivalent environment variable and CLI flag names
please see the section on how non top-level configuration properties are
[transformed](./config/#all-other-properties).

## Activating the Driver
To activate the EC2 driver please follow the instructions for
[activating storage drivers](/user-guide/config#activating-storage-drivers),
using `ec2` as the driver name.

## Examples
The following is an example of a working configuration for the `AWS EC2` driver.

```yaml
storageDrivers:
    - ec2
aws:
  accessKey: <your_access_key>
  secretKey: <your_secret_key>
```
