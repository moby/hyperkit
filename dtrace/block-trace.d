#!/usr/sbin/dtrace -C -s
/*
 * block-trace.d - Trace all block device accesses
 *
 * USAGE: sudo block-trace.d -p <pid of hyperkit>
 */

#pragma D option quiet

dtrace:::BEGIN
{
    printf("Tracing... Hit Ctrl-C to end.\n");
    printf("%13s  %-15s  Arguments\n", "Time(us)", "Function")
}

hyperkit$target:::block-preadv,
hyperkit$target:::block-pwritev,
hyperkit$target:::block-delete
{
    printf("%13d  %-15s  %#016x %#016x\n", timestamp/1000, probefunc, arg0, arg1)
}

dtrace:::END
{
}
