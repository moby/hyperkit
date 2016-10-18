The majority of code is derived from `bhyve` and therefore the general
coding style should follow the
[FreeBSD coding style guidelines](https://www.freebsd.org/cgi/man.cgi?query=style&sektion=9)
(which are also accessible via `man 9 style` on Mac OS X).

You may use tools like
[uncrustify](http://uncrustify.sourceforge.net/) with this
[config file](https://github.com/freebsd/pkg/blob/master/freebsd.cfg)
for *new* code, though the result may not be perfect.

Keep in mind that, especially for most of the `bhyve` derived code, it
is more important to try to keep the code in line with `bhyve` if
possible.
