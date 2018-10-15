package main

import(
	"fmt"
	"syscall"
	"os"
	"strings"
	"bufio"
	"time"
	"strconv"
	"container/heap"
)

const(
	_splitFiilePath = "dataset/tmpFile/split_text_"
	_externalOutputPath = "dataset/external_result.rec"
)

type Node struct {
	fileIdx int
	chunkIdx int
	value string
}

func externalMergeSort(inputFile string, userArgs RsortArgs, chunk int) {
	fileNum, err := generateKFile(inputFile, userArgs, chunk)
	fmt.Printf("File N = %d and user sets chunk = %d\n", fileNum, chunk)
	if err != nil {
		fmt.Println(err)
	} else {
		err = generateRun(fileNum, chunk, userArgs)
		if err != nil {
			fmt.Println(err)
		} 
	}
}

// Maximum memory available in my computer is 5916029232
func generateKFile(path string, userArgs RsortArgs, chunk int) (int, error) {
	var text []Text; var tmpText Text
	memUsed, fileIdx := 0, 1

	memLimit, err := getMemLimit()
	if err != nil {
		return -1, err
	}

	inputFile, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "@url") {
			tmpText.url = strings.TrimLeft(scanner.Text(), "@url:")
		} else if strings.HasPrefix(scanner.Text(), "@published:") {
			tmpText.published = strings.TrimLeft(scanner.Text(), "@published:")
		} else if strings.HasPrefix(scanner.Text(), "@title:") {
			tmpText.title = strings.TrimLeft(scanner.Text(), "@title:")
		} else if strings.HasPrefix(scanner.Text(), "@content:") {
			tmpText.content = strings.TrimLeft(scanner.Text(), "@content:")
		} else if strings.HasPrefix(scanner.Text(), "@author:") {
			tmpText.author = strings.TrimLeft(scanner.Text(), "@author:")
		} else if strings.HasPrefix(scanner.Text(), "@favoriteCount:") {
			tmpText.favoriteCount = strings.TrimLeft(scanner.Text(), "@favoriteCount:")
		} else if strings.HasPrefix(scanner.Text(), "@viewCount:") {
			tmpText.viewCount = strings.TrimLeft(scanner.Text(), "@viewCount:")
		} else if strings.HasPrefix(scanner.Text(), "@res:") {
			tmpText.res = strings.TrimLeft(scanner.Text(), "@res:")
		} else if strings.HasPrefix(scanner.Text(), "@duration:") {
			tmpText.duration = strings.TrimLeft(scanner.Text(), "@duration:")
		} else if strings.HasPrefix(scanner.Text(), "@category:") {
			tmpText.category = strings.TrimLeft(scanner.Text(), "@category:")
			text = append(text, tmpText)
			memUsed += len(tmpText.url + tmpText.title + tmpText.content + tmpText.published + tmpText.author + tmpText.favoriteCount + tmpText.viewCount + tmpText.res + tmpText.duration + tmpText.category)
			if memUsed >= memLimit {
				fileName := _splitFiilePath + strconv.Itoa(fileIdx)

				err = writeKFile(text, fileName, chunk, userArgs)
				if err != nil {
					return fileIdx, err
				}

				fileIdx++; text = []Text{}; memUsed = 0
			}
		}
	}

	if len(text) != 0 {
		fileName := _splitFiilePath + strconv.Itoa(fileIdx)
		err = writeKFile(text, fileName, chunk, userArgs)
		if err != nil {
			return fileIdx, err
		}
	} else {
		fileIdx = fileIdx - 1
	}

	fmt.Printf("Completely generate %d files\n", fileIdx * chunk)
	return fileIdx, scanner.Err()
}

func getMemLimit() (int, error) {
	var stat syscall.Statfs_t

	wd, err := os.Getwd()
	if err != nil {
		return -1, err
	}

	syscall.Statfs(wd, &stat)
	memLimit := int(stat.Bavail * uint64(stat.Bsize)) / 200
	fmt.Printf("Available memory = %d\n", memLimit)
	return memLimit, nil
}

