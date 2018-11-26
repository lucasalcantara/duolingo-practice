package main

import (
	"bufio"
	"duolingo"
	"fmt"
	"os"
	"strings"
)

func main() {
	bio := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("Duolingo username> ")
		username, err := bio.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command: ", err)
			continue
		}
		username = strings.TrimSpace(username)

		fmt.Printf("language (i.e. de, pt, en)> ")
		lng, err := bio.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command: ", err)
			continue
		}
		lng = strings.TrimSpace(lng)

		d, err := duolingo.New(username, lng)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		words := removeLearnedWords(d.Words())

		fmt.Printf("translate words to (i.e. de, pt, en)> ")
		target, err := bio.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command: ", err)
			continue
		}
		target = strings.TrimSpace(target)

		translateWords, err := d.TranslateWords(lng, target, words)
		if err != nil {
			fmt.Println(err)
			continue
		}

		practice(translateWords, lng, target)
	}
}

func removeLearnedWords(words []string) []string {
	unknownWords := []string{}
	learned, err := learnedWords()
	if err != nil {
		panic(err)
	}

	for _, w := range words {
		if _, ok := learned[w]; !ok {
			unknownWords = append(unknownWords, w)
		}
	}

	return unknownWords
}

func learnedWords() (map[string]int, error) {
	file, err := os.Open("words.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	words := make(map[string]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		values := strings.Split(scanner.Text(), "=")
		if values[1] == "1" {
			words[values[0]] = 0
		}
	}
	return words, scanner.Err()
}

func practice(learnedWords map[string][]string, lng, target string) {
	bio := bufio.NewReader(os.Stdin)
loop:
	for {
		fmt.Printf("type 1 to %s->%s or 2 to %s->%s or exit > ", lng, target, target, lng)
		cmd, err := bio.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command: ", err)
			continue
		}

		switch strings.TrimSpace(cmd) {
		case "1":
			lngToTarget(learnedWords)
		case "2":
			targetToLng(learnedWords)
		case "exit":
			break loop
		default:
			fmt.Println("Unknow command.")

		}
	}
}

func lngToTarget(learnedWords map[string][]string) {
	bio := bufio.NewReader(os.Stdin)
	for w, a := range learnedWords {
		fmt.Printf("%s > ", w)
		bio.ReadString('\n')
		fmt.Println("answer(s)> ", strings.Join(a, " or "))
	}
}

func targetToLng(learnedWords map[string][]string) {
	bio := bufio.NewReader(os.Stdin)
	for w, a := range learnedWords {
		fmt.Printf("%s > ", a[0])
		bio.ReadString('\n')
		fmt.Println("answer> ", w)
	}
}
