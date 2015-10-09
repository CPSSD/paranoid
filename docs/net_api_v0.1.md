# Network API #

## The Architecture ##

Architecture: multiple clients and one central server. The server's only purpose is to echo messages received to all clients.

## Server Messages

None.

## Client Messages
Each message is a JSON object with 2 mandatory fields
* `sender` - (string) A unique identifier for the sender address of the sender of the message
* `type` (string)
  - `creat`
  - `write`
  - `link`
  - `unlink`
  - `truncate`

The rest of the fields of a message depend on the message type

---

### write ###
Writes to a file specified in `name` and Base64 encoded data
* __name__ - (string) The name of the file to write to
* __offset__ - (int) The offset of the write in bytes
* __length__ - (int) The length of the write in bytes
* __data__ - (Base64 Encoded String) The data to write

```
{
  "sender":1
  "type:"write",
  "name":"helloworld.txt",
  "offset":0,
  "length":8
  "data":"aGVsbG8="
}

```

---
### creat ###
Creates a file with the name specified in the `name` field
* __name__ - (string) The name of the file to create

```
{
  "sender":1,
  "messageID":1,
  "type:"creat",
  "name":"helloworld.txt"
}
```

---
### link ###
Creates a link to a file.
* __name__ - (string) The name of the file which will be the link
* __target__ - (string) The file that the link is pointing to

```
{
  "sender":1,
  "type":"link",
  "name":"hello.txt"         // This is the new file name.
  "target":"helloworld.txt"  // This is the existing file name.
}
```

---

### unlink ###
Removes the file that is the link to target as specified in `link`
* __name__ - (string) The name of the file to unlink

```
{
  "sender":1,
  "type":"unlink",
  "name":"hello.txt"
}
```

---
### truncate ###
Removes the data from a file.
* __name__ - (string) The name of the file to truncate
* __offset__ - (int) the offset in bytes the file will be after the truncate

```
{
  "sender":1,
  "type":"link",
  "name":"hello.txt"
  "offset":1
}
```

