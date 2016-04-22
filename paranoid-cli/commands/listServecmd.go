package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	pb "github.com/cpssd/paranoid/proto/fileserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"os/user"
	"strconv"
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

	serverClient := pb.NewFileserverClient(connection)
	response, err := serverClient.ListServer(context.Background(),
		&pb.ListServeRequest{
			Uuid: uuid,
			Pool: pool,
		})
	if err != nil {
		fmt.Println("Unable to send File to Discovery Share Server")
		Log.Fatal("Couldn't message Discovery Share Server", err)
	}
	if len(response.ServedFiles) == 0 {
		fmt.Println("You Have no files being Currently served")
	} else {
		fmt.Printf("%s%s%s%s\n\n",
			fmt.Sprintf("%-60s", "File Path")[:60],
			fmt.Sprintf("%-35s", "FileHash")[:35],
			fmt.Sprintf("%-20s", "Times Accessed")[:20],
			fmt.Sprintf("%-20s", "Expires")[:20],
		)
		for _, files := range response.ServedFiles {
			i, _ := strconv.ParseInt(files.ExpirationTime, 10, 64)
			expiration := time.Unix(i, 0)
			fmt.Printf("%s%s%s%s\n",
				fmt.Sprintf("%-60s", files.FilePath)[:60],
				fmt.Sprintf("%-35s", files.FileHash)[:35],
				fmt.Sprintf("%-20d", files.AccessLimit)[:20],
				fmt.Sprintf("%-20s", expiration)[:20],
			)
		}
	}
}
