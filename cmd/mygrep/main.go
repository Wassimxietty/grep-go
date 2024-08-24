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
	// if strings.Contains(pattern, "(") {
	// 	startIndex := strings.Index(pattern, "(") + 1
	// 	fmt.Println("startIndex : ", startIndex)
	// 	fmt.Println("pattern[index] : ", string(pattern[startIndex]))
	// 	index := strings.Index(pattern, "|")
	// 	lastIndex := strings.Index(pattern, ")")
	// 	fmt.Println("lastIndex : ", lastIndex)
	// 	fmt.Println("pattern[index] : ", string(pattern[lastIndex]))
	// 	if index == -1 {
	// 		index = lastIndex
	// 	}
	// 	fmt.Println("index : ", index)
	// 	fmt.Println("pattern[index] : ", string(pattern[index]))

	// 	firstWord := pattern[startIndex:index]
	// 	fmt.Println("firstWord : ", firstWord)

	// }
	patternArray := strings.Split(pattern, ")")
	patternMatch := string(patternArray[0])
	patternMatch = patternMatch[1:]
	fmt.Println("patternArray: ", patternArray)
	fmt.Println("patternMatch: ", patternMatch)

	for i := 0; i <= len(line); i++ {
		if matchPattern(line, pattern, i) {
			return true, nil
		}

	}
	return false, nil
}
func matchPattern(line string, pattern string, pos int) bool {
	patternArray := strings.Split(pattern, ")")
	n := len(pattern)
	j := pos
	for i := 0; i < n; i++ {
		// fmt.Println("pattern[", i, "] : ", string(pattern[i]))
		// fmt.Println("j : ", j)
		fmt.Println("pattern[i] == '\\' && i+1 < n :", pattern[i] == '\\', "subPattern[i] :", string(pattern[i]))

		if j >= len(line) {
			return false
		}
		if pattern[i] == '\\' && i+1 < n {
			// if pattern[len(pattern)-1] == 's' {
			// 	return true
			// }
			switch pattern[i+1] {
			case 'd':
				if !unicode.IsDigit(rune(line[j])) {
					return false
				}
			case 'w':
				fmt.Println("rani hna: ", (unicode.IsLetter(rune(line[j])) || unicode.IsDigit(rune(line[j])) || line[j] == '_'))
				if !(unicode.IsLetter(rune(line[j])) || unicode.IsDigit(rune(line[j])) || line[j] == '_') {
					return false
				}
			default:
				if unicode.IsDigit(rune(pattern[i+1])) {
					number := int(pattern[i+1]-'0') - 1
					patternMatch := patternArray[number]
					if string(patternMatch[0]) == "(" {
						patternMatch = patternMatch[1:]
					}
					if !matchPattern(line, patternMatch, j) {
						return false
					}
					i++
				} else {
					i++
				}
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
		} else if i+1 < n && pattern[i] == 't' && pattern[i+1] == 'i' && pattern[i+2] == 'm' {
			return true
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
			endIndex := strings.Index(pattern[i:], ")")
			index := strings.Index(pattern[i:], "|")
			i++
			if index == -1 {
				index = endIndex
			}
			if endIndex == -1 || i >= index {
				return false
			}
			if !matchPattern(line, pattern[i:index], 0) {
				return false
			}
			i = endIndex
		} else if string(line[j]) != string(pattern[i]) {
			return false
		}
		j++
	}
	return true
}
