package conf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
*
* # 注释
* [section]
* key1 = val
* key2 = val
*
 */

type sectionUnit struct {
	comment     string
	sectionName string
	keyNameList []string
	confData    map[string][]string
}

func newSection(name, comments string) *sectionUnit {
	return &sectionUnit{
		comment:     comments,
		sectionName: name,
		keyNameList: make([]string, 0, 10),
		confData:    make(map[string][]string)}
}

type Configer struct {
	confFileName    string
	sectionNameList []string
	sectionData     map[string]*sectionUnit
}

func newConfiger(file string) (*Configer, error) {
	return constructerConfiger(file), nil
}

func constructerConfiger(file string) *Configer {
	cfg := &Configer{confFileName: file,
		sectionNameList: make([]string, 0, 10),
		sectionData:     make(map[string]*sectionUnit)}
	return cfg
}

func (cfg *Configer) readFile() error {
	fd, err := os.Open(cfg.confFileName)
	if err != nil {
		return err
	}

	defer fd.Close()
	return cfg.read(fd)
}

func (cfg *Configer) read(fd io.Reader) error {
	reader := bufio.NewReader(fd)

	var commentsTmp string
	var sectionTmp string

	for {

		line, err := reader.ReadString('\n')
		if nil != err && err != io.EOF {
			return err
		}

		if io.EOF == err {
			if len(line) <= 0 {
				return nil
			}
		}

		line = strings.TrimSpace(line)

		if 0 == len(line) {
			continue
		}

		switch line[0] {
		case '[':

			if line[len(line)-1] != ']' {
				return errors.New("section format lost")
			}

			sectionTmp = line[1 : len(line)-1]
			if err := cfg.setSection(sectionTmp, commentsTmp); nil != err {
				return err
			}
			commentsTmp = ""

		case ';':
			fallthrough
		case '#':
			commentsTmp = fmt.Sprintf("%s%s%s", commentsTmp, lineDelimiter, line)

		default:

			key, value, err := parseSingleLine(line)
			if nil != err {
				return err
			}
			if err := cfg.setValue(sectionTmp, key, value, commentsTmp); nil != err {
				return err
			}
			commentsTmp = ""
		}

	}
}

func (cfg *Configer) setSection(section, comment string) error {

	sec, exist := cfg.sectionData[section]
	if exist {
		return errors.New("can't load multi same secion " + section)
	}

	if nil == sec {
		sec = newSection(section, comment)
	}

	cfg.sectionNameList = append(cfg.sectionNameList, section)
	cfg.sectionData[section] = sec
	return nil
}

func (cfg *Configer) setValue(section, key, value, comment string) error {

	if 0 == len(section) || 0 == len(key) || 0 == len(value) {
		return errors.New("set value failed")
	}

	sec, exist := cfg.sectionData[section]
	if false == exist {
		return errors.New("not found section[" + section + "]for set value")
	}

	//sec.comment = comment
	if len(sec.confData[key]) == 0 {
		sec.confData[key] = append(sec.confData[key], comment)
		sec.keyNameList = append(sec.keyNameList, key)
	}

	sec.confData[key] = append(sec.confData[key], value)
	cfg.sectionData[section] = sec
	return nil
}

func (cfg *Configer) getValue(section, key string) ([]string, error) {

	sec := cfg.sectionData[section]
	if nil == sec {
		//	fmt.Println("@@@@@@@@@@@@@", section, key)
		return nil, errors.New(p_SECTION_NOT_FOUND)
	}

	if values := sec.confData[key]; len(values) > 0 {
		return values, nil
	}
	//	fmt.Println("!!!@@@@@@@@@@@@@", section, key)
	return nil, errors.New(p_KEY_NOT_FOUND)
}

func (cfg *Configer) GetValueByString(section, key string) (string, error) {

	values, err := cfg.getValue(section, key)
	if err != nil {
		return "", err
	}

	return values[1], nil
}

func (cfg *Configer) GetValueByStringCarryDefault(section, key, defaultValue string) string {

	values, err := cfg.getValue(section, key)
	if err != nil {
		return defaultValue
	}

	return values[1]
}

func (cfg *Configer) GetValueByBool(section, key string) (bool, error) {

	values, err := cfg.getValue(section, key)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(values[1])
}

func (cfg *Configer) GetValueByBoolCarryDefault(section, key string, defaultValue bool) bool {

	values, err := cfg.getValue(section, key)
	if err != nil {
		return defaultValue
	}

	ans, err := strconv.ParseBool(values[1])
	if err != nil {
		return defaultValue
	}
	return ans
}

func (cfg *Configer) GetValueByFloat64(section, key string) (float64, error) {
	values, err := cfg.getValue(section, key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(values[1], 64)
}

func (cfg *Configer) GetValueByInt(section, key string) (int, error) {
	values, err := cfg.getValue(section, key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(values[1])
}

func (cfg *Configer) GetValueByIntCarryDefault(section, key string, defaultValue int) int {
	values, err := cfg.getValue(section, key)
	if err != nil {
		return defaultValue
	}

	num, err := strconv.Atoi(values[1])
	if err != nil {
		return defaultValue
	}

	return num
}

func (cfg *Configer) GetValueByInt64(section, key string) (int64, error) {
	values, err := cfg.getValue(section, key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(values[1], 10, 64)
}

func (cfg *Configer) GetValueByStringList(section, key string) ([]string, error) {
	values, err := cfg.getValue(section, key)
	if err != nil {
		return nil, err
	}
	return values[1:], nil
}
