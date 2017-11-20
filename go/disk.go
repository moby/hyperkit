package hyperkit

import (
	"fmt"
	"os"
	"path/filepath"
)

// DiskConfig describes a disk image.
type DiskConfig struct {
	// Path specifies where the image file will be.
	Path string `json:"path"`
	// Size specifies the size of the disk image.  Used if the image needs to be created.
	Size int `json:"size"`
	// Format is passed as-is to the driver.
	Format string `json:"format"`
	// Driver is the name of the disk driver, "ahci-hd" or "virtio-blk".
	Driver string `json:"driver"`
}

// GetPath is the location of the disk image file.
func (d *DiskConfig) GetPath() string {
	return d.Path
}

// SetPath changes the location of the disk image file.
func (d *DiskConfig) SetPath(p string) {
	d.Path = p
}

// GetSize is the desired size of the disk image file.
func (d *DiskConfig) GetSize() int {
	return d.Size
}

// String returns the path.
func (d *DiskConfig) String() string {
	return d.Path
}

// Ensure create the image file if it does not exist.
func (d *DiskConfig) Ensure() error {
	if !d.exists() {
		return d.create()
	}
	return nil
}

// exists if the image file exists.
func (d *DiskConfig) exists() bool {
	_, err := os.Stat(d.Path)
	return err == nil
}

// create creates an empty file suitable for use as a disk image for a hyperkit VM.
func (d *DiskConfig) create() error {
	if d.Size == 0 {
		return fmt.Errorf("Disk image %s not found and unable to create it as size is not specified", d.Path)
	}
	diskDir := filepath.Dir(d.Path)
	if err := os.MkdirAll(diskDir, 0755); err != nil {
		return err
	}

	f, err := os.Create(d.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	return f.Truncate(int64(d.Size) * int64(1024) * int64(1024))
}

// AsArgument returns the command-line option to pass after `-s <slot>:0,` to hyperkit for this disk.
func (d *DiskConfig) AsArgument() string {
	res := fmt.Sprintf("%s,%s", defaultString(d.Driver, "virtio-blk"), d.Path)
	if d.Format != "" {
		res += ",format=" + d.Format
	}
	return res
}
