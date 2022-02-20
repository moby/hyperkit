// Package hyperkit provides a Go wrapper around the hyperkit
// command. It currently shells out to start hyperkit with the
// provided configuration.
//
// Most of the arguments should be self explanatory, but console
// handling deserves a mention. If the Console is configured with
// ConsoleStdio, the hyperkit is started with stdin, stdout, and
// stderr plumbed through to the VM console. If Console is set to
// ConsoleFile hyperkit the console output is redirected to a file and
// console input is disabled. For this mode StateDir has to be set and
// the interactive console is accessible via a 'tty' file created
// there.
//
// Currently this module has some limitations:
// - Only supports zero or one disk image
// - Only support zero or one network interface connected to VPNKit
// - Only kexec boot
//
// This package is currently implemented by shelling out a hyperkit
// process. In the future we may change this to become a wrapper
// around the hyperkit library.
package hyperkit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
)

const (
	// ConsoleStdio configures console to use Stdio (deprecated)
	ConsoleStdio = iota
	// ConsoleFile configures console to a tty and output to a file (deprecated)
	ConsoleFile
	// ConsoleLog configures console to a tty and sends its contents to the logs (deprecated)
	ConsoleLog

	legacyVPNKitSock  = "Library/Containers/com.docker.docker/Data/s50"
	defaultVPNKitSock = "Library/Containers/com.docker.docker/Data/vpnkit.eth.sock"

	defaultCPUs   = 1
	defaultMemory = 1024 // 1G

	defaultVSockGuestCID = 3

	jsonFile = "hyperkit.json"
	pidFile  = "hyperkit.pid"
)

var defaultHyperKits = []string{"hyperkit",
	"com.docker.hyperkit",
	"/usr/local/bin/hyperkit",
	"/Applications/Docker.app/Contents/Resources/bin/hyperkit",
	"/Applications/Docker.app/Contents/MacOS/com.docker.hyperkit"}

// Socket9P contains a unix domain socket path and 9p tag
type Socket9P struct {
	Path string `json:"path"`
	Tag  string `json:"tag"`
}

// HyperKit contains the configuration of the hyperkit VM
type HyperKit struct {
	// HyperKit is the path to the hyperkit binary.
	HyperKit string `json:"hyperkit"`
	// Argv0 is the name to declare as argv[0].  If left empty, argv[0] is left untouched.
	Argv0 string `json:"argv0"`

	// StateDir is the directory where runtime state is kept. If left empty, no state will be kept.
	StateDir string `json:"state_dir"`
	// VPNKitSock is the location of the VPNKit socket used for networking.
	VPNKitSock string `json:"vpnkit_sock"`
	// VPNKitUUID is a string containing a UUID, it can be used in conjunction with VPNKit to get consistent IP address.
	VPNKitUUID string `json:"vpnkit_uuid"`
	// VPNKitPreferredIPv4 is a string containing an IPv4 address, it can be used to request a specific IP for a UUID from VPNKit.
	VPNKitPreferredIPv4 string `json:"vpnkit_preferred_ipv4"`
	// UUID is a string containing a UUID, it sets BIOS DMI UUID for the VM (as found in /sys/class/dmi/id/product_uuid on Linux).
	UUID string `json:"uuid"`
	// Disks contains disk images to use/create.
	Disks []Disk `json:"disks"`
	// ISOImage is the (optional) path to a ISO image to attach.
	ISOImages []string `json:"iso"`

	// VSock enables the virtio-socket device and exposes it on the host.
	VSock bool `json:"vsock"`
	// VSockDir specifies where the unix domain sockets will be created. Defaults to StateDir when empty.
	VSockDir string `json:"vsock_dir"`
	// VSockPorts is a list of guest VSock ports that should be exposed as sockets on the host.
	VSockPorts []int `json:"vsock_ports"`
	// VSock guest CID
	VSockGuestCID int `json:"vsock_guest_cid"`

	// VMNet is whether to create vmnet network.
	VMNet bool `json:"vmnet"`

	// Sockets9P holds the 9P sockets.
	Sockets9P []Socket9P `json:"9p_sockets"`

	// Kernel is the path to the kernel image to boot.
	Kernel string `json:"kernel"`
	// Initrd is the path to the initial ramdisk to boot off.
	Initrd string `json:"initrd"`
	// Bootrom is the path to a boot rom eg for UEFI boot.
	Bootrom string `json:"bootrom"`

	// CPUs is the number CPUs to configure.
	CPUs int `json:"cpus"`
	// Memory is the amount of megabytes of memory for the VM.
	Memory int `json:"memory"`

	// Console defines where the console of the VM should be connected to. (deprecated)
	Console int `json:"console"`

	// Serials defines what happens to the I/O on the serial ports. If this is not nil
	// it overrides the Console setting.
	Serials []Serial `json:"serials"`

	// Below here are internal members, but they are exported so
	// that they are written to the state json file, if configured.

	// Pid of the hyperkit process
	Pid int `json:"pid"`
	// Arguments used to execute the hyperkit process
	Arguments []string `json:"arguments"`
	// CmdLine is a single string of the command line
	CmdLine string `json:"cmdline"`

	process *os.Process
}

