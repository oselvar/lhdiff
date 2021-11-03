package lhdiff

import (
	"bytes"
	"fmt"
	"github.com/ianbruene/go-difflib/difflib"
	levenshtein "github.com/ka-weihe/fast-levenshtein"
	"github.com/sourcegraph/go-diff/diff"
	"math"
	"regexp"
	"sort"
	"strings"
)

type LineInfo struct {
	lineNumber int32
	content    string
	context    string
}

type LinePair struct {
	left  LineInfo
	right LineInfo
}

func (linePair LinePair) contentNormalizedLevenshteinSimilarity() float64 {
	distance := levenshtein.Distance(linePair.left.content, linePair.right.content)
	normalizedLevenhsteinDistance := float64(distance) / math.Max(float64(len(linePair.left.content)), float64(len(linePair.right.content)))
	return 1 - normalizedLevenhsteinDistance
}

func (linePair LinePair) contextTfIdfCosineSimilarity() float64 {
	return TfIdfCosineSimilarity(linePair.left.context, linePair.right.context)
}

func (linePair LinePair) combinedSimilarity() float64 {
	contentSimilarity := linePair.contentNormalizedLevenshteinSimilarity()
	if contentSimilarity <= 0.5 {
		return 0.0
	}
	contextSimilarity := linePair.contextTfIdfCosineSimilarity()
	return ContentSimilarityFactor*contentSimilarity + ContextSimilarityFactor*contextSimilarity
}

type ByCombinedSimilarity []LinePair

func (a ByCombinedSimilarity) Len() int { return len(a) }
func (a ByCombinedSimilarity) Less(i, j int) bool {
	return a[j].combinedSimilarity() < a[i].combinedSimilarity()
}
func (a ByCombinedSimilarity) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

const ContextSimilarityFactor = 0.4
const ContentSimilarityFactor = 0.6
const SimilarityThreshold = 0.45

func Lhdiff(left string, right string, contextSize int) (map[int32]LinePair, int32, []int) {
	leftLines := ConvertToLinesWithoutNewLine(left)
	rightLines := ConvertToLinesWithoutNewLine(right)

	mappedRightLines := make(map[int32]bool)
	allPairs := make(map[int32]LinePair, 0)

	diffScript, _ := difflib.GetUnifiedDiffString(difflib.LineDiffParams{
		A:        leftLines,
		B:        rightLines,
		FromFile: "left",
		ToFile:   "right",
		Context:  3,
	})
	// fmt.Println(diffScript)
	if diffScript != "" {
		fileDiff, _ := diff.ParseFileDiff([]byte(diffScript))

		unchangedDiffPairs, leftLineNumbers, rightLineNumbers := LineNumbersFromDiff(fileDiff, leftLines, rightLines, contextSize)
		for _, unchangedDiffPair := range unchangedDiffPairs {
			allPairs[unchangedDiffPair.left.lineNumber] = unchangedDiffPair
			mappedRightLines[unchangedDiffPair.right.lineNumber] = true
		}

		leftLineInfos := MakeLineInfos(leftLineNumbers, leftLines, contextSize)
		rightLineInfos := MakeLineInfos(rightLineNumbers, rightLines, contextSize)

		for _, rightLineInfo := range rightLineInfos {
			var similarPairCandidates []LinePair
			for _, leftLineInfo := range leftLineInfos {
				pair := LinePair{
					left:  leftLineInfo,
					right: rightLineInfo,
				}
				similarPairCandidates = append(similarPairCandidates, pair)
			}
			sort.Sort(ByCombinedSimilarity(similarPairCandidates))
			if len(similarPairCandidates) > 0 {
				mostSimilarPair := similarPairCandidates[0]
				if mostSimilarPair.combinedSimilarity() > SimilarityThreshold {
					allPairs[mostSimilarPair.left.lineNumber] = mostSimilarPair
					mappedRightLines[mostSimilarPair.right.lineNumber] = true
				}
			}
		}
	} else {
		// The files are identical
		for leftLineNumber := range leftLines {
			lineInfo := MakeLineInfo(int32(leftLineNumber), leftLines, 4)
			allPairs[int32(leftLineNumber)] = LinePair{
				left:  lineInfo,
				right: lineInfo,
			}
			mappedRightLines[int32(leftLineNumber)] = true
		}
	}
	rightLineNumbers := make([]int, 0)
	for rightLineNumber, _ := range rightLines {
		_,mapped := mappedRightLines[int32(rightLineNumber)]
		if !mapped {
			rightLineNumbers = append(rightLineNumbers, rightLineNumber)
		}
	}
	return allPairs, int32(len(leftLines)), rightLineNumbers
}

