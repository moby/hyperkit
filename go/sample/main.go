// +build darwin

// This mostly serves as sample and test code for the hyperkit package
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/docker/hyperkit/go"
)

func main() {
	hk := flag.String("hyperkit", "", "HyperKit binary to use")
	statedir := flag.String("state", "", "Directory to keep state in")
	vpnkitsock := flag.String("vpnkitsock", "auto", "Path to VPNKit socket")
	disk := flag.String("disk", "", "Path to disk image")

	kernel := flag.String("kernel", "", "Kernel to boot")
	initrd := flag.String("initrd", "", "Initial RAM Disk")

	bg := flag.Bool("background", false, "Start VM in the background")

	cpus := flag.Int("cpus", 1, "Number of CPUs")
	mem := flag.Int("mem", 1024, "Amount of memory in MB")
	diskSz := flag.Int("disk-size", 0, "Size of Disk in MB")
	vsock := flag.Bool("vsock", false, "Enable virtio-sockets")

	data := flag.String("data", "", "User data to pass to VM via ISO")

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
			err = h.Remove(*disk != "")
			if err != nil {
				log.Fatalln("Error removing state: ", err)
			}
			return
		default:
			log.Fatalln("Unknown command: ", cmd[0])
		}
	}

	h, err := hyperkit.New(*hk, *statedir, *vpnkitsock, *disk)
	if err != nil {
		log.Fatalln("Error creating hyperkit: ", err)
	}

	h.Kernel = *kernel
	h.Initrd = *initrd

	h.CPUs = *cpus
	h.Memory = *mem
	h.DiskSize = *diskSz
	h.VSock = *vsock

	h.UserData = *data

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
