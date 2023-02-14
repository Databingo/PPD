package main

import (
  "os"
  "fmt"
  "log"
  "net"
  "sync"
  "time"
  "runtime"
  "strings"
  "net/http"
  "io/ioutil"
  "io"
  "bufio"
  "github.com/databingo/webview"
  "github.com/gorilla/websocket"
  "github.com/siongui/gojianfan"
  "gopkg.in/igm/sockjs-go.v2/sockjs"

 )



const DefaultPolicy = "<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\"/></cross-domain-policy>\x00"
type Server struct {
   b []byte
}
func NewServer() *Server {
 s:= &Server{}
 s.b = []byte(DefaultPolicy)
 return s
}
func (s *Server) ListenAndServe() {
    addr, err := net.ResolveTCPAddr("tcp", ":843")
    if err != nil {
         println(err)
	}
    l, err := net.ListenTCP("tcp", addr)
    if err != nil {
         println(err)
	}
    for {
       conn, err := l.AcceptTCP()
       if err != nil {
             println(err)
       }
       println("receive...TCP")
       go s.handleConnection(conn)
    }
    l.Close()
}

func (s *Server) handleConnection(conn *net.TCPConn){
    defer conn.Close()
    conn.SetReadDeadline(time.Now().Add(time.Second))
    r := bufio.NewReader(conn)
    i, err := r.ReadString('\x00')
    println("read")
    println(i)
    println("---")
    if err != nil {
	println("err!= nil")
        if err != io.EOF {
	    println("err")
	    println(err)
	    log.Fatal(err)
	}
    return
    }
    println("Write(s.b)---", string(s.b))
    conn.Write(s.b)
    println("fsp: sent policy flie to", conn.RemoteAddr().String())
   }





func main() {
    s := NewServer()
    go s.ListenAndServe()
   }
