package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/buraksekili/storage/proto/pb"

	"google.golang.org/grpc"

	"github.com/buraksekili/storage/service"
)

var path = *flag.String("path", "./img/test.png", "define path of the image.")
var addr = *flag.String("addr", "localhost", "define address of the server.")
var port = *flag.String("port", "8080", "define port")

func main() {
	flag.Parse()

	servAddr := fmt.Sprintf("%s:%s", addr, port)

	if !strings.Contains(path, "/") {
		path = "./" + path
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot read cwd: %v", err)
	}

	l := log.New(os.Stdout, "storage: ", log.LstdFlags)
	storagePath := filepath.Join(cwd, "./img/test-server")
	imgStorage := service.NewLocalImgStorage(storagePath, l)
	us := service.NewUploaderServer(imgStorage)

	grpcServer := grpc.NewServer()
	pb.RegisterImageUploaderServer(grpcServer, us)

	lis, err := net.Listen("tcp", servAddr)
	if err != nil {
		log.Fatalf("cannot start the server: %v", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("cannot start the server: ", err)
	}

	log.Println("[INFO] started listening on: ", servAddr)

}
