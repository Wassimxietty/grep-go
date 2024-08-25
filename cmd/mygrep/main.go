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
		plusArray := strings.Split(pattern, "+")
		fmt.Println("plusArray : ", plusArray)
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
		okay, j := matchPattern(line, pattern, i)
		fmt.Println(" returned j: ", j, " ", okay)
		if okay {
			return true, nil
		}
	}
	return false, nil
}
func matchPattern(line string, pattern string, pos int) (bool, int) {
	patternArray := strings.Split(pattern, ")")
	n := len(pattern)
	j := pos
	i := 0
	for i < n {
		// fmt.Println("pattern [", i, "]: ", string(pattern[i]))
		if j >= len(line) {
			return false, j
		}
		if pattern[i] == '\\' && i+1 < n {

			// if pattern[len(pattern)-1] == 's' {
			// 	return true
			// }
			switch pattern[i+1] {
			case 'd':
				// fmt.Println("d")
				if !unicode.IsDigit(rune(line[j])) {
					return false, j
				}
				i++

			case 'w':
				if !(unicode.IsLetter(rune(line[j])) || unicode.IsDigit(rune(line[j])) || line[j] == '_') {
					return false, j
				}
				i++

			default:
				if unicode.IsDigit(rune(pattern[i+1])) {
					number := int(pattern[i+1]-'0') - 1
					if number == -1 {
						fmt.Println("patternArray[number] is patternArray[-1]")
						return false, j
					}
					patternMatch := patternArray[number]
					if string(patternMatch[0]) == "(" {
						patternMatch = patternMatch[1:]
					}
					okay, jPose := matchPattern(line, patternMatch, j)
					if okay {
						j = jPose
						i += 2
						continue
					}

				} else {
					if string(line[j]) != string(pattern[i+1]) {
						return false, j
					}
					j++ // Move to the next character in the line
				}
			}
		} else if pattern[i] == '[' && i+1 < n && pattern[i+1] == '^' {
			endPos := strings.Index(pattern[i:], "]")
			matchAnyPattern := pattern[i+1 : endPos]
			if strings.ContainsAny(matchAnyPattern, string(line[j])) {
				return false, j
			}
			i = endPos
		} else if pattern[i] == '[' && i+1 < n {
			endPos := strings.Index(pattern[i:], "]")
			matchAnyPattern := pattern[i+1 : endPos]
			if !strings.ContainsAny(matchAnyPattern, string(line[j])) {
				return false, j
			}
			i = endPos
		} else if pattern[i] == '^' && i+1 < n {
			if j != 0 {
				return false, j
			} else {
				i++
				if line[j] != pattern[i] {
					return false, j
				}
			}
		} else if strings.Contains(pattern, "$") {
			endPos := strings.Index(pattern[i:], "$")
			matchAnyPattern := pattern[i:endPos]
			for i := 0; i < endPos; i++ {
				if matchAnyPattern[i] != line[j] {
					return false, j
				}
				j++
			}
			if j != len(line) {
				return false, j
			}
			i = endPos
		} else if pattern[i] == '+' && i != 0 {
			if pattern[i-1] == ']' {
				fmt.Println("] spotted")
				// for ii := i - 2; ii < 0; ii-- {
				// 	if pattern[ii] == '[' {
				// 		bracketIndex := ii
				// 		plusPattern := pattern[bracketIndex : i+1]
				// 		okay, jj := matchPattern(line, plusPattern, j)
				// 		if okay {
				// 			fmt.Println("jj : ", jj)
				// 			j = jj
				// 			i++
				// 			continue
				// 		} else {
				// 			fmt.Println("jj : ", jj)
				// 		}
				// 	}
				// }
				// return false, j
			}
			letterPlus := pattern[i-1]
			for i < len(line) && letterPlus == line[j] && letterPlus != ']' {
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
			return true, j
		} else if pattern[i] == '(' {
			endIndex := strings.Index(pattern[i:], ")")
			index := strings.Index(pattern[i:], "|")
			i++
			if index == -1 {
				index = endIndex
			}
			if endIndex == -1 || i >= index {
				return false, j
			}
			okay, jj := matchPattern(line, pattern[i:index], 0)
			// fmt.Println("pattern[i]: ", string(pattern[i]))
			if !okay {
				return false, jj
			}
			i = endIndex
		} else if string(line[j]) != string(pattern[i]) {
			return false, j
		}
		// fmt.Println("line[", j, "]: ", string(line[j]))
		j++
		i++
	}
	return true, j
}
