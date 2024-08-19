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
	for i := 0; i <= len(line); i++ {
		if matchPattern(line, pattern, i) {
			return true, nil
		}

	}
	return false, nil
	// // You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")
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
}
func matchPattern(line string, pattern string, pos int) bool {
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
			if !(j == 0) {
				return false
			} else {
				i++
				if line[j] != pattern[i] {
					return false
				}
			}
		} else if strings.Contains(pattern, "$") {
			endPos := strings.Index(pattern[i:], "$")
			for i := 0; i < endPos; i++ {
				if pattern[i] != line[j] {
					return false
				}
				j++
			}
			return true
		} else {
			if line[j] != pattern[i] {
				return false
			}
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
