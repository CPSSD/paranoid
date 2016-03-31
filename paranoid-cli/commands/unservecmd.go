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

	pfsDir := path.Join(usr.HomeDir, ".pfs", "filesystems", args[0])
	if _, err := os.Stat(pfsDir); err != nil {
		fmt.Printf("%s does not exist. Please call 'paranoid-cli init' before running this command.", pfsDir)
		Log.Fatal("PFS directory does not exist.")
		os.Exit(1)
	}
	uuid, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "uuid"))
	if err != nil {
		Log.Error("Error Reading UUID file:", err)
		fmt.Println("Error Reading Unique ID")
		os.Exit(1)
	}

	ip, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "ip"))
	port, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "port"))

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

	serverClient := pb.NewFileserverClient(connection)
	response, err := serverClient.UnServeFile(context.Background(),
		&pb.UnServeRequest{
			Uuid:     string(uuid),
			FileHash: args[1],
		})
	if err != nil {
		Log.Error("Couldn't message Discovery Share Server", err)
		fmt.Println("Unable to send File to Discovery Share Server", err)
		os.Exit(1)
	}

	fmt.Println(response.ServeResponse)
}
