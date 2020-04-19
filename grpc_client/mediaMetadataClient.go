package grpc_client

import (
	"context"
	"fmt"
	"go-publisher-worker/Models"
	"google.golang.org/grpc"
	"log"
	pbMediaMetadata "go-publisher-worker/proto/media_metadata"
)


type MediaMetadataClient struct {
	Conn *grpc.ClientConn
	client pbMediaMetadata.MediaMetadataClient
}

func (mediaMetadataClient *MediaMetadataClient) UpdateMediaMetadata() (*pbMediaMetadata.MediaMetadataResponse, error)  {
	response, err := mediaMetadataClient.client.UpdateMediaMetadata(context.Background(), &pbMediaMetadata.UpdateMediaRequest{
		MediaId:                  0,
		Name:                     "",
		SiteName:                 "",
		Length:                   0,
		Status:                   0,
		Thumbnail:                "",
		ProjectId:                0,
		AwsBucketWholeMedia:      "",
		AwsStorageNameWholeMedia: "",
		CreatedAt:                0,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (mediaMetadataClient *MediaMetadataClient) CreateMediaMetadata (newMedia *Models.NewMediaModel) (*pbMediaMetadata.MediaMetadataResponse, error)  {

	response, err := mediaMetadataClient.client.NewMediaMetadata(context.Background(), &pbMediaMetadata.CreateNewMediaMetadataRequest{
		Name:                     newMedia.Name,
		SiteName:                 newMedia.SiteName,
		Length:                   newMedia.Length,
		Status:                   newMedia.Status,
		Thumbnail:                newMedia.Thumbnail,
		ProjectId:                newMedia.ProjectId,
		AwsBucketWholeMedia:      newMedia.AwsBucketWholeMedia,
		AwsStorageNameWholeMedia: newMedia.AwsStorageNameWholeMedia,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}


func InitMediaMetadataGrpcClient() *MediaMetadataClient  {
	env := Models.GetEnvStruct()
	fmt.Println("CONNECTING")
	conn, err := grpc.Dial(env.MediaMetadataGrpcServer + ":" + env.MediaMetadataGrpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("END CONNECTION")

	client := pbMediaMetadata.NewMediaMetadataClient(conn)
	return &MediaMetadataClient{
		Conn:    conn,
		client:  client,
	}

}
