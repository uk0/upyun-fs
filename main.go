package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"upyun-fs/fsv2"
)

//打开数据库
func init() {
	var err error
	//数据存储路径和一些初始文件
	db, err := leveldb.OpenFile("./levelDB/db", nil)
	if err != nil {
		log.Fatalln(err)
	}
	fsv2.FSsysTemp.DB = db
}
func main() {

	cg := fsv2.ParseFromCmd()
	fsv2.Run(cg);
}
