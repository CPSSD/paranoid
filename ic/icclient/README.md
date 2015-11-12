# icclient

## Importing
```
import "github.com/cpssd/paranoid/ic/icclient"
```

## Sending messages with no data
To send a message to the server that does not contain any date i.e. creat or rename etc.. use `icclient.SendMessage`. The first parameter is pfsDirectory.

**Example :** To send a rename command do the following
```
icclient.SendMessage(pfsDirectory, "rename", arguments)
```

## Sending messages with data
To send a message to the server that contains data i.e. write
use `icclient.SendMessageWithData`. The first parameter is pfsDirectory.

**Example :** To send a write command
```
icclient.SendMessageWithData(pfsDirectory, "write", arguments, data)
```
