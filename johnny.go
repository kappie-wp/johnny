package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Walk the path, check for .go files, read the file line by line
	// If TODO is found and matched then print comment up to one complete line
	// Perhaps cater for multiline comments

	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println(args)
		log.Println("Need a single path argument")
		return
	}
	start := args[0]
	count := 0

	err := filepath.Walk(start,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if !info.IsDir() && strings.HasSuffix(path, ".go") {
				todos, err := findTodos(path)
				if err != nil {
					fmt.Println("File reading error", err)
					return err
				}

				if len(todos) > 0 {
					count += len(todos)
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Total TODO's : ", count)
}

func findTodos(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var todos []string
	var reg = regexp.MustCompile(`(?m)^.*\/\/.*TODO(?:.*\\\s*\n)*.*`)
	s := bufio.NewScanner(file)

	for s.Scan() {
		todo := reg.FindString(s.Text())
		if todo != "" {
			todos = append(todos, todo)
		}
	}

	return todos, s.Err()
}

