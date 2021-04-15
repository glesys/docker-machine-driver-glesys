[![Build Status](https://travis-ci.org/glesys/docker-machine-driver-glesys.svg?branch=master)](https://travis-ci.org/glesys/docker-machine-driver-glesys)

# GleSYS driver for Docker Machine

This is the official Docker Machine driver to create Docker machines with
[GleSYS](https://glesys.com).

## Installation

### From a Release

Binaries are available for Linux, OS X and Windows on the
[releases](https://github.com/glesys/docker-machine-driver-glesys/releases)
page.

Download the binary for your platform, make sure to add it to your `PATH` and
make it executable.

Rename the file to `docker-machine-driver-glesys`.

### Building from Source

Make sure you have Go and `GOPATH` setup correctly. Then run:

```shell
go get github.com/glesys/docker-machine-driver-glesys
cd $GOPATH/src/github.com/glesys/docker-machine-driver-glesys
go install ./cmd/docker-machine-driver-glesys
```

## Obtaining Credentials

To use this driver you need a GleSYS Cloud account and a valid API key. You can
sign up for a free account at https://glesys.com/signup then visit
https://customer.glesys.com to create an API key for any of your projects.

## Using the Driver

You can ensure that `docker-machine` can find the GleSYS driver by asking for
the driver help:

```shell
docker-machine create -d glesys | grep glesys

  --glesys-project                             GleSYS project (e.g. CL12345) [$GLESYS_PROJECT]
  --glesys-api-key                             GleSYS API key [$GLESYS_API_KEY]

  --glesys-bandwidth "100"                     Bandwidth in MBit/s
  --glesys-campaign-code                       Campaign code to use for the machine
  --glesys-cpu "2"                             Number of CPU cores
  --glesys-data-center "Falkenberg"            Data center to place the machine in
  --glesys-memory "2048"                       Memory in MB
  --glesys-root-password                       Root password to use for the machine. If omitted, a random password will be generated (VMware only)
  --glesys-username-kvm "docker-machine"       Username to use in KVM platform
  --glesys-ssh-key-path                        Path to the SSH key file you want to use. If omitted, a new key will be generated
  --glesys-storage "20"                        Storage in GB
  --glesys-template "Ubuntu 16.04 LTS 64-bit"  Template to use for the machine
  --glesys-platform "VMware"                   Virtualization platform (VMware or KVM)
```

To create a machine you need to specify a project and an API key:

```shell
docker-machine create -d glesys --glesys-project=CL12345 --glesys-api-key=my-api-key example-host

Running pre-create checks...
Creating machine...
(example-host) Waiting for machine to come online...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with ubuntu(systemd)...
Installing Docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine,
run: docker-machine env example-host
```

## Help

If you need any help using this driver feel free to send an email to
support@glesys.com.

## License

The contents of this repository are distributed under the MIT license, see
[LICENSE](LICENSE).
