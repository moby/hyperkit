* Why is this change needed or what are the use cases?

   Currently inspecting a hyperkit running kernel requires detailed knowledge
   of hyperkit or adding additional compilation flags for tracing. Ideally it
   should be possible to debug a hosted kernel as if using hardware debugger or
   similar to how qemu/bochs expose a gdbserver.

* What are the requirements this change should meet?

   Provide a mechanism to enable connecting a debugger (lldb, gdb) to the hosted kernel.

* What are some ways to design/implement this feature?

   Both gdb and lldb support using the gdb protocol for remote target debugging.

    * [GDB Remote Protocol](https://sourceware.org/gdb/current/onlinedocs/gdb/Remote-Protocol.html)
    * [lldb remote](http://lldb.llvm.org/remote.html)
    * [lldb protocol extensions](http://llvm.org/viewvc/llvm-project/lldb/trunk/docs/lldb-gdb-remote.txt?view=markup)

    Support can be added by exposing a port that acts as a gdbserver allowing both debuggers to connect options are. 

    * TCP/IP `target remote localhost:1234`
    * Debug serial port for debugging (eg COM2) `target remote com.docker.hyperkit/debugtty`
    * Support both TCP/IP and serial

* Which design/implementation do you think is best and why?

   Supporting gdb and lldb by using the common protocol without extensions. Whilst lldb is the default debugger for OS X users may wish to use a more familiar debugger with their target kernel which is running under hyperkit. In addition
   supporting both TCP/IP and serial give the widest availability for users to debug.

* What are the risks or limitations of your proposal?

  * There could be performance impact so the gdb hooks would need to be guarded by lightweight checks in order to minimize overhead.
  * The impact of OS X security features such as sandboxd and System Integrity Protection would need to be tested to ensure debugging is possible in a signed build.
    * [Overview of SIP](https://derflounder.wordpress.com/2015/10/01/system-integrity-protection-adding-another-layer-to-apples-security-model/)
