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
pfi [flags] <PfsInitLocation> <MountLocation>
```

* MountLocation - The path to the folder where you wish to mount the paranoid file system.
* PfsInitLocation - The path to the folder where you initialised pfs with `pfs init initlocation`.

example :
```
pfi -v ~/coding/initDir /home/mladen/coding/mountDir
```

###Flags
-v

Prints logs in standard output.
