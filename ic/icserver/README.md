# icserver
## Importing
Import like this.
```
import "github.com/cpssd/paranoid/ic/icserver"
```

## Starting the server
To start the server use `icserver.RunServer`.

`RunServer` takes 2 parameters: pfsDirectory and verboseLogging

Specify true for verbose logging if you wish for the server to log actions.
```
go icserver.RunServer("/home/mladen/pfsDir", true)
```

## Listening for messages
To get messages from the server as they come in, attach a listener to the `icserver.MessageChan` channel. The channel is of type `icserver.FileSystemMessage`
```
for {
    select {
        case newMessage := <- icserver.MessageChan :
            // newMessage is of type icserver.FileSystemMessage
            // do something with newMessage
    }
}
```
