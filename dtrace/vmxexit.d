#!/usr/sbin/dtrace -C -s
/*
 * vmxexit.d - report VMX exits for particular VM
 *
 * USAGE: sudo vmxexit.d -p <pid of hyperkit>
 *        sudo vmxexit.d -D TOTAL -p <pid of hyperkit>
 *        sudo vmxexit.d -D INTERVAL=1 -p <pid of hyperkit>
 *
 * If '-D INTERVAL=<seconds>' is specified periodically print a
 * summary of VM Exits per vCPU.
 *
 * It seems that tracing starts before the end of dtrace:::BEGIN, i.e.
 * before the 'reasons' array is fully initialised. This results in
 * some empty/corrupted exit reasons to be printed. According to the
 * documentation, this should not happen, so it might be a bug in High
 * Sierra.
 *
 * As a workaround, specifying '-D·NUMERIC' stores and prints the VM
 * Exit reasons numbers instead of names.
 */

#pragma D option quiet

string reasons[int];

dtrace:::BEGIN
{
        start = timestamp;

        reasons[0]  = "EXCEPTION";
        reasons[1]  = "EXT_INTR";
        reasons[2]  = "TRIPLE_FAULT";
        reasons[3]  = "INIT";
        reasons[4]  = "SIPI";
        reasons[5]  = "IO_SMI";
        reasons[6]  = "SMI";
        reasons[7]  = "INTR_WINDOW";
        reasons[8]  = "NMI_WINDOW";
        reasons[9]  = "TASK_SWITCH";
        reasons[10] = "CPUID";
        reasons[11] = "GETSEC";
        reasons[12] = "HLT";
        reasons[13] = "INVD";
        reasons[14] = "INVLPG";
        reasons[15] = "RDPMC";
        reasons[16] = "RDTSC";
        reasons[17] = "RSM";
        reasons[18] = "VMCALL";
        reasons[19] = "VMCLEAR";
        reasons[20] = "VMLAUNCH";
        reasons[21] = "VMPTRLD";
        reasons[22] = "VMPTRST";
        reasons[23] = "VMREAD";
        reasons[24] = "VMRESUME";
        reasons[25] = "VMWRITE";
        reasons[26] = "VMXOFF";
        reasons[27] = "VMXON";
        reasons[28] = "CR_ACCESS";
        reasons[29] = "DR_ACCESS";
        reasons[30] = "INOUT";
        reasons[31] = "RDMSR";
        reasons[32] = "WRMSR";
        reasons[33] = "INVAL_VMCS";
        reasons[34] = "INVAL_MSR";
        /* 35 not documented */
        reasons[36] = "MWAIT";
        reasons[37] = "MTF";
        /* 38 not documented */
        reasons[39] = "MONITOR";
        reasons[40] = "PAUSE";
        reasons[41] = "MCE_DURING_ENTRY";
        /* 42 not documented */
        reasons[43] = "TPR";
        reasons[44] = "APIC_ACCESS";
        reasons[45] = "VIRTUALIZED_EOI";
        reasons[46] = "GDTR_IDTR";
        reasons[47] = "LDTR_TR";
        reasons[48] = "EPT_FAULT";
        reasons[49] = "EPT_MISCONFIG";
        reasons[50] = "INVEPT";
        reasons[51] = "RDTSCP";
        reasons[52] = "VMX_PREEMPT";
        reasons[53] = "INVVPID";
        reasons[54] = "WBINVD";
        reasons[55] = "XSETBV";
        reasons[56] = "APIC_WRITE";

        printf("Tracing... Hit Ctrl-C to end.\n");

        #ifdef INTERVAL
        secs = INTERVAL;
        printf("\n\n");
        printf("Per CPU VM Exits\n");
        #endif
}

hyperkit$target:::vmx-exit
{
        #ifdef INTERVAL
        /* Per vCPU count for periodic reporting */
        @total[arg0] = count();
        #endif

        /* Per Reason per vCPU counts for summary */
        #ifdef NUMERIC
        @num[arg1, arg0] = count();
        #else
        @num[reasons[arg1], arg0] = count();
        #endif
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
        printf("%16s %-4s %8s\n", "REASON", "vCPU", "COUNT");
        #else
        printf("%16s %-4s %8s\n", "REASON", "vCPU", "RATE (1/s)");
        normalize(@num, (timestamp - start) / 1000000000);
        #endif
        #ifdef NUMERIC
        printa("              %2d %-4d %@8d\n", @num);
        #else
        printa("%16s %-4d %@8d\n", @num);
        #endif
}
