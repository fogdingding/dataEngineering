package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"math/rand"
	"strconv"
	"time"
)

const (
	_inputPath = "dataset/ettoday.rec"
	_outputPath = "dataset/result.rec"
)

type Text struct {
	idx string // Test rsort -n
	url string
	title string
	content string
	published string
	author string
	favoriteCount string
	viewCount string
	res string
	duration string
	category string
}

type RsortArgs struct {
	dIdx int 
	kIdx int
	iIdx int
	rIdx int
	sIdx int
	nIdx int
	parallelIdx int
	externalIdx int
	chunkIdx int
}

func main() {
	userArgs := getArgsIdx()
	if userArgs.externalIdx != -1 {
		chunk := 1
		if userArgs.chunkIdx != -1 {
			chunk, _ = strconv.Atoi(os.Args[userArgs.chunkIdx+1])
		}
		start := time.Now()
		externalMergeSort(os.Args[userArgs.externalIdx+1], userArgs, chunk)
		duration := time.Since(start).Seconds()
		fmt.Printf("External rsort spends %.3f sec\n", duration)	
	} else {
		if userArgs.dIdx == -1 {
			text, err := readNewsData(_inputPath)
			if err == nil {
				fmt.Printf("Total news data = %d\n", len(text))
				start := time.Now()
				text = mergeSortText(text, userArgs)
				duration := time.Since(start).Seconds()
				fmt.Printf("Rsort spends %.3f sec\n", duration)
				err = writeNewsInFile(_outputPath, text, userArgs)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		} else {
			content, err := readNewsBySplitN(_inputPath, os.Args[userArgs.dIdx+1])
			if err == nil {
				fmt.Printf("Total split string = %d\n", len(content))		
				start := time.Now()
				content = mergeSortString(content, userArgs)
				duration := time.Since(start).Seconds()
				fmt.Printf("Rsort spends %.3f sec\n", duration)
				err = writeStrInFile(_outputPath, content, userArgs)
				if err != nil {
					fmt.Println(err)
				}	
			} else {
				fmt.Println(err)
			}
		}
	}
}

//***************************************************************************
// read / write functions 
//***************************************************************************
func readNewsData(path string) ([]Text, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return []Text{}, err
	}
	defer inputFile.Close()

	var text []Text
	var tmp Text
	contentFlag := 0
	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "@U:") {
			tmp.url = strings.TrimLeft(scanner.Text(), "@U:")
		} else if strings.HasPrefix(scanner.Text(), "@T:") {
			tmp.title = strings.TrimLeft(scanner.Text(), "@T:")
		} else if strings.HasPrefix(scanner.Text(), "@B:") {
			contentFlag = 1
		} else if contentFlag == 1 {
			tmp.idx = strconv.Itoa(rand.Intn(1000000)) // Test rsort -n
			tmp.content = strings.TrimSpace(scanner.Text())
			text = append(text, tmp)
			contentFlag = 0
		}
	}

	return text, scanner.Err()
}

func writeNewsInFile(path string, text []Text, userArgs RsortArgs) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)
	for idx, _ := range text {
		w.WriteString("@GAISRec:\n")
		if userArgs.rIdx == -1 {
			w.WriteString("@idx:" + text[idx].idx + "\n") // Test rsort -n
			w.WriteString("@U:" + text[idx].url + "\n")
			w.WriteString("@T:" + text[idx].title + "\n")
			w.WriteString("@B:\n" + text[idx].content + "\n")
		} else {
			w.WriteString("@idx:" + text[len(text) - idx - 1].idx + "\n") // Test rsort -n
			w.WriteString("@U:" + text[len(text) - idx - 1].url + "\n")
			w.WriteString("@T:" + text[len(text) - idx - 1].title + "\n")
			w.WriteString("@B:\n" + text[len(text) - idx - 1].content + "\n")
		}
		w.WriteString("\n")
	}
	return w.Flush()
}

//***************************************************************************
// @Purpose : rsort -d [delimiter] ( only for news content )
// @Description : split contents by delimiters
//***************************************************************************
func readNewsBySplitN(path string, delimiter string) ([]string, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	var content []string
	contentFlag := 0
	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "@B:") {
			contentFlag = 1
		} else if contentFlag == 1 {
			size := len(scanner.Text())
			var splitStr []string

			if delimiter != "--all" {
				splitStr = strings.SplitN(scanner.Text(), delimiter, size)
			} else {
				splitStr = strings.FieldsFunc(scanner.Text(), splitByDefault)
			}

			for strIdx, _ := range splitStr {
				substr := strings.TrimSpace(splitStr[strIdx])
				if substr != "" {
					content = append(content, substr)
				}
			}

			contentFlag = 0
		}
	}

	return content, scanner.Err()
}

//***************************************************************************
// @Purpose : rsort -d --all
// @Description : set default delimiters
//***************************************************************************
func splitByDefault(word rune) bool {
	return word == '，' || word == '；' || word == '。' || word == '？' || word == '！'
}

func writeStrInFile(path string, content []string, userArgs RsortArgs) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)
	size := len(content) 
	for idx, _ := range content {
		if userArgs.rIdx == -1 {
			w.WriteString(content[idx] + "\n")
		} else {
			w.WriteString(content[size - idx - 1] + "\n")
		}
	}
	return w.Flush()
}

//***************************************************************************
// Rsort argument functions
// @Purpose : get rsort argument
// @Description : 
// if argument exist return there index, otherwise return -1
//***************************************************************************

func newArgs() RsortArgs {
	return RsortArgs{-1, -1, -1, -1, -1, -1, -1, -1, -1}
}

func getArgsIdx() (RsortArgs) {
	userArgs := newArgs()
	for idx, args := range os.Args {
		if args == "-d" {
			userArgs.dIdx = idx
		} else if args == "-k" {
			userArgs.kIdx = idx
		} else if args == "-i" {
			userArgs.iIdx = idx
		} else if args == "-r" {
			userArgs.rIdx = idx
		} else if args == "-s" {
			userArgs.sIdx = idx
		} else if args == "-n" {
			userArgs.nIdx = idx
		} else if args == "--parallel" {
			userArgs.parallelIdx = idx
		} else if args == "--external" {
			userArgs.externalIdx = idx
		} else if args == "--chunk" && userArgs.externalIdx != -1 {
			userArgs.chunkIdx = idx
		}
	}
	return userArgs
}