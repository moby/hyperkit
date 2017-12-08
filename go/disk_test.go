package hyperkit

import (
	"os"
	"os/exec"
	"reflect"
	"testing"
)

// checkEqual allows to avoid importing testify.
func checkEqual(t *testing.T, expected, effective interface{}) {
	if !reflect.DeepEqual(expected, effective) {
		t.Errorf("FAIL:\n expected: %#v\neffective: %#v", expected, effective)
	}
}

func TestQcowDisk(t *testing.T) {
	if _, err := exec.LookPath("qcow-tool"); err != nil {
		t.Skip("cannot find qcow-tool: %v", err)
	}
	disk := QcowDisk{
		Path: "test.qcow",
		Size: 1,
	}
	os.Remove(disk.Path)

	checkEqual(t, "virtio-blk,file://test.qcow?sync=&buffered=1,format=qcow,qcow-config=discard=false;compact_after_unmaps=0;keep_erased=0;runtime_asserts=false", disk.AsArgument())
	checkEqual(t, nil, disk.Ensure())
	// Running twice is fine.
	checkEqual(t, nil, disk.Ensure())
	checkEqual(t, 1, disk.GetSize())
	// Rerunning resizes.
	disk.Size = 2
	checkEqual(t, nil, disk.Ensure())
	checkEqual(t, 2, disk.GetSize())
}

func TestRawDisk(t *testing.T) {
	disk := RawDisk{
		Path: "test.raw",
		Size: 1,
	}
	os.Remove(disk.Path)

	checkEqual(t, "virtio-blk,test.raw", disk.AsArgument())
	checkEqual(t, nil, disk.Ensure())
	// Running twice is fine.
	checkEqual(t, nil, disk.Ensure())
	{
		checkEqual(t, 1, disk.GetSize())
		s, err := os.Stat(disk.Path)
		checkEqual(t, nil, err)
		checkEqual(t, 1024*1024, int(s.Size()))
	}
	// Rerunning resizes.
	disk.Size = 2
	checkEqual(t, nil, disk.Ensure())
	{
		checkEqual(t, 2, disk.GetSize())
		s, err := os.Stat(disk.Path)
		checkEqual(t, nil, err)
		checkEqual(t, 2*1024*1024, int(s.Size()))
	}
}

func TestRawDiskTrim(t *testing.T) {
	disk := RawDisk{
		Path: "test.raw",
		Size: 1,
		Trim: true,
	}
	checkEqual(t, "ahci-hd,test.raw", disk.AsArgument())
}

func newDisk(t *testing.T, p string, s int) Disk {
	res, err := NewDisk(p, s)
	checkEqual(t, nil, err)
	return res
}

func TestNewDisk(t *testing.T) {
	{
		var ref Disk = &QcowDisk{Path: "/test.qcow2", Size: 1}
		checkEqual(t, ref, newDisk(t, "/test.qcow2", 1))
		checkEqual(t, ref, newDisk(t, "file:///test.qcow2", 1))
	}
	{
		var ref Disk = &RawDisk{Path: "/test.raw", Size: 1, Trim: true}
		checkEqual(t, ref, newDisk(t, "/test.raw", 1))
		checkEqual(t, ref, newDisk(t, "file:///test.raw", 1))
	}
}
