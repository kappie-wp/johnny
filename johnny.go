package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-echarts/go-echarts"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Println("Needs two program arguments - usage: go run johnny.go /path/to/src .go")
		return
	}

	start := args[0]
	ext := args[1]
	count := 0

	err := filepath.Walk(start,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, ext) {
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
