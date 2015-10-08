Paranoid Virtual Disk V0.1
=====================

#### File System Structure

The file system will initially be a flat file system. (no dirs only files)

In this sprint the i-node and the file data are the same thing. The node
contains the file data as well as a UUID for referencing the file.

#### Generic structure
```
-                 Root
-           Node        Node
```
### Links
  nodeUUID: 123456789
  linkname: adsf

Node Structure
```
{
  nodeUUID: 1234567890
  file_name: asdf
  file_metadata: []
  length: 12
  offset: 1
  file_content: asdf123
}
```

For the file system the UUID and filename must be unique.

#### Supported Functions

- creat (filename)

      Create a file with a new UUID and Name

- write (filename data)

      Update data in a currently existing file

- link (file nodeID)

      Creates a link between a node and a file

- unlink (filename)

      Removes Link to a file. If no other links to a file exist file
      it is removed

- truncate (filename size)

      Cuts the file down to a specific size

- link (filename )

      creates a link between a file and a node

