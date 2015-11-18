paranoid-cli(8) -- interact with a paranoid filesystem
=================================================

## SYNOPSIS

`paranoid-cli` `init` `<pfs-name>`<br>
`paranoid-cli` `mount` `<port>` `<discovery-server-address>` `<pfs-name>` `<mountpoint>`<br>
`paranoid-cli` `automount` `<pfs-name>`<br>
`paranoid-cli` `unmount` `<pfs-name>`<br>

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

## OPTIONS

  * `--verbose`:
    This enables debug logging to standard error.

  * `--networkoff`:
	This disables networking for testing purposes
