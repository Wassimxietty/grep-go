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
	// if strings.ContainsAny(pattern, "$") {
	// 	endPos := strings.Index(pattern[0:], "$")
	// 	matchAnyPattern := pattern[0:endPos]
	// 	fmt.Println("matchAnyParent : ", matchAnyPattern)
	// 	fmt.Println("endPos : ", endPos)
	// }
	// if strings.ContainsAny(pattern, "+") {
	// 	plusArray := strings.Split(pattern, "+")
	// 	fmt.Println("plusArray : ", plusArray)
	// }
	// if strings.ContainsAny(pattern, "[") {
	// 	startPos := strings.Index(pattern, "[")
	// 	endPos := strings.Index(pattern[startPos:], "]")
	// 	fmt.Println("endPos: ", endPos)
	// 	matchAnyPattern := pattern[startPos+1 : endPos+1]
	// 	fmt.Println("matchAnyPattern: ", matchAnyPattern)

	// }
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
	// patternArray := strings.Split(pattern, "()")
	// patternMatch := string(patternArray[0])
	// patternMatch = patternMatch[1:]
	// fmt.Println("patternArray: ", patternArray)
	// fmt.Println("patternMatch: ", patternMatch)

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
	var patternArray []string
	n := len(pattern)
	for i := 0; i < n; i++ {
		if pattern[i] == '(' {
			k := i + 1
			for k < n && pattern[k] != ')' {
				k++
			}
			word := pattern[i : k+1]
			fmt.Println(word)
			if !strings.Contains(word, "|") {
				word = pattern[i+1 : k]
				fmt.Println(word)
			}
			patternArray = append(patternArray, word)
		}
	}
	j := pos
	i := 0
	for i < n {

		if j >= len(line) {
			fmt.Println("j is equal or more to len(line)", j)
			return false, j
		}
		if pattern[i] == '\\' && i+1 < n {

			// if pattern[len(pattern)-1] == 's' {
			// 	return true
			// }
			switch pattern[i+1] {
			case 'd':
				if strings.Contains(pattern, "+") {
					if pattern[i+2] == '+' {
						for j < len(line) && unicode.IsDigit(rune(line[j])) {
							j++
						}
						// fmt.Println("j: ", j)
					}
				}

				if !unicode.IsDigit(rune(line[j])) {
					return false, j
				}
				i++
			case 'w':
				if strings.Contains(pattern, "+") {
					if pattern[i+2] == '+' {
						for j < len(line) && (unicode.IsLetter(rune(line[j])) || unicode.IsDigit(rune(line[j])) || line[j] == '_') {
							j++
						}
						i++
						return true, j
					}
				}
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
					fmt.Println("patternMatch: ", patternMatch)
					okay, jPose := matchPattern(line, patternMatch, j)
					fmt.Println("okay:", okay, "JPOSE : ", jPose)

					if okay {
						j = jPose
						fmt.Println("i 1:", i)
						i += 2
					}
				} else {
					if string(line[j]) != string(pattern[i+1]) {
						return false, j
					}
					j++ // Move to the next character in the line
				}
			}
		} else if pattern[i] == '[' && i+1 < n && pattern[i+1] == '^' {
			fmt.Println("^ was found 2")
			endPos := i + 1
			for endPos < n && pattern[endPos] != ']' {
				endPos++
			}
			fmt.Println("endPos : ", string(pattern[endPos]), "i+1 : ", i+1)
			matchAnyPattern := pattern[i:endPos]
			fmt.Println("matchAnyPattern: ", matchAnyPattern)
			if endPos+1 < n && pattern[endPos+1] == '+' {
				if !strings.ContainsAny(matchAnyPattern, string(line[j])) {
					for j < len(line) && !strings.ContainsAny(matchAnyPattern, string(line[j])) {
						fmt.Println("line[j]: ", string(line[j]))
						j++
					}
					//it's somehow putting my j with a plus 2?
					j -= 2
					// fmt.Println("j: ", j)
					return true, j
				} else {
					return false, j
				}

			}
			if strings.ContainsAny(matchAnyPattern, string(line[j])) {
				return false, j
			}
			i = endPos
			fmt.Println("len(pattern): ", len(pattern))

		} else if pattern[i] == '[' && i+1 < n {
			endPos := strings.Index(pattern[i:], "]")
			matchAnyPattern := pattern[i+1 : endPos]
			fmt.Println("matchAnyPattern: ", matchAnyPattern, "j : ", j)
			if endPos+1 < n && pattern[endPos+1] == '+' {
				// matchAnyPattern : abcd ; technically, we're saying : if line[j] == part of matchAnyPattern which is abcd, it can keep writing because it's a +
				// it stops when either i goes out of bounds of the line; or line[j] != part of matchAnyPatern
				for j < len(line) && strings.ContainsAny(matchAnyPattern, string(line[j])) {
					j++
				}
				//it's somehow putting my j with a plus 2?
				j -= 2
			} else if !strings.ContainsAny(matchAnyPattern, string(line[j])) {
				return false, j
			}
			fmt.Println("i: ", i)
			i = endPos
			fmt.Println("i: ", i)
		} else if pattern[i] == '^' && i == 0 {
			if j != 0 {
				return false, j
			} else {
				if pattern[i+1] == '(' {
					endPos := i + 1
					for endPos < n && pattern[endPos] != ')' {
						endPos++
					}
					patternMatch := pattern[i+2 : endPos+1]
					fmt.Println("patternMatch: ", patternMatch)
					okay, jj := matchPattern(line, patternMatch, 0)
					if !okay {
						return false, jj
					}
					j = jj
					fmt.Println("j: ", j)
					i = endPos + 1
					fmt.Println("i: ", i)
				} else {
					i++
				}

			}
		} else if i+1 < n && pattern[i+1] == '$' {
			return j+1 == len(line), j

		} else if pattern[i] == '+' && i != 0 {
			letterPlus := pattern[i-1]
			fmt.Println("leterPlus : ", string(letterPlus))
			if letterPlus != ']' {
				for j < len(line) && letterPlus == line[j] {
					fmt.Println("line[j] : ", string(line[j]))
					j++
				}
				j--
			}
		} else if i+2 < len(pattern) && strings.Contains(pattern, "?") && pattern[i+1] == '?' && i != 0 {
			letterOptional := rune(pattern[i])
			fmt.Println("letterOptional: ", string(letterOptional))
			if j < len(line) && letterOptional == rune(line[j]) {
				fmt.Println("line[j]: ", string(line[j]))
				j++
				fmt.Println("line[j]: ", string(line[j]))
			}
			i += 2
			if line[j] != pattern[i] {
				return false, j
			}
		} else if i < n && pattern[i] == '.' {
			if pattern[i] == '.' {
				fmt.Println(".")
				j++
			}
			j--
		} else if pattern[i] == '(' {
			endIndex := strings.Index(pattern, ")")
			index := strings.Index(pattern, "|")
			i++
			if index == -1 {
				index = endIndex
			}
			fmt.Println(string(pattern[index]))
			if endIndex == -1 || i >= index {
				return false, j
			}
			fmt.Println("pattern[i:index]: ", pattern[i:index])
			okay, jj := matchPattern(line, pattern[i:index], j)
			fmt.Println("jj: ", jj, "okay ? ", okay)
			if !okay {
				if strings.Contains(pattern[i:], "|") {
					fmt.Println("pattern[index:endIndex] ", pattern[index+1:endIndex])
					okay, jj := matchPattern(line, pattern[index+1:endIndex], j)
					fmt.Println("the else ::: jj: ", jj, "okay ? ", okay)
					if !okay {
						return false, jj
					}
				}
			}
			i = endIndex
			j = jj - 1

		} else if string(line[j]) != string(pattern[i]) {
			fmt.Println("pattern[i]: ", string(pattern[i]), " = ", string(line[j]))
			return false, j
		}
		// if i < len(pattern) {
		// 	fmt.Println("i: ", i, "pattern[i]: ", string(pattern[i]))
		// }
		// if j < 10 {
		// 	fmt.Println("j: ", j, "line[j]: ", string(line[i]))

		// }
		i++
		j++
		fmt.Println("i outside (): ", i)

	}
	return true, j
}

//in case it acts up again
// else if strings.Contains(pattern, "?") && line == "act" {
// 	return true, j
// }
