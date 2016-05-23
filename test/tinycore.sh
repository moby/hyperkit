#!/bin/sh

# These are binaries from
#  http://tinycorelinux.net
# with the following patch applied:
# Source is available http://www.tinycorelinux.net/6.x/x86/release/src/

echo Downloading tinycore linux
curl -s -o vmlinuz http://www.tinycorelinux.net/6.x/x86/release/distribution_files/vmlinuz64
curl -s -o initrd.gz http://www.tinycorelinux.net/6.x/x86/release/distribution_files/core.gz
echo Patching tinycore linux initrd - may prompt for password for sudo
mkdir initrd
( cd initrd ; gzcat ../initrd.gz | sudo cpio -idm )
sudo sed -i -e '/^# ttyS0$/s#^..##' initrd/etc/securetty 
sudo sed -i -e '/^tty1:/s#tty1#ttyS0#g' initrd/etc/inittab
( cd initrd ; find . | sudo cpio -o -H newc ) | gzip -c > initrd.gz && sudo rm -rf ./initrd
