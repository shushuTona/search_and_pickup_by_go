package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	cli "github.com/urfave/cli/v2"
)

const (
	SEARCH_FILE = "./search.csv"
	RESULT_FILE = "./result.csv"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dir",
				Value: "",
				Usage: "target dir name for search",
			},
			&cli.StringFlag{
				Name:  "extension",
				Value: "",
				Usage: "target extension for search",
			},
			&cli.StringFlag{
				Name:  "ignore",
				Value: "",
				Usage: "ignore dir name for search",
			},
            &cli.BoolFlag{
                Name:  "exists",
                Value: false,
                Usage: "check without target file path",
            },
		},
		Action: func(cCtx *cli.Context) error {
			// dir arg
			if len(cCtx.String("dir")) == 0 {
				return errors.New("error : require dir arg")
			}
			tagetRootDir := cCtx.String("dir")

			// extension arg
			if len(cCtx.String("extension")) == 0 {
				return errors.New("error : require extension arg")
			}
			var fileExtension string
			if cCtx.String("extension")[0:1] == "." {
				fileExtension = cCtx.String("extension")
			} else {
				fileExtension = "." + cCtx.String("extension")
			}

			// ignore arg
			ignoreDirNameList := []string{}
			if len(cCtx.String("ignore")) > 0 {
				ignoreDirNameList = strings.Split(cCtx.String("ignore"), ":")
			}

			// existsFlag
			existsFlag := cCtx.Bool("exists")

			err := search(tagetRootDir, fileExtension, ignoreDirNameList, existsFlag)

			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func search(tagetRootDir string, fileExtension string, ignoreDirNameList []string, existsFlag bool) error {
	fmt.Println("ignoreDirNameList")
	fmt.Println(ignoreDirNameList)

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
		return err
	}

	fmt.Println(filePathList)

	// 検索ワード一覧を取得
	searchFile, err := os.Open(SEARCH_FILE)
	if err != nil {
		return err
	}
	defer searchFile.Close()

	// 検索結果を書き込むCSVファイルを生成 or 読み込み
	resultFile, err := os.Create(RESULT_FILE)
	if err != nil {
		return err
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
				// existsFlagが有効の場合、1つでも対象ファイルが存在する場合繰り返し処理を終了する。
				if existsFlag && len(targetList) > 0 {
					break
				}

				if strings.Contains(scanner.Text(), searchWord) {
					targetList = targetList + f.Name() + ":" + strconv.Itoa(line) + "\n"
				}

				line++
			}

			// if err := scanner.Err(); err != nil {
			// 	// Handle the error
			// }
		}

		// existsFlagが有効の場合、ファイルパスの代わりに真偽値を追加する。
		if existsFlag {
			flag := "false"
			if len(targetList) > 0 {
				flag = "true"
			}
			resultRecord = append(resultRecord, flag)
		} else {
			resultRecord = append(resultRecord, targetList)
		}

		resultFileWriter.Write(resultRecord)
	}

	return nil
}
