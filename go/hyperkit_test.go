package hyperkit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyConsole(t *testing.T) {
	h, err := New("sh", "", "state-dir")
	require.Nil(t, err)

	h.Console = ConsoleFile
	h.buildArgs("")
	assert.EqualValues(t, []string{"-A", "-u", "-F", "state-dir/hyperkit.pid", "-c", "1", "-m", "1024M", "-s", "0:0,hostbridge", "-s", "31,lpc", "-s", "1,virtio-rnd", "-l", "com1,autopty=state-dir/tty,log=state-dir/console-ring", "-f", "kexec,,,earlyprintk=serial "}, h.Arguments)
}

func TestNewSerial(t *testing.T) {
	h, err := New("sh", "", "state-dir")
	require.Nil(t, err)

	h.Serials = []Serial{
		{
			InteractiveConsole: TTYInteractiveConsole,
			LogToRingBuffer:    true,
		},
	}
	h.buildArgs("")
	assert.EqualValues(t, []string{"-A", "-u", "-F", "state-dir/hyperkit.pid", "-c", "1", "-m", "1024M", "-s", "0:0,hostbridge", "-s", "31,lpc", "-s", "1,virtio-rnd", "-l", "com1,autopty=state-dir/tty1,log=state-dir/console-ring", "-f", "kexec,,,earlyprintk=serial "}, h.Arguments)
}

func TestNullSerial(t *testing.T) {
	h, err := New("sh", "", "state-dir")
	require.Nil(t, err)

	h.Serials = []Serial{
		{
			LogToRingBuffer: true,
		},
	}
	h.buildArgs("")
	assert.EqualValues(t, []string{"-A", "-u", "-F", "state-dir/hyperkit.pid", "-c", "1", "-m", "1024M", "-s", "0:0,hostbridge", "-s", "31,lpc", "-s", "1,virtio-rnd", "-l", "com1,null,log=state-dir/console-ring", "-f", "kexec,,,earlyprintk=serial "}, h.Arguments)
}
