#!/usr/bin/env python3
"""Exercise the real EPT decoder and IOAPIC device with stale VMCS lengths."""
import json
import os
from pathlib import Path
import pty
import resource
import select
import struct
import subprocess
import sys
import threading


def extract_rom(obj, destination):
    data = obj.read_bytes()
    magic, cpu, _, _, commands, _, _, _ = struct.unpack_from('<IiiIIIII', data)
    if magic != 0xfeedfacf or cpu != 0x1000007:
        raise RuntimeError('Expected an x86_64 Mach-O object')
    offset = 32
    for _ in range(commands):
        command, size = struct.unpack_from('<II', data, offset)
        if command == 0x19:
            sections = struct.unpack_from('<I', data, offset + 64)[0]
            for section in range(sections):
                at = offset + 72 + 80 * section
                if data[at:at + 16].rstrip(b'\0') == b'__text':
                    length, position = struct.unpack_from('<QI', data, at + 40)
                    relocations = struct.unpack_from('<I', data, at + 60)[0]
                    if length != 65536 or relocations:
                        raise RuntimeError('ROM must be 64 KiB and relocation-free')
                    destination.write_bytes(data[position:position + length])
                    return
        offset += size
    raise RuntimeError('No ROM text section found')


def run_case(binary, rom, library, length, output):
    master, slave = pty.openpty()
    console = bytearray()
    done = threading.Event()

    def read_console():
        while not done.is_set():
            if select.select([master], [], [], 0.1)[0]:
                try:
                    chunk = os.read(master, 65536)
                except OSError:
                    return
                if not chunk:
                    return
                console.extend(chunk)

    reader = threading.Thread(target=read_console)
    env = dict(os.environ)
    env.pop('DYLD_INSERT_LIBRARIES', None)
    env.pop('HYPERKIT_TEST_EPT_LENGTH', None)
    if length is not None:
        env['DYLD_INSERT_LIBRARIES'] = str(library)
        env['HYPERKIT_TEST_EPT_LENGTH'] = str(length)
    args = [str(binary), '-A', '-m', '16M', '-c', '1',
            '-s', '0:0,hostbridge', '-s', '31,lpc', '-l', 'com1,stdio',
            '-f', 'bootrom,' + str(rom)]
    process = subprocess.Popen(args, env=env, stdin=slave, stdout=subprocess.PIPE,
                               stderr=subprocess.STDOUT)
    reader.start()
    timed_out = False
    try:
        try:
            log, _ = process.communicate(timeout=10)
        except subprocess.TimeoutExpired:
            timed_out = True
            process.kill()
            log, _ = process.communicate()
    finally:
        os.close(slave)
        reader.join(timeout=1)
        done.set()
        reader.join()
        os.close(master)
    name = 'native' if length is None else str(length)
    (output / ('length-' + name + '.log')).write_bytes(log)
    (output / ('length-' + name + '.serial')).write_bytes(console)
    injected = log.count(b'EPT_TEST_LENGTH=')
    passed = (not timed_out and process.returncode == 0 and
              b'EPT_MMIO_PASS' in console and b'EPT_MMIO_FAIL' not in console and
              (length is None or injected >= 2))
    result = {'length': name, 'passed': passed, 'exit': process.returncode,
              'timeout': timed_out, 'injected_mmio_exits': injected}
    print(json.dumps(result), flush=True)
    if not passed:
        print(log.decode(errors='replace')[-1800:], flush=True)
        print(console.decode(errors='replace')[-800:], flush=True)
    return result


def main():
    resource.setrlimit(resource.RLIMIT_CORE, (0, 0))
    root = Path(__file__).resolve().parent.parent
    binary = Path(sys.argv[1]).resolve() if len(sys.argv) > 1 else root / 'build/hyperkit'
    output = root / 'build/ept-mmio'
    output.mkdir(parents=True, exist_ok=True)
    obj = output / 'guest.o'
    subprocess.run(['clang', '-arch', 'x86_64', '-c',
                    str(root / 'test/ept_mmio_guest.S'), '-o', str(obj)], check=True)
    rom = output / 'guest.rom'
    extract_rom(obj, rom)
    library = output / 'stale-length.dylib'
    subprocess.run(['clang', '-arch', 'x86_64', '-dynamiclib', '-Wall', '-Wextra',
                    '-Werror', '-Wno-deprecated-declarations', '-framework', 'Hypervisor',
                    str(root / 'test/ept_stale_length.c'), '-o', str(library)], check=True)
    results = [run_case(binary, rom, library, length, output)
               for length in (None, 0, 1, 2, 6, 15)]
    (output / 'results.json').write_text(json.dumps(results, indent=2) + '\n')
    return 0 if all(case['passed'] for case in results) else 1


if __name__ == '__main__':
    sys.exit(main())
