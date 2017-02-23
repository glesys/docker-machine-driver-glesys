package glesys

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/glesys/glesys-go"
)

const (
	defaultBandwidth  = 100
	defaultCPU        = 2
	defaultDataCenter = "Falkenberg"
	defaultMemory     = 2048
	defaultStorage    = 20
	defaultTemplate   = "Ubuntu 16.04 LTS 64-bit"
)

const (
	apiKeyFlag       = "glesys-api-key"
	bandwidthFlag    = "glesys-bandwidth"
	campaignCodeFlag = "glesys-campaign-code"
	cpuFlag          = "glesys-cpu"
	dataCenterFlag   = "glesys-data-center"
	memoryFlag       = "glesys-memory"
	projectFlag      = "glesys-project"
	rootPasswordFlag = "glesys-root-password"
	sshKeyPathFlag   = "glesys-ssh-key-path"
	storageFlag      = "glesys-storage"
	templateFlag     = "glesys-template"
)

// Driver for GleSYS
type Driver struct {
	*drivers.BaseDriver
	APIKey       string
	Bandwidth    int
	CampaignCode string
	CPU          int
	DataCenter   string
	Memory       int
	Project      string
	RootPassword string
	ServerID     string
	Storage      int
	Template     string
}

// NewDriver creates a new driver
func NewDriver(hostName, storePath string) drivers.Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
}

func stringFlagToEnvVar(s string) string {
	return strings.ToUpper(strings.Replace(s, "-", "_", -1))
}

// GetCreateFlags defines all flags and environment variables that can be used with docker-machine create
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: stringFlagToEnvVar(projectFlag),
			Name:   projectFlag,
			Usage:  "GleSYS project (e.g. CL12345)",
		},
		mcnflag.StringFlag{
			EnvVar: stringFlagToEnvVar(apiKeyFlag),
			Name:   apiKeyFlag,
			Usage:  "GleSYS API key",
		},
		mcnflag.IntFlag{
			Name:  memoryFlag,
			Usage: "Memory in MB",
			Value: defaultMemory,
		},
		mcnflag.IntFlag{
			Name:  cpuFlag,
			Usage: "Number of CPU cores",
			Value: defaultCPU,
		},
		mcnflag.IntFlag{
			Name:  storageFlag,
			Usage: "Storage in GB",
			Value: defaultStorage,
		},
		mcnflag.IntFlag{
			Name:  bandwidthFlag,
			Usage: "Bandwidth in MBit/s",
			Value: defaultBandwidth,
		},
		mcnflag.StringFlag{
			Name:  dataCenterFlag,
			Usage: "Data center to place the machine in",
			Value: defaultDataCenter,
		},
		mcnflag.StringFlag{
			Name:  rootPasswordFlag,
			Usage: "Root password to use for the machine. If omitted, a random password will be generated",
		},
		mcnflag.StringFlag{
			Name:  templateFlag,
			Usage: "Template to use for the machine",
			Value: defaultTemplate,
		},
		mcnflag.StringFlag{
			Name:  campaignCodeFlag,
			Usage: "Campaign code to use for the machine",
			Value: "",
		},
		mcnflag.StringFlag{
			Name:  sshKeyPathFlag,
			Usage: "Path to the SSH key file you want to use. If omitted, a new key will be generated",
			Value: "",
		},
	}
}

// SetConfigFromFlags configures the driver with the object that was returned
// by RegisterCreateFlags
func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.Project = opts.String(projectFlag)
	d.APIKey = opts.String(apiKeyFlag)

	if d.Project == "" {
		return fmt.Errorf("glesys driver requires the --%v option", projectFlag)
	}

	if d.APIKey == "" {
		return fmt.Errorf("glesys driver require the --%v option", apiKeyFlag)
	}

	d.Bandwidth = opts.Int(bandwidthFlag)
	d.CampaignCode = opts.String(campaignCodeFlag)
	d.CPU = opts.Int(cpuFlag)
	d.DataCenter = opts.String(dataCenterFlag)
	d.Memory = opts.Int(memoryFlag)
	d.RootPassword = opts.String(rootPasswordFlag)
	d.SSHKeyPath = opts.String(sshKeyPathFlag)
	d.Storage = opts.Int(storageFlag)
	d.Template = opts.String(templateFlag)

	return nil
}

