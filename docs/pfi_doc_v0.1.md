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
pfi [flags] <MountLocation> <PfsInitLocation> <PfsBinaryPath>
```

* MountLocation - The path to the folder where you wish to mount the paranoid file system.
* PfsInitLocation - The path to the folder where you initialised pfs with `pfs init initlocation`.
* PfsBinaryPath - The path to the pfs binary.

example :
```
pfi -v /home/mladen/coding/mountDir ~/coding/initDir bin/pfs
```

###Flags
-v

Prints logs in standard output.
