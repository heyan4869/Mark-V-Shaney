// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Write by Yan He from course 02-601

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"strconv"
	"math/rand"
	"time"
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a map of suffixes and their frequency
type Chain struct {
	chain     map[string]map[string]int
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string]map[string]int), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {

		br := bufio.NewReader(r)
		p := make(Prefix, c.prefixLen)
		for {
			var s string
			if _, err := fmt.Fscan(br, &s); err != nil {
				break
			}
			key := p.String()
			_, ok := c.chain[key]
			if !ok {
				var suffix map[string]int = make(map[string]int)
				suffix[s] = 1
				c.chain[key] = suffix

			} else {
				c.chain[key][s]++
			}
			p.Shift(s)
		}
	}

// Write the frequency table into a file
func (c *Chain) WirteToFile(outfileName io.Writer) {

	// Write the length of the prefix into the file
	fmt.Fprintln(outfileName, c.prefixLen)

	// Calculate the number of empty space in the prefix
	for prefix, _ := range c.chain {
		var prefixstring string
		var empty string = "\"\""
		var checkempty []string = make([]string, c.prefixLen)
		for n := 0; n < c.prefixLen; n++ {
			checkempty[n] = ""
		}
		checkemptystring := strings.Join(checkempty, " ")

		// Replace the empty space with "" in the prefix
		if prefix == checkemptystring {
			var emptyitem []string = make([]string, c.prefixLen)
			for m := 0; m < c.prefixLen; m++ {
				emptyitem[m] = empty
			}
			prefixstring = strings.Join(emptyitem, " ")
		} else {
			var items []string = strings.Split(prefix, " ")
			var length int = len(items)
			var num int
			for i := 0; i < length; i++{
				if items[i] != "" {
					num = i
					break
				}
			}
			for j := 0; j < num; j++ {
				items[j] = empty
			}
			prefixstring = strings.Join(items, " ")
		}

		// Write the prefix into the file
		fmt.Fprint(outfileName, prefixstring, " ")		

		// Write the suffix and its frequency into the file
    	for suffix, freq := range c.chain[prefix] {
			fmt.Fprint(outfileName, suffix, " ", freq, " ")
		}

		// Change to another line after writing each line in the frequency table
		fmt.Fprintln(outfileName, "")
	}
}

// Read the model file and build a new chain
func (cGen *Chain) ReadModelfile(lines []string) {

	// Get the length of prefix from the first line
	prefixLength, err := strconv.Atoi(lines[0])
	if err != nil {
        fmt.Println("Error:  Number of prefixes should be a positive number")
        return
    }

    var key string
	var length int = len(lines)
	var prefixarray []string = make([]string, prefixLength)

	for i := 1; i < length; i++ {
		var item []string = strings.Split(lines[i], " ")
		for j := 0; j < prefixLength; j++ {
			if item[j] == "\"\"" {
				item[j] = ""
			}
			prefixarray[j] = item[j]		
		}
		key = strings.Join(prefixarray, " ")
		var suffix map[string]int = make(map[string]int)
		var suffixword string		
		for k := prefixLength; k < len(item) - 1; k = k + 2 {
			suffixword = item[k]
			frequent, err := strconv.Atoi(item[k+1])
			if err != nil {
        		fmt.Println("Error:  Number of frequency should be a positive number")
        		return
    		}
    		suffix[suffixword] = frequent		
		}
		cGen.chain[key] = suffix
	}
}

// Generate returns a string of at most n words generated from Chain
func (cGen *Chain) Generate(n int) string {
	p := make(Prefix, cGen.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := cGen.chain[p.String()]
		var freqsum int
		var candidate []string
		var next string
		for suffix, freq := range choices {
			freqsum = freqsum + freq
			for j := 0; j < freq; j++ {
				candidate = append(candidate, suffix)
			}
			
			next = candidate[rand.Intn(freqsum)]
			
		}
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

// Read model file that given by the user input
func ReadModel(modelfilename string) []string {
	tablefile, err := os.Open(modelfilename)
	if err != nil {
        fmt.Println("Error: Something wrong with the file.")
     	os.Exit(3)
   	}

   	var lines []string = make([]string, 0)
   	scanner := bufio.NewScanner(tablefile)
   	for scanner.Scan() {
    	// Append the content of the file to the lines slice
       	lines = append(lines, scanner.Text())
   	}
   	if scanner.Err() != nil {
       	fmt.Println("Sorry: There was some kind of error during the file reading")
		os.Exit(3)
	}
	tablefile.Close()
	return lines
}

// Main function of this program
func main() {

	// Check if the command is "read" or "generate"
	if os.Args[1] == "read" {
		prefixLen, err := strconv.Atoi(os.Args[2])
		if err != nil || prefixLen <= 0 {
        	fmt.Println("Error:  Number of prefixes should be a positive number")
        	return
    	}
		outfileName := os.Args[3]
		input := os.Args[4:]
		var inputlen int = len(input)	
		c := NewChain(prefixLen)
		for i := 0; i < inputlen; i++ {
			var inputfilename string = input[i]
     		inputfile, err2 := os.Open(inputfilename)
    		if err2 != nil {
        		fmt.Println("Error: Something wrong with the file.")
     		}
			c.Build(inputfile)
		}
		var outputfilename string = outfileName
		outputfile, err3 := os.Create(outputfilename)
		if err3 != nil {
        		fmt.Println("Error: Something wrong with the file.")
     		}
		c.WirteToFile(outputfile)

	} else if os.Args[1] == "generate" {
		modelfile := os.Args[2]
		var lines []string = ReadModel(modelfile)
		prefixLength, err := strconv.Atoi(lines[0])
		if err != nil {
    	    fmt.Println("Error:  Number of prefixes should be a positive number")
        	return
    	}
    	rand.Seed(time.Now().UnixNano())
		cGen := NewChain(prefixLength)
		cGen.ReadModelfile(lines)
		number, err := strconv.Atoi(os.Args[3])
		if err != nil || number <= 0 {
			fmt.Println("Error:  Number of words should be a positive number")
        	return
		}
		text := cGen.Generate(number)
		fmt.Println(text)
	}
}