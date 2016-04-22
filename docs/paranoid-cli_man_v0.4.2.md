paranoid-cli(1) -- interact with a paranoid filesystem
======================================================

## SYNOPSIS

`paranoid-cli` `init` `<pfs-name>`<br>
`paranoid-cli` `mount` `<port>` `<discovery-server-address>` `<pfs-name>` `<mountpoint>`<br>
`paranoid-cli` `automount` `<pfs-name>`<br>
`paranoid-cli` `unmount` `<pfs-name>`<br>
`paranoid-cli` `list`<br>
`paranoid-cli` `secure` `<pfs-name>`<br>
`paranoid-cli` `delete` `<pfs-name>`<br>
`paranoid-cli` `history` `<pfs-name>||<log-directory>`<br>
`paranoid-cli` `buildfs` `<pfs-name>` `<log-directory>`<br>
`paranoid-cli` `restart` `<pfs-name>`<br>
`paranoid-cli` `serve` `<pfs-name>` `<file>` <br>
`paranoid-cli` `unserve` `<pfs-name>` `<file>` <br>
`paranoid-cli` `list-serve` `<pfs-name>` <br>

## DESCRIPTION

**paranoid-cli** is used by a user to interact with a paranoid file system. It can be used to init,
mount or unmount a paranoid file system.

## COMMANDS

* `init`:
    Create a new filesystem with the indicated name.

* `mount`:
	This is used to mount the indicated filesystem to mountpoint. As part of the cluster servered by the server at the `discovery-server-address`. It will instruct pfsd to run on port `port`

* `automount`:
	This is used to mount the indicated filesystem with the previous settings used.

* `unmount`:
	This is used to unmount the paranoid filesystem mounted at the given mountpoint.

* `list`:
    This lists all of the currently-existing paranoid filesystems available to mount.

* `secure`:
    Generate TLS/SSL certificate files for a previously-unsecured paranoid filesystem.

* `delete`:
    Permanently delete a paranoid filesystem.

* `history`:
    View the history of a paranoid filesystem or of the specified log-directory

* `buildfs`:
    builds a filesystem with the given <pfs-name> from the logfiles whos location is specified by <log-directory>

* `restart`:
    Restart the network services of the named paranoid filesystem.

* `serve`:
    Set a file to be served from the Paranoid File Server

* `unserve`:
    Remove a file that has been previously served from the paranoid file server

* `list-serve`:
    List all files that are currently on the server and are available to be served

## OPTIONS

  * `--verbose`:
    This enables debug logging to standard error.

  * `--networkoff`:
	This disables networking for testing purposes