// Serial port.
type Serial struct {
	// InteractiveConsole allows a user to connect to a live VM serial console.
	InteractiveConsole InteractiveConsole
	// LogToRingBuffer will write console output to a fixed size ring buffer file.
	LogToRingBuffer bool
	// LogToASL will write console output to the Apple System Log.
	LogToASL bool
}

// InteractiveConsole is an optional interactive VM console.
type InteractiveConsole int

const (
	// NoInteractiveConsole disables the interactive console.
	NoInteractiveConsole = InteractiveConsole(iota)
	// StdioInteractiveConsole creates a console on stdio.
	StdioInteractiveConsole
	// TTYInteractiveConsole creates a console on a TTY.
	TTYInteractiveConsole
)

// New creates a template config structure.
// - If hyperkit can't be found an error is returned.
// - If vpnkitsock is empty no networking is configured. If it is set
//   to "auto" it tries to re-use the Docker for Mac VPNKit connection.
// - If statedir is "" no state is written to disk.
func New(hyperkit, vpnkitsock, statedir string) (*HyperKit, error) {
	h := HyperKit{}
	var err error

	h.HyperKit, err = checkHyperKit(hyperkit)
	if err != nil {
		return nil, err
	}
	h.StateDir = statedir
	h.VPNKitSock, err = checkVPNKitSock(vpnkitsock)
	if err != nil {
		return nil, err
	}

	h.CPUs = defaultCPUs
	h.Memory = defaultMemory

	h.VSockGuestCID = defaultVSockGuestCID

	h.Console = ConsoleStdio

	return &h, nil
}

// Run the VM with a given command line until it exits.
func (h *HyperKit) Run(cmdline string) error {
	errCh, err := h.Start(cmdline)
	if err != nil {
		return err
	}
	return <-errCh
}

// Start the VM with a given command line in the background.  On
// success, returns a channel on which the result of Wait'ing for
// hyperkit is sent.  Return failures to start the process.
func (h *HyperKit) Start(cmdline string) (chan error, error) {
	log.Debugf("hyperkit: Start %#v", h)
	if err := h.check(); err != nil {
		return nil, err
	}
	h.buildArgs(cmdline)
	cmd, err := h.execute()
	if err != nil {
		return nil, err
	}
	errCh := make(chan error, 1)
	// Make sure we reap the child when it exits
	go func() {
		log.Debugf("hyperkit: Waiting for %#v", cmd)
		errCh <- cmd.Wait()
	}()
	return errCh, nil
}

func (h *HyperKit) checkLegacyConsole() error {
	switch h.Console {
	case ConsoleFile, ConsoleLog:
		if h.StateDir == "" {
			return fmt.Errorf("If ConsoleFile or ConsoleLog is set, StateDir must be specified")
		}
	case ConsoleStdio:
		if !isTerminal(os.Stdout) && h.StateDir == "" {
			return fmt.Errorf("If ConsoleStdio is set but stdio is not a terminal, StateDir must be specified")
		}
	}
	return nil
}

func (h *HyperKit) checkSerials() error {
	stdioConsole := -1
	for i, serial := range h.Serials {
		if serial.LogToRingBuffer && h.StateDir == "" {
			return fmt.Errorf("If VM is to log to a ring buffer, StateDir must be specified")
		}
		if serial.InteractiveConsole == StdioInteractiveConsole {
			if !isTerminal(os.Stdout) {
				return fmt.Errorf("If StdioInteractiveConsole is set, stdio must be a TTY")
			}
			if stdioConsole != -1 {
				return fmt.Errorf("Only one serial port can be nominated as the stdio interactive console")
			}
			stdioConsole = i
		}
		if serial.InteractiveConsole == TTYInteractiveConsole && h.StateDir == "" {
			return fmt.Errorf("If TTYInteractiveConsole is set, StateDir must be specified ")
		}
		if serial.LogToRingBuffer && h.StateDir == "" {
			return fmt.Errorf("If LogToRingBuffer is set, StateDir must be specified")
		}
	}
	return nil
}

