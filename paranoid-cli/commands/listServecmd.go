package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	pb "github.com/cpssd/paranoid/proto/fileserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"os/user"
	"time"
)

//Adds files from Paranoid File Server
func ListServe(c *cli.Context) {
	args := c.Args()

	if len(args) < 1 {
		cli.ShowCommandHelp(c, "list-serve")
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Unable to get information on current user:", err)
		Log.Fatal("Could not get user information:", err)
	}

	ip, port, uuid := getFsMeta(usr, args[0])
	address := ip + ":" + port

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(2*time.Second))
	opts = append(opts, grpc.WithInsecure())
	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		fmt.Println("Failed to Connect to Discovery Share Server")
		Log.Fatal("Unable to Connect to Discovery Share Server", err)
	}
	defer connection.Close()

	serverClient := pb.NewFileserverClient(connection)
	response, err := serverClient.ListServer(context.Background(),
		&pb.ListServeRequest{
			Uuid: uuid,
		})
	if err != nil {
		fmt.Println("Unable to send File to Discovery Share Server")
		Log.Fatal("Couldn't message Discovery Share Server", err)
	}
	if len(response.ServedFiles) == 0 {
		fmt.Println("You Have no files being Currently served")
	} else {
		fmt.Println("File Path", "\t", "FileHash", "\t", "Times Accessed", "\t", "Expiration Time", "\n")
		for _, files := range response.ServedFiles {
			fmt.Println(files.FilePath, "\t", files.FileHash, "\t", files.AccessLimit, "\t", files.ExpirationTime)
		}
	}
}
