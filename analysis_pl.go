package main

import (
	"os"
	"fmt"
	"sync"
	"strings"
	"net/http"
	"path/filepath"
	"github.com/grokify/html-strip-tags-go"
	"github.com/shurcooL/httpfs/vfsutil"
)

var an_dir = "/"

var an_public_ctn = make(chan string)
var words_map = make(map[string]int)
var an_wg sync.WaitGroup

func anaylsiser(private_ctn chan string, an_dir string, word string) {
	an_wg.Add(1)
	fs := assets
	go an_walk_paths_recursive(private_ctn, fs, an_dir, word)
	an_wg.Wait()
	close(an_public_ctn)
}

func an_walk_paths_recursive(private_ctn chan string, fs http.FileSystem, an_dir string, word string) (err error) {
	defer an_wg.Done()
	visit := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != an_dir {
			an_wg.Add(1)
			go an_walk_paths_recursive(private_ctn, fs, path, word)
			return filepath.SkipDir
		}
		if info == nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		an_wg.Add(1)
		go an_parse_file(private_ctn, fs, path, word)
		return err
	}
	vfsutil.Walk(fs, an_dir, visit)
	return err
}

func an_parse_file(private_ctn chan string, fs http.FileSystem, path string, word string) error {
	defer an_wg.Done()
	data, err := vfsutil.ReadFile(fs, path)
	if err != nil {
		return err
	}
	switch filepath.Ext(path) {
	case ".htm":
		println("-------read .htm")
		an_wg.Add(1)
		an_parse_zcj(private_ctn, path, data, word)
		return nil
	}
	return nil
}

func an_parse_zcj(private_ctn chan string, path string, data []byte, word string) {
	defer an_wg.Done()
	content_str := string(data)
	words_string := strip.StripTags(content_str)
	an_public_ctn <- words_string
}

func count_words() map[string]int {
	long_str := ""
	go func() {
		for {
			str, ok := <-an_public_ctn
			long_str += str
			if !ok {
				break
			}

		}
	}()
	anaylsiser(an_public_ctn, an_dir, "ok")
	words := strings.Fields(long_str)
	for _, word := range words {
	        word = strings.Replace(word, ",", "", -1)
	        word = strings.Replace(word, ".", "", -1)
	        word = strings.Replace(word, "-", "", -1)
	        word = strings.Replace(word, "–", "", -1)
	        word = strings.Replace(word, "“", "", -1)
	        word = strings.Replace(word, "”", "", -1)
	        word = strings.Replace(word, "…", "", -1)
		words_map[word] = words_map[word] + 1

	}
	long_str = ""
	fmt.Println("finish count words")
	return words_map

}
