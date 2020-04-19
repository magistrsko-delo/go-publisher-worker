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
	sequenceGrpcClient *grpc_client.SequenceServiceClient
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

			if err != nil{
				log.Println(err)
			}

			log.Println("new media rsp: ", newMediaRsp)

			sequenceTimeShiftData, err := worker.timeShiftGrpcClient.GetSequenceChunkInfo(publishInput.SequenceId)
			if err != nil{
				log.Println(err)
			}

			sequenceChunksResolution := sequenceTimeShiftData.GetData()
			for i := 0; i < len(sequenceChunksResolution); i++ {
				resolution := sequenceChunksResolution[i].GetResolution()

				chunks := sequenceChunksResolution[i].GetChunks()
				for j := 0; j < len(chunks); j++ {
					_, err := worker.mediaChunksGrpcClient.LinkMediaChunks(newMediaRsp.GetMediaId(), int32(j), resolution, chunks[j].GetChunkId())
					if err != nil {
						log.Println(err)
					}
				}
			}

			newMediaRsp.Status = 3
			// TODO for next version download new sequence chunks.. join then and create new media for download on aws...


			_, err = worker.mediaMetadataGrpcClient.UpdateMediaMetadata(newMediaRsp)
			if err != nil {
				log.Println(err)
			}

			sequenceData, err := worker.sequenceGrpcClient.GetSequenceMedia(sequenceTimeShiftData.GetSequenceId())

			if err != nil {
				log.Println(err)
			}
			_, err = worker.sequenceGrpcClient.UpdateSequence(sequenceData)

			if err != nil {
				log.Println(err)
			}

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
		sequenceGrpcClient: 		grpc_client.InitSequenceServiceMetadata(),
	}

}