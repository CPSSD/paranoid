Paranoid Fuse Interface v0.1
=============================
##Installation
Before installing PFI, fuse must first be installed on the machine.

[FUSE](http://fuse.sourceforge.net/)

To install PFI go to the directory paranoid/pfi and run
```
go build -i
```
##Usage
###General
```
pfi [flags] <PfsDirectory> <MountPoint>
```

* PfsDirectory - The path to the folder where the pfs file system is located.
* MountPoint - The path to the folder where you wish to mount the paranoid file system.

example :
```
pfi -v ~/coding/initDir /home/mladen/coding/mountDir
```

###Flags
-v

Prints logs in standard output.
