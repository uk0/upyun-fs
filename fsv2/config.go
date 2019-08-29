package fsv2

import (
	"flag"
	"fmt"
	"os"
)

type UpxConfig struct {
	Bucket  string
	Operator  string
	Password  string

}


type Config struct {
	Mountpoint  string
	Debug                bool // 是否是debug模式
	NotExistCacheTimeout int  // 文件不存在会缓存的时间，单位秒

	Upx UpxConfig
}

var config = Config{}

func init() {
	flag.StringVar(&config.Mountpoint, "mp", "", "mountpoint")
	flag.StringVar(&config.Upx.Bucket, "upyun_bkt", "", "Bucket")
	flag.StringVar(&config.Upx.Password, "upyun_pass", "", "Password")
	flag.StringVar(&config.Upx.Operator, "upyun_op", "", "Operator")

	flag.Parse()
}

func ParseFromCmd() Config {

	if config.Mountpoint == "" {
		fmt.Println("Please input mountpoint!")
		os.Exit(-1)
	}

	stat, err := os.Stat(config.Mountpoint)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if !stat.IsDir() {
		fmt.Println("Mountpoint is not a directory!")
		os.Exit(-1)
	}

	if config.Upx.Operator == "" {
		fmt.Println("Please input Upx Operator!")
		os.Exit(-1)

	}
	if config.Upx.Password == "" {
		fmt.Println("Please input Password !")
		os.Exit(-1)
	}

	return config
}
