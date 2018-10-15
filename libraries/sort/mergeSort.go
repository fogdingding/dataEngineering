package main

import(
	"os"
	"strings"
	"strconv"
	"sync"
	"runtime"
)

func mergeSortText(inputArr []Text, userArgs RsortArgs) []Text {
	size := len(inputArr)
	
	if size <= 1 {
		return inputArr
	}
	
	if runtime.NumGoroutine() >= runtime.NumCPU() || userArgs.parallelIdx == -1 {
		if size < 20 {
			return insertionSortText(inputArr, size, userArgs)
		} else {
			return mergeText(mergeSortText(inputArr[:size/2], userArgs), mergeSortText(inputArr[size/2:], userArgs), userArgs)
		}
	} else {
		var left, right []Text
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {left = mergeSortText(inputArr[:size/2], userArgs); wg.Done()}()
		go func() {right = mergeSortText(inputArr[size/2:], userArgs); wg.Done()}()
		wg.Wait()		
		return mergeText(left, right, userArgs)
	}
}

func mergeText(left []Text, right []Text, userArgs RsortArgs) []Text {
	result := make([]Text, len(left) + len(right))
	leftIdx, rightIdx := 0, 0

	for idx := 0; idx < len(result); idx++ {
		if leftIdx < len(left) && rightIdx < len(right) {
			leftCmpStr, rightCmpStr := setCompareString(left[leftIdx], right[rightIdx], userArgs)
			if setCompareFunc(leftCmpStr, rightCmpStr, userArgs) {
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

func insertionSortText(inputArr []Text, size int, userArgs RsortArgs) []Text {
	var i, j int
	for i = 1; i < size; i++ {
		tmp := inputArr[i]
		for j = i-1; j >= 0; j-- {
			left, right := setCompareString(inputArr[j], tmp, userArgs)
			if !setCompareFunc(left, right, userArgs) {
				inputArr[j+1] = inputArr[j]
			} else {
				break
			}
		}
		inputArr[j+1] = tmp
	}
	return inputArr
}

//***************************************************************************
// @Purpose : rsort -k [ specific content ] and rsort -i 
// @Description : return two strings using to compare 
// case -k : specific content ex : @U, @T and @B etc...
// case -i : all strings in lower case ( only for english )  
// default : @U + @T + @B strings
//***************************************************************************
func setCompareString(left, right Text, userArgs RsortArgs) (string, string) {
	leftCmpStr, rightCmpStr := "", ""

	if userArgs.kIdx == -1 {
			leftCmpStr, rightCmpStr = left.url + left.title + left.content + left.published  + left.author + left.favoriteCount + left.viewCount + left.res + left.duration + left.category, right.url + right.title + right.content + right.published + right.author + right.favoriteCount + right.viewCount + right.res + right.duration + right.category
	} else {
		if os.Args[userArgs.kIdx + 1] == "@U" {
			leftCmpStr, rightCmpStr = left.url, right.url
		} else if os.Args[userArgs.kIdx + 1] == "@T" {
			leftCmpStr, rightCmpStr = left.title, right.title
		} else if os.Args[userArgs.kIdx + 1] == "@B" {
			leftCmpStr, rightCmpStr = left.content, right.content
		} else if os.Args[userArgs.kIdx + 1] == "@idx" {
			leftCmpStr, rightCmpStr = left.idx, right.idx // Test rsort -n
		}
	}

	if userArgs.iIdx != -1 {
		leftCmpStr, rightCmpStr = strings.ToLower(leftCmpStr), strings.ToLower(rightCmpStr)
	}

	return leftCmpStr, rightCmpStr
}

//***************************************************************************
// @Purpose : rsort -n and rsort -s
// @Description : set compare function in merge
// case -n : convert strings to integer and compare them (numerical comparison)
// case -s : compare news by there sizes (size oreder comparison)
// default : compare news by there ascii 
//***************************************************************************
func setCompareFunc(leftCmpStr, rightCmpStr string, userArgs RsortArgs) bool {
	if userArgs.nIdx != -1 {
		leftCmpInt, _ := strconv.Atoi(leftCmpStr)
		rightCmpInt, _ := strconv.Atoi(rightCmpStr)
		return leftCmpInt <= rightCmpInt
	} else if userArgs.sIdx != -1 {
		return len(leftCmpStr) <= len(rightCmpStr)
	} else {
		return leftCmpStr <= rightCmpStr
	}
}

func mergeSortString(inputArr []string, userArgs RsortArgs) []string {
	size := len(inputArr)
	
	if size <= 1 {
		return inputArr
	}
	
	if runtime.NumGoroutine() >= runtime.NumCPU() || userArgs.parallelIdx == -1 {
		if size < 20 {
			return insertionSortString(inputArr, size)
		} else {
			return mergeString(mergeSortString(inputArr[:size/2], userArgs), mergeSortString(inputArr[size/2:], userArgs))
		}
	} else {
		var left, right []string
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {left = mergeSortString(inputArr[:size/2], userArgs); wg.Done()}()
		go func() {right = mergeSortString(inputArr[size/2:], userArgs); wg.Done()}()
		wg.Wait()		
		return mergeString(left, right)
	}
}

func mergeString(left, right []string) []string {
	result := make([]string, len(left) + len(right))
	leftIdx, rightIdx := 0, 0

	for idx := 0; idx < len(result); idx++ {
		if leftIdx < len(left) && rightIdx < len(right) {
			if left[leftIdx] <= right[rightIdx] {
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

func insertionSortString(inputArr []string, size int) []string {
	var i, j int
	for i = 1; i < size; i++ {
		tmp := inputArr[i]
		for j = i-1; j >= 0 && inputArr[j] > tmp; j-- {
			inputArr[j+1] = inputArr[j]
		}
		inputArr[j+1] = tmp
	}
	return inputArr
}