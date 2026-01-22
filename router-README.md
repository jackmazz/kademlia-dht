# Router Server

This is a correct router server that you can use for your k-DHT Part 2
implementation.

It includes router implementations to support:
* Linux on x86-64 (`kdht-router-linux-x86_64`)
* Windows with WSL2 (`kdht-router-linux-x86_64`)
* macOS on x86-64 (`kdht-router-macos-x86_64`)
* macOS on Apple Silicon (`kdht-router-macos-arm64`)

You will have to install the appropriate version for your computer,
and follow these directions _very carefully_ to ensure that it works.

## Usage

In order to use this router, you need to do four things:

* Put the appropriate `kdht-router` somewhere in your `PATH`
* Put the `given/router` directory in your Git workspace's `given`
  directory
* Put the `api/kdht/router.pb.go` file in your Git workspace's
  `api/kdht` directory
* Return the given `router.New(node, k)` from your
  `impl.NewRoutingTable(node, k)` function

It is _very important that you do exactly these things, exactly as
described_.  If you do anything slightly differently, the most likely
result is that the router does not work at all, or your program does
not compile.  For some reason students often edit commands that they
are given, ignore whitespace, etc., which will certainly cause this to
fail.  If possible, open this README in Emacs on your computer, then
use Emacs to copy the commands you need to run (one at a time) and
paste them into your terminal.  If you get an error, _stop,_ look to
see what you did wrong, and fix it before proceeding.

For the following instructions, we will assume that your k-DHT node
sources are in the directory `$CHECKOUT` (which I treat as a shell
variable for Unix commands).  You can cause this to be the case by, in
your terminal, changing to the top level of your Git checkout (where
you would run `make test`), and running **exactly** the command
(spaces are critical, the quotes must be backquotes!):

```
CHECKOUT=`pwd`
```

After you do that, change to the directory that contains this tar
file in your terminal, and follow these steps:

### Install `kdht-router` in your path

You will need to identify the `kdht-router` executable that matches
your system; there is a list at the beginning of this document that
gives the corresponding filename for each supported operating system
and architecture.  If you are on macOS and are not sure of your
architecture, the command `uname -m` at the prompt will tell you
whether you are on x86-64 or ARM64 (Apple Silicon).  Once you have
determined the correct binary for your system, replace `$ROUTER` in
the below commands with its full name; as with the `CHECKOUT`
environment variable, you can set `ROUTER=kdht-router-XXX`, where XXX
is the appropriate text for your system, and simply copy and paste the
commands here.

Run these commands:

```
mkdir -p ~/bin
cp $ROUTER ~/bin/kdht-router
```

These should complete silently.  Any output indicates an error.

### Copy the required Go files into your sources

Run this command:

```
rsync -r given api $CHECKOUT
```

This should complete silently.

### Use the given router

Open your `impl/routing.go` file, and make these changes:

* Remove the import "errors"
* Add the import `"cse586.kdht/given/router"`
* Change the body of the `NewRoutingTable()` function to contain
  exactly the line `return router.New(node, k)`

### Load the new PATH

Run the command:

```
kdht-router
```

If you get the message `usage: kdht-router sockaddr`, you do not need
to do anything else to put this in your `PATH`, and you can move on to
testing the installation.

If you get a message like `bash: kdht-router: command not found` (it
might be from `zsh` on macOS), you did not already have a local `bin`
directory and it is not in your path.  Try logging out and back in to
your computer, and see if the above command gives the usage message.
If it does, you can move on to testing the installation.

If logging out and back in does **not** cause the command
`kdht-router` to emit a usage message, you need to edit your `PATH` in
your login configuration.  To do this, you need to know what your
shell is; there are two possibilities that we will support: GNU bash
and zsh.  If you are running some other shell, presumably you know how
to configure it, and you should do so.  Run this command:

```
echo $SHELL
```

You should see a command that ends in either `bash` (in which case you
are running bash) or `zsh` (in which case you are running zsh).
Follow the appropriate instructions as follow.

#### For bash users

Run this command:

```
echo 'PATH=$PATH:$HOME/bin' >> ~/.bashrc
```

Again, it is **VERY IMPORTANT** that you run that command _exactly_ as
written.  In particular, you must make sure to use single quotes and a
double angle bracket, and to get the spacing exactly right.

#### For zsh users

Run this command:

```
echo 'PATH=$PATH:$HOME/bin' >> ~/.zshrc
```

Again, it is **VERY IMPORTANT** that you run that command _exactly_ as
written.  In particular, you must make sure to use single quotes and a
double angle bracket, and to get the spacing exactly right.

#### After updating your login files

Once you have updated your login files, whether you are running bash
or zsh, you must open a new terminal window for this change to take
effect.  (Opening a new terminal will not change the configuration of
any existing terminals.)  You may find that logging out and back in is
necessary in some circumstances.  In any case, after (possibly)
logging out and back in, you should find that running `kdht-router` at
the prompt gives you the usage message above.

## Testing the Installation

Once you have installed the router in your sources (and restarted your
VM or logged out and back in, if necessary), you can run some simple
tests on it.  **If running `kdht-router` at the prompt with no
arguments does not result in the message `usage: kdht-router socket`,
please return to the previous section and fix your installation.**
Test your installation by running, from your Git workspace top-level
directory:

```
make
go test ./given/router
```

If the above `go test` command does not compile and run, you did not
copy the sources from this directory into your git workspace
correctly.  If it compiles and runs but the tests report errors, there
is some other kind of error that the error messages may help you
identify (please read them carefully).  If the tests run and work, the
given router works, but your `impl` might still be problematic.  Run
`make test` and check the given router-related tests from Part 1 of
this project.  If they fail, check the Usage section of this document
for how to install the router in your `impl`.  If they pass, you
should be good to go.
