module github.com/glesys/docker-machine-driver-glesys

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/docker/docker v1.13.2-0.20170601211448-f5ec1e2936dc // indirect
	github.com/docker/machine v0.16.2
	github.com/glesys/glesys-go/v2 v2.4.1
	github.com/sirupsen/logrus v1.5.0 // indirect
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20200423211502-4bdfaf469ed5 // indirect
)

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.5.0
	github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190822205725-ed20165a37b4
)

go 1.14
