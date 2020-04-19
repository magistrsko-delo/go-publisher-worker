package Worker

import (
	"encoding/json"
	"go-publisher-worker/Models"
	"go-publisher-worker/grpc_client"
	"log"
)

type Worker struct {
	RabbitMQ *RabbitMqConnection
	env *Models.Env
	mediaMetadataGrpcClient *grpc_client.MediaMetadataClient
	mediaChunksGrpcClient *grpc_client.MediaChunksClient
	timeShiftGrpcClient *grpc_client.TimeShiftClient
}

func (worker *Worker) Work()  {
	forever := make(chan bool)

	go func() {
		for d := range worker.RabbitMQ.msgs {
			log.Printf("Received a message: %s", d.Body)

			var publishInput *Models.InputPublishSequence
			err := json.Unmarshal(d.Body, &publishInput)
			if err != nil{
				log.Println(err)
			}

			newMedia := &Models.NewMediaModel{
				MediaId:                  0,
				Name:                     publishInput.Name,
				SiteName:                 publishInput.SiteName,
				Length:                   0,
				Status:                   0,
				Thumbnail:                "",
				ProjectId:                -1,
				AwsBucketWholeMedia:      "",
				AwsStorageNameWholeMedia: "",
			}

			newMediaRsp, err :=  worker.mediaMetadataGrpcClient.CreateMediaMetadata(newMedia)

			log.Println("new media rsp: ", newMediaRsp)


			log.Printf("Done")
			_ = d.Ack(false)
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func InitWorker() *Worker  {
	return &Worker{
		RabbitMQ: 					initRabbitMqConnection(Models.GetEnvStruct()),
		env:      					Models.GetEnvStruct(),
		mediaMetadataGrpcClient: 	grpc_client.InitMediaMetadataGrpcClient(),
		mediaChunksGrpcClient: 		grpc_client.InitChunkMetadataClient(),
		timeShiftGrpcClient: 		grpc_client.InitTimeShiftClient(),
	}

}