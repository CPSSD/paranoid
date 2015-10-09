pfs(1) -- issue commands to a paranoid filesystem
=================================================

## SYNOPSIS

`pfs` `init` `<pfs-directory>`<br>
`pfs` `mount` `<pfs-directory>` `<mountpoint>`<br>
`pfs` [`-f`|`--fuse`] `stat` `<pfs-directory>` `<file>`<br>
`pfs` [`-f`|`--fuse`] `read` `<pfs-directory>` `<file>` [`<offset>` `<length>`]<br>
`pfs` [`-n`|`--net`|`-f`|`--fuse`] `creat` `<pfs-directory>` `<file>`<br>
`pfs` [`-n`|`--net`|`-f`|`--fuse`] `write` `<pfs-directory>` `<file>` [`<offset>` `<length>]`<br>

## DESCRIPTION

**pfs** is the control system for the paranoid file storage system. It handles
communication between the network layers and FUSE and the virtual file system.
It can also be used to test the file system by omitting the
`-n` or `-f` flags.

## COMMANDS

* `init`:
    Create a new filesystem in the indicated directory.  The directory must already exist and must be empty.

* `mount`:
    Mount the indicated paranoid filesystem with FUSE. It will be mounted at `<mountpoint>`.

* `stat`:
    Writes stat information for the indicated file to standard output (in a format to be determined, but must initially include at least the length).

* `read`:
    Reads the file `<file>` and prints it to standard output.  If `<offset>` and `<length>` are omitted, then output all of the file.

* `creat`:
    Create a new file in the filesystem and create a hard link to it called `<file>`.

* `write`:
    Write the data piped in through standard input to the file referenced by `<file>` starting at `<offset>` and
    writing `<length>` bytes of data.  If `<offset>` and `<length>` are omitted, then the file is first truncated and the length of the write
    is determined by the amount of data on standard input.

## OPTIONS

These options specify the source of the command.

  * `-n`, `--net`:
    The source of the command is the network. It is coming from another node. This
    flag instructs `pfs` not to send a message out on the network after performing the
    operation.

  * `-f`, `--fuse`:
    The source of the command is the FUSE layer. This flag instructs `pfs` to send a message
    out on the network after performing the operation.

These options specify the output of the program.

  * `-v`, `--verbose`:
    This enables debug logging to standard error.

These options are miscellaneous options.

  * `--version`:
    Show pfs version and exit.

## EXAMPLES

Create new file with primary link name of `helloworld.txt` and write the text "Hello World!"

    $ pfs creat ~/.pfs helloworld.txt
    $ echo "Hello World!" | pfs write ~/.pfs helloworld.txt 0 13
    $ pfs read ~/.pfs helloworld.txt
    Hello World!

