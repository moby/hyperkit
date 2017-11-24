package hyperkit

import (
	"os"
	"reflect"
	"testing"
)

// checkEqual allows to avoid importing testify.
func checkEqual(t *testing.T, expected, effective interface{}) {
	if !reflect.DeepEqual(expected, effective) {
		t.Errorf("FAIL:\n expected: %v\neffective: %v", expected, effective)
	}
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
		s, err := os.Stat(disk.Path)
		checkEqual(t, nil, err)
		checkEqual(t, 1024*1024, int(s.Size()))
	}
	// Rerunning resizes.
	disk.Size = 2
	checkEqual(t, nil, disk.Ensure())
	{
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
