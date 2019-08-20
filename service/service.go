package service

import (
	"log"
	"upyun-fs/config"
	"upyun-fs/fs"
)

func Service(cg config.Config) {
	if err := fs.Run(cg); err != nil {
		log.Fatal(err)
	}
}