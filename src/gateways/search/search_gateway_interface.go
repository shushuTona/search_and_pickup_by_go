package search

import "encoding/csv"

type ISearchGateway interface {
	GetTargetFilePath() ([]string, error)
	GetSearchWordList(search_file_path string) *csv.Reader
}
