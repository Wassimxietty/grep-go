package main

import (
	// Uncomment this to pass the first stage
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
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

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}
	var ok bool

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	fmt.Println(pattern)
	wordPattern := strings.Trim(pattern, "\\d")
	wordPattern = strings.Trim(wordPattern, "\\w")
	subslice := []byte(wordPattern)
	fmt.Println(subslice)
	fmt.Println("IF IT CONTAINS THE WORD PATTERN: ", bytes.Contains(line, subslice))
	ok1 := bytes.ContainsAny(line, "0123456789")
	fmt.Println("OK1 (IF IT ENTERS THE OK1 CONDITION) :", ok1)

	// \d apple
	// sally has 1 orange

	if pattern == "\\d" {
		ok1 := bytes.ContainsAny(line, "0123456789")
		if ok1 {
			wordPattern := strings.Trim(pattern, "\\d")
			wordPattern = strings.Trim(wordPattern, "\\w")
			subslice := []byte(wordPattern)
			ok = bytes.Contains(line, subslice)
		}
	} else if pattern == "\\w" {
		ok1 := bytes.ContainsAny(line, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_")
		if ok1 {
			wordPattern := strings.Trim(pattern, "\\d")
			wordPattern = strings.Trim(wordPattern, "\\w")
			subslice := []byte(wordPattern)
			ok = bytes.Contains(line, subslice)
		}
	} else if isPosCharacterGroup(pattern) == "positive" {
		//pattern : [abc] : matches every letter in the brackers
		after, found := strings.CutPrefix(pattern, "[")
		if found {
			after1, found1 := strings.CutSuffix(after, "]")
			if found1 {
				ok = bytes.ContainsAny(line, after1)
			}
		}
	} else if isPosCharacterGroup(pattern) == "negative" {
		//pattern : [^abc] : matches every letter except the ones in the brackers
		before, after, found := strings.Cut(pattern, "[^")
		if found && before == "" {
			after1, found1 := strings.CutSuffix(after, "]")
			if found1 {
				ok = !bytes.ContainsAny(line, after1)
			}
		}
	}
	fmt.Println("LAST OK :", ok, nil)
	return ok, nil
}

func isPosCharacterGroup(pattern string) string {
	// Check if the pattern length is at least 2 (to contain at least "[]")
	if len(pattern) < 2 {
		return "none"
	}
	if pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
		if pattern[1] == '^' {
			return "negative"
		} else {
			return "positive"
		}
	}
	return "none"

}

//TEST IN GIT COPY AND PASTE
// git add .
// git commit --allow-empty -m 'pass 1st stage'
// git push origin master
