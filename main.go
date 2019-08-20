package main

import (
	"upyun-fs/config"
	"upyun-fs/service"
)

func main() {

	cg := config.ParseFromCmd()
	service.Service(cg);

}