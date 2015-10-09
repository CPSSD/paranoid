Other Interfaces/APIs
=====================

## User-Facing FUSE Control

Mounts a PFS file system using FUSE and starts the PFS network client.

Notes:

- This is the **only** user-facing utility (??).

- The `mount` functionality should be removed from `pfs`.

#### Synopsis

    mount.pfs -o server=1.2.3.4,port=1234 pfs-directory mountpoint

`pfs-directory` must exist.  If the directory is empty, then runs `pfs -f init <pfs-directory>`.

(Aside... At some point we will need an `fsck.pfs` utility to first verify/fix the PFS file-system structure.)

## Network Client

#### Synopsis 1

    pfs-network-client --client <pfs-directory> <server-ip> <server-port>

Starts the network client, connects to the server and runs `pfs` commands in response to messages received from the server.

The process runs until it receives an appropriate signal (e.g. `INT` or `TERM`).  It exits if the connection with the server is lost.

#### Synopsis 2

    pfs-network-client --send <server-ip> <server-port>

Connects to the server, reads and sends its standard input to the server (stdin is assumed to be a valid message), then exits.

## Network Server

#### Synopsis

    pfs-network-server <server-ip> <server-port>

Starts the server.