func PrintLinePairs(linePairs map[int32]LinePair, leftLineCount int32, newRightLines []int, lines bool) {
	for leftLineNumber := int32(0); leftLineNumber < leftLineCount; leftLineNumber++ {
		pair, exists := linePairs[leftLineNumber]
		if !exists {
			fmt.Printf("%d,_\n", leftLineNumber+1)
		} else {
			if lines {
				fmt.Printf("%d:%s,%d:%s\n", leftLineNumber+1, strings.TrimSpace(pair.left.content), pair.right.lineNumber+1, strings.TrimSpace(pair.right.content))
			} else {
				fmt.Printf("%d,%d\n", leftLineNumber+1, pair.right.lineNumber+1)
			}
		}
	}
	for _, rightLine := range newRightLines {
		fmt.Printf("_,%d\n", rightLine+1)
	}
}

func MakeLineInfos(lineNumbers []int32, lines []string, contextSize int) []LineInfo {
	lineInfos := make([]LineInfo, len(lineNumbers))
	for i, lineNumber := range lineNumbers {
		lineInfos[i] = MakeLineInfo(lineNumber, lines, contextSize)
	}
	return lineInfos
}

func MakeLineInfo(lineNumber int32, lines []string, contextSize int) LineInfo {
	content := lines[lineNumber]
	context := GetContext(lineNumber, lines, contextSize)
	lineInfo := LineInfo{
		lineNumber: lineNumber,
		context:    context,
		content:    content,
	}
	return lineInfo
}

var /* const */ brackets = regexp.MustCompile("^[{()}]$")

// GetContext returns a string consisting of (up to) contextSize context lines above and below lineIndex.
// a line is considered to be a context line if it is not an "insignificant" line, i.e. either blank
// or just a curly brace or parenthesis (whitespace trimmed).
func GetContext(lineNumber int32, lines []string, contextSize int) string {
	var context []string

	i := int(lineNumber) - 1

	for j := 0; i >= 0 && j < contextSize; {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if len(trimmed) != 0 && !brackets.MatchString(trimmed) {
			context = append([]string{line}, context...)
			j++
		}
		i--
	}

	i = int(lineNumber) + 1
	for j := 0; i < len(lines) && j < contextSize; {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if len(trimmed) != 0 && !brackets.MatchString(trimmed) {
			context = append(context, line)
			j++
		}
		i++
	}

	return strings.Join(context, "")
}

// LineNumbersFromDiff returns two slices:
// 1: a slice of removed line numbers in left
// 2: a slice of added line numbers in right
// 3:
func LineNumbersFromDiff(fileDiff *diff.FileDiff, leftLines []string, rightLines []string, contextSize int) ([]LinePair, []int32, []int32) {
	var unchangedPairs []LinePair
	// Deleted from left
	var leftLineNumbers []int32
	// Added to right
	var rightLineNumbers []int32

	previousLeftLineNumber := int32(0)
	previousRightLineNumber := int32(0)
	for _, hunk := range fileDiff.Hunks {
		unchangedHunkPairs, leftLineNumbersHunk, rightLineNumbersHunk := LineNumbersFromHunk(hunk, leftLines, rightLines, previousLeftLineNumber, previousRightLineNumber, contextSize)
		leftLineNumbers = append(leftLineNumbers, leftLineNumbersHunk...)
		rightLineNumbers = append(rightLineNumbers, rightLineNumbersHunk...)
		unchangedPairs = append(unchangedPairs, unchangedHunkPairs...)
		previousLeftLineNumber = hunk.OrigStartLine - 1 + hunk.OrigLines
		previousRightLineNumber = hunk.NewStartLine - 1 + hunk.NewLines
	}
	// Add unchanged lines after last hunk
	leftLineNumber := previousLeftLineNumber
	rightLineNumber := previousRightLineNumber
	for int(leftLineNumber) < len(leftLines) {
		leftLineInfo := MakeLineInfo(leftLineNumber, leftLines, contextSize)
		rightLineInfo := MakeLineInfo(rightLineNumber, rightLines, contextSize)
		unchangedPairs = append(unchangedPairs, LinePair{
			left:  leftLineInfo,
			right: rightLineInfo,
		})
		leftLineNumber++
		rightLineNumber++
	}
	return unchangedPairs, leftLineNumbers, rightLineNumbers
}

func LineNumbersFromHunk(hunk *diff.Hunk, leftLines []string, rightLines []string, previousLeftLineNumber int32, previousRightLineNumber int32, contextSize int) ([]LinePair, []int32, []int32) {
	var unchangedPairs []LinePair
	var leftLineNumbers []int32
	var rightLineNumbers []int32

	leftLineNumber := previousLeftLineNumber
	rightLineNumber := previousRightLineNumber
	for leftLineNumber < hunk.OrigStartLine-1 {
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

func ConvertToLinesWithoutNewLine(text string) []string {
	lines := strings.SplitAfter(text, "\n")
	return Map(lines, RemoveMultipleSpaceAndTrim)
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func RemoveMultipleSpaceAndTrim(s string) string {
	re := regexp.MustCompile("[ \t]+")
	return strings.TrimSpace(re.ReplaceAllString(s, " ")) + "\n"
}
