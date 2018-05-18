package conf

/***
使用开源项目：
https://github.com/msbranco/goconfig
读取配置文件。和windows ini 文件一样.

*/

import (
	"fmt"
	"github.com/msbranco/goconfig"
	//"log"
	"os"
)

/*
/search/push_server
	--bin/push_srv
	--conf/push_conf
	--log
	--
*/
var PushConfig *goconfig.ConfigFile

func init() {
	file, _ := os.Getwd()
	fmt.Println("push server current path:", file)
	//读取当前路径文件。
	config, err := goconfig.ReadConfigFile(file + "/conf/b8push.conf")
	PushConfig = config
	fmt.Println("load config, file=", file+"/conf/b8push.conf")
	if err != nil {
		fmt.Println("loaded config error, err=", err)
	} else {
		fmt.Println("loaded config, config=", PushConfig)
	}
}

//按照类型获得value.
func GetVal(conf_type string, key string) string {
	str, err := PushConfig.GetString(conf_type, key)
	if err != nil {
		panic("Config error: " + err.Error())
	}
	fmt.Println("get config value, type=", conf_type, ", key=", key, ", val=", str)
	return str
}
