package main

import (
	"fmt"
	"go-publisher-worker/Worker"
)

func init()  {
}

func main()  {
	fmt.Println("main")
	Worker.InitRabbitMqConnection()
}
