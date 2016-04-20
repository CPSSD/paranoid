package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cpssd/paranoid/paranoid-cli/tls"
	pb "github.com/cpssd/paranoid/proto/fileserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

//Adds files from Paranoid File Server
func Serve(c *cli.Context) {
	args := c.Args()

	if len(args) < 2 {
		cli.ShowCommandHelp(c, "serve")
		os.Exit(1)
	}
	var requestLimit int
	var requestTimeout int
	var err error
	if len(args) < 3 {
		requestLimit = 0
		requestTimeout = 0
	} else if len(args) < 4 {
		requestLimit, err = strconv.Atoi(args[2])
		requestTimeout = 0
		if err != nil {
			fmt.Println("Unable to parse optional paramaters")
			Log.Fatal("Unable to parse optional paramaters:", err)
		}
	} else if len(args) < 5 {
		requestLimit, err = strconv.Atoi(args[2])
		requestTimeout, err = strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("Unable to parse optional paramaters")
			Log.Fatal("Unable to parse optional paramaters:", err)
		}
	}

	file := args[1]

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

	serveFilePath, err := filepath.Abs(file)
	serveData, err := ioutil.ReadFile(serveFilePath)
	uuid, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "uuid"))
	if err != nil {
		Log.Error("Error Reading supplied file", err)
		fmt.Println("Error Reading supplied file:", file)
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
	if tls.CertExists(pfsDir) {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	connection, err := grpc.Dial(address, opts...)
	if err != nil {
		fmt.Println("Failed to Connect to Discovery Share Server")
		Log.Fatal("Unable to Connect to Discovery Share Server", err)
	}
	defer connection.Close()

	serverClient := pb.NewFileserverClient(connection)
	response, err := serverClient.ServeFile(context.Background(),
		&pb.ServeRequest{
			Uuid:     string(uuid),
			FilePath: serveFilePath,
			FileData: serveData,
			Timeout:  int32(requestTimeout),
			Limit:    int32(requestLimit),
		})
	if err != nil {
		fmt.Println("Unable to send File to Discovery Share Server")
		Log.Fatal("Couldn't message Discovery Share Server", err)
	}
	fmt.Println("File now avaliable at:", "http://"+string(ip)+response.ServerPort+"/"+response.ServeResponse)
}
