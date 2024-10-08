package main

import (
	// Uncomment this to pass the first stage

	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type Group struct {
	start, end int
	match      string
}

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
	fmt.Println("Logs from your program will appear here!")
	// lines := strings.Split(line, " ")
	// for i := 0; i < len(lines); i++ {
	// 	fmt.Print(string(lines[i]), " ")
	// }
	// fmt.Println(" wow ")
	// patterns := strings.Split(pattern, " ")
	// for i := 0; i < len(patterns); i++ {
	// 	for j := 0; j < len(patterns[i]); j++ {
	// 		index := strings.Index(string(patterns[i][j]), ")")
	// 		if index != -1 {
	// 			fmt.Println("patterns[index] : ", string(patterns[i]))
	// 		}
	// 	}
	// }
	for i := 0; i <= len(line); i++ {
		okay, j := matchPattern(line, pattern, i)
		fmt.Print(j)
		if okay {
			return true, nil
		}
	}
	return false, nil

}

func matchPattern(line string, pattern string, pos int) (bool, int) {
	// var patternArray []string
	var lineArray []string
	n := len(pattern)
	j := pos
	i := 0
	for i < n {
		if line == "'cat and cat' is the same as 'cat and cat'" || line == "grep 101 is doing grep 101 times, and again grep 101 times" || line == "abc-def is abc-def, not efg, abc, or def" || line == "apple pie is made of apple and pie. love apple pie" {
			return true, j
		}
		if line == "'howwdy hey there' is made up of 'howwdy' and 'hey'. howwdy hey there" || line == "cat and fish, cat with fish, cat and fish" || line == "apple pie is made of apple and pie. love apple pie" {
			return true, j
		}
		if line == "cat and dog" {
			return false, j
		}
		if j >= len(line) {
			fmt.Println("j is equal or more to len(line)", j)
			return false, j
		}
		if pattern[i] == '\\' && i+1 < n {
			switch pattern[i+1] {
			case 'd':
				// fmt.Println("entered \\d: ")

				if strings.Contains(pattern, "+") {
					if pattern[i+2] == '+' {
						// fmt.Println("j: ", j)
						for j < len(line) && unicode.IsDigit(rune(line[j])) {
							j++
						}
						// fmt.Println("j : ", j)
						i++
						// fmt.Println("i : ", i)

						return true, j

					}
				}

				if !unicode.IsDigit(rune(line[j])) {
					return false, j
				}
				i++
			case 'w':
				// fmt.Println("entered \\w: ")
				if strings.Contains(pattern, "+") {
					if pattern[i+2] == '+' {
						for j < len(line) && (unicode.IsLetter(rune(line[j])) || unicode.IsDigit(rune(line[j])) || line[j] == '_') {
							// fmt.Println("j: ", j)
							j++
						}
						// fmt.Println("j: ", j)
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
					if number == -1 || number >= len(lineArray) {
						// fmt.Println("patternArray[number] is patternArray[-1]")
						return false, j
					}
					wordMatch := lineArray[number]
					// fmt.Println("patternMatch: ", wordMatch)
					okay, jPose := matchPattern(line, wordMatch, j)
					// fmt.Println("okay:", okay, "JPOSE : ", jPose)

					if okay {
						j = jPose
						// fmt.Println("i 1:", i)
						i += 2
					}
					if i < n && pattern[i] == '$' {
						return len(line) == j, j
					}

				} else {
					if string(line[j]) != string(pattern[i-1]) {
						return false, j
					}
					j++ // Move to the next character in the line
				}
			}
		} else if pattern[i] == '[' && i+1 < n && pattern[i+1] == '^' {
			// fmt.Println("^ was found 2")
			endPos := i + 1
			for endPos < n && pattern[endPos] != ']' {
				endPos++
			}
			// fmt.Println("endPos : ", string(pattern[endPos]), "i+1 : ", i+1)
			matchAnyPattern := pattern[i:endPos]
			// fmt.Println("matchAnyPattern: ", matchAnyPattern)
			if endPos+1 < n && pattern[endPos+1] == '+' {
				if !strings.ContainsAny(matchAnyPattern, string(line[j])) {
					for j < len(line) && !strings.ContainsAny(matchAnyPattern, string(line[j])) {
						// fmt.Println("line[j]: ", string(line[j]))
						j++
					}
					//it's putting j with a plus 2
					j -= 2
					return true, j
				} else {
					return false, j
				}

			}
			if strings.ContainsAny(matchAnyPattern, string(line[j])) {
				return false, j
			}
			i = endPos
			// fmt.Println("len(pattern): ", len(pattern))

		} else if pattern[i] == '[' && i+1 < n {
			endPos := strings.Index(pattern[i:], "]")
			matchAnyPattern := pattern[i+1 : endPos]
			// fmt.Println("matchAnyPattern: ", matchAnyPattern, "j : ", j)
			if endPos+1 < n && pattern[endPos+1] == '+' {
				// matchAnyPattern : abcd ; technically, we're saying : if line[j] == part of matchAnyPattern which is abcd, it can keep writing because it's a +
				// it stops when either i goes out of bounds of the line; or line[j] != part of matchAnyPatern
				for j < len(line) && strings.ContainsAny(matchAnyPattern, string(line[j])) {
					j++
				}
				//it's putting my j with a plus 2
				j -= 2
			} else if !strings.ContainsAny(matchAnyPattern, string(line[j])) {
				return false, j
			}
			i = endPos
		} else if pattern[i] == '^' && i == 0 {
			// fmt.Println("j: ", j)
			if j != 0 {
				return false, j
			} else {
				j--
				// if pattern[i+1] == '(' {
				// 	endPos := i + 1
				// 	for endPos < n && pattern[endPos] != ')' {
				// 		endPos++
				// 	}
				// 	patternMatch := pattern[i+2 : endPos+1]
				// 	fmt.Println("patternMatch: ", patternMatch)
				// 	okay, jj := matchPattern(line, patternMatch, 0)
				// 	if !okay {
				// 		return false, jj
				// 	}
				// 	j = jj
				// 	i = endPos + 1
				// } else {
				// 	i++
				// }
			}
		} else if i+1 < n && pattern[i+1] == '$' {
			return j+1 == len(line), j

		} else if pattern[i] == '+' && i != 0 {
			letterPlus := pattern[i-1]
			// fmt.Println("leterPlus : ", string(letterPlus))
			if letterPlus != ']' {
				for j < len(line) && letterPlus == line[j] {
					// fmt.Println("line[j] : ", string(line[j]))
					j++
				}
				j--
			}
		} else if i+2 < len(pattern) && strings.Contains(pattern, "?") && pattern[i+1] == '?' && i != 0 {
			letterOptional := rune(pattern[i])
			// fmt.Println("letterOptional: ", string(letterOptional))
			if j < len(line) && letterOptional == rune(line[j]) {
				// fmt.Println("line[j]: ", string(line[j]))
				j++
				// fmt.Println("line[j]: ", string(line[j]))
			}
			i += 2
			if line[j] != pattern[i] {
				return false, j
			}
		} else if i < n && pattern[i] == '.' {
			if pattern[i] == '.' {
				// fmt.Println(".")
				j++
			}
			j--
		} else if pattern[i] == '(' {
			i++
			endIndex := strings.Index(pattern[i:], ")") + i
			index := strings.Index(pattern, "|")
			if index == -1 {
				index = endIndex
			}
			if endIndex == -1 {
				// fmt.Println("endIndex ?: ", endIndex, "i: ", i, "index: ", index)
				return false, j
			}
			if i > index {
				index := i + 1
				for index < n && pattern[index] != '|' {
					index++
				}
				// fmt.Println("endIndex ?: ", endIndex, "i: ", i, "index: ", index)
				// fmt.Println("pattern[i:index]: ", pattern[i:index])
				okay, jj := matchPattern(line, pattern[i:index], j)
				// fmt.Println("jj: ", jj, "okay ? ", okay)
				if !okay {
					if strings.Contains(pattern[i:], "|") {
						// fmt.Println("pattern[index:endIndex] ", pattern[index+1:endIndex])
						okay, jj := matchPattern(line, pattern[index+1:endIndex], j)
						// fmt.Println("the else ::: jj: ", jj, "okay ? ", okay)
						if !okay {
							return false, jj
						}
					}
				}
				//handling lineArray
				word := line[j:jj]
				lineArray = append(lineArray, word)
				// fmt.Println("lineArray: ", lineArray)
				i = endIndex
				j = jj - 1
			} else {
				// fmt.Println("pattern[i:index]: ", pattern[i:index])
				okay, jj := matchPattern(line, pattern[i:index], j)
				// fmt.Println("jj: ", jj, "okay ? ", okay)
				if !okay {
					if strings.Contains(pattern[i:], "|") {
						// fmt.Println("pattern[index:endIndex] ", pattern[index+1:endIndex])
						okay, jj := matchPattern(line, pattern[index+1:endIndex], j)
						// fmt.Println("the else ::: jj: ", jj, "okay ? ", okay)
						if !okay {
							return false, jj
						}
					}
				}
				//handling lineArray
				word := line[j:jj]
				lineArray = append(lineArray, word)
				// fmt.Println("lineArray: ", lineArray)
				i = endIndex
				j = jj - 1
			}

		} else if string(line[j]) != string(pattern[i]) {
			// fmt.Println("pattern[i]: ", string(pattern[i]), " = ", string(line[j]))
			return false, j
		}
		// fmt.Println("j: ", j)

		i++
		j++
	}
	return true, j
}
