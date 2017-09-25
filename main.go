package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/kelseyhightower/envconfig"
	pb "github.com/nhite/pb-backend"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type configuration struct {
	ListenAddress  string `envconfig:"LISTEN_ADDRESS" required:"true" default:"127.0.0.1:1234"`
	MaxMessageSize int    `envconfig:"MAX_RECV_MSG_SIZE" required:"true" default:"16500545"`
	CertFile       string `envconfig:"CERT_FILE" required:"true"`
	KeyFile        string `envconfig:"KEY_FILE" required:"true"`
	WorkingDir     string `envconfig:"WORKING_DIR" default:"/tmp" required:"true"`
	FileExtension  string `envconfig:"FILE_EXTENSION" default:".nhite" required:"false"`
	PermMode       uint32 `envconfig:"PERM_MODE" default:"0600" required:"true"`
}

const envPrefix = "N_LOCAL_BACKEND"

var (
	config configuration
	// Build date
	Build string
	// Version number
	Version     string
	versionFlag bool
)

func main() {
	flag.BoolVar(&versionFlag, "v", false, "Display version then exit")
	flag.Parse()
	if versionFlag {
		if Version == "" {
			Version = "dev"
		}
		fmt.Printf("%v version %v, build %v\n", os.Args[0], Version, Build)
		os.Exit(0)
	}
	if len(os.Args) > 1 {
		envconfig.Usage(envPrefix, &config)
		os.Exit(1)

	}
	err := envconfig.Process(envPrefix, &config)
	if err != nil {
		envconfig.Usage(envPrefix, &config)
		fmt.Println(err)
		os.Exit(1)
	}

	log.Println("Listening on " + config.ListenAddress)
	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
	if err != nil {
		log.Fatal("could not load TLS keys: ", err)
	}
	server := grpc.NewServer(grpc.Creds(creds), grpc.MaxRecvMsgSize(config.MaxMessageSize))

	pb.RegisterBackendServer(server, &backend{
		workingDir:    config.WorkingDir,
		fileExtension: config.FileExtension,
		permission:    os.FileMode(config.PermMode),
	})
	log.Fatal(server.Serve(listener))
}
