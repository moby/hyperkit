#!/usr/sbin/dtrace -C -s
/*
 * eptfault.d - report all EPT faults for particular VM
 *
 * USAGE: sudo eptfault.d -p <pid of hyperkit>
 *        sudo eptfault.d -D TOTAL -p <pid of hyperkit>
 *        sudo eptfault.d -D INTERVAL=1 -p <pid of hyperkit>
 *
 * If '-D INTERVAL=<seconds>' is specified periodically print a
 * summary of EPT Faults per vCPU.
 *
 * It seems that tracing starts before the end of dtrace:::BEGIN, i.e.
 * before the 'lapic_map' array is fully initialised. This results in
 * some empty/corrupted Local APIC register names being
 * printed. According to the documentation, this should not happen, so
 * it might be a bug in High Sierra.
 *
 * As a workaround, specifying '-D·NUMERIC' stores and prints the Local
 * APIC register offsets instead of names.
 */

#pragma D option quiet

string lapic_map[uint16_t];

dtrace:::BEGIN
{
        start = timestamp;

        /* See Chapter: "Advanced Programmable Interrupt Controller (APIC)" in
         * Intel 64 and IA-32 Architectures Software Developer's Manual Vol 3 */
        lapic_map[0x000] = "RES0x0";
        lapic_map[0x010] = "RES0x1";
        lapic_map[0x020] = "ID";
        lapic_map[0x030] = "VER";
        lapic_map[0x040] = "RES0x4";
        lapic_map[0x050] = "RES0x5";
        lapic_map[0x060] = "RES0x6";
        lapic_map[0x070] = "RES0x7";
        lapic_map[0x080] = "TPR";
        lapic_map[0x090] = "APR";
        lapic_map[0x0A0] = "PPR";
        lapic_map[0x0B0] = "EOI";
        lapic_map[0x0C0] = "RRR";
        lapic_map[0x0D0] = "LDR";
        lapic_map[0x0E0] = "DFR";
        lapic_map[0x0F0] = "SVR";
        lapic_map[0x100] = "ISR0";
        lapic_map[0x110] = "ISR1";
        lapic_map[0x120] = "ISR2";
        lapic_map[0x130] = "ISR3";
        lapic_map[0x140] = "ISR4";
        lapic_map[0x150] = "ISR5";
        lapic_map[0x160] = "ISR6";
        lapic_map[0x170] = "ISR7";
        lapic_map[0x180] = "TMR0";
        lapic_map[0x190] = "TMR1";
        lapic_map[0x1A0] = "TMR2";
        lapic_map[0x1B0] = "TMR3";
        lapic_map[0x1C0] = "TMR4";
        lapic_map[0x1D0] = "TMR5";
        lapic_map[0x1E0] = "TMR6";
        lapic_map[0x1F0] = "TMR7";
        lapic_map[0x200] = "IPR0";
        lapic_map[0x210] = "IRR1";
        lapic_map[0x220] = "IRR2";
        lapic_map[0x230] = "IRR3";
        lapic_map[0x240] = "IRR4";
        lapic_map[0x250] = "IRR5";
        lapic_map[0x260] = "IRR6";
        lapic_map[0x270] = "IRR7";
        lapic_map[0x280] = "ESR";
        lapic_map[0x290] = "RES0x29";
        lapic_map[0x2A0] = "RES0x2A";
        lapic_map[0x2B0] = "RES0x2B";
        lapic_map[0x2C0] = "RES0x2C";
        lapic_map[0x2D0] = "RES0x2D";
        lapic_map[0x2E0] = "RES0x2E";
        lapic_map[0x2F0] = "CMCI_LVT";
        lapic_map[0x300] = "ICR_LOW";
        lapic_map[0x310] = "ICR_HI";
        lapic_map[0x320] = "LVT_TIMER";
        lapic_map[0x330] = "LVT_THERM";
        lapic_map[0x340] = "LVT_PERF";
        lapic_map[0x350] = "LVT_LINT0";
        lapic_map[0x360] = "LVT_LINT1";
        lapic_map[0x370] = "LVT_ERROR";
        lapic_map[0x380] = "TIMER_ICR";
        lapic_map[0x390] = "TIMER_CCR";
        lapic_map[0x3A0] = "RES0x3A";
        lapic_map[0x3B0] = "RES0x3B";
        lapic_map[0x3C0] = "RES0x3C";
        lapic_map[0x3D0] = "RES0x3D";
        lapic_map[0x3E0] = "TIMER_DCR";
        lapic_map[0x3F0] = "SELF_IPI";

        printf("Tracing... Hit Ctrl-C to end.\n");

        #ifdef INTERVAL
        secs = INTERVAL;
        printf("\n\n");
        printf("Per CPU VM Exits\n");
        #endif
}

hyperkit$target:::vmx-ept-fault
/(arg1 & 0xfffff000) == 0xfee00000/
{
        // LAPIC FAULTS
        // The Local APIC has a 4K register space with registers 128bit aligned

        #ifdef NUMERIC
        @lapic_faults[arg1 & 0xfff, arg0] = count();
        #else
        @lapic_faults[lapic_map[arg1 & 0xfff], arg0] = count();
        #endif
}

hyperkit$target:::vmx-ept-fault
{
        #ifdef INTERVAL
        /* Per vCPU count for periodic reporting */
        @total[arg0] = count();
        #endif

        @all_faults[arg1, arg0] = count();
}

#ifdef INTERVAL
/* timer */
profile:::tick-1sec
{
       secs--;
}

/* Periodically print per vCPU VM Exits */
profile:::tick-1sec
/secs == 0/
{
        printa(@total);
        clear(@total);
        secs = INTERVAL;
}
#endif

dtrace:::END
{
        #ifdef TOTAL
        printf("%18s %-4s %10s\n", "ADDRESS", "vCPU", "COUNT");
        #else
        printf("%18s %-4s %10s\n", "ADDRESS", "vCPU", "RATE (1/s)");
        normalize(@all_faults, (timestamp - start) / 1000000000);
        #endif
        printa("%18x %-4d %@10d\n", @all_faults);

        #ifdef TOTAL
        printf("%18s %-4s %10s\n", "LAPIC REGISTER", "vCPU", "COUNT");
        #else
        printf("%18s %-4s %10s\n", "LAPIC REGISTER", "vCPU", "RATE (1/s)");
        normalize(@lapic_faults, (timestamp - start) / 1000000000);
        #endif
        #ifdef NUMERIC
        printa("             0x%03x %-4d %@10d\n", @lapic_faults);
        #else
        printa("%18s %-4d %@10d\n", @lapic_faults);
        #endif
}
