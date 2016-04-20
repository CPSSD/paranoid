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
	pfsDir := path.Join(usr.HomeDir, ".pfs", "filesystems", args[0])
	if _, err := os.Stat(pfsDir); err != nil {
		fmt.Printf("%s does not exist. Please call 'paranoid-cli init' before running this command.", pfsDir)
		Log.Fatal("PFS directory does not exist.")
	}

	uuid, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "uuid"))
	if err != nil {
		Log.Error("Error Reading supplied file", err)
		fmt.Println("Error Reading UUID")
		os.Exit(1)
	}

	ip, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "ip"))
	port, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "port"))

	if err != nil {
		fmt.Println("Could not find Ip address of the file server")
		Log.Fatal("Unable to read Ip and Port of discovery server", err)
	}

	address := string(ip) + ":" + string(port)

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
			Uuid: string(uuid),
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
