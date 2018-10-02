package main

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	_inputPath = "ettoday.rec"
	_outputPath = "result.rec"
)

type Text struct {
	idx string // Test rsort -n
	url string
	title string
	content string
}

// Set buffered channel to control the number of concurrent units
var sem = make(chan struct{}, getSyncUnits())

func main() {
	text, _ := readFile(_inputPath)
	fmt.Printf("Total news data = %d\n", len(text))

	start := time.Now()
	text = mergeSortSync(text)
	duration := time.Since(start).Seconds()
	fmt.Printf("Rsort spends %.3f sec\n", duration)

	_ = writeFile(_outputPath, text)
}

func mergeSort(inputArr []Text) []Text {
	if len(inputArr) <= 1 {
		return inputArr
	}
	left, right, _ := split(inputArr)
	left = mergeSort(left)
	right = mergeSort(right)
	return merge(left, right)
}

func mergeSortSync(inputArr []Text) []Text {
	size := len(inputArr)

	if size <= 1 {
		return inputArr
	}
	
	left, right, _ := split(inputArr)
	var wg sync.WaitGroup
	wg.Add(2)

	select {
		case sem <- struct{}{}:
			go func() {
				left = mergeSortSync(left);
				<-sem
				wg.Done()
			}()
		default:
			left = mergeSort(left)
			wg.Done()
	}

	select {
		case sem <- struct{}{}:
			go func() {
				right = mergeSortSync(right)
				<-sem
				wg.Done()
			}()
		default:
			right = mergeSort(right)
			wg.Done()
    }

	wg.Wait()
	return merge(left, right)
}

func split(inputArr []Text) ([]Text, []Text, int) {
	mid := len(inputArr) / 2
	return inputArr[0:mid], inputArr[mid:], mid
}

func merge(left []Text, right []Text) []Text {
	result := make([]Text, len(left) + len(right))
	leftIdx, rightIdx := 0, 0

	for idx := 0; idx < len(result); idx++ {
		if leftIdx < len(left) && rightIdx < len(right) {
			leftCmpStr, rightCmpStr := setCompareString(left[leftIdx], right[rightIdx])
			if setCompareFunc(leftCmpStr, rightCmpStr) {
				result[idx] = left[leftIdx]
				leftIdx++
			} else {
				result[idx] = right[rightIdx]
				rightIdx++
			}
		} else {
			if leftIdx < len(left) {
				result[idx] = left[leftIdx]
				leftIdx++
			} else if rightIdx < len(right) {
				result[idx] = right[rightIdx]
				rightIdx++
			}
		}
	}

	return result
}

func readFile(path string) ([]Text, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return []Text{}, err
	}
	defer inputFile.Close()

	var text []Text
	var tmp Text
	idx, contentFlag := 0, 0
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
			tmp.content = scanner.Text()
			text = append(text, tmp)
			contentFlag = 0
			idx++
		}
	}

	return text, scanner.Err()
}

// Handling rsort argument -r (reverse order)
func writeFile(path string, text []Text) error {
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)
	for idx, _ := range text {
		argsIdx := checkArgsExist("-r")
		w.WriteString("@GAISRec:\n")
		if argsIdx == -1 {
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

func checkArgsExist(targetArgs string) int {
	for idx, args := range os.Args {
		if args == targetArgs {
			return idx
		}
	}
	return -1
}

// Handling rsort argument -k (key pat) & -i (case insensitive)
func setCompareString(left, right Text) (string, string) {
	leftCmpStr, rightCmpStr := "", ""
	kIdx, iIdx := checkArgsExist("-k"), checkArgsExist("-i")

	if kIdx == -1 {
			leftCmpStr, rightCmpStr = left.url + left.title + left.content, right.url + right.title + right.content
	} else {
		if os.Args[kIdx + 1] == "@U" {
			leftCmpStr, rightCmpStr = left.url, right.url
		} else if os.Args[kIdx + 1] == "@T" {
			leftCmpStr, rightCmpStr = left.title, right.title
		} else if os.Args[kIdx + 1] == "@B" {
			leftCmpStr, rightCmpStr = left.content, right.content
		} else if os.Args[kIdx + 1] == "@idx" {
			leftCmpStr, rightCmpStr = left.idx, right.idx // Test rsort -n
		}
	}

	if iIdx != -1 {
		leftCmpStr, rightCmpStr = strings.ToLower(leftCmpStr), strings.ToLower(rightCmpStr)
	}

	return leftCmpStr, rightCmpStr
}

// Handling rsort argument -s (size order) & -n (numerical comparison)
func setCompareFunc(leftCmpStr string, rightCmpStr string) bool {
	if checkArgsExist("-n") != -1 {
		leftCmpInt, _ := strconv.Atoi(leftCmpStr)
		rightCmpInt, _ := strconv.Atoi(rightCmpStr)
		return leftCmpInt <= rightCmpInt
	} else if checkArgsExist("-s") != -1 {
		return len(leftCmpStr) <= len(rightCmpStr)
	} else {
		return leftCmpStr <= rightCmpStr
	}
}