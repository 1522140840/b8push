package log

// import (
// 	kafka "apollod/kafka"
// )

var logger = InitLog("static")

const (
	ACK_LOG      = "ack_log"
	BIND_LOG     = "bind_log"
	UNBIND_LOG   = "unbind_log"
	ACTIVE_LOG   = "active_log"
	INACTIVE_LOG = "inactive_log"
	CLICK_LOG    = "click_log"
	LOGIN_LOG    = "login_log"
	UPSTREAM_LOG = "upstream_log"
	LOGOUT_LOG   = "logout_log"
	MSG_LOG      = "msg_log"
)

func WriteToLog(logType string, data ...string) {

	res := ""
	for _, v := range data {
		if v == "" {
			v = "-"
		}
		res = res + v + "\t"
	}
	logger.Info(logType + "\t" + res)

	//将写kafka的工作转移到/search/sogou_push/monitorFileToKafka下
	// kafka.SendToKafka("push_static", logType+"\t"+res)
}
