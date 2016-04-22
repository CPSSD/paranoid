package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	pb "github.com/cpssd/paranoid/proto/fileserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"os/user"
	"path"
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
		Log.Error("Could not get user information:", err)
		fmt.Println("Unable to get information on current user:", err)
		os.Exit(1)
	}
	ip, port, uuid := getFsMeta(args[0])

	if err != nil {
		Log.Error("Unable to read Ip and Port of discovery server", err)
		fmt.Println("Could not find Ip address of the file server")
		os.Exit(1)
	}

	address := string(ip) + ":" + string(port)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(2*time.Second))
	opts = append(opts, grpc.WithInsecure())
	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		Log.Error("Unable to Connect to Discovery Share Server", err)
		fmt.Println("Failed to Connect to Discovery Share Server")
		os.Exit(1)
	}
	defer connection.Close()
	filePath, err := filepath.Abs(args[1])
	if err != nil {
		Log.Error("Failed to get path to file", err)
		fmt.Println("Could Not get path to file", args[1])
	}
	serverClient := pb.NewFileserverClient(connection)
	response, err := serverClient.UnServeFile(context.Background(),
		&pb.UnServeRequest{
			Uuid:     uuid,
			FilePath: filePath,
		})
	if err != nil {
		Log.Error("Couldn't remove file from Discovery Share Server", err)
		fmt.Println("Unable to remove to Discovery Share Server")
		os.Exit(1)
	}

	fmt.Println(response.ServeResponse)
}