// PreCreateCheck allows for pre-create operations to make sure a driver is ready for creation
func (d *Driver) PreCreateCheck() error {
	if d.RootPassword == "" {
		d.RootPassword = generatePassword(64)
	}

	if d.SSHKeyPath != "" {
		if _, err := os.Stat(d.SSHKeyPath); os.IsNotExist(err) {
			return fmt.Errorf("SSH key file does not exist: %q", d.SSHKeyPath)
		}
	}

	return nil
}

// Create a host using the driver's config
func (d *Driver) Create() error {
	client := d.getClient()

	if d.SSHKeyPath == "" {
		if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
			return err
		}
	}

	publicKey, err := ioutil.ReadFile(d.GetSSHKeyPath() + ".pub")
	if err != nil {
		return err
	}

	server, err := client.Servers.Create(context.Background(), glesys.CreateServerParams{
		Bandwidth:    d.Bandwidth,
		CampaignCode: d.CampaignCode,
		CPU:          d.CPU,
		DataCenter:   d.DataCenter,
		Hostname:     d.GetMachineName(),
		IPv4:         "any",
		IPv6:         "any",
		Memory:       d.Memory,
		Password:     d.RootPassword,
		Platform:     "VMware",
		PublicKey:    string(publicKey),
		Storage:      d.Storage,
		Template:     d.Template,
	})

	if err != nil {
		return fmt.Errorf("Failed to create machine: %v", err)
	}

	d.ServerID = server.ID

	log.Info("Waiting for machine to come online...")

	for {
		server, err = client.Servers.Details(context.Background(), d.ServerID)
		if err != nil {
			return err
		}

		if server.IsLocked() == false {
			break
		}
		time.Sleep(1 * time.Second)
	}

	d.IPAddress = server.IPList[0].Address

	return nil
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return "glesys"
}

// GetSSHHostname returns hostname for use with ssh
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// GetSSHPort returns port for use with ssh
func (d *Driver) GetSSHPort() (int, error) {
	return 22, nil
}

// GetSSHUsername returns username for use with ssh
func (d *Driver) GetSSHUsername() string {
	return "root"
}

// GetURL returns a Docker compatible URL for connecting to this machine
// e.g. tcp://1.2.3.4:2376
func (d *Driver) GetURL() (string, error) {
	if err := drivers.MustBeRunning(d); err != nil {
		return "", err
	}

	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, "2376")), nil
}

// GetState returns the current state of a machine
func (d *Driver) GetState() (state.State, error) {

	s, err := d.getClient().Servers.Details(context.Background(), d.ServerID)
	if err != nil {
		return state.Error, err
	}
	switch s.State {
	case "locked":
		return state.None, nil
	case "running":
		return state.Running, nil
	case "stopped":
		return state.Stopped, nil
	}

	return state.None, nil
}

// Kill stops a machine forcefully
func (d *Driver) Kill() error {
	return d.getClient().Servers.Stop(context.Background(), d.ServerID, glesys.StopServerParams{
		Type: "hard",
	})
}

func generatePassword(passlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, passlen)
	for i := 0; i < passlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

// Remove a machine
func (d *Driver) Remove() error {
	if d.ServerID == "" {
		log.Warn("Can not remove the machine, ID is missing")
		return nil
	}
	return d.getClient().Servers.Destroy(context.Background(), d.ServerID, glesys.DestroyServerParams{
		KeepIP: false,
	})
}

// Restart a machine.
func (d *Driver) Restart() error {
	return d.getClient().Servers.Stop(context.Background(), d.ServerID, glesys.StopServerParams{
		Type: "reboot",
	})
}

// Start a machine
func (d *Driver) Start() error {
	return d.getClient().Servers.Start(context.Background(), d.ServerID)
}

// Stop a machine gracefully
func (d *Driver) Stop() error {
	return d.getClient().Servers.Stop(context.Background(), d.ServerID, glesys.StopServerParams{
		Type: "soft",
	})
}

func (d *Driver) getClient() *glesys.Client {
	return glesys.NewClient(d.Project, d.APIKey, "docker-machine-driver-glesys/1.0.0")
}
