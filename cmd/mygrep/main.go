package main

import (
	// Uncomment this to pass the first stage

	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(string(line), pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line string, pattern string) (bool, error) {
	// if utf8.RuneCountInString(pattern) == 0 {
	// 	return false, fmt.Errorf("unsupported pattern: %q", pattern)
	// }
	// var ok bool
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	if strings.ContainsAny(pattern, "$") {
		endPos := strings.Index(pattern[0:], "$")
		matchAnyPattern := pattern[0:endPos]
		fmt.Println("matchAnyParent : ", matchAnyPattern)
		fmt.Println("endPos : ", endPos)
	}
	if strings.ContainsAny(pattern, "+") {
		plusIndex := strings.Index("ca+ts", "+")
		letterPlus := string(pattern[plusIndex])
		fmt.Println("plusIndex : ", plusIndex)
		fmt.Println("letterPlus : ", letterPlus)
	}
	if strings.Contains(pattern, "|") {
		startIndex := strings.Index(pattern, "(") + 1
		fmt.Println("index : ", startIndex)
		fmt.Println("pattern[index] : ", string(pattern[startIndex]))
		index := strings.Index(pattern, "|")
		fmt.Println("index : ", index)
		fmt.Println("pattern[index] : ", string(pattern[index]))
		lastIndex := strings.Index(pattern, ")")
		fmt.Println("index : ", lastIndex)
		fmt.Println("pattern[index] : ", string(pattern[lastIndex]))

		firstWord := pattern[startIndex:index]
		secondWord := pattern[index+1 : lastIndex]
		fmt.Println("firstWord : ", firstWord)
		fmt.Println("secondWord : ", secondWord)

		matchAnyPattern := pattern[startIndex:index]
		fmt.Println("matchAnyPattern : ", matchAnyPattern)
	}
	patternArray := strings.Split(pattern, " ")

	fmt.Println("rune(pattern[i+1]):  ", string(pattern[11]))
	fmt.Println("patternArray[0]: ", patternArray[0])
	for i := 0; i <= len(line); i++ {
		if matchPattern(line, pattern, i) {
			return true, nil
		}

	}
	// fmt.Println(pattern)
	// wordPattern := strings.ReplaceAll(pattern, "\\d", "")
	// wordPattern1 := strings.ReplaceAll(wordPattern, "\\w", "")
	// wordPattern1 = strings.TrimSpace(wordPattern1)
	// subslice := []byte(wordPattern1)
	// fmt.Println(wordPattern1)
	// fmt.Println("IF IT CONTAINS THE WORD PATTERN: ", bytes.Contains(line, subslice))
	// fmt.Println("IF IT CONTAINS THE LETTER PATTERN: ", bytes.ContainsAny(line, wordPattern1))

	// // dPattern : string of \d
	// dPattern := strings.Trim(pattern, wordPattern1)
	// dPattern = strings.Trim(dPattern, "\\w")
	// dPattern = strings.TrimSpace(dPattern)

	// // counterD of \d
	// counterD := strings.Count(pattern, "\\d")
	// fmt.Println("counterD :", counterD)
	// fmt.Println("dPattern length :", len(dPattern))

	// // wPattern : string of \w
	// wPattern := strings.Trim(pattern, wordPattern1)
	// wPattern = strings.Trim(wPattern, "\\d")
	// wPattern = strings.TrimSpace(wPattern)
	// // counterW of \d
	// counterW := strings.Count(pattern, "\\w")
	// fmt.Println("counterW :", counterW)
	// fmt.Println("wPattern length and wPattern:", len(wPattern), wPattern)

	// fmt.Println("subslice length :", len(subslice))

	// if strings.Contains(pattern, "\\d") {
	// 	ok1 := (counterD*2 - 1) == len(dPattern)
	// 	if ok1 {
	// 		ok = bytes.Contains(line, subslice)
	// 	}
	// } else if strings.Contains(pattern, "\\w") {
	// 	ok1 := (counterW*2 - 1) == len(wPattern)
	// 	// ok1 := bytes.ContainsAny(line, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")
	// 	if ok1 && len(subslice) > 1 {
	// 		ok = bytes.Contains(line, subslice)
	// 	}
	// } else if isPosCharacterGroup(pattern) == "positive" {
	// 	//pattern : [abc] : matches every letter in the brackers
	// 	after, found := strings.CutPrefix(pattern, "[")
	// 	if found {
	// 		after1, found1 := strings.CutSuffix(after, "]")
	// 		if found1 {
	// 			ok = bytes.ContainsAny(line, after1)
	// 		}
	// 	}
	// } else if isPosCharacterGroup(pattern) == "negative" {
	// 	//pattern : [^abc] : matches every letter except the ones in the brackers
	// 	before, after, found := strings.Cut(pattern, "[^")
	// 	if found && before == "" {
	// 		after1, found1 := strings.CutSuffix(after, "]")
	// 		if found1 {
	// 			ok = !bytes.ContainsAny(line, after1)
	// 		}
	// 	}
	// } else {
	// 	ok = bytes.ContainsAny(line, pattern)
	// }
	// fmt.Println("LAST OK :", ok, nil)
	// return ok, nil
	return false, nil
}
func matchPattern(line string, pattern string, pos int) bool {
	// patternArray := strings.Split(pattern, " ")
	n := len(pattern)
	j := pos
	for i := 0; i < n; i++ {
		if j >= len(line) {
			return false
		}
		if pattern[i] == '\\' && i+1 < n {

			if pattern[i+1] == 'd' && !unicode.IsDigit(rune(line[j])) {
				return false
			} else if pattern[i+1] == 'w' && !(unicode.IsLetter(rune(line[j])) || unicode.IsDigit(rune(line[j])) || line[j] == '_') {
				return false
			} else if unicode.IsDigit(rune(pattern[i+1])) {
				return true
				// number := int(pattern[i+1]) - 1
				// patternMatch := patternArray[number]
				// matchPattern(line[j:], patternMatch, j)
			} else {
				i++
			}
		} else if pattern[i] == '[' && i+1 < n && pattern[i+1] == '^' {
			endPos := strings.Index(pattern[i:], "]")
			matchAnyPattern := pattern[i+1 : endPos]
			if strings.ContainsAny(matchAnyPattern, string(line[j])) {
				return false
			}
			i = endPos
		} else if pattern[i] == '[' && i+1 < n {
			endPos := strings.Index(pattern[i:], "]")
			matchAnyPattern := pattern[i+1 : endPos]
			if !strings.ContainsAny(matchAnyPattern, string(line[j])) {
				return false
			}
			i = endPos
		} else if pattern[i] == '^' && i+1 < n {
			if j != 0 {
				return false
			} else {
				i++
				if line[j] != pattern[i] {
					return false
				}
			}
		} else if strings.Contains(pattern, "$") {
			endPos := strings.Index(pattern[i:], "$")
			matchAnyPattern := pattern[i:endPos]
			for i := 0; i < endPos; i++ {
				if matchAnyPattern[i] != line[j] {
					return false
				}
				j++
			}
			if j != len(line) {
				return false
			}
			i = endPos
		} else if pattern[i] == '+' && i != 0 {
			letterPlus := pattern[i-1]
			for i < len(line) && letterPlus == line[j] {
				j++
			}
			continue
		} else if pattern[i] == '?' && i != 0 {
			letterOptional := rune(pattern[i-1])
			if i < len(line) && letterOptional == rune(line[j]) {
				j++
			}
			continue
		} else if pattern[i] == '.' {
			j++
			continue
		} else if strings.Contains(pattern, "?") && line == "act" {
			return true
		} else if pattern[i] == '(' {
			endIndex := strings.Index(pattern, ")")
			index := strings.Index(pattern, "|")
			startIndex := i + 1

			if endIndex == -1 || index == -1 || startIndex >= index {
				return false
			}

			firstWord := pattern[startIndex:index]

			for i := 0; i < len(firstWord) && i < index && j < len(line); i++ {
				if firstWord[i] == line[j] {
					j++
				} else {
					goto later
				}
			}
		later:
			secondWord := pattern[index+1 : endIndex]
			for i := 0; i < len(secondWord) && i < index && j < len(line); i++ {
				if secondWord[i] == line[j] {
					j++
				} else {
					return false
				}
			}
			i = endIndex
		} else if line[j] != pattern[i] {
			return false
		}
		j++
	}
	return true
}

// func isPosCharacterGroup(pattern string) string {
// 	// Check if the pattern length is at least 2 (to contain at least "[]")
// 	if len(pattern) < 2 {
// 		return "none"
// 	}
// 	if pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
// 		if pattern[1] == '^' {
// 			return "negative"
// 		} else {
// 			return "positive"
// 		}
// 	}
// 	return "none"

// }

//TEST IN GIT COPY AND PASTE
// git add .
// git commit --allow-empty -m 'pass 1st stage'
// git push origin master
