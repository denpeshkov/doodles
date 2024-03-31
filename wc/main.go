package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	log.SetFlags(0)

	var (
		bytesFlag bool
		linesFlag bool
		charsFlag bool
		wordsFlag bool
	)

	fs := flag.NewFlagSet("wc", flag.ExitOnError)
	fs.BoolVar(&bytesFlag, "c", false, "The number of bytes in each input file is written to the standard output. This will cancel out any prior usage of the -m option.")
	fs.BoolVar(&linesFlag, "l", false, "The number of lines in each input file is written to the standard output.")
	fs.BoolVar(&charsFlag, "m", false, "The number of characters in each input file is written to the standard output. If the current locale does not support multibyte characters, this is equivalent to the -c option. This will cancel out any prior usage of the -c option.")
	fs.BoolVar(&wordsFlag, "w", false, "The number of words in each input file is written to the standard output.")
	fs.Parse(os.Args[1:])

	fname := fs.Arg(0)

	if fs.NFlag() == 0 {
		bytesFlag = true
		linesFlag = true
		wordsFlag = true
	}
	if charsFlag {
		bytesFlag = false
	}

	// if only bytesFlag is set use optimized Stat call for file
	if fs.NFlag() == 1 && bytesFlag && fname != "" {
		fi, err := os.Stat(fname)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stdout, "%d %s", fi.Size(), fname)
		return
	}

	var in io.Reader
	if fname == "" {
		in = os.Stdin
	} else {
		f, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
	}

	bytes, chars, words, lines, err := stat(bufio.NewReader(in))
	if err != nil {
		log.Fatal(err)
	}

	var out []string
	if linesFlag {
		out = append(out, strconv.Itoa(lines))
	}
	if wordsFlag {
		out = append(out, strconv.Itoa(words))
	}
	if bytesFlag {
		out = append(out, strconv.Itoa(bytes))
	}
	if charsFlag {
		out = append(out, strconv.Itoa(chars))
	}
	fmt.Printf("%s %s", strings.Join(out, " "), fname)
}

func stat(r io.RuneReader) (bytes int, chars int, words int, lines int, err error) {
	var pr rune
	for {
		ch, sz, err := r.ReadRune()
		if err == io.EOF {
			if unicode.IsLetter(pr) && unicode.IsSpace(ch) {
				words++
			}
			return bytes, chars, words, lines, nil
		}
		if err != nil {
			return 0, 0, 0, 0, err
		}

		bytes += sz
		chars++
		if !unicode.IsSpace(pr) && unicode.IsSpace(ch) {
			words++
		}
		if ch == '\n' {
			lines++
		}
		pr = ch
	}
}
