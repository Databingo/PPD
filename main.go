package main

import (
	"fmt"
	"sync"
	"strings"
	"net/http"
	"github.com/databingo/webview"
	"github.com/gorilla/websocket"
	"github.com/siongui/gojianfan"
	"github.com/tidwall/gjson"
)

var w webview.WebView

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(assets))
	go http.ListenAndServe(":9995", mux)

	go count_words()

	go websocket_center()
	http.HandleFunc("/search", search_Handler)

	w = webview.New(webview.Settings{
		Title:                  "Pali Parallel Dictionary",
		Width:                  1000,
		Height:                 600,
		URL:                    "http://localhost:9995/index.html",
		//URL:                    "http://localhost:8888/myapp/index.html",
		Debug:                  true,
		Resizable:              false,
		ExternalInvokeCallback: callback,
	})
	w.Run()

}

func callback(w webview.WebView, data string) {
	switch {
	case strings.HasPrefix(data, "inner_data:"):
	case (data == "register"):
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

type worker struct {
	source chan string
	conn   *websocket.Conn
	quite  chan struct{}
	stop   chan string
        word_record   chan string
}

func (w *worker) Start() {
	w.source = make(chan string)
	w.stop = make(chan string)
	w.word_record = make(chan string)
	go func() {
		var count = 0
		for {
			_, p, err := w.conn.ReadMessage()
			if err != nil {
				fmt.Println("read from  client failed, close this worker's conn \n")
				w.conn.Close()
				threadSafeMaper.Delete(w)
				//close(w.source)
				break
			}
			word := gojianfan.S2T(string(p))
			if word != "" {
				select {
				case w.word_record <- word:
					fmt.Printf("receive %s\n", word)
					if count > 0 {
						close(w.stop)
					}
					w.stop = make(chan string)
					go searcher(w.stop, w.source, dir, word)
					go analys(w.stop, w.source, word)
					count++
				}
			}
		}
	}()

	go func() {
		var word_asking string
		for {
			select {
			case msg, _ := <-w.word_record:
				word_asking = string(msg)
			case msg, _ := <-w.source:
				word := gjson.Get(string(msg), "word").String()
				if word == word_asking {
					err := w.conn.WriteMessage(websocket.TextMessage, []byte(string(msg)))
					if err != nil {
						fmt.Println("push to client failed, close this worker's conn \n")
						w.conn.Close()
						threadSafeMaper.Delete(w)
						//close(w.source)
					}
				}
			}
		}
	}()
}

type threadSafeMap struct {
	sync.Mutex
	workers map[*worker]bool
}

func (slice *threadSafeMap) Put(w *worker) {

	slice.Lock()
	defer slice.Unlock()
	slice.workers[w] = true
}
func (slice *threadSafeMap) Delete(w *worker) {
	slice.Lock()
	defer slice.Unlock()
	delete(slice.workers, w)
}

func (slice *threadSafeMap) Iter(routine func(*worker)) {
	slice.Lock()
	defer slice.Unlock()
	for worker, _ := range slice.workers {
		routine(worker)
	}
}

var threadSafeMaper = threadSafeMap{workers: make(map[*worker]bool)}

func search_Handler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	wr := &worker{conn: conn}
	wr.Start()
	threadSafeMaper.Put(wr)
}

func websocket_center() {
	wg := &sync.WaitGroup{}

	//------------push_to_clients--------------
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			msg, _ := <-public_ctn
			threadSafeMaper.Iter(func(w *worker) { w.source <- msg })
		}
	}(wg)

	http.ListenAndServe(":8300", nil)
	wg.Wait()
}
func analys(stop chan string, private_ctn chan string, word string) {
	select {
	case <-stop:
		return
	default:
		var new_ls []map[string]int
		ls := make([]map[string]int, 10000)
		for i := range ls {
			ls[i] = make(map[string]int)
		}
		for k, v := range words_map {
			if strings.Contains(k, word) {
				ls[v][k] = v
			}
		}

		for _, e := range ls {
			if len(e) != 0 {
				new_ls = append([]map[string]int{e}, new_ls...)
			}

		}
		ls = nil
		var c = 0
		for _, a := range new_ls {
			for k, v := range a {
				if c < 10 {
					s := fmt.Sprintf(`{"type":1, "word":"%s", "hit":"%s", "count":%d}`, word, strings.Replace(k, `"`, `'`, -1), v)
					private_ctn <- s
					c++

				}
			}
		}

	}
}
