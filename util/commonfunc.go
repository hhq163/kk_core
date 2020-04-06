package util

import (
	"crypto/md5"
	"fmt"
	"io"
)

//公共方法

//是否为密码规定字符
func IsHaveInvalidPwdChar(strChar string) bool {
	bReuslt := false
	ilen := len(strChar)

	for i := 0; i < ilen; i++ {
		if (strChar[i] < '0' || strChar[i] > '9') && (strChar[i] < 'A' || strChar[i] > 'Z') && (strChar[i] < 'a' || strChar[i] > 'z') {
			bReuslt = true
			break
		}
	}
	return bReuslt
}

/**
 * 对字符串进行MD5哈希
 * @param data string 要加密的字符串
 */
func MD5(data string) string {
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//将游戏前端登录类型转换成web后台维护类型
func ConvertLoginDeviceId(loginType uint8) int {
	deviceType := 0
	switch loginType {
	case 0: // PC_FLASH
		deviceType = 1
	case 3:
		deviceType = 2
	case 1:
		deviceType = 3
	case 5, 6: // NEW_APP密码登录和手势登录
		deviceType = 3
	case 4: // PC_H5
		deviceType = 1
	}
	return deviceType
}
