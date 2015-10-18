Paranoid Fuse Interface v0.1
=============================
##Installation
Before installing PFI, fuse must first be installed on the machine.

[FUSE](http://fuse.sourceforge.net/ "Fuse Homepage")

To install PFI run, from the directory paranoid/pfi run

go build -i

##Usage
###General
```
pfi (flags) MountLocation PfsInitLocation PfsBinaryPath
```

* MountLocation - The path to the folder where you wish to mount the paranoid file system. This **does not** need to be an absolute path, it my be relative.
* PfsInitLocation - The path to the folder where you initialised pfs with `pfs init initlocation`.  This **does not** need to be an absolute path, it my be relative.
* PfsBinaryPath - The path to the pfs binary.  This **does not** need to be an absolute path, it my be relative.

example :
```
pfi -log /home/mladen/coding/mountDir ~/coding/initDir bin/pfs
```

###Flags
-log

Prints logs in standard output.
