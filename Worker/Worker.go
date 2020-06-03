package Worker

import (
	"encoding/json"
	"fmt"
	"go-publisher-worker/Http"
	"go-publisher-worker/Models"
	"go-publisher-worker/execCommand"
	"go-publisher-worker/grpc_client"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Worker struct {
	RabbitMQ *RabbitMqConnection
	env *Models.Env
	mediaMetadataGrpcClient *grpc_client.MediaMetadataClient
	mediaChunksGrpcClient *grpc_client.MediaChunksClient
	timeShiftGrpcClient *grpc_client.TimeShiftClient
	sequenceGrpcClient *grpc_client.SequenceServiceClient
	awsStorageGrpcClient *grpc_client.AwsStorageClient
	mediaDownloaderClient *Http.MediaDownloader
	execCommand *execCommand.ExecCommand
	rabbitMQMessagePublisher *RabbitMQPublisher
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

			sequenceData, err := worker.sequenceGrpcClient.GetSequenceMedia(publishInput.SequenceId)

			if err != nil {
				log.Println(err)
			}

			_, err = worker.sequenceGrpcClient.UpdateSequence(sequenceData)

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
				fmt.Println(resolution)
				chunks := sequenceChunksResolution[i].GetChunks()
				for j := 0; j < len(chunks); j++ {
					_, err := worker.mediaChunksGrpcClient.LinkMediaChunks(newMediaRsp.GetMediaId(), int32(j), resolution, chunks[j].GetChunkId())
					if err != nil {
						log.Println(err)
						continue
					}

					if resolution == "1920x1080" {
						err = worker.mediaDownloaderClient.DownloadFile("./assets/chunks/" + strconv.Itoa(j) + "_" + chunks[j].GetAwsStorageName(), chunks[j].GetChunksUrl())
						if err != nil {
							log.Println(err)
							continue
						}
					}

				}
			}

			newMediaChunksPaths, err := worker.getFilesPathsInDirectory()

			if err != nil {
				log.Println(err)
			}

			newMediaName := strings.Replace(publishInput.Name, " ", "_", -1)
			err = worker.execCommand.CreateFilesConcatFile(newMediaChunksPaths, "./assets/concatFileList.txt")
			err = worker.execCommand.ExecFFmpegCommand([]string{"-f", "concat", "-safe", "0", "-i", "assets/concatFileList.txt", "-c", "copy", "assets/" + newMediaName + ".ts"})
			err = worker.execCommand.ExecFFmpegCommand([]string{"-i", "assets/" + newMediaName + ".ts", "-acodec", "copy", "-vcodec", "copy", "assets/" + newMediaName + ".mp4"})

			if err != nil {
				log.Println(err)
			}

			awsBucketRsp, err := worker.awsStorageGrpcClient.CreateBucket(strings.Replace(newMediaName, "_", "-", -1))

			if err != nil {
				log.Println(err)
			}
			_, err = worker.awsStorageGrpcClient.UploadMedia("assets/" + newMediaName + ".mp4", awsBucketRsp.GetBucketname(), newMediaName + ".mp4")
			if err != nil {
				log.Println(err)
			}

			videoAnalysisQueueData, err := json.Marshal(Models.VideoAnalysisMessage{MediaId: int(newMediaRsp.MediaId)})
			if err != nil {
				log.Println(err)
			}
			err = worker.rabbitMQMessagePublisher.publishMessageOnQueue(videoAnalysisQueueData)
			if err != nil {
				log.Println(err)
			}

			thumbnailMedia, err := worker.getMediaScreenShot(newMediaName, newMediaRsp.GetMediaId())

			if err != nil {
				log.Println(err)
			}

			////////////////////////////////////////
			newMediaRsp.Status = 3
			newMediaRsp.AwsBucketWholeMedia = awsBucketRsp.GetBucketname()
			newMediaRsp.AwsStorageNameWholeMedia = newMediaName + ".mp4"
			newMediaRsp.Thumbnail = thumbnailMedia
			_, err = worker.mediaMetadataGrpcClient.UpdateMediaMetadata(newMediaRsp)
			if err != nil {
				log.Println(err)
			}

			removePaths := newMediaChunksPaths
			removePaths = append(removePaths, "./assets/concatFileList.txt")
			removePaths = append(removePaths, "assets/" + newMediaName + ".ts")
			removePaths = append(removePaths, "assets/" + newMediaName + ".mp4")
			worker.removeGeneratedFiles(removePaths)

			log.Printf("Done")
			_ = d.Ack(false)
		}

	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (worker *Worker) getMediaScreenShot(newMediaName string, newMediaId int32) (string, error) {
	log.Println("CREATING MEDIA SCREENSHOT")
	imageName := strconv.Itoa(int(newMediaId)) + "-" + strconv.Itoa(rand.Intn(1000000000000)) + "-" + newMediaName + ".jpg"
	err := worker.execCommand.ExecFFmpegCommand([]string{"-ss", "00:00:01", "-i", "./assets/" + newMediaName + ".mp4", "-vframes", "1", "-g:v", "2", "./assets/" + imageName})
	if err != nil {
		log.Println(err)
		return "" , err
	}
	_, err = worker.awsStorageGrpcClient.UploadMedia( "./assets/" + imageName, "mag20-images", imageName)  // TODO for later add this to configuration maybe..
	if err != nil {
		log.Println(err)
		return "" , err
	}

	worker.removeFile("./assets/" + imageName)
	return "v1/mediaManager/mag20-images/" + imageName, nil
}


func (worker *Worker) getFilesPathsInDirectory() ([]string, error) {
	var files []string
	root := "./assets/chunks"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if  !strings.Contains(path, ".gitkeep") && strings.Contains(path, ".ts") {
			files = append(files, path)
		}

		return err
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return files, nil
}

func (worker *Worker) removeGeneratedFiles(paths []string)  {
	for _, path := range paths {
		worker.removeFile(path)
	}
}

func (worker *Worker) removeFile(path string)  {
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}

func InitWorker() *Worker  {
	return &Worker{
		RabbitMQ: 					initRabbitMqConnection(Models.GetEnvStruct()),
		rabbitMQMessagePublisher:   InitRabbitMQVideoAnalysisPublisher(),
		env:      					Models.GetEnvStruct(),
		mediaMetadataGrpcClient: 	grpc_client.InitMediaMetadataGrpcClient(),
		mediaChunksGrpcClient: 		grpc_client.InitChunkMetadataClient(),
		timeShiftGrpcClient: 		grpc_client.InitTimeShiftClient(),
		sequenceGrpcClient: 		grpc_client.InitSequenceServiceMetadata(),
		awsStorageGrpcClient: 		grpc_client.InitAwsStorageGrpcClient(),
		mediaDownloaderClient: 		&Http.MediaDownloader{},
		execCommand: 				&execCommand.ExecCommand{},
	}

}