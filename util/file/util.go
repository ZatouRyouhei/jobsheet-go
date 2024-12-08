package file

import (
	"os"
)

// ディレクトリの存在確認
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
