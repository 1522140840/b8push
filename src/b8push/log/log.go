package log

import (
	"fmt"
	log "github.com/cihub/seelog"
	"os"
)

//var logMap map[string]LoggerInterface
/**
wiki:
https://github.com/cihub/seelog/wiki/Constraints-and-exceptions
https://github.com/cihub/seelog/wiki/Dispatchers-and-receivers
https://github.com/cihub/seelog/wiki/Formatting
https://github.com/cihub/seelog/wiki/Logger-types

日志写入文件配置，标示同步，异步，循环异步。
https://github.com/cihub/seelog/wiki/Logger-types-reference
:%Func 带全路径，会造成日志文件比较大，不显示.
其他日志参数，参考：
https://github.com/cihub/seelog-examples/blob/master/formats_main.go

*/

//初始化日志文件.
func InitLog(name string) log.LoggerInterface {
	//name = "amq"
	rootPath, _ := os.Getwd()
	//os.Stderr.WriteString("**********root Path '" + rootPath + "'\n")
	fmt.Println("root path=", rootPath)
	if rootPath != "" {
		logger, err := log.LoggerFromConfigAsFile(rootPath + "/conf/" + name + "_log.xml")

		//os.Stderr.WriteString("**********load Log Config '" + rootPath + "/conf/" + name + "_log.xml\n")
		fmt.Println("load log config, file=", rootPath+"/conf/"+name+"_log.xml")

		if err != nil {
			//os.Stderr.WriteString("Logger " + name + " init failed: " + err.Error() + "\n")
			fmt.Println("load log config error, error=", err)
			return nil
		}

		return logger

	} else {
		logger, err := log.LoggerFromConfigAsString(`<seelog minlevel="debug">
        <outputs formatid="common">
            <console/>
        </outputs>
        <formats>
        <format id="common" format="%LEV %Date/%Time %Msg%n" />
        </formats>
        </seelog>`)

		if err != nil {
			os.Stderr.WriteString("Logger " + name + " init failed: " + err.Error() + "\n")
			return nil
		}

		return logger
	}

}
