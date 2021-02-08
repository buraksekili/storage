package service

import (
	"fmt"
	"io"
	"log"

	"github.com/buraksekili/storage/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UploaderServer struct {
	Storage Storage
}

func NewUploaderServer(storage Storage) *UploaderServer {
	return &UploaderServer{storage}
}

func (us *UploaderServer) UploadImage(stream pb.ImageUploader_UploadImageServer) error {
	r, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Unknown, "receive response error: %v", err)
	}

	imgName, ext := r.GetInfo().GetImageName(), r.GetInfo().GetImageExtension()
	if !validate(imgName, ext) {
		return status.Errorf(codes.InvalidArgument, "invalid input image name %s ,extension %s: %v", imgName, ext, err)
	}

	is := 0
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			log.Println("[INFO] reached end of file")
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive data: %v", err)
		}

		ic := r.GetImageChunk()
		lic := len(ic)
		is += lic
		log.Printf("chunk received with size: %d", lic)

		err = us.Storage.Save(ic, ext, imgName)
		if err != nil {
			return status.Errorf(codes.Unknown, "[ERROR] cannot save data: %v", err)
		}
	}

	res := &pb.UploadImageResponse{ImageSize: int32(is)}

	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: %v", err)
	}

	fmt.Printf("[INFO] image with name %s and size %d is saved", imgName, is)
	return nil
}
