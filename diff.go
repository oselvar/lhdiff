package lhdiff

import (
	"bytes"

	"github.com/sourcegraph/go-diff/diff"
)

func LineNumbersFromDiff(fileDiff *diff.FileDiff, leftLines []string, rightLines []string, contextSize int) ([]LinePair, []int, []int) {
	var unchangedPairs []LinePair
	// Deleted from left
	var leftDeletedLineNumbers []int
	// Added to right
	var rightDeletedLineNumbers []int

	previousLeftLineNumber := 0
	previousRightLineNumber := 0
	for _, hunk := range fileDiff.Hunks {
		unchangedHunkPairs, leftLineNumbersHunk, rightLineNumbersHunk := LineNumbersFromHunk(hunk, leftLines, rightLines, previousLeftLineNumber, previousRightLineNumber, contextSize)
		leftDeletedLineNumbers = append(leftDeletedLineNumbers, leftLineNumbersHunk...)
		rightDeletedLineNumbers = append(rightDeletedLineNumbers, rightLineNumbersHunk...)
		unchangedPairs = append(unchangedPairs, unchangedHunkPairs...)
		previousLeftLineNumber = int(hunk.OrigStartLine - 1 + hunk.OrigLines)
		previousRightLineNumber = int(hunk.NewStartLine - 1 + hunk.NewLines)
	}
	// Add unchanged lines after last hunk
	leftLineNumber := previousLeftLineNumber
	rightLineNumber := previousRightLineNumber
	for leftLineNumber >= 0 && leftLineNumber < len(leftLines) {
		leftLineInfo := MakeLineInfo(leftLineNumber, leftLines, contextSize)
		rightLineInfo := MakeLineInfo(rightLineNumber, rightLines, contextSize)
		unchangedPairs = append(unchangedPairs, LinePair{
			left:  leftLineInfo,
			right: rightLineInfo,
		})
		leftLineNumber++
		rightLineNumber++
	}
	return unchangedPairs, leftDeletedLineNumbers, rightDeletedLineNumbers
}

func LineNumbersFromHunk(hunk *diff.Hunk, leftLines []string, rightLines []string, previousLeftLineNumber int, previousRightLineNumber int, contextSize int) ([]LinePair, []int, []int) {
	var unchangedPairs []LinePair
	leftLineNumbers := make([]int, 0)
	rightLineNumbers := make([]int, 0)

	leftLineNumber := previousLeftLineNumber
	rightLineNumber := previousRightLineNumber
	for leftLineNumber < int(hunk.OrigStartLine)-1 {
		leftLineInfo := MakeLineInfo(leftLineNumber, leftLines, contextSize)
		rightLineInfo := MakeLineInfo(rightLineNumber, rightLines, contextSize)
		unchangedPairs = append(unchangedPairs, LinePair{
			left:  leftLineInfo,
			right: rightLineInfo,
		})
		leftLineNumber++
		rightLineNumber++
	}

	lines := bytes.Split(hunk.Body, []byte{'\n'})

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '-':
			leftLineNumbers = append(leftLineNumbers, leftLineNumber)
			leftLineNumber++
		case '+':
			rightLineNumbers = append(rightLineNumbers, rightLineNumber)
			rightLineNumber++
		default:
			unchangedPairs = append(unchangedPairs, LinePair{
				left:  MakeLineInfo(leftLineNumber, leftLines, contextSize),
				right: MakeLineInfo(rightLineNumber, rightLines, contextSize),
			})
			leftLineNumber++
			rightLineNumber++
		}
	}
	return unchangedPairs, leftLineNumbers, rightLineNumbers
}
