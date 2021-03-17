package conf

import (
	"errors"
	"strings"
)

func LoadConfig(file string) error {

	if len(strings.TrimSpace(file)) <= 0 {
		file = defaultConfiger.confFileName
	}

	cfg := constructerConfiger(file)
	err := cfg.readFile()
	if nil != err {
		return err
	}

	installConfiger(cfg)
	return err
}

func Reload() error {
	return LoadConfig(defaultConfiger.confFileName)
}

func LoadAndGetConfiger(file string) (*Configer, error) {
	cfg := constructerConfiger(file)
	err := cfg.readFile()

	return cfg, err
}

func parseSingleLine(line string) (string, string, error) {

	strs := strings.Split(line, confDelimiter)
	if confLineCount > len(strs) {
		return "", "", errors.New("config " + line + " parse failed")
	}

	if confLineCount == len(strs) {
		return strings.TrimSpace(strs[0]), strings.TrimSpace(strs[1]), nil
	}

	key := strings.TrimSpace(strs[0])
	value := strings.Join(strs[1:], confDelimiter)

	return key, strings.TrimSpace(value), nil
}
