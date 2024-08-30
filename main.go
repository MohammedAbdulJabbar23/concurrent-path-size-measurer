package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)


func main() {
	var wg sync.WaitGroup;
	var filePath string;
	ch := make(chan int64,4);

	fmt.Println("Enter a path: ");
	fmt.Scanln(&filePath);

	wg.Add(1);
	
	GetPathSize(filePath,ch,&wg);
	
	go func() {
		wg.Wait();
		close(ch);
	}()
	
	var size int64 = 0;
	
	for val := range ch {
		size += val;
	}

	fmt.Printf("Total size: %.2f MB\n", float64(size)/(1024*1024));
}



func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path);
	if err != nil {
		return false, err;
	}
	return fileInfo.IsDir(),err;
}

func GetPathSize(path string, ch chan int64, wg *sync.WaitGroup) {
	defer wg.Done();
	isDir, _ := IsDirectory(path);  

	if  isDir {
		entries, err := os.ReadDir(path);
		if err != nil {
			log.Fatal(err);
		}
		for _, e := range entries {
			wg.Add(1);
			go GetPathSize(path+"/"+e.Name(), ch,wg);
		}
	}else {
		fi, err := os.Stat(path);
		if err != nil {
			fmt.Println("No perrmissions to", path);
			return;
		}
		ch <- fi.Size();
	}
}