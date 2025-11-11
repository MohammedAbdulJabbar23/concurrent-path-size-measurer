package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

func main() {
	var filePath string
	fmt.Println("Enter a path: ")
	fmt.Scanln(&filePath)

	size := GetPathSize(filePath)
	fmt.Printf("Total size: %.2f MB\n", float64(size)/(1024*1024))
}

func GetPathSize(path string) int64 {
	var totalSize atomic.Int64
	var wg sync.WaitGroup
	
	sem := make(chan struct{}, 100)

	var walk func(string)
	walk = func(currentPath string) {
		defer wg.Done()

		entries, err := os.ReadDir(currentPath)
		if err != nil {
			return
		}

		for _, entry := range entries {
			fullPath := filepath.Join(currentPath, entry.Name())
			
			if entry.IsDir() {
				wg.Add(1)
				sem <- struct{}{}
				go func(p string) {
					defer func() { <-sem }()
					walk(p)
				}(fullPath)
			} else {
				if info, err := entry.Info(); err == nil {
					totalSize.Add(info.Size())
				}
			}
		}
	}

	wg.Add(1)
	sem <- struct{}{}
	go func() {
		defer func() { <-sem }()
		walk(path)
	}()

	wg.Wait()
	return totalSize.Load()
}
