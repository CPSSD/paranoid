# Net API
## The Architecture
The sprint follows a simple network architecture i.e, multiple clients and one central server. The server's only purpose is to echo messages received to all clients.

## Messages
Each message is a JSON object with 3 mandatory fields
* SenderID - (String) A unique identifier for the sender address of the sender of the message
* MessageID - (String) A unique identifier for the message
* MessageType - (String) "write"/"read"/"create"/"link"/"unlink"/"truncate"

The rest of the fields of a message depend on the message type

#### "write"
* FileName - (String) The name of the file to write to
* Offset - (Int) The offset of the write in bytes
* Length - (Int) The length of the write in bytes
* WriteData - (Base64 String) The date to write as a base64 String

#### "read"
* FileName - (String) The name of the file to read

#### "create"
* FileName - (String) The name of the file to create

#### "link"
* FileName - (String) The name of the file to link
* Target - (String) The Target

#### "unlink"
* FileName - (String) The name of the file to unlink

#### "truncate"
* FileName - (String) The name of the file to truncate
* Length - (Int) the length (bytes) the file will be after the truncate
