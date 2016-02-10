Key Management in PFS
=====================

## The Key Struct ##

A key is, in essence, a byte array of lengths of 16, 24, or 32 (corresponding to AES-128, 
AES-192, and AES-256, respectively). A key of invalid size will cause most functions to return 
a KeySizeError.

## Key Distribution ##

Our key distribution strategy is the one described in [this paper](https://www.cs.jhu.edu/~sdoshi/crypto/papers/shamirturing.pdf).
The algorithm requires two variables along with the key to be distributed: the total number
of pieces *n*, and the number of pieces required to piece together the original key *k*.

A basic description of the process is:

* Choose a prime *p* larger than both the key and *n*. This prime will publicly known as it is required to
    reconstruct the key.
* Create a random polynomial *q* of degree *k* - 1, with the zeroth coefficient being the key and the
    others being randomly chosen from the range [0, *p*).
* Now, *q*(0) is the key and all other values of *q* (mod *p*) from 1 up to *n* are the key "chunks".
* One chunk is distributed to each node on the network.
* Finally, to calculate the key from the "chunks" we use Lagrange interpolation to calculate the polynomial *q*
    based on what we know about each of the "chunks". If we have at least *k* "chunks", the key will once again be *q*(0).
    Otherwise, we will not have enough data to reliably interpolate the key, and we will get the wrong answer.
