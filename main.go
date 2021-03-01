package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Complete the balancedForest function below.
func balancedForest(c []int32, edges [][]int32) int64 {
	neighbors := make([][]int32, len(c))
	for _, e := range edges {
		neighbors[e[0]-1] = append(neighbors[e[0]-1], e[1]-1)
		neighbors[e[1]-1] = append(neighbors[e[1]-1], e[0]-1)
	}
	markers := make([]bool, len(c))
	parents := make([]int32, len(c))
	sizeMap := make(map[int64][]int32)
	sizes := make([]int64, len(c))

	var size int64
	for _, s := range c {
		size += int64(s)
	}

	var buildTree func(int32) int64
	buildTree = func(nodeIndex int32) int64 {
		markers[nodeIndex] = true
		nodeSize := int64(c[nodeIndex])
		for _, n := range neighbors[nodeIndex] {
			if markers[n] {
				continue
			}
			parents[n] = nodeIndex
			nodeSize += buildTree(n)
		}

		sizes[nodeIndex] = nodeSize
		if isValidSubTree(size, nodeSize) {
			sizeMap[nodeSize] = append(sizeMap[nodeSize], nodeIndex)
		}
		return nodeSize
	}

	buildTree(0)

	min := int64(-1)

	getValidSpread := func(nodeSize int64, nodes []int32) int64 {
		if isLargerSubTree(size, nodeSize) {
			upper := nodeSize
			lower := size - nodeSize*2
			spread := upper - lower
			if lower == 0 {
				return spread
			}
			if len(nodes) >= 2 {
				return spread
			}
			for _, n := range nodes {
				if hasParentSize(n, parents, sizes, upper+lower) || hasParentSize(n, parents, sizes, upper+upper) {
					return spread
				}
				if lower == upper {
					break
				}
				for _, l := range sizeMap[lower] {
					if !hasParentNode(l, parents, n) {
						return spread
					}
				}
			}
		}
		if isSmallerSubTree(size, nodeSize) {
			lower := nodeSize
			upper := (size - nodeSize) / 2
			spread := upper - lower
			for _, n := range sizeMap[nodeSize] {
				if hasParentSize(n, parents, sizes, upper+lower) {
					return spread
				}
			}
		}
		return -1
	}

	for nodeSize, nodes := range sizeMap {
		spread := getValidSpread(nodeSize, nodes)
		if min == -1 || (spread != -1 && spread < min) {
			min = spread
		}
	}

	return min
}

func hasParentNode(node int32, parents []int32, expectedParent int32) bool {
	parent := parents[node]
	for parent != 0 {
		if parent == expectedParent {
			return true
		}
		parent = parents[parent]
	}
	return false
}

func hasParentSize(node int32, parents []int32, sizes []int64, expectedSize int64) bool {
	parent := parents[node]
	for parent != 0 {
		if sizes[parent] == expectedSize {
			return true
		}
		parent = parents[parent]
	}
	return false
}

func isLargerSubTree(size, nodeSize int64) bool {
	return nodeSize*2 <= size && size-nodeSize*2 <= nodeSize
}

func isSmallerSubTree(size, nodeSize int64) bool {
	return nodeSize%2 == size%2 && (size-nodeSize)/2 >= nodeSize
}

func isValidSubTree(size, nodeSize int64) bool {
	return isLargerSubTree(size, nodeSize) || isSmallerSubTree(size, nodeSize)
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	stdout := os.Stdout

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	qTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	q := int32(qTemp)

	for qItr := 0; qItr < int(q); qItr++ {
		nTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
		checkError(err)
		n := int32(nTemp)

		cTemp := strings.Split(readLine(reader), " ")

		var c []int32

		for i := 0; i < int(n); i++ {
			cItemTemp, err := strconv.ParseInt(cTemp[i], 10, 64)
			checkError(err)
			cItem := int32(cItemTemp)
			c = append(c, cItem)
		}

		var edges [][]int32
		for i := 0; i < int(n)-1; i++ {
			edgesRowTemp := strings.Split(readLine(reader), " ")

			var edgesRow []int32
			for _, edgesRowItem := range edgesRowTemp {
				edgesItemTemp, err := strconv.ParseInt(edgesRowItem, 10, 64)
				checkError(err)
				edgesItem := int32(edgesItemTemp)
				edgesRow = append(edgesRow, edgesItem)
			}

			if len(edgesRow) != 2 {
				panic("Bad input")
			}

			edges = append(edges, edgesRow)
		}

		result := balancedForest(c, edges)

		fmt.Fprintf(writer, "%d\n", result)
	}

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
