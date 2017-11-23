package hyperkit

import (
	"fmt"
	"os"
)

const (
	mib = int64(1024 * 1024)
)

/*-------.
| Disk.  |
`-------*/

// Disk in an interface for qcow2 and raw disk images.
type Disk interface {
	// GetPath returns the location of the disk image file.
	GetPath() string
	// SetPath changes the location of the disk image file.
	SetPath(p string)
	// GetSize returns the desired size of the disk image file.
	GetSize() int
	// String returns the path.
	String() string

	// Exists is true iff the disk image file exists.
	Exists() bool
	// Ensure creates the disk image if needed, and resizes it if needed.
	Ensure() error
	// Stop can be called when hyperkit has quit.  It performs sanity checks, compaction, etc.
	// It can be passed a lock file which is not-used, but kept alive even if the parent
	// process died.
	Stop(lockFile *os.File) error

	// AsArgument returns the command-line option to pass after `-s <slot>:0,` to hyperkit for this disk.
	AsArgument() string

	create() error
	getFileSize() (int, error)
	resize() error
}

// exists if the image file exists.
func exists(d Disk) bool {
	_, err := os.Stat(d.GetPath())
	return err == nil
}

// ensure creates the disk image if needed, and resizes it if needed.
func ensure(d Disk) error {
	current, err := d.getFileSize()
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return d.create()
	}
	if current < d.GetSize() {
		log.Infof("Attempting to resize %q from %dMiB to %dMiB", d, current, d.GetSize())
		return d.resize()
	}
	if d.GetSize() < current {
		log.Errorf("Cannot safely shrink %q from %dMiB to %dMiB", d, current, d.GetSize())
	}
	return nil
}

// diskDriver to use.
//
// Dave Scott writes:
//
// > Regarding TRIM and raw disks
// > (https://github.com/docker/pinata/pull/8235/commits/0e2c7c2e21114b4ed61589bd42b720f7d88c0d8e):
// > it works like this: the `ahci-hd` virtual hardware in hyperkit
// > exposes the `ATA_SUPPORT_DSM_TRIM` capability
// > (https://github.com/moby/hyperkit/blob/81fa6279fcb17e8435f3cec0978e9aa3af02e63b/src/lib/pci_ahci.c#L996)
// > if the `fcntl(F_PUNCHHOLE)`
// > (https://github.com/moby/hyperkit/blob/81fa6279fcb17e8435f3cec0978e9aa3af02e63b/src/lib/block_if.c#L276)
// > API works on the raw file (it's dynamically detected so on HFS+ it's
// > disabled and on APFS it's enabled) -> TRIM on raw doesn't need any
// > special flags set in the Go code; the special flags are only for the
// > TRIM on qcow implementation. When images are deleted in the VM the
// > `trim-after-delete`
// > (https://github.com/linuxkit/linuxkit/tree/master/pkg/trim-after-delete)
// > daemon calls `fstrim /var/lib/docker` which causes Linux to emit the
// > TRIM commands to hyperkit, which calls `fcntl`, which tells macOS to
// > free the space in the file, visible in `ls -sl`.
// >
// > Unfortunately the `virtio-blk` protocol doesn't support `TRIM`
// > requests at all so we have to use `ahci-hd` (if you try to run
// > `fstrim /var/lib/docker` with `virtio-blk` it'll give an `ioctl`
// > error).
func diskDriver(trim bool) string {
	if trim {
		return "ahci-hd"
	}
	return "virtio-blk"
}

/*----------.
| RawDisk.  |
`----------*/

// RawDisk describes a raw disk image file.
type RawDisk struct {
	// Path specifies where the image file will be.
	Path string `json:"path"`
	// Size specifies the size of the disk image.  Used if the image needs to be created.
	Size int `json:"size"`
	// Format is passed as-is to the driver.
	Format string `json:"format"`
	// Trim specifies whether we should trim the image file.
	Trim bool `json:"trim"`
}

// GetPath returns the location of the disk image file.
func (d *RawDisk) GetPath() string {
	return d.Path
}

// SetPath changes the location of the disk image file.
func (d *RawDisk) SetPath(p string) {
	d.Path = p
}

// GetSize returns the desired size of the disk image file.
func (d *RawDisk) GetSize() int {
	return d.Size
}

// String returns the path.
func (d *RawDisk) String() string {
	return d.Path
}

// Exists if the image file exists.
func (d *RawDisk) Exists() bool {
	return exists(d)
}

// Ensure creates the disk image if needed, and resizes it if needed.
func (d *RawDisk) Ensure() error {
	return ensure(d)
}

// Create a disk.
func (d *RawDisk) create() error {
	f, err := os.Create(d.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	return f.Truncate(int64(d.Size) * mib)
}

// Query the current virtual size of the disk in MiB
func (d *RawDisk) getFileSize() (int, error) {
	// Return a failure if the file doesn't exist yet
	fileinfo, err := os.Stat(d.Path)
	if err != nil {
		return 0, err
	}
	return int(fileinfo.Size() / mib), nil
}

// Resize the virtual size of the disk
func (d *RawDisk) resize() error {
	return os.Truncate(d.Path, int64(d.Size)*mib)
}

// Stop cleans up this disk when we are quitting.
func (d *RawDisk) Stop(lockFile *os.File) error {
	return nil
}

// AsArgument returns the command-line option to pass after `-s <slot>:0,` to hyperkit for this disk.
func (d *RawDisk) AsArgument() string {
	res := fmt.Sprintf("%s,%s", diskDriver(d.Trim), d.Path)
	if d.Format != "" {
		res += ",format=" + d.Format
	}
	return res
}
