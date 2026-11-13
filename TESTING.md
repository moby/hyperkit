## Testing

Current testing is limited and under development. Proposals should be discussed
via github issues/pull requests following 

https://docs.docker.com/opensource/workflow/advanced-contributing/

## EPT MMIO instruction length regression

Run `make HAVE_OCAML_QCOW=NO test-ept-mmio` on an Intel Mac with Python 3 and the existing macOS
build tools. No downloaded kernel, disk image, network service, or OCaml
dependency is needed. The test builds a small reset ROM and uses the real
HyperKit EPT exit, decoder, and IOAPIC device paths. The ROM copies its test
payload to guest RAM so that the instruction emulator can fetch it.

The guest writes and reads an IOAPIC register using six-byte instructions,
checks the returned value, and powers off. Halt guards detect incorrect RIP
advancement. A test-only interposer supplies VMCS instruction lengths 0, 1, 2,
6, and 15 for these MMIO exits; other exits are unchanged. The native path is
also exercised without the interposer. Each case requires a success marker,
normal poweroff, and, for injected cases, both MMIO exits to have been reached.
The timeout is a failure. Logs and a JSON summary are in `build/ept-mmio`.

## Expect based test

The current integration test uses expect to drive it. Some tips for diagnosing tests for people new to expect.

* [A basic introductory guide to expect](https://gist.github.com/Fluidbyte/6294378)
* Turn on internal diagnostics - add `exp_internal 1` to the top of the expect script. This will allow you to see the matches expect is using and how it sees the strings.
* Interact with a running test - you can use the `interact` to give control of the current process to the user. 
