Paranoid Virtual Disk V0.01
=====================

#### File System Structure

The file system will initially be a flat file system. (no dirs only files)

In this sprint the i-node and the file data are the same thing. The node
contains the file data as well as a UUID for referencing the file.

#### Generic structure
```
-                 Root
-       File UUID      File UUID2
-       File Data      File Data2
```

File Structure
```
{
  nodeUUID: 1234567890
  file_name: asdf
  file_metadata: []
  file_content: (base64 encoded data?)
}
```

For the file system the UUID and filename must be unique.

#### Supported Functions

- create

      Create a file with a new UUID and Name

- write

      Update data in a currently existing file

- unlink

      Remove the file from the FS