// check validates `h`.  It also creates the disks if needed.
func (h *HyperKit) check() error {
	log.Debugf("hyperkit: check %#v", h)
	var err error
	// Sanity checks on configuration
	if h.Serials == nil {
		if err := h.checkLegacyConsole(); err != nil {
			return err
		}
	} else {
		if err := h.checkSerials(); err != nil {
			return err
		}
	}
	for _, image := range h.ISOImages {
		if _, err = os.Stat(image); os.IsNotExist(err) {
			return fmt.Errorf("ISO %s does not exist", image)
		}
	}

	if h.VSock {
		if h.VSockDir == "" {
			if h.StateDir == "" {
				return fmt.Errorf("If virtio-sockets are enabled, VSockDir or StateDir must be specified")
			}
			h.VSockDir = h.StateDir
		}
	} else if len(h.VSockPorts) > 0 {
		return fmt.Errorf("To forward vsock ports VSock must be enabled")
	}

	if h.Bootrom == "" {
		if _, err = os.Stat(h.Kernel); os.IsNotExist(err) {
			return fmt.Errorf("Kernel %s does not exist", h.Kernel)
		}
		if h.Initrd != "" {
			if _, err = os.Stat(h.Initrd); os.IsNotExist(err) {
				return fmt.Errorf("initrd %s does not exist", h.Initrd)
			}
		}
	} else {
		if _, err = os.Stat(h.Bootrom); os.IsNotExist(err) {
			return fmt.Errorf("Bootrom %s does not exist", h.Bootrom)
		}
	}
	if h.VPNKitPreferredIPv4 != "" {
		if ip := net.ParseIP(h.VPNKitPreferredIPv4); ip == nil {
			return fmt.Errorf("Invalid VPNKit IP: %s", h.VPNKitPreferredIPv4)
		}
	}

	// Create directories.
	for _, dir := range []string{h.VSockDir, h.StateDir} {
		if dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("Cannot create directory %q: %v", dir, err)
			}
		}
	}

	for idx, disk := range h.Disks {
		if disk.GetPath() == "" {
			if h.StateDir == "" {
				return fmt.Errorf("Unable to create disk image when neither path nor state dir is set")
			}
			if disk.GetSize() <= 0 {
				return fmt.Errorf("Unable to create disk image when size is 0 or not set")
			}
			disk.SetPath(filepath.Clean(filepath.Join(h.StateDir, fmt.Sprintf("disk%02d.img", idx))))
			h.Disks[idx] = disk
		}
		if err := disk.Ensure(); err != nil {
			return err
		}
	}
	return nil
}

// Stop the running VM
func (h *HyperKit) Stop() error {
	if h.process == nil {
		return fmt.Errorf("hyperkit process not known")
	}
	if !h.IsRunning() {
		return nil
	}
	err := h.process.Kill()
	if err != nil {
		return err
	}

	return nil
}

// IsRunning returns true if the hyperkit process is running.
func (h *HyperKit) IsRunning() bool {
	// os.FindProcess on Unix always returns a process object even
	// if the process does not exists. There does not seem to be
	// a call to find out if the process is running either, so we
	// use another package to find out.
	proc, err := ps.FindProcess(h.Pid)
	if err != nil {
		return false
	}
	if proc == nil {
		return false
	}
	return true
}

// isDisk checks if the specified path is used as a disk image
func (h *HyperKit) isDisk(path string) bool {
	for _, disk := range h.Disks {
		if filepath.Clean(path) == filepath.Clean(disk.GetPath()) {
			return true
		}
	}
	return false
}

// Remove deletes all statefiles if present.
// This also removes the StateDir if empty.
// If keepDisk is set, the disks will not get removed.
func (h *HyperKit) Remove(keepDisk bool) error {
	if h.IsRunning() {
		return fmt.Errorf("Can't remove state as process is running")
	}
	if h.StateDir == "" {
		// If there is not state directory we don't mess with the disk
		return nil
	}

	if !keepDisk {
		return os.RemoveAll(h.StateDir)
	}

	files, _ := ioutil.ReadDir(h.StateDir)
	for _, f := range files {
		fn := filepath.Clean(filepath.Join(h.StateDir, f.Name()))
		if !h.isDisk(fn) {
			if err := os.Remove(fn); err != nil {
				return err
			}
		}
	}
	return nil
}

// Convert to json string
func (h *HyperKit) String() string {
	s, err := json.Marshal(h)
	if err != nil {
		return err.Error()
	}
	return string(s)
}

func intArrayToString(i []int, sep string) string {
	if len(i) == 0 {
		return ""
	}
	s := make([]string, len(i))
	for idx := range i {
		s[idx] = strconv.Itoa(i[idx])
	}
	return strings.Join(s, sep)
}

