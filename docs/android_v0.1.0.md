Android Research
================

In Sprint 6 one of our tasks was to implement a fully-featured Android client for Paranoid.
Ultimately, we decided to indefinitely postpone this task due its lack of viability. This document
is an assessment of the solutions we researched and their outcomes.

## Cross-compilation of Go to the Android NDK ##

Since all of the existing Paranoid code is written Go, it made sense to attempt to cross-compile
it for the Android NDK, as we would have otherwise had to rewrite it all in Java from scratch.
Go has an official package for mobile development, available [here](https://godoc.org/golang.org/x/mobile).
It describes two ways in which you can build an app with Go:

1.  Using pure native code. In this case all interfaces to the mobile device itself go through the
    mobile package. We decided that this solution was insufficient due to the `mobile` package being
    very immature and unstable, lacking many key features.
2.  Using a hybrid of native Go code compiled for the Android NDK and Java code compiled for the SDK.
    This enables us to write all of the backend code for the project as a single dynamically-linked
    library. This solution is preferable to the first because it means that all of the frontend code
    can be written in Java while it works identically to the desktop client in the background.
    However, this solution is still far from ideal. It only supports libraries, whereas all of our backend
    code that we want to port compiles to a single large binary executable. All of our code would have
    to be heavily refactored into a library, which would likely need to move into a separate repo.
    Another problem is that `gomobile`, the program used to generate Java bindings for Go code, is
    also very immature and places heavy restrictions on the API of our library.

## FUSE in Android ##

Aside from the above difficulties, we discovered that it is impossible to use our backend code
as it currently stands either way, due to the fact that non-rooted Android does not support FUSE.
Many distributions of the Android kernel do not even include the FUSE module, and those that do
lock it down such that the only way to use it is to root the device. We are in unanimous agreement
that we do not want our app only to be available for rooted devices, so we have no choice but to
use a different method for implementing our filesystem backend.

We spent some time researching direct alternatives to FUSE in Android, but none seem to exist.

## Syncing files ##

An alternative method to using FUSE to create a filesystem is to simply use the host filesystem instead.
In this method we would have our own branch of the host's directory tree for our "filesystems".
All files inside would be monitored by `inotify` and any detected changes would be propagated across
the network.

When it comes to failure tolerance, however, it becomes that this solution is not viable.
The benefit of using FUSE is that we can intercept all filesystem calls made and interpret
them however we like. With this method, we are only notified of a filesystem call being made 
after it has already been committed to the host filesystem.

As an example, let's say there is a file on a desktop Paranoid filesystem with read and write permissions
containing the text `Hello, World!`. If this node loses connection to the network for any reason,
our FUSE program will intercept a chmod command changing the permissions to read-only and
-- since we cannot commit the change across the network -- will return with an error
and the file will remain unchanged.

If the same situation were to occur on an Android Paranoid filesystem which has been disconnected
from the network, the host Android filesystem will happily make the change locally before `inotify`
informs our program that a change has been made. Since we cannot commit that change across the network
we have a split-brain scenario. The other nodes on the network have a read-writable version of the file,
while the Android system has a read-only version.

Now if that Android system is reconnected to the network and a change from `Hello, World!` to `Hello, Planet!`
is made on another computer, the Android system will attempt to change the contents of the file, but will
fail due to not having write permissions on the file.

## Read-only client ##

Writing a read-only client for Android is also something we considered. In this, an Android client would
only be able to view the files on a Paranoid network, and would not have permission to edit them.
We decided that this is not worth doing, as, at the time of writing, we are working on a feature
which lets a user host a subset of a node's filesystem as a read-only webserver. Since the two features
are incredibly similar, we decided replicating it as a standalone Android client would be redundant.
