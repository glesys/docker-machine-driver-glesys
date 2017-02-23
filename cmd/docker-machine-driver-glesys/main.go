package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/glesys/docker-machine-driver-glesys"
)

func main() {
	plugin.RegisterDriver(glesys.NewDriver("default", "path"))
}
