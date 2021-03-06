package Models

import (
	"fmt"
	"os"
)

var envStruct *Env

type Env struct {
	RabbitUser string
	RabbitPassword string
	RabbitQueueWorker string
	RabbitMqVideoAnalysis string
	RabbitHost string
	RabbitPort string
	Env string
	MediaMetadataGrpcServer string
	MediaMetadataGrpcPort string
	ChunkMetadataGrpcServer string
	ChunkMetadataGrpcPort string
	TimeShiftGrpcServer string
	TimeShiftGrpcPort string
	SequenceServiceServer string
	SequenceServicePort string
	AwsStorageGrpcServer string
	AwsStorageGrpcPort string
}

func InitEnv()  {
	envStruct = &Env{
		RabbitUser:       			os.Getenv("RABBIT_USER"),
		RabbitPassword:   			os.Getenv("RABBIT_PASSWORD"),
		RabbitQueueWorker:      	os.Getenv("RABBIT_QUEUE_WORKER"),
		RabbitMqVideoAnalysis:      os.Getenv("RABBIT_QUEUE_VIDEO_ANALYSIS"),
		RabbitHost:       			os.Getenv("RABBIT_HOST"),
		RabbitPort: 				os.Getenv("RABBIT_PORT"),
		Env: 			  			os.Getenv("ENV"),
		MediaMetadataGrpcServer: 	os.Getenv("MEDIA_METADATA_GRPC_SERVER"),
		MediaMetadataGrpcPort:   	os.Getenv("MEDIA_METADATA_GRPC_PORT"),
		ChunkMetadataGrpcServer:  	os.Getenv("CHUNK_METADATA_GRPC_SERVER"),
		ChunkMetadataGrpcPort:		os.Getenv("CHUNK_METADATA_GRPC_PORT"),
		TimeShiftGrpcServer:		os.Getenv("TIMESHIFT_GRPC_SERVER"),
		TimeShiftGrpcPort:			os.Getenv("TIMESHIFT_GRPC_PORT"),
		SequenceServiceServer:		os.Getenv("SEQUENCE_SERVICE_GRPC_SERVER"),
		SequenceServicePort:		os.Getenv("SEQUENCE_SERVICE_GRPC_PORT"),
		AwsStorageGrpcServer: 		os.Getenv("AWS_STORAGE_GRPC_SERVER"),
		AwsStorageGrpcPort:			os.Getenv("AWS_STORAGE_GRPC_PORT"),
	}
	fmt.Println(envStruct)
}

func GetEnvStruct() *Env  {
	return  envStruct
}