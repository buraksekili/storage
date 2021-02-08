package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/buraksekili/storage/proto/pb"
	"google.golang.org/grpc"
)

func uploadImage(c pb.ImageUploaderClient, path, imgName, ext string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("couldn't get cwd: %v", err)
	}

	// open image and store it in f
	path = filepath.Join(cwd, path)
	fmt.Println("path is: ", path)

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("cannot open image: %v", err)
	}
	defer f.Close()

	stream, err := c.UploadImage(context.Background())
	if err != nil {
		log.Fatalf("cannot upload image: %v", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{ImageName: imgName, ImageExtension: ext},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatalf("cannot send image info req: %v", err)
	}

	// read image (fn)
	reader := bufio.NewReader(f)
	chunk := make([]byte, 1000000)
	for {
		n, err := reader.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("[ERROR] cannot read chunk: %v", err)
		}
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ImageChunk{ImageChunk: chunk[:n]},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatalf("[ERROR] cannot send chunk: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("[ERROR] cannot receive close res: %v", err)
	}

	log.Printf("[INFO] image with size %v uploaded", res.GetImageSize())
}

var path = flag.String("path", "./img/test-client/test.png", "define path of the image.")
var addr = flag.String("addr", "localhost", "define address of the server.")
var port = flag.String("port", "8080", "define port")

func main() {
	flag.Parse()

	servAddr := fmt.Sprintf("%s:%s", *addr, *port)

	if !strings.Contains(*path, "/") {
		*path = "./" + *path
	}

	i := strings.LastIndex(*path, "/")
	if i == -1 {
		log.Fatalf("invalid image path as: %s", *path)
	}

	dotIdx := strings.LastIndex(*path, ".")
	if i == -1 {
		log.Fatalf("invalid image path as: %s", *path)
	}

	imgName := (*path)[i+1 : dotIdx]
	ext := (*path)[dotIdx+1:]

	cc, err := grpc.Dial(servAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[ERROR] couldn't establish connection: %v", err)
	}
	fmt.Println("dial: ", servAddr)

	lc := pb.NewImageUploaderClient(cc)

	uploadImage(lc, *path, imgName, ext)

}
