package service_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buraksekili/storage/proto/pb"
	"google.golang.org/grpc"

	"github.com/buraksekili/storage/service"
)

func TestUploaderServer_UploadImage(t *testing.T) {
	t.Parallel()

	testStoragePath := "../img/test-client"
	testOutputPath := "../img/test-server"

	l := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	storage := service.NewLocalImgStorage(testOutputPath, l)

	testServAddr := startTestServer(t, storage)
	client := getTestClient(t, testServAddr)

	imgName := "test"
	ext := "png"
	fullImageName := fmt.Sprintf("%s.%s", imgName, ext)

	cwd, err := os.Getwd()
	require.NoError(t, err)
	fullTestImagePath := filepath.Join(cwd, testStoragePath, fullImageName)
	fullOutputImagePath := filepath.Join(cwd, testOutputPath, fullImageName)

	f, err := os.Open(fullTestImagePath)
	require.NoError(t, err)

	defer f.Close()
	require.FileExists(t, fullTestImagePath)

	stream, err := client.UploadImage(context.Background())
	require.NoError(t, err)

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{ImageName: imgName, ImageExtension: ext},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	send := 0

	// read image (fn)
	reader := bufio.NewReader(f)
	chunk := make([]byte, 1000000)
	for {
		n, err := reader.Read(chunk)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		send += n

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ImageChunk{ImageChunk: chunk[:n]},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetImageSize())
	require.EqualValues(t, send, res.GetImageSize())
	require.FileExists(t, fullOutputImagePath)
	require.NoError(t, os.Remove(fullOutputImagePath))
}

func getTestClient(t *testing.T, addr string) pb.ImageUploaderClient {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewImageUploaderClient(cc)
}

func startTestServer(t *testing.T, storage *service.LocalImgStorage) string {
	us := service.NewUploaderServer(storage)

	gRPCServ := grpc.NewServer()
	pb.RegisterImageUploaderServer(gRPCServ, us)

	lis, err := net.Listen("tcp", "")
	require.NoError(t, err)

	go gRPCServ.Serve(lis)

	return lis.Addr().String()
}
