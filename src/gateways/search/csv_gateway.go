package search

import (
	"encoding/csv"
	"os"
)

type CSVSearchGateway struct {
	tagetRootDir      string
	fileExtension     string
	ignoreDirNameList []string
	searchFile        *os.File
}

func NewSearchGateway(tagetRootDir string, fileExtension string, ignoreDirNameList []string) ISearchGateway {
	return &CSVSearchGateway{
		tagetRootDir:      tagetRootDir,
		fileExtension:     fileExtension,
		ignoreDirNameList: ignoreDirNameList,
	}
}

func (csg *CSVSearchGateway) GetTargetFilePath() ([]string, error) {
	filePathList, err := getTargetFilePath(csg.tagetRootDir, csg.fileExtension, csg.ignoreDirNameList)
	if err != nil {
		return nil, err
	}

	return filePathList, nil
}

func (csg *CSVSearchGateway) GetSearchWordList(search_file_path string) *csv.Reader {
	searchFile, err := os.Open(search_file_path)
	if err != nil {
		panic(err)
	}
	csg.searchFile = searchFile

	csvReader := csv.NewReader(searchFile)

	return csvReader
}
