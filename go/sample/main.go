// +build darwin

// This mostly serves as sample and test code for the hyperkit package
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/moby/hyperkit/go"
)

func stringToIntArray(l string, sep string) ([]int, error) {
	var err error
	if l == "" {
		return []int{}, nil
	}
	a := strings.Split(l, sep)
	r := make([]int, len(a))
	for idx := range a {
		if r[idx], err = strconv.Atoi(a[idx]); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func main() {
	hk := flag.String("hyperkit", "", "HyperKit binary to use")
	statedir := flag.String("state", "", "Directory to keep state in")
	vpnkitsock := flag.String("vpnkitsock", "auto", "Path to VPNKit socket")
	vpnkituuid := flag.String("vpnkituuid", "", "VPNKit UUID. Allows VMs to reconnect and get the same network configuration.")
	vpnkitip := flag.String("vpnkitip", "", "Preferred IPv4 address in VPNKit range. Requires an unused UUID.")
	var disks disks
	flag.Var(&disks, "disk", "Can be specified multiple times. Format: {file=}PATH{,size=SIZE_IN_MB} or just size=SIZE_IN_MB if -state is set. If PATH doesn't exist and size is specified the image will automatically be created.")

	kernel := flag.String("kernel", "", "Kernel to boot")
	initrd := flag.String("initrd", "", "Initial RAM Disk")

	bg := flag.Bool("background", false, "Start VM in the background")

	cpus := flag.Int("cpus", 1, "Number of CPUs")
	mem := flag.Int("mem", 1024, "Amount of memory in MB")
	vsock := flag.Bool("vsock", false, "Enable virtio-sockets")
	vsockports := flag.String("vsock-ports", "", "Comma separated list of ports to expose as sockets from guest")

	_9psock := flag.String("9p-socket", "", "9P unix domain socket to forward to the guest. Format: socket_path,9p_tag")

	iso := flag.String("iso", "", "ISO image to pass to the VM (not for booting from)")

	flag.Parse()
	cmd := flag.Args()

	if len(cmd) != 0 {
		if *statedir == "" {
			log.Fatalln("Specify existing state directory for: ", cmd[0])
		}
		h, err := hyperkit.FromState(*statedir)
		if err != nil {
			log.Fatalln("Error getting hyperkit: ", err)
		}
		switch cmd[0] {
		case "info":
			fmt.Println("Running: ", h.IsRunning())
			s, _ := json.MarshalIndent(h, "", "  ")
			fmt.Println(string(s))
			return
		case "stop", "kill":
			err := h.Stop()
			if err != nil {
				log.Fatalln("Error stopping hyperkit: ", err)
			}
			err = h.Remove(false)
			if err != nil {
				log.Fatalln("Error removing state: ", err)
			}
			return
		default:
			log.Fatalln("Unknown command: ", cmd[0])
		}
	}

	h, err := hyperkit.New(*hk, *vpnkitsock, *statedir)
	if err != nil {
		log.Fatalln("Error creating hyperkit: ", err)
	}

	h.Kernel = *kernel
	h.Initrd = *initrd

	h.CPUs = *cpus
	h.Memory = *mem
	h.VSock = *vsock
	if h.VSock {
		ports, err := stringToIntArray(*vsockports, ",")
		if err != nil {
			log.Fatalln("Unable to parse vsockports: ", err)
		}
		h.VSockPorts = ports
	}
	if *iso != "" {
		h.ISOImages = []string{*iso}
	}

	h.Disks = disks

	if *_9psock != "" {
		p := strings.Split(*_9psock, ",")
		if len(p) != 2 {
			log.Fatalln("9psock requires two parameters: path,tag")
		}
		h.Sockets9P = []hyperkit.Socket9P{{Path: p[0], Tag: p[1]}}
	}

	h.VPNKitUUID = *vpnkituuid
	h.VPNKitPreferredIPv4 = *vpnkitip

	if *bg {
		h.Console = hyperkit.ConsoleFile
		err = h.Start("console=ttyS0")

	} else {
		err = h.Run("console=ttyS0")
	}

	if err != nil {
		fmt.Println(h.Arguments)
		log.Fatalln("Error creating hyperkit: ", err)
	}

}

type disks []hyperkit.DiskConfig

func (d *disks) String() string {
	return fmt.Sprintf("%v", *d)
}

// Set parses a string with a disk configuration passed as a command line option.
func (d *disks) Set(v string) error {
	if v == "" {
		return fmt.Errorf("Empty disk config")
	}

	var config hyperkit.DiskConfig
	for _, kv := range strings.Split(v, ",") {
		p := strings.SplitN(kv, "=", 2)
		if len(p) == 1 { // Assume no key is a path
			if config.Path != "" {
				return fmt.Errorf("Invalid disk config, path already set")
			}
			config.Path = p[0]
		} else {
			switch p[0] {
			case "size":
				var err error
				if config.Size, err = strconv.Atoi(p[1]); err != nil {
					return err
				}
			case "file":
				config.Path = p[1]
			default:
				return fmt.Errorf("Unrecognised disk config key: %s", p[0])
			}
		}
	}
	*d = append(*d, config)
	return nil
}
