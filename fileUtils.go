package tLog

import (
	"io/fs"
	"os"
)

// 创建目录
func createDir(dir string, perm fs.FileMode) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, perm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
