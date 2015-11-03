pfsm(8) -- issue commands to a paranoid filesystem
=================================================

## SYNOPSIS

`pfsm` `init` `<pfs-directory>`<br>
`pfsm` `mount` `<pfs-directory>` `<server-ip>` `<server-port>`<br>
`pfsm` [`-f`|`--fuse`] `stat` `<pfs-directory>` `<file>`<br>
`pfsm` [`-n`|`--net`|`-f`|`--fuse`] `utimes` `<pfs-directory>` `<file>`<br>
`pfsm` [`-f`|`--fuse`] `access` `<pfs-directory>` `<file>` `<filemode>`<br>
`pfsm` [`-n`|`--net`|`-f`|`--fuse`] `chmod` `<pfs-directory>` `<file>` `<permflags>`<br>
`pfsm` [`-f`|`--fuse`] `read` `<pfs-directory>` `<file>` [`<offset>` `<length>`]<br>
`pfsm` [`-f`|`--fuse`] `readdir` `<pfs-directory>`<br>
`pfsm` [`-n`|`--net`|`-f`|`--fuse`] `creat` `<pfs-directory>` `<file>` `<permflags>` <br>
`pfsm` [`-n`|`--net`|`-f`|`--fuse`] `link` `<pfs-directory>` `<file>` `<target>` <br>
`pfsm` [`-n`|`--net`|`-f`|`--fuse`] `unlink` `<pfs-directory>` `<file>` <br>
`pfsm` [`-n`|`--net`|`-f`|`--fuse`] `rename` `<pfs-directory>` `<file>` `<newname>` <br>
`pfsm` [`-n`|`--net`|`-f`|`--fuse`] `truncate` `<pfs-directory>` `<file>` `<length>` <br>
`pfsm` [`-n`|`--net`|`-f`|`--fuse`] `write` `<pfs-directory>` `<file>` [`<offset>` `<length>]`<br>

## DESCRIPTION

**pfsm** is the control system for the paranoid file storage system. It handles
communication between the network layers and FUSE and the paranoid file system.
It can also be used to test the file system by omitting the
`-n` or `-f` flags.

## COMMANDS

* `init`:
    Create a new filesystem in the indicated directory.  The directory must already exist and must be empty. Also generates the file-system's UUID.

* `mount`:
    This informs that the indicated file system has been mounted; specifically, it provides the server IP address and port to use.  These are written into the relevant meta-data files.

* `stat`:
    Writes stat information for the indicated file to standard output (in a format to be determined, but must initially include at least the length).

* `utimes`:
	Updates the atime and mtime of `<file>` to those given as JSON on stdin

* `access`:
	Check if the file `<file>` can be acessed in `<filemode>`

* `chmod`:
	Change the permissions of `<file>` to `<permflags>`

* `read`:
    Reads the file `<file>` and prints it to standard output.  If `<offset>` and `<length>` are omitted, then output all of the file.

* `readdir`:
    Returns a list of all files in the filesystem to standard output one per line. 

* `creat`:
    Create a new file in the filesystem and create a hard link to it called `<file>`.

* `link`:
    Create a new link in the filesystem `<file>` to `<target>`.

* `unlink`:
	Delete the link `<file>` and decrement the reference count of it's file.
	If the reference count is 0, delete the file.

* `rename`:
	Change the name of `<file>` to `<newname>`

* `truncate`:
	Change the size of `<file>` to `<length>` bytes

* `write`:
    Write the data piped in through standard input to the file referenced by `<file>` starting at `<offset>` and
    writing `<length>` bytes of data.  If `<offset>` and `<length>` are omitted, then the file is first truncated and the length of the write
    is determined by the amount of data on standard input.

## OPTIONS

These options specify the source of the command.

  * `-n`, `--net`:
    The source of the command is the network. It is coming from another node. This
    flag instructs `pfsm` not to send a message out on the network after performing the
    operation.

  * `-f`, `--fuse`:
    The source of the command is the FUSE layer. This flag instructs `pfsm` to send a message
    out on the network after performing the operation.

These options specify the output of the program.

  * `-v`, `--verbose`:
    This enables debug logging to standard error.

These options are miscellaneous options.

  * `--version`:
    Show pfs version and exit.

## EXAMPLES

Create new file with primary link name of `helloworld.txt` and write the text "Hello World!"

    $ pfsm creat ~/.pfs helloworld.txt
    $ echo "Hello World!" | pfsm write ~/.pfs helloworld.txt 0 13
    $ pfsm read ~/.pfs helloworld.txt
    Hello World!

