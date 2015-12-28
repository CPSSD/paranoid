Glossary
=====================

- `pfs` : Paranoid File System. This is how the data is stored by paranoid.
- `libpfs` : Paranoid File System Manager. This is the liabary that manages operations on pfs.
- `FUSE` : File System in User Space. This is a kernel module for Unix operating systems
that allows users create their own file systems without editing kernel code.
- `pfi` : Paranoid Fuse Interface. This is the package that handles communication with FUSE.
- `pfsd` : Paranoid File System Daemon. This is the binary that handles backaground funcionality of the system such as interfacing with fuse and networking with other nodes.
- `paranoid-cli` : Paranoid Command Line Interface. This is the only user-facing element of the system.
Used to run commands like mount and unmount for a given paranoid file system.
- `paranoid-directory` : This is the directory where the paranoid file system data is stored.
- `mountpoint` : This is the directory where the data from pfs is to be displayed by FUSE.
