package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	pb "github.com/cpssd/paranoid/proto/fileserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

//Removes files from Paranoid File Server
func Unserve(c *cli.Context) {
	args := c.Args()

	if len(args) < 2 {
		cli.ShowCommandHelp(c, "unserve")
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Unable to get information on current user:", err)
		Log.Fatal("Could not get user information:", err)
	}

	ip, port, uuid, pool := getFsMeta(usr, args[0])

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
	filePath, err := filepath.Abs(args[1])
	if err != nil {
		fmt.Println("Could Not get path to file", args[1])
		Log.Fatal("Failed to get path to file", err)
	}
	serverClient := pb.NewFileserverClient(connection)
	response, err := serverClient.UnServeFile(context.Background(),
		&pb.UnServeRequest{
			Uuid:     uuid,
			Pool:     pool,
			FilePath: filePath,
		})
	if err != nil {
		fmt.Println("Unable to remove to Discovery Share Server")
		Log.Fatal("Couldn't remove file from Discovery Share Server", err)
	}

	fmt.Println(response.ServeResponse)
}
