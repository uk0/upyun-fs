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
	db1, err := leveldb.OpenFile("./levelDB/info", nil)
	if err != nil {
		log.Fatalln(err)
	}
	//数据存储路径和一些初始文件
	db2, err := leveldb.OpenFile("./levelDB/table", nil)
	if err != nil {
		log.Fatalln(err)
	}
	fsv2.FSsysTemp.DBInfo = db1
	fsv2.FSsysTemp.DBTable = db2
}
func main() {

	cg := fsv2.ParseFromCmd()
	fsv2.Run(cg);
}
