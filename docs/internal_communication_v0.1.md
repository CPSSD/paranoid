# Internal Communication
This document outlines how PFSM talks to PFSD

## Technology
PFSM will talk to PFSD through [Unix Domain Sockets](https://en.wikipedia.org/wiki/Unix_domain_socket).

The messages will be encoded with JSON.

## How it works
PFSM receives messages from FUSE in the form of arguments and data.
an example is the write command : The arguments would be something like `write -f /home/mladen/pfsmInitDirectory file.txt 0 2` and the data would be a base 64 encoded string which when decoded will look something like this `"hi"`.


PFSM will need to package this information in a JSON string and send it to PFSD though a unix domain socket. once packaged the JSON will look like this

```
{
    "command" : "write",
    "args" : ["file.txt", "0", "2"],
    "data" : "hi" // base 64 encoded
}
```

Notice that the `pfsmInitDirectory` and `"-f"` arguments are left out. This is because other nodes will have a different pfsmInitDirectory and there is no need for `"-f"`

Different messages may contain no `data` and it will just be `"data" : null`.
