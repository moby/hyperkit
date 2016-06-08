#!/bin/sh -e

# These are binaries from a mirror of
#  http://tinycorelinux.net
# with the following patch applied:
# Upstream source is available http://www.tinycorelinux.net/6.x/x86/release/src/
#BASE_URL="http://www.tinycorelinux.net/"

BASE_URL="http://distro.ibiblio.org/tinycorelinux/"

echo Downloading tinycore linux
curl -O -s "${BASE_URL}/6.x/x86/release/distribution_files/vmlinuz64"
mv vmlinuz64 vmlinuz
curl -O -s "${BASE_URL}/6.x/x86/release/distribution_files/core.gz"
mv core.gz initrd.gz
echo Patching tinycore linux initrd - may prompt for password for sudo
mkdir initrd
( cd initrd ; gzcat ../initrd.gz | sudo cpio -idm )
sudo sed -i -e '/^# ttyS0$/s#^..##' initrd/etc/securetty 
sudo sed -i -e '/^tty1:/s#tty1#ttyS0#g' initrd/etc/inittab
( cd initrd ; find . | sudo cpio -o -H newc ) | gzip -c > initrd.gz && sudo rm -rf ./initrd
