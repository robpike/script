Script is an interactive driver for a line-at-a-time command language such as
the shell. It takes no arguments.

It reads each line of the script file, waiting for a newline on
standard input to proceed. After receiving a newline, it prints the
next line of of the script and also feeds it to a single running
command instance, which prints the output of the command to standard
output. Thus script serves as a way to control the sequential input
to the command and thus supervise its activities.

The argument file is in whatever format the command normally accepts.
Comments, in whatever single-line form the command accepts, will
be passed to the command but can serve as documentation during
execution.  There is no facility to send more than one line of input
at a time.

Normally the user types an empty line to proceed, advancing to the
next line of the script. But a non-empty line is passed to the
command without advanding the script, allowing the user to inject
extra commands.
