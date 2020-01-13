package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/chenjiandongx/go-echarts/charts"
)

var todos = map[string]int{}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Println("Needs two program arguments - usage: go run johnny.go /path/to/src .go")
		return
	}

	start := args[0]
	ext := args[1]

	err := filepath.Walk(start,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, ext) {
				err = findTodos(path)
				if err != nil {
					fmt.Println("File reading error", err)
					return err
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	buildChart()
}

func buildChart() {
	names, nums := extractNames()
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "TODOs"})
	bar.AddXAxis(names)
	bar.AddYAxis("Hover for names", nums)

	f, err := os.Create("graph.html")
	if err != nil {
		log.Println(err)
	}

	bar.Render(f)
}

func extractNames() ([]string, []int) {
	names := make([]string, 0, len(todos))
	nums := make([]int, 0, len(todos))
  for name, num := range todos {
    names = append(names, name)
		nums = append(nums, num)
  }

	return names, nums
}

func findTodos(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var reg = regexp.MustCompile(`(?m)^.*\/\/.*TODO(?:.*\\\s*\n)*.*`)

	s := bufio.NewScanner(file)
	for s.Scan() {
		todo := reg.FindString(s.Text())
		if todo != "" {
			dev := findDev(todo)
			if dev == "" {
				dev = "unclaimed"
			}
			todos[dev] += 1
		}
	}

	return s.Err()
}

func findDev(s string) string {
    i := strings.Index(s, "TODO(")
		pad := 5
		if i == -1 {
			i = strings.Index(s, "TODO (")
			pad = 6
		}
    if i >= 0 {
        j := strings.Index(s[i:], ")")
        if j >= 0 {
            return s[i+pad : j+i]
        }
    }
    return ""
}
