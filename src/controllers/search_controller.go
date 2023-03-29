package controllers

import "fmt"

type SearchController struct {
	tagetRootDir      string
	fileExtension     string
	ignoreDirNameList []string
}

func (sc *SearchController) Search() error {
	fmt.Println(sc.tagetRootDir)
	fmt.Println(sc.fileExtension)
	fmt.Println(sc.ignoreDirNameList)

	return nil
}
