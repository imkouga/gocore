package kits

import (
	"reflect"
	"runtime"
	"strings"
)

// 获取i定义模块的模块名(包名)
func GetPacketName(i interface{}) string {
	return getPacketNameByFuncName(i)
}

func getPacketNameByFuncName(i interface{}) string {

	val := reflect.ValueOf(i)
	if k := val.Kind(); k != reflect.Func {
		return ""
	}

	// 获取函数长名称(包含包名) 如:
	//"dcommon/functioner.TestRuntimeFunInnerMostPacketName"
	//"dcommon/functioner.TestGetPacketName.func1"
	//"strings.Compare"
	fn := runtime.FuncForPC(val.Pointer()).Name()

	// strings.Index函数没找到对应字符串的idx时, 返回-1
	dotIdx := strings.Index(fn, dotStr)
	if dotIdx == -1 {
		return ""
	}
	slashIdx := strings.Index(fn, slashStr)
	return fn[slashIdx+1 : dotIdx]
}

// 获取调用函数的函数全名
func RuntimeFuncName() string {

	ptr := make([]uintptr, 1)
	runtime.Callers(2, ptr)

	fn := runtime.FuncForPC(ptr[0])
	return fn.Name()
}

// 获取调用函数的函数短名称
func RuntimeFunNameByShortcut() string {

	ptr := make([]uintptr, 1)
	runtime.Callers(2, ptr)

	fn := runtime.FuncForPC(ptr[0])
	fnameSlice := strings.Split(fn.Name(), dotStr)
	return fnameSlice[len(fnameSlice)-1]
}

// 获取调用函数所在的顶层模块名称(顶层包名)
func RuntimeFunTopPacketName() string {
	ptr := make([]uintptr, 1)
	runtime.Callers(2, ptr)

	fn := runtime.FuncForPC(ptr[0])
	fnameSlice := strings.Split(fn.Name(), slashStr)
	return fnameSlice[0]
}

// 获取调用函数所在的内层模块名称(最内层包名)
func RuntimeFunInnermostPacketName() string {
	ptr := make([]uintptr, 1)
	runtime.Callers(2, ptr)

	fn := runtime.FuncForPC(ptr[0])
	fnameSlice := strings.Split(fn.Name(), slashStr)
	pSlice := strings.Split(fnameSlice[len(fnameSlice)-1], dotStr)
	return pSlice[0]
}
