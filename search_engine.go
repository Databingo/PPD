package main

import (
	"os"
	"fmt"
	"bytes"
	"strconv"
	"strings"
	"net/http"
	"path/filepath"
	"github.com/grokify/html-strip-tags-go"
	"github.com/shurcooL/httpfs/vfsutil"
	"github.com/siongui/gojianfan"
)

var dir = "/zh_sutta"
var public_ctn = make(chan string)

func searcher(stop chan string, private_ctn chan string, dir string, word string) {
	select {
	case <-stop:
	default:
		fs := assets
		go walk_paths_recursive(stop, private_ctn, fs, dir, word)
	}

}

func walk_paths_recursive(stop chan string, private_ctn chan string, fs http.FileSystem, dir string, word string) (err error) {
	select {
	case <-stop:
	default:
		visit := func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && path != dir {
				go walk_paths_recursive(stop, private_ctn, fs, path, word)
				return filepath.SkipDir
			}
			if info == nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			go parse_file(stop, private_ctn, fs, path, word)
			return err
		}
		vfsutil.Walk(fs, dir, visit)
	}
	return err
}

func parse_file(stop chan string, private_ctn chan string, fs http.FileSystem, path string, word string) error {
	select {
	case <-stop:
	default:
		data, err := vfsutil.ReadFile(fs, path)
		if err != nil {
			return err
		}
		if !bytes.Contains(data, []byte(word)) {
			return nil
		} else {
			switch filepath.Ext(path) {
			case ".htm":
				parse_zcj(stop, private_ctn, path, data, word)
			}
		}
	}
	return nil
}

func parse_zcj(stop chan string, private_ctn chan string, path string, data []byte, word string) {
	select {
	case <-stop:
	default:
		content_str := string(data)
		sen_list := strings.Split(content_str, "<br>")
		for _, n := range sen_list {
			// get raw number extracted from pali lines
			word_list := strings.Split(n, ".")
			no := strings.Join(strings.Fields(word_list[0]), "") // remove space
			if _, err := strconv.Atoi(no); err == nil {          // check if number
				no_pl := " " + no + "."
				no_cn := "<!" + no + ">"
				content_str = strings.Replace(content_str, no_pl, "^-^"+ no +".pl.", -1)
				content_str = strings.Replace(content_str, no_cn, "^-^"+ no +".cn.", -1)
			}
		}

		no := ""
		nu_para_list := strings.Split(content_str, "^-^")
		// 454.cn ...
		for _, p := range nu_para_list {
			// find hit
			if strings.Contains(strip.StripTags(p), word) {
				word_list := strings.Split(p, ".")
				// find num
				no = strings.Join(strings.Fields(word_list[0]), "") // remove space
				//check if number
				if _, err := strconv.Atoi(no); err != nil {
					//println("sentence prefix is not number, continue \n")
					continue
				}
				// find sutta
				sutta_no := strings.Split(filepath.Base(path), ".")[0] + "."

				// find para
				for _, n := range nu_para_list {
					select {
					case <-stop:
					default:
						word_list_n := strings.Split(n, ".")
						// find num
						no_n := strings.Join(strings.Fields(word_list_n[0]), "") // remove space
						if no_n == no && n != p {
							//println(no_n)
							hit_lg := word_list[1]
							hit_para := sutta_no + strings.Join(strings.Fields(strip.StripTags(p)), " ") //remove duplicate space
							//println(hit_para)
							hit_para = strings.Replace(hit_para, word, "<item class='word'>"+word+"</item>", -1)
							hit_para = strings.Replace(hit_para, `"`, `'`, -1)
							parallel_lg := word_list_n[1]
							parallel_para := sutta_no + strings.Join(strings.Fields(strip.StripTags(n)), " ") //remove duplicate space
							parallel_para = strings.Replace(parallel_para, `"`, `'`, -1)
							json := fmt.Sprintf(`{"type":0, "word":"%s", "lg":"%s", "%s":"%s", "%s":"%s"}`,
								gojianfan.T2S(word), hit_lg, hit_lg, hit_para, parallel_lg, parallel_para)
							private_ctn <- json
							//println("put2ch --" + json[:100])
							break //now just for two kind language ok
						}
					}
				}
			}
		}
	}
}
