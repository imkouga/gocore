package filex

import (
	"io/ioutil"
	"os"
	"path"
)

// 写data数据到完整的文件, 如果文件不存在则创建，如果存在则覆写
func WriteFullFile(fileName string, data []byte) error {

	if IsFileExist(fileName) {
		fd, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}

		defer fd.Close()
		if err := fd.Truncate(0); err != nil {
			return err
		}
	}

	fileDir := path.Dir(fileName)
	if false == IsPathExist(fileDir) {
		if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(fileName, data, os.ModePerm)
}

func IsPathExist(path string) bool {
	return isPathExist(path)
}

func IsFileExist(file string) bool {
	return isPathExist(file)
}

func isPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
