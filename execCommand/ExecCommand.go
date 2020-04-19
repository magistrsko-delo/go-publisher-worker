package execCommand

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ExecCommand struct {
}

func (ffmpeg *ExecCommand) ExecFFmpegCommand(arguments []string) error  {
	cmd := exec.Command("ffmpeg", arguments...)
	fmt.Println(cmd.String())
	err := cmd.Run()

	if err != nil {
		fmt.Println("error: ")
		fmt.Println(err)
		return err
	}
	return nil
}

func (ffmpeg *ExecCommand) CreateFilesConcatFile(filePaths []string, filePathOut string) error {

	file, err := os.OpenFile(filePathOut, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
		return  err
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range filePaths {
		_, _ = datawriter.WriteString("file " +  strings.Replace(strings.Replace(data, "\\", "/", -1), "assets/", "", 1)  + "\n")
	}

	_ = datawriter.Flush()
	_ = file.Close()

	return nil
}