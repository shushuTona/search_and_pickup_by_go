package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// 対象ディレクトリ
	tagetRootDir := "/go/src/testdata"

	// 対象ファイル拡張子
	fileExtension := ".js"

	// 対象ディレクトリ配下に存在する対象拡張子のファイルを再帰的に確認して対象ファイルパス一覧を生成
	var filePathList []string
	err := filepath.Walk(tagetRootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == fileExtension {
			filePathList = append(filePathList, path)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// 検索ワード一覧を取得
	searchFile, err := os.Open("./search.csv")
	if err != nil {
		panic(err)
	}
	defer searchFile.Close()

	// 検索結果を書き込むCSVファイルを生成 or 読み込み
	resultFile, err := os.Create("./result.csv")
	if err != nil {
		panic(err)
	}
	resultFileWriter := csv.NewWriter(resultFile)
	defer resultFileWriter.Flush()

	csvReader := csv.NewReader(searchFile)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			// 最後まで読み出した場合break
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		searchWord := record[0]
		fmt.Printf("%#v\n", searchWord)

		resultRecord := []string{}
		// 1カラム目に検索単語を追加する
		resultRecord = append(resultRecord, searchWord)

		var targetList string
		// 対象ファイル内に存在検索ワードを確認
		for _, filePath := range filePathList {
			f, err := os.Open(filePath)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			// Splits on newlines by default.
			scanner := bufio.NewScanner(f)

			line := 1
			// https://golang.org/pkg/bufio/#Scanner.Scan
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), searchWord) {
					targetList = targetList + f.Name() + ":" + strconv.Itoa(line) + "\n"
				}

				line++
			}

			// if err := scanner.Err(); err != nil {
			// 	// Handle the error
			// }
		}
		resultRecord = append(resultRecord, targetList)

		resultFileWriter.Write(resultRecord)
	}
}
