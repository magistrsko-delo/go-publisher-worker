package grpc_client

import (
	"fmt"
	"go-publisher-worker/Models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"

	pbMediaChunks "go-publisher-worker/proto/media_chunks_metadata"
)

type MediaChunksClient struct {
	Conn *grpc.ClientConn
	client pbMediaChunks.MediaMetadataClient
}

func (mediaChunksClient *MediaChunksClient) LinkMediaChunks() (*pbMediaChunks.LinkMediaChunkResponse, error)  {

	response, err := mediaChunksClient.client.LinkMediaWithChunk(context.Background(), &pbMediaChunks.LinkMediaWithChunkRequest{
		MediaId:              0,
		Position:             0,
		Resolution:           "",
		ChunkId:              0,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}


func InitChunkMetadataClient() *MediaChunksClient  {
	env := Models.GetEnvStruct()
	fmt.Println("CONNECTING chunks metadata")

	conn, err := grpc.Dial(env.ChunkMetadataGrpcServer + ":" + env.ChunkMetadataGrpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("END CONNECTION chunk metadata")

	client := pbMediaChunks.NewMediaMetadataClient(conn)
	return &MediaChunksClient{
		Conn:    conn,
		client:  client,
	}

}
