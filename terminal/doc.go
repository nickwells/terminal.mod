/*
Package terminal provides a way of controlling POSIX terminal behaviour. This
includes representing the desired attributes for a piece of text. These
attributes can then be used to generate terminal control strings to display
the text suitably decorated.

When it comes to using the available attributes restraint is advised. While
the package makes it possible to display the text in bold, fast-blinking
italic with the text crossed out and doubly underlined, the foreground colour
set to barbie pink and the background set to lime green, you might not want
this.

Note also that no checks are made that the foreground and background colours
are different or that they offer good contrast so take care when deciding
which colours to use.

Note also that not all terminals support all the attributes offered by this
package. This applies especially to colours where the full range of RGB
encoded colours may not be available.
*/
package terminal
