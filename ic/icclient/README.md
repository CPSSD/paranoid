# icclient

## Importing
```
import "github.com/cpssd/paranoid/ic/icclient"
```

## Sending messages with no data
To send a message to the server that does not contain any date i.e. creat or rename etc..
use `icclient.SendMessage`. if you wanted to send a rename command :
```
icclient.SendMessage("rename", arguments)
```

## Sending messages with data
To send a message to the server that contains data i.e. write
use `icclient.SendMessageWithData`/ if you want to send a write command :
```
icclient.SendMessageWithData("write", arguments, data)
``` 
