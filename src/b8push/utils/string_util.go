package util

import (
	"strings"
)

/**
 *用于校验字符串是否是空串
 */
func StrIsBlank(str string) bool {
	if len(strings.TrimSpace(str)) == 0 {
		return true
	}
	return false
}
