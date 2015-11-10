Glossary
=====================

- `pfs` : Paranoid File System. This is how the data is stored by paranoid.
- `pfsm` : Paranoid File System Manager. This is the binary that manages operations on pfs.
- `FUSE` : File System in User Space. This is a kernel module for Unix operating systems
that allows users create their own file systems without editing kernel code.
- `pfi` : Paranoid Fuse Interface. This is the binary that manages communication between FUSE and pfsm.
- `pfsd` : Paranoid File System Daemon. This handles networking between paranoid nodes.
- `ic` : Internal Communication. The packages that aid in communication between pfi and pfsd
- `paranoid-cli` : Paranoid Command Line Interface. This is the only user-facing element of the system.
Used to run commands like mount and unmount for a given paranoid file system.
- `paranoid-directory` : This is the directory where the paranoid file system data is stored.
- `mountpoint` : This is the directory where the data from pfs is to be displayed by FUSE.