func (h *HyperKit) legacyConsoleArgs() []string {
	cfg := "com1"
	if h.Console == ConsoleStdio && isTerminal(os.Stdout) {
		cfg += fmt.Sprintf(",stdio")
	} else {
		cfg += fmt.Sprintf(",autopty=%s/tty", h.StateDir)
	}
	if h.Console == ConsoleLog {
		cfg += fmt.Sprintf(",asl")
	} else {
		cfg += fmt.Sprintf(",log=%s/console-ring", h.StateDir)
	}
	return []string{"-l", cfg}
}

func (h *HyperKit) serialArgs() []string {
	results := []string{}
	for i, serial := range h.Serials {
		cfg := fmt.Sprintf("com%d", i+1)
		switch serial.InteractiveConsole {
		case NoInteractiveConsole:
			cfg += ",null"
		case StdioInteractiveConsole:
			cfg += fmt.Sprintf(",stdio")
		case TTYInteractiveConsole:
			cfg += fmt.Sprintf(",autopty=%s/tty%d", h.StateDir, i+1)
		}
		if serial.LogToASL {
			cfg += fmt.Sprintf(",asl")
		}
		if serial.LogToRingBuffer {
			cfg += fmt.Sprintf(",log=%s/console-ring", h.StateDir)
		}
		results = append(results, "-l", cfg)
	}
	return results
}

func (h *HyperKit) buildArgs(cmdline string) {
	a := []string{"-A", "-u"}
	if h.StateDir != "" {
		a = append(a, "-F", filepath.Join(h.StateDir, pidFile))
	}

	a = append(a, "-c", fmt.Sprintf("%d", h.CPUs))
	a = append(a, "-m", fmt.Sprintf("%dM", h.Memory))

	a = append(a, "-s", "0:0,hostbridge")
	a = append(a, "-s", "31,lpc")

	nextSlot := 1

	if h.VPNKitSock != "" {
		var uuid string
		if h.VPNKitUUID != "" {
			uuid = fmt.Sprintf(",uuid=%s", h.VPNKitUUID)
		}
		var preferredIPv4 string
		if h.VPNKitPreferredIPv4 != "" {
			preferredIPv4 = fmt.Sprintf(",preferred_ipv4=%s", h.VPNKitPreferredIPv4)
		}
		a = append(a, "-s", fmt.Sprintf("%d:0,virtio-vpnkit,path=%s%s%s", nextSlot, h.VPNKitSock, uuid, preferredIPv4))
		nextSlot++
	}

	if h.VMNet {
		a = append(a, "-s", fmt.Sprintf("%d:0,virtio-net", nextSlot))
		nextSlot++
	}

	if h.UUID != "" {
		a = append(a, "-U", h.UUID)
	}

	for _, disk := range h.Disks {
		a = append(a, "-s", fmt.Sprintf("%d:0,%s", nextSlot, disk.AsArgument()))
		nextSlot++
	}

	if h.VSock {
		l := fmt.Sprintf("%d,virtio-sock,guest_cid=%d,path=%s", nextSlot, h.VSockGuestCID, h.VSockDir)
		if len(h.VSockPorts) > 0 {
			l = fmt.Sprintf("%s,guest_forwards=%s", l, intArrayToString(h.VSockPorts, ";"))
		}
		a = append(a, "-s", l)
		nextSlot++
	}

	for _, image := range h.ISOImages {
		a = append(a, "-s", fmt.Sprintf("%d,ahci-cd,%s", nextSlot, image))
		nextSlot++
	}

	a = append(a, "-s", fmt.Sprintf("%d,virtio-rnd", nextSlot))
	nextSlot++

	for _, p := range h.Sockets9P {
		a = append(a, "-s", fmt.Sprintf("%d,virtio-9p,path=%s,tag=%s", nextSlot, p.Path, p.Tag))
		nextSlot++
	}

	// -l: LPC device configuration.
	if h.Serials == nil {
		a = append(a, h.legacyConsoleArgs()...)
	} else {
		a = append(a, h.serialArgs()...)
	}

	if h.Bootrom == "" {
		kernArgs := fmt.Sprintf("kexec,%s,%s,earlyprintk=serial %s", h.Kernel, h.Initrd, cmdline)
		a = append(a, "-f", kernArgs)
	} else {
		kernArgs := fmt.Sprintf("bootrom,%s,,", h.Bootrom)
		a = append(a, "-f", kernArgs)
	}

	h.Arguments = a
	h.CmdLine = h.HyperKit + " " + strings.Join(a, " ")
	log.Debugf("hyperkit: Arguments: %#v", h.Arguments)
	log.Debugf("hyperkit: CmdLine: %#v", h.CmdLine)
}

