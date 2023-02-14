package main

import (
	"fmt"
	"github.com/shurcooL/httpfs/vfsutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var sen_dir = "/zh_sutta"
var sen_public_ctn = make(chan interface{})
var sen_wg sync.WaitGroup
var sen_list []string
var page_m = make(map[string]string)

func sen_searcher(private_ctn chan interface{}, sen_dir string, word string) {
	sen_wg.Add(1)
	fs := assets
	go sen_walk_paths_recursive(private_ctn, fs, sen_dir, word)
	sen_wg.Wait()
	close(sen_public_ctn)
	runtime.GC()
}

func sen_walk_paths_recursive(private_ctn chan interface{}, fs http.FileSystem, sen_dir string, word string) (err error) {
	defer sen_wg.Done()
	visit := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != sen_dir {
			sen_wg.Add(1)
			go sen_walk_paths_recursive(private_ctn, fs, path, word)
			return filepath.SkipDir
		}
		if info == nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		sen_wg.Add(1)
		go sen_parse_file(private_ctn, fs, path, word)
		return err
	}
	vfsutil.Walk(fs, sen_dir, visit)
	return err
}

func sen_parse_file(private_ctn chan interface{}, fs http.FileSystem, path string, word string) error {
	defer sen_wg.Done()
	defer runtime.GC()
	data, err := vfsutil.ReadFile(fs, path)
	if err != nil {
		return err
	}
	switch filepath.Ext(path) {
	case ".htm":
		m := make(map[string]string)
		m[path] = string(data)
		println("read page: ", path)
		sen_public_ctn <- m
	}
	return nil
}

func prepare_pages() {
	sen_searcher(sen_public_ctn, sen_dir, "")
	for {
		page, ok := <-sen_public_ctn
		if page != nil {
			for k, v := range page.(map[string]string) {
				page_m[k] = v
			}
		}

		if !ok {
			println(len(sen_list))
			break
		}
	}
	fmt.Println("finish prepare pages")
}
