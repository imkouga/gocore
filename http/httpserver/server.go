package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	KDefaultPath = "/"
	KDefaultPort = 18080
	KDefaultIP   = "0.0.0.0"
	KDefaultAddr = "0.0.0.0:18080"
)
const (
	ipIndex         = 0
	portIndex       = 1
	addrIndex       = 0
	pathIndex       = 1
	crtFileIndex    = 0
	crtKeyFileIndex = 1
)

var (
	invalidListenAddr         = errors.New("invalid listen address")
	httpServerStartFailed     = errors.New("http server listen start failed")
	invalidCertificateKeyFile = errors.New("invalid certificate key info")
)

func Start(addr, path string, options ...string) error {

	if len(strings.TrimSpace(addr)) <= 0 {
		return invalidListenAddr
	}

	if len(strings.TrimSpace(path)) <= 0 {
		path = KDefaultPath
	}

	if len(options) > 0 {
		return httpsStart(addr, path, dispatchHandler(), options...)
	}

	return httpStart(addr, path, dispatchHandler())
}

// http服务启动入口, 提供监听地址(第一参数)和path(第二参数)
// 若参数为空地址默认使用 "0.0.0.0:18080" 和 "/"
func HttpServerStart(addr ...string) error {

	var (
		serverAddr = KDefaultAddr
		serverPath = KDefaultPath
	)

	switch len(addr) {
	case 0:
	case 1:
		serverAddr = addr[addrIndex]
	case 2:
		serverAddr = addr[addrIndex]
		serverPath = addr[pathIndex]
	}

	if len(strings.TrimSpace(serverPath)) <= 0 {
		serverPath = KDefaultPath
	}
	serverAddr = parseAddressArgument(serverAddr)

	return Start(serverAddr, serverPath)
}

// https服务启动入口, 提供监听地址(第一参数)和证书信息(第二参数)
// 若参数为空地址默认使用 "0.0.0.0:18080" 和 私有本地自签证书
func HttpsServerStart(addr string, options ...string) error {
	var (
		serverAddr = KDefaultAddr
	)

	addr = strings.TrimSpace(addr)
	if len(addr) > 0 {
		serverAddr = addr
	}

	serverAddr = parseAddressArgument(serverAddr)

	if len(options) < 2 {
		if err := dumnCommonCertificateFile(); err != nil {
			return err
		}
		return Start(serverAddr, KDefaultPath, realDumpCertificateFile, realDumpCertificateKeyFile)
	}

	return Start(serverAddr, KDefaultPath, options...)
}

func httpStart(addr, path string, finderHandle func(w http.ResponseWriter, r *http.Request)) error {

	go func() {
		http.HandleFunc(path, finderHandle)
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(fmt.Errorf("%s:%s", httpServerStartFailed, err.Error()))
		}
	}()

	return nil
}

func httpsStart(addr, path string, finderHandle func(w http.ResponseWriter, r *http.Request), options ...string) error {

	if len(options) < 2 {
		return invalidCertificateKeyFile
	}

	crtFile := strings.TrimSpace(options[crtFileIndex])
	keyFile := strings.TrimSpace(options[crtKeyFileIndex])
	if len(crtFile) <= 0 || len(keyFile) <= 0 {
		return invalidCertificateKeyFile
	}

	http.HandleFunc(path, finderHandle)
	go func() {
		if err := http.ListenAndServeTLS(addr, crtFile, keyFile, nil); err != nil {
			panic(fmt.Sprintf("%s:%s", httpServerStartFailed, err))
		}
	}()

	return nil
}

func parseAddressArgument(addr string) string {

	addr = strings.TrimSpace(addr)

	// 当参数为空时, 直接使用默认的地址进行监听
	if len(addr) <= 0 {
		return KDefaultAddr
	}

	// 当参数不包含:时, 默认优先包含端口处理
	if false == strings.Contains(addr, ":") {
		return fmt.Sprintf("%s:%s", KDefaultIP, addr)
	}

	addrs := strings.Split(addr, ":")
	ip := strings.TrimSpace(addrs[ipIndex])
	port := strings.TrimSpace(addrs[portIndex])
	if len(ip) <= 0 {
		ip = KDefaultIP
	}

	if len(port) <= 0 {
		port = fmt.Sprintf("%d", KDefaultPort)
	}

	return fmt.Sprintf("%s:%s", ip, port)
}
