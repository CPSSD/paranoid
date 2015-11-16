paranoid-cli(8) -- interact with a paranoid filesystem
=================================================

## SYNOPSIS

`paranoid-cli` `init` `<pfs-directory>`<br>
`paranoid-cli` `mount` `port` `<discovery-server-address>` `<paranoid-directory>` `<mountpoint>`<br>
`paranoid-cli` `unmount` `<mountpoint>`<br>

## DESCRIPTION

**paranoid-cli** is used by a user to interact with a paranoid file system. It can be used to init,
mount or unmount a paranoid file system.

## COMMANDS

* `init`:
    Create a new filesystem in the indicated directory.  The directory must already exist and must be empty.

* `mount`:
	This is used to mount the indicated filesystem to mountpoint. As part of the cluster servered by the server at the `discovery-server-address`. It will instruct pfsd to run on port `port`

* `unmount`:
	This is used to unmount the paranoid filesystem mounted at the given mountpoint.

## OPTIONS

  * `--verbose`:
    This enables debug logging to standard error.
