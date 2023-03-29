package search

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getTargetFilePath(tagetRootDir string, fileExtension string, ignoreDirNameList []string) ([]string, error) {
	// 対象ディレクトリ配下に存在する対象拡張子のファイルを再帰的に確認して対象ファイルパス一覧を生成
	var filePathList []string
	err := filepath.Walk(tagetRootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		ignoreFlag := false
		for _, ignore := range ignoreDirNameList {
			if strings.Contains(path, ignore) {
				ignoreFlag = true
				break
			}
		}
		if ignoreFlag {
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == fileExtension {
			filePathList = append(filePathList, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filePathList, nil
}
