## [HyperKit](http://github.com/moby/hyperkit)

![Build Status macOS](https://circleci.com/gh/moby/hyperkit.svg?style=shield&circle-token=cf8379b302eab2bbf33821cafe164dbefb71982d)

*HyperKit* is a toolkit for embedding hypervisor capabilities in your application. It includes a complete hypervisor, based on [xhyve](https://github.com/mist64/xhyve)/[bhyve](http://bhyve.org), which is optimized for lightweight virtual machines and container deployment.  It is designed to be interfaced with higher-level components such as the [VPNKit](https://github.com/moby/vpnkit) and [DataKit](https://github.com/moby/datakit).

HyperKit currently only supports macOS using the [Hypervisor.framework](https://developer.apple.com/library/mac/documentation/DriversKernelHardware/Reference/Hypervisor/index.html). It is a core component of Docker Desktop for Mac.


## Requirements

* OS X 10.10.3 Yosemite or later
* a 2010 or later Mac (i.e. a CPU that supports EPT)

## Reporting Bugs

If you are using a version of Hyperkit which is embedded into a higher level application (e.g. [Docker Desktop for Mac](https://github.com/docker/for-mac)) then please report any issues against that higher level application in the first instance. That way the relevant team can triage and determine if the issue lies in Hyperkit and assign as necessary.

If you are using Hyperkit directly then please report issues against this repository.

## Usage

    $ hyperkit -h

## Building

    $ git clone https://github.com/moby/hyperkit
    $ cd hyperkit
    $ make

The resulting binary will be in `build/hyperkit`

To enable qcow support in the block backend an OCaml [OPAM](https://opam.ocaml.org) development
environment is required with the qcow module available. A
suitable environment can be setup by installing `opam` and `libev`
via `brew` and using `opam` to install the appropriate libraries:

    $ brew install opam libev
    $ opam init
    $ eval `opam env`
    $ opam pin add qcow.0.11.0 git://github.com/mirage/ocaml-qcow -n
    $ opam pin add qcow-tool.0.11.0 git://github.com/mirage/ocaml-qcow -n
    $ opam install uri qcow.0.11.0 conduit.2.1.0 lwt.5.3.0 qcow-tool mirage-block-unix.2.12.0 conf-libev logs fmt mirage-unix prometheus-app

Notes:

- `opam config env` must be evaluated each time prior to building
  hyperkit so the build will find the ocaml environment.
- Any previous pin of `mirage-block-unix` or `qcow`
  should be removed with the commands:
  
  ```sh
  $ opam update
  $ opam pin remove mirage-block-unix
  $ opam pin remove qcow
  ```

## Tracing

HyperKit defines a number of static DTrace probes to simplify investigation of
performance problems. To list the probes supported by your version of HyperKit,
type the following command while HyperKit VM is running:

     $ sudo dtrace -l -P 'hyperkit$target' -p $(pgrep hyperkit)

Refer to scripts in dtrace/ directory for examples of possible usage and
available probes.

### Relationship to xhyve and bhyve

HyperKit includes a hypervisor derived from [xhyve](http://www.xhyve.org), which in turn
was derived from [bhyve](http://www.bhyve.org). See the [original xhyve
README](README.xhyve.md) which incorporates the bhyve README.

We try to avoid deviating from these upstreams unnecessarily in order
to more easily share code, for example the various device
models/emulations should be easily shareable.

### Reporting security issues

The maintainers take security seriously. If you discover a security issue,
please bring it to their attention right away!

Please **DO NOT** file a public issue, instead send your report privately to
[security@docker.com](mailto:security@docker.com).

Security reports are greatly appreciated and we will publicly thank you for it.
We also like to send gifts&mdash;if you're into Docker schwag, make sure to let
us know. We currently do not offer a paid security bounty program, but are not
ruling it out in the future.
