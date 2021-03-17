package httpserver

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const (
	slashTag = '/'
	slashStr = "/"
)

var (
	forwarder *forwardManager = nil
)

type forwardEntry func(w http.ResponseWriter, r *http.Request)

type prefixRouter struct {
	keyPath string
	regexp  *regexp.Regexp
	fn      forwardEntry
}

func newPrefixRouter(path string, reg *regexp.Regexp, fn forwardEntry) *prefixRouter {
	return &prefixRouter{
		keyPath: path,
		regexp:  reg,
		fn:      fn,
	}
}

func (pfr *prefixRouter) match(str string) bool {
	return pfr.regexp.MatchString(str)
}

type forwardManager struct {
	pathTable       map[string]forwardEntry
	pathPrefixTable []*prefixRouter
	tableLock       *sync.RWMutex
}

func newForwarder() *forwardManager {

	return &forwardManager{
		pathTable:       make(map[string]forwardEntry),
		pathPrefixTable: make([]*prefixRouter, 0, 1),
		tableLock:       new(sync.RWMutex),
	}
}

func init() {
	if nil == forwarder {
		forwarder = newForwarder()
	}
}

func isPathValid(path string) (bool, error) {
	if path[0] != slashTag {
		return false, fmt.Errorf("path must start with a slash, got %q", path)
	}

	return true, nil
}

func compilePrefixPathRegexp(template string) (*regexp.Regexp, error) {
	regStr := fmt.Sprintf("^%s", template)
	if strings.HasSuffix(regStr, slashStr) == false {
		regStr = fmt.Sprintf("%s%s", regStr, slashStr)
	}

	return regexp.Compile(regStr)
}

func HandlePrefixPathFunc(path string, forward func(w http.ResponseWriter, r *http.Request)) error {

	if valid, err := isPathValid(path); !valid {
		return err
	}

	reg, err := compilePrefixPathRegexp(path)
	if err != nil {
		return err
	}

	return forwarder.register(path, forward, reg)
}

func HandleFunc(keyWord string, forward func(w http.ResponseWriter, r *http.Request)) error {
	return RegisterForward(keyWord, forward)
}

func RegisterForward(path string, forward func(w http.ResponseWriter, r *http.Request)) error {

	if valid, err := isPathValid(path); !valid {
		return err
	}

	return forwarder.register(path, forward, nil)
}

func UnRegisterForward(keyWord string) error {

	if nil == forwarder {
		forwarder = newForwarder()
	}

	return forwarder.unRegister(keyWord)
}

func (fm *forwardManager) register(path string, forward func(w http.ResponseWriter, r *http.Request), reg *regexp.Regexp) error {

	fm.tableLock.Lock()
	defer fm.tableLock.Unlock()

	if reg != nil {
		fm.pathPrefixTable = append(fm.pathPrefixTable, newPrefixRouter(path, reg, forward))
		return nil
	}

	_, ok := fm.pathTable[path]
	if ok {
		return fmt.Errorf("path:%s had be registered", path)
	}

	fm.pathTable[path] = forward
	return nil
}

func (fm *forwardManager) unRegister(keyWord string) error {

	fm.tableLock.Lock()
	defer fm.tableLock.Unlock()

	delete(fm.pathTable, keyWord)
	return nil
}

func (fm *forwardManager) router(keyWord string, w http.ResponseWriter, r *http.Request) {

	handle, find := fm.getPathTableHandle(keyWord)
	if !find {
		handle, find = fm.getPathPrefixTableHandle(keyWord)
		if !find {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	handle(w, r)
}

func (fm *forwardManager) getPathTableHandle(path string) (forwardEntry, bool) {
	fm.tableLock.RLock()
	defer fm.tableLock.RUnlock()
	handle, find := fm.pathTable[path]
	return handle, find
}

func (fm *forwardManager) getPathPrefixTableHandle(path string) (forwardEntry, bool) {

	fm.tableLock.RLock()
	defer fm.tableLock.RUnlock()
	handle, find := forwardEntry(nil), false
	for _, preUnit := range fm.pathPrefixTable {
		if preUnit.match(path) {
			handle, find = preUnit.fn, true
			break
		}
	}
	return handle, find
}

func dispatchHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqTag := r.URL.Path
		if len(reqTag) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		forwarder.router(reqTag, w, r)
	}
}
