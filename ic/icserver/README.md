# icserver
## Importing
to import use
```
import "github.com/cpssd/paranoid/ic/icserver"
```

## Starting the server
please start the server in another gorutine.
if you want the server to use verbose logging specify true as a parameter to `icserver.RunServer`
```
go icserver.RunServer(true)
```

## Listening for messages
To get messages from the server as they come in attach a listener to the `icserver.MessageChan` channel. The channel is of type `icserver.FileSystemMessage`
```
for {
    select {
        case newMessage := <- icserver.MessageChan :
            // do something with newMessage
            // newMessage is of type icserver.FileSystemMessage
    }
}
```
