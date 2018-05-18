package conf

/***
使用开源项目：
https://github.com/msbranco/goconfig
读取配置文件。和windows ini 文件一样.

*/

import (
	. "../conf"
	"testing"
)

func TestGet(t *testing.T) {
	//获得当前路径。
	//	file, _ := os.Getwd()
	//	fmt.Println("current path:", file)
	//	//读取当前路径文件。
	//	config, err := goconfig.ReadConfigFile(file + "/conf/push.conf")
	//	fmt.Println(err)
	//	if err == nil {
	//		fmt.Println(config.GetString("service-1", "url"))
	//	}
	port := GetVal("activemq", "port")
	if port != "61613" {
		t.Fatalf("port is not 61613 is:" + port)
	}
}
