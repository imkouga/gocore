package kits

import (
	"crypto/md5"
	"encoding/hex"
	"net"
	"strconv"
	"strings"
)

const (
	DefaultSepChar = "@"
	dotStr         = "."
	slashStr       = "/"
	errorJoinChar  = "|"
)

// string类型的数组合并为一个字符串
func Array2String(arr []string, sep string) string {
	if len(arr) <= 0 {
		return ""
	}

	if len(sep) <= 0 {
		sep = DefaultSepChar
	}

	return strings.Join(arr, sep)
}

// string类型的数组转换成map string:bool
func StringArray2MapStringBool(arr []string) map[string]bool {

	if len(arr) <= 0 {
		return nil
	}

	arrMap := make(map[string]bool)

	for _, str := range arr {
		str = strings.TrimSpace(str)
		if len(str) <= 0 {
			continue
		}
		arrMap[str] = true
	}

	return arrMap
}

// 判断字符串是否为空
func StrIsEmpty(str string) bool {
	str = strings.TrimSpace(str)
	return len(str) == 0
}

// 判断是否为合法的IPv4地址字符串
func IsIPv4Legal(addr string) bool {
	addr = strings.TrimSpace(addr)

	ip := net.ParseIP(addr)
	if nil == ip {
		return false
	}

	return true
}

// 判断是否为合法的端口号
func IsPortLegal(port string) bool {

	port = strings.TrimSpace(port)

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	if portNum > 65536 || portNum < 1 {
		return false
	}

	return true
}

// 判断给定的字符串是否在字符串切片中存在
func IsInSlice(slice []string, key string) bool {

	for _, value := range slice {
		if key == value {
			return true
		}
	}
	return false
}

// 判断给定的Int是否在Int切片中存在
func IsInSliceInt(slice []int, key int) bool {

	for _, value := range slice {
		if key == value {
			return true
		}
	}
	return false
}

// 获取两个字符串切片的交集;
func SliceIntersect(slice1, slice2 []string) []string {

	inter := make([]string, 0)
	for _, v := range slice1 {
		if IsInSlice(slice2, v) {
			inter = append(inter, v)
		}
	}
	return inter
}

// 判断两个字符串切片是否有交集
func IsIntersectInSlices(slice1, slice2 []string) bool {

	for _, v := range slice1 {
		if IsInSlice(slice2, v) {
			return true
		}
	}
	return false
}

// 删除字符串切片内重复的元素
func RemoveRepeatedElementAtStringSlice(arr []string) []string {

	newArr := make([]string, 0)
	for i := 0; i < len(arr); i++ {

		repeat := false

		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}

		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}

	return newArr
}

// 计算生成指定字节序列对应的MD5
func GenerateMD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// 计算生成指定字符串对应的MD5
func GenerateMD5ForString(str string) string {
	return GenerateMD5([]byte(str))
}

//  获取两个字符串切片合并后的副本
func MergeAndCopyStringSlice(a, b []string) []string {
	tmp := make([]string, 0, len(a)+len(b))
	tmp = append(tmp, a...)
	tmp = append(tmp, b...)
	return tmp
}

// 创建字符串map副本
func CopyStringMap(src map[string]string) map[string]string {
	dst := make(map[string]string)
	for key, val := range src {
		dst[key] = val
	}

	return dst
}

// 创建字符串slice副本
func CopyStringSlice(src []string) []string {
	dst := make([]string, 0, len(src))
	copy(dst, src)
	return dst
}

// 错误信息拼接
func JoinAllErrInfo(errs []error) string {

	errStrs := make([]string, 0, len(errs))
	for _, err := range errs {
		errStrs = append(errStrs, err.Error())
	}

	return strings.Join(errStrs, errorJoinChar)
}
