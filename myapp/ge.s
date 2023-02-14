package main


import (
 "github.com/shurcooL/vfsgen"
 "net/http"

)

func main(){
   fs := http.Dir("./myapp")
   err := vfsgen.Generate(fs, vfsgen.Options{})
   if err != nil {
    println(err)
   }

}
