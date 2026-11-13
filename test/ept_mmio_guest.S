/* SPDX-License-Identifier: BSD-2-Clause */

/* A reset ROM with six-byte IOAPIC accesses and no external firmware. */
.text
.code16
Lstart:
	cli
	lgdtl %cs:(Lgdt_pointer - Lstart)
	movl %cr0, %eax
	orl $1, %eax
	movl %eax, %cr0
	ljmpl $8, $(0xffff0000 + Lprotected_mode - Lstart)

.code32
Lprotected_mode:
	movw $16, %ax
	movw %ax, %ds
	movw %ax, %es
	movw %ax, %ss
	movl $0x8000, %esp
	/* The emulator fetches instructions from guest RAM, not ROM mappings. */
	movl $(0xffff0000 + Lpayload - Lstart), %esi
	movl $0x100000, %edi
	movl $(Lpayload_end - Lpayload), %ecx
	cld
	rep movsb
	ljmpl $8, $0x100000
Lpayload:
	movl $0xfec00000, %edi
	movl $7, %eax
	/* movl %eax, 0(%edi), with a disp32 to require six bytes. */
	.byte 0x89, 0x87
	.long 0
	jmp Lafter_write
	/* A stale length of 15 must not accidentally reach the success path. */
	.fill 15, 1, 0xf4
Lafter_write:
	movl $0xdeadbeef, %eax
	/* movl 0(%edi), %eax, also using a disp32. */
	.byte 0x8b, 0x87
	.long 0
	jmp Lafter_read
	.fill 15, 1, 0xf4
Lafter_read:
	cmpl $7, %eax
	jne Lfailed
	movl $(0x100000 + Lpassed_message - Lpayload), %esi
	jmp Lprint
Lfailed:
	movl $(0x100000 + Lfailed_message - Lpayload), %esi
Lprint:
	movw $0x3f8, %dx
Lnext_char:
	lodsb
	testb %al, %al
	jz Lpoweroff
	outb %al, %dx
	jmp Lnext_char
Lpoweroff:
	/* HyperKit's PM1 control port: SLP_TYP=S5, SLP_EN=1. */
	movw $0x404, %dx
	movw $0x3400, %ax
	outw %ax, %dx
Lstopped:
	hlt
	jmp Lstopped

Lpassed_message:
	.asciz "EPT_MMIO_PASS\n"
Lfailed_message:
	.asciz "EPT_MMIO_FAIL\n"
Lpayload_end:

.p2align 3
Lgdt:
	.quad 0
	.quad 0x00cf9b000000ffff
	.quad 0x00cf93000000ffff
Lgdt_end:
Lgdt_pointer:
	.word Lgdt_end - Lgdt - 1
	.long 0xffff0000 + Lgdt - Lstart

.code16
.org 0xfff0, 0xff
	/* A 16-bit near jump wraps IP from 0xfff0 back to offset zero. */
	.byte 0xe9, 0x0d, 0x00
.org 0x10000, 0xff
