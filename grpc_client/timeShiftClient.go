package grpc_client

import (
	"fmt"
	"go-publisher-worker/Models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pbTimeshift "go-publisher-worker/proto/timeshift_service"
	"log"
)

type TimeShiftClient struct {
	Conn *grpc.ClientConn
	client pbTimeshift.TimeshiftClient
}

func (timeShiftClient *TimeShiftClient) GetSequenceChunkInfo(sequenceId int32) (*pbTimeshift.TimeShiftSequenceResponse, error)  {
	response, err := timeShiftClient.client.GetSequenceChunkInformation(context.Background(), &pbTimeshift.TimeShiftSequenceRequest{
		SequenceId:           sequenceId,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}


func InitTimeShiftClient() *TimeShiftClient  {
	env := Models.GetEnvStruct()
	fmt.Println("CONNECTING timeshift client")
	conn, err := grpc.Dial(env.TimeShiftGrpcServer + ":" + env.TimeShiftGrpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("END CONNECTION timeshift client")

	client := pbTimeshift.NewTimeshiftClient(conn)
	return &TimeShiftClient{
		Conn:    conn,
		client:  client,
	}
}