func writeKFile(text []Text, fileName string, chunk int, userArgs RsortArgs) error {
	start := time.Now()
	text = mergeSortText(text, userArgs)
	duration := time.Since(start).Seconds()
	fmt.Printf("Rsort spends %.3f sec\n", duration)

	size := len(text)
	for chunkIdx := 1; chunkIdx <= chunk; chunkIdx++ {
		lower := ((chunkIdx-1) * size) / chunk
		upper := (chunkIdx * size) / chunk
		writeText := text[lower:upper]

		outputFile, err := os.Create(fileName + "_" + strconv.Itoa(chunkIdx) + ".rec")
		if err != nil {
			return err
		}
		defer outputFile.Close()
	
		w := bufio.NewWriter(outputFile)
		for idx, _ := range writeText {
			w.WriteString("@url:" + writeText[idx].url)
			w.WriteString("@published:" + writeText[idx].published)
			w.WriteString("@title:" + writeText[idx].title)
			w.WriteString("@content:" + writeText[idx].content)
			w.WriteString("@author:" + writeText[idx].author)
			w.WriteString("@favoriteCount:" + writeText[idx].favoriteCount)
			w.WriteString("@viewCount:" + writeText[idx].viewCount)
			w.WriteString("@res:" + writeText[idx].res)
			w.WriteString("@duration:" + writeText[idx].duration)
			w.WriteString("@category:" + writeText[idx].category + "\n")
		}
		
		if w.Flush() != nil {
			return w.Flush()
		}
	} 
	return nil
}

func generateRun(fileNum int, chunk int, userArgs RsortArgs) error {
	files := make(map[int]*bufio.Reader, fileNum)
	var pq Heap
	heap.Init(&pq)
	for i := 1; i <= fileNum; i++ {
		fileName := _splitFiilePath + strconv.Itoa(i) + "_1.rec"

		file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		defer file.Close()

		files[i-1] = bufio.NewReaderSize(file, 66000)
		text, _ := files[i-1].ReadString('\n')
		heap.Push(&pq, &Node{i, 1, text})
	}

	outputFile, err := os.Create(_externalOutputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	w := bufio.NewWriter(outputFile)
	for true {
		if pq.Len() == 0 {
			break
		}

		minNode := heap.Pop(&pq).(*Node)
		w.WriteString(minNode.value)

		text, _ := files[minNode.fileIdx-1].ReadString('\n')
		if len(text) != 0 {
			heap.Push(&pq, &Node{minNode.fileIdx, minNode.chunkIdx, text})
		} else {
			if minNode.chunkIdx != chunk {
				newChunk := _splitFiilePath + strconv.Itoa(minNode.fileIdx) + "_" + strconv.Itoa(minNode.chunkIdx+1) + ".rec"

				file, err := os.OpenFile(newChunk, os.O_RDONLY, 0)
				if err != nil {
					return err
				}
				defer file.Close()
		
				files[minNode.fileIdx-1] = bufio.NewReaderSize(file, 66000)
				newText, _ := files[minNode.fileIdx-1].ReadString('\n')
				heap.Push(&pq, &Node{minNode.fileIdx, minNode.chunkIdx+1, newText})
			}
			fileName := _splitFiilePath + strconv.Itoa(minNode.fileIdx) + "_" + strconv.Itoa(minNode.chunkIdx) + ".rec"
			os.Remove(fileName)
			fmt.Printf("Sucessfully merge %s\n", fileName)		
		}
	}
	return w.Flush()
}

//***************************************************************************
// Heap functions
//***************************************************************************

type Heap []*Node

func (pq Heap) Len() int {
	return len(pq)
}

func (pq Heap) Less(i, j int) bool {
	return pq[i].value < pq[j].value
}

func (pq Heap) Swap(i,j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *Heap) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}

func (pq *Heap) Pop() interface{} {
	oldNodes := *pq
    size := len(oldNodes)
    node := oldNodes[size-1]
    *pq = oldNodes[0 : size-1]
    return node
}