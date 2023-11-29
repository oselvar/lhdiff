package lhdiff

import (
	"fmt"
	"strconv"

	"github.com/mongodb-forks/go-difflib/difflib"
	"github.com/sourcegraph/go-diff/diff"
)

const ContextSimilarityFactor = 0.4
const ContentSimilarityFactor = 0.6
const SimilarityThreshold = 0.45

/**
 * Returns a list of mappings between the lines of the left and right file.
 * Each mapping is a pair of line numbers, where -1 indicates that the line is not present in the file.
 * The mappings are sorted by the line number of the left file.
 * If includeIdenticalLines is true, then lines that are identical in both files are included in the mappings.
 * Otherwise, only lines that are not identical are included.
 * The contextSize parameter determines how many lines of context are used to determine the similarity of lines.
 * The context lines are not included in the mappings.
 * The context lines are lines that are not blank and do not consist of only curly braces or parenthesis.
 * The context lines are used to determine the similarity of lines.
 * The similarity of lines is determined by a combination of the normalized Levenshtein distance of the content of the lines and the cosine similarity of the context of the lines.
 * The similarity of lines is only considered if it is above a certain threshold.
 * The mappings are determined by first finding the unchanged lines using the difflib library.
 * Then, for each line in the right file, the most similar line in the left file is found.
 * The most similar line is the line with the highest combined similarity.
 */
func Lhdiff(left string, right string, contextSize int, includeIdenticalLines bool) ([][]uint32, error) {
	leftLines := ConvertToLinesWithoutNewLine(left)
	rightLines := ConvertToLinesWithoutNewLine(right)

	diffScript, err := difflib.GetUnifiedDiffString(difflib.LineDiffParams{
		A:        leftLines,
		B:        rightLines,
		FromFile: "left",
		ToFile:   "right",
		Context:  3,
	})
	//fmt.Println(diffScript)
	if err != nil {
		return nil, err
	}

	mappedByRightLineNumber := make(map[int]bool)
	linePairByLeftLineNumber := make(map[int]LinePair, 0)

	if diffScript == "" {
		// The files are identical
		for leftLineNumber := range leftLines {
			lineInfo := MakeLineInfo(leftLineNumber, leftLines, 4)
			linePairByLeftLineNumber[leftLineNumber] = LinePair{
				left:  lineInfo,
				right: lineInfo,
			}
			mappedByRightLineNumber[leftLineNumber] = true
		}
	} else {
		fileDiff, err := diff.ParseFileDiff([]byte(diffScript))
		if err != nil {
			return nil, err
		}

		Algo1(fileDiff, leftLines, rightLines, contextSize, &mappedByRightLineNumber, &linePairByLeftLineNumber)
	}
	return makeLineMappings(linePairByLeftLineNumber, mappedByRightLineNumber, leftLines, rightLines, includeIdenticalLines), nil
}

func makeLineMappings(linePairs map[int]LinePair, mappedByRightLineNumber map[int]bool, leftLines []string, rightLines []string, includeIdenticalLines bool) [][]uint32 {
	leftLineCount := len(leftLines)
	rightLineNumbers := make([]int, 0)
	for rightLineNumber := range rightLines {
		_, mapped := mappedByRightLineNumber[rightLineNumber]
		if !mapped {
			rightLineNumbers = append(rightLineNumbers, rightLineNumber)
		}
	}

	lineMappings := make([][]uint32, 0)
	for leftLineNumber := 0; leftLineNumber < leftLineCount; leftLineNumber++ {
		pair, exists := linePairs[leftLineNumber]
		if !exists {
			lineMappings = append(lineMappings, []uint32{uint32(leftLineNumber + 1), 0})
		} else {
			if includeIdenticalLines || !(pair.left.content == pair.right.content && leftLineNumber == pair.right.lineNumber) {
				lineMappings = append(lineMappings, []uint32{uint32(leftLineNumber + 1), uint32(pair.right.lineNumber + 1)})
			}
		}
	}
	for _, rightLineNumber := range rightLineNumbers {
		lineMappings = append(lineMappings, []uint32{0, uint32(rightLineNumber + 1)})
	}
	return lineMappings
}

func PrintMappings(mappings [][]uint32) error {
	for _, mapping := range mappings {
		_, err := fmt.Printf("%s,%s\n", toString(mapping[0]), toString(mapping[1]))
		if err != nil {
			return err
		}
	}
	return nil
}

func toString(i uint32) string {
	var left string
	if i == 0 {
		left = "_"
	} else {
		left = strconv.FormatUint(uint64(i), 10)
	}
	return left
}