// openTTY opens the tty files for reading, and returns it.
func (h *HyperKit) openTTY(filename string) *os.File {
	path := path.Join(h.StateDir, filename)
	for {
		if res, err := os.OpenFile(path, os.O_RDONLY, 0); err != nil {
			log.Infof("hyperkit: openTTY: %v, retrying", err)
			time.Sleep(10 * time.Millisecond)
		} else {
			log.Infof("hyperkit: openTTY: got %v", path)
			saneTerminal(res)
			setRaw(res)
			return res
		}
	}
}

func (h *HyperKit) findStdioTTY() string {
	if h.Serials != nil {
		for i, serial := range h.Serials {
			if serial.InteractiveConsole == StdioInteractiveConsole {
				return fmt.Sprintf("tty%d", i)
			}
		}
		return ""
	}
	if h.Console == ConsoleStdio {
		return "tty"
	}
	return ""
}

// execute forges the command to run hyperkit, runs and returns it.
// It also plumbs stdin/stdout/stderr.
func (h *HyperKit) execute() (*exec.Cmd, error) {

	cmd := exec.Command(h.HyperKit, h.Arguments...)
	if h.Argv0 != "" {
		cmd.Args[0] = h.Argv0
	}
	cmd.Env = os.Environ()

	// Plumb in stdin/stdout/stderr.
	//
	// If ConsoleStdio is configured and we are on a terminal,
	// just plugin stdio. If we are not on a terminal we have a
	// StateDir (as per checks above) and have configured HyperKit
	// to use a PTY in the statedir. In this case, we just open
	// the PTY slave and copy it to stdout (and ignore
	// stdin). This allows for redirecting of the VM output, used
	// for testing.
	//
	// If a logger is specified, use it for stdout/stderr
	// logging. Otherwise use the default /dev/null.
	filename := h.findStdioTTY()
	if filename != "" {
		if isTerminal(os.Stdout) {
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else {
			go func() {
				tty := h.openTTY(filename)
				defer tty.Close()
				io.Copy(os.Stdout, tty)
			}()
		}
	} else if log != nil {
		log.Debugf("hyperkit: Redirecting stdout/stderr to logger")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return nil, err
		}
		go logStream(stdout, "stdout")
		go logStream(stderr, "stderr")
	}

	log.Debugf("hyperkit: Starting %#v", cmd)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	h.Pid = cmd.Process.Pid
	log.Debugf("hyperkit: Pid is %v", h.Pid)
	h.process = cmd.Process
	err = h.writeState()
	if err != nil {
		log.Debugf("hyperkit: Cannot write state: %v, killing %v", err, h.Pid)
		h.process.Kill()
		return nil, err
	}
	return cmd, nil
}

// writeState write the state to a JSON file
func (h *HyperKit) writeState() error {
	if h.StateDir == "" {
		// This is not an error
		return nil
	}

	s, err := json.Marshal(h)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(h.StateDir, jsonFile), []byte(s), 0644)
}

// logStream redirects a file to the logs.
func logStream(r io.ReadCloser, name string) {
	defer r.Close()
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Warnf("hyperkit: failed to read %v: %v", name, err)
			break
		}
		log.Infof("hyperkit: %v: %v", name, line)
	}
}

// checkHyperKit tries to find and/or validate the path of hyperkit
func checkHyperKit(hyperkit string) (string, error) {
	if hyperkit != "" {
		p, err := exec.LookPath(hyperkit)
		if err != nil {
			return "", fmt.Errorf("Could not find hyperkit executable %s: %s", hyperkit, err)
		}
		return p, nil
	}

	// Look in a number of default locations
	for _, hyperkit := range defaultHyperKits {
		p, err := exec.LookPath(hyperkit)
		if err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("Could not find hyperkit executable")
}

// checkVPNKitSock tries to find and/or validate the path of the VPNKit socket
func checkVPNKitSock(vpnkitsock string) (string, error) {
	if vpnkitsock == "auto" {
		vpnkitsock = filepath.Join(getHome(), defaultVPNKitSock)
		if _, err := os.Stat(vpnkitsock); err != nil {
			vpnkitsock = filepath.Join(getHome(), legacyVPNKitSock)
		}
	}
	if vpnkitsock == "" {
		return "", nil
	}

	vpnkitsock = filepath.Clean(vpnkitsock)
	_, err := os.Stat(vpnkitsock)
	if err != nil {
		return "", err
	}
	return vpnkitsock, nil
}

func getHome() string {
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}
