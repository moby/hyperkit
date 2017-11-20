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
// - Only support zero or one network interface connected to VPNKit
//
// This package is currently implemented by shelling out a hyperkit
// process. In the future we may change this to become a wrapper
// around the hyperkit library.
package hyperkit
