Key Distribution Planning
=========================

*This document is the product of a collaboration between Terry Bolt and Conor Griffin.*

Our key distribution plan has to cover 3 things:

* Physically transferring the key pieces.
* Keeping track of what key pieces were shared successfully, who has the pieces, and to whom they belong.
* Dealing with a failure to distribute enough keys.


## Key Piece Transferral ##

This part of key distribution is already complete. We use Shamir’s Secret Sharing algorithm to create a set of shares
(we call them KeyPieces) which can be interpolated together to get the original key. That way chunking keys is
a relatively simple process, as the alternative requires a huge number of keys and shares.

## Tracking Shared Pieces ##

Physically transferring the shares works fine, but there is currently no way to track what shares were
successfully transferred and no way to know if enough shares have been distributed to later unlock the system.
We propose the creation of a state machine designed solely to track what pieces are held by what nodes.
This will be integrated into Raft so that there is a record of what machines have what shares in the Raft log.
It will not track the contents of the shares, as including those in the log would be a security issue.
Instead it will only track what node the share belongs to, what node is currently holding it, and how many
of its sibling shares are also being held.

## Failure Handling ##
Mostly discussed in the Continuous Key Sharing Plan. If a node fails to distribute its keys correctly then it is
not allowed to join the network and a new generation of shares is not created.

## Continuous Key Sharing Plan ##

If we share keys continuously, this greatly reduces the complexity of failure handling.
This means a node losing connection to the internet does not mean it can no longer share its keys and loses data.
This plan means that there is a lot of additional complexity whenever a node attempts to join
the cluster, however during normal operation it works pretty simply. 

Steps a node must complete to join the network:

1. Distribute its key chunks to a majority of nodes.
2. Receive key chunks from a majority of nodes
3. Start a new key generation that a majority of nodes must share a majority of keys in. Previous key generations are kept by nodes until all nodes share their keys in a new generation. 

Once a node rejoins the cluster, as it catches up with the raft log, it will catch up on key generations
and distribute its keys as part of the current generation.

A problem with this method is that nodes do not always need a majority of keys to unlock their contents,
as if a node is offline during multiple nodes joining the cluster, they will only need a majority of their
previously distributed keys. However I do not think this problem is solvable without requiring all nodes
to agree on a configuration change.
