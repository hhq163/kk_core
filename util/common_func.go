package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/hhq163/logger"
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

func CRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

/**
 * 对字符串进行SHA1哈希
 * @param data string 要加密的字符串
 */
func SHA1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//获取本机IP
func GetServerIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", errors.New("Can not find the client ip address!")

}

func Ip2long(ipstr string) (ip uint32) {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return
	}
	ips := reg.FindStringSubmatch(ipstr)
	if ips == nil {
		return
	}

	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return
	}

	ip += uint32(ip1 * 0x1000000)
	ip += uint32(ip2 * 0x10000)
	ip += uint32(ip3 * 0x100)
	ip += uint32(ip4)

	return
}

func Long2ip(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip>>24, ip<<8>>24, ip<<16>>24, ip<<24>>24)
}

func SafeGo(f interface{}, clog logger.Logger, args ...interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				clog.Error("SafeGo func=", f, ",args=", args, ", err=", err, "stack()=", string(debug.Stack()))

				time.Sleep(10 * time.Second)
				os.Exit(0)
			}
		}()

		v1 := reflect.ValueOf(f)
		params := make([]reflect.Value, len(args))
		for k, v := range args {
			params[k] = reflect.ValueOf(v)
		}
		v1.Call(params)
	}()
}

//安全启动一个go协程，带回调函数
func ReliableGo(f interface{}, callback func(), clog logger.Logger, args ...interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				clog.Error("SafeGo func=", f, ",args=", args, ", err=", err, "stack()=", string(debug.Stack()))

				if callback != nil {
					callback()
				}
				time.Sleep(10 * time.Second)
				os.Exit(0)
			}
		}()

		v1 := reflect.ValueOf(f)
		params := make([]reflect.Value, len(args))
		for k, v := range args {
			params[k] = reflect.ValueOf(v)
		}
		v1.Call(params)
	}()
}

//启动一个go协程，不记录日志到文件
func InstanceGo(f interface{}, args ...interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Fatalln("SafeGo func=", f, ",args=", args, ", err=", err, "stack()=", string(debug.Stack()))
			}
		}()

		v1 := reflect.ValueOf(f)
		params := make([]reflect.Value, len(args))
		for k, v := range args {
			params[k] = reflect.ValueOf(v)
		}
		v1.Call(params)
	}()
}

//安全启动一个go协程，带回调函数，不记录日志到文件
func CallbackGo(f interface{}, callback func(), args ...interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if callback != nil {
					callback()
				}
				time.Sleep(5 * time.Second)

				log.Fatalln("SafeGo func=", f, ",args=", args, ", err=", err, "stack()=", string(debug.Stack()))
			}
		}()

		v1 := reflect.ValueOf(f)
		params := make([]reflect.Value, len(args))
		for k, v := range args {
			params[k] = reflect.ValueOf(v)
		}
		v1.Call(params)
	}()
}

func StrToMd5(source string) string {
	h := md5.New()
	h.Write([]byte(source))
	return hex.EncodeToString(h.Sum(nil))
}

//生成长度为n的随机字符串
func RandStr(n int) string {
	result := make([]byte, n/2)
	rand.Read(result)
	return hex.EncodeToString(result)
}

// 获取当前程序运行目录
func GetExecpath() string {
	execpath, _ := os.Executable() // 获得程序路径
	path := filepath.Dir(execpath)
	return strings.Replace(path, "\\", "/", -1)
}
