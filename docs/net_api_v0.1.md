# Network API #

## The Architecture ##
The sprint follows a simple network architecture i.e, multiple clients and one central server. The server's only purpose is to echo messages received to all clients.

## Messages
Each message is a JSON object with 2 mandatory fields
* SenderID - (string) A unique identifier for the sender address of the sender of the message
* MessageType (string)
  - write
  - creat
  - link
  - unlink
  - truncate

The rest of the fields of a message depend on the message type

---

### write ###
When a client issues a `write` command it is broadcasted to all connected clients. They then write the data to the files specified
* __fileName__ - (string) The name of the file to write to
* __offset__ - (int) The offset of the write in bytes
* __length__ - (int) The length of the write in bytes
* __writeData__ - (Base64 Encoded String) The date to write as a base64 String

```
{
  "senderID":1
  "messageType:"write",
  "fileName":"helloworld.txt",
  "offset":0,
  "length":8
  "fileData":"aGVsbG8="
}

```

---
### creat ###
Broadcasted by the server to all connected clients and tells them to `creat` a file with the name specified in the `fileName` field
* __fileName__ - (string) The name of the file to create

```
{
  "senderID":1,
  "messageID":1,
  "messageType:"creat",
  "fileName":"helloworld.txt"
}
```

---
### link ###
Creates a link to a file.
* __fileName__ - (string) The name of the file which will be the link
* __target__ - (string) The file that the link is pointing to

```
{
  "senderID":1,
  "messageType":"link",
  "fileName":"hello.txt"
  "target":"helloworld.txt"
}
```

---

### unlink ###
Removes the file that is the link to target as specified in `link`
* __fileName__ - (string) The name of the file to unlink

```
{
  "senderID":1,
  "messageType":"unlink",
  "fileName":"hello.txt"
}
```

---
### truncate ###
Removes the data from a file.
* __fileName__ - (string) The name of the file to truncate
* __length__ - (int) the length in bytes the file will be after the truncate

```
{
  "senderID":1,
  "messageType":"link",
  "fileName":"hello.txt"
  "length":1
}
```
