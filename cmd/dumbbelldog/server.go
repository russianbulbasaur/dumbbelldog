package main

import (
    "dumbbelldog/internals/server"
    "dumbbelldog/internals/response"
    "dumbbelldog/internals/request"
    "strings"
    "strconv"
    "bytes"
    "os"
    "fmt"
)

func main() {
	server := server.NewServer("127.0.0.1",8000)
    server.AddPath("index",index)
    server.AddPath("echo",echo)
    server.AddPath("user-agent",userAgent)
    server.AddPath("files",files)
    server.AddPath("not_found",notFound)
    server.Run()
}


func index(req *request.Request) []byte {
  responseStruct := response.NewResponse(200,"OK",[]byte{},
   make(map[string]string,0))
  return responseStruct.Marshal() 
}

func echo(req *request.Request) []byte {
    target := req.RLine.Target
    echo := strings.TrimPrefix(target,"/echo/")
    headers := make(map[string]string,0)
    headers["Content-Type"] = "text/plain"
    headers["Content-Length"] = strconv.Itoa(len(echo))
    responseStruct := response.NewResponse(200,"OK",[]byte(echo),headers)
    return responseStruct.Marshal()
}


func userAgent(req *request.Request) []byte{
    output := []byte(req.Headers["User-Agent"].Val)
    fmt.Println(string(output))
    headers := make(map[string]string,0)
    headers["Content-Type"] = "text/plain"
    headers["Content-Length"] = strconv.Itoa(len(output))
    responseStruct := response.NewResponse(200,"OK",output,headers)
    return responseStruct.Marshal()
}


func files(req *request.Request) []byte {
    target := req.RLine.Target
    fileName := strings.TrimPrefix(target,"/files/")
    var buffer *bytes.Buffer = bytes.NewBuffer(make([]byte,0))
    file,err := os.Open(fileName)
    if err != nil {
      fmt.Println(err)
      return notFound(req)
    }
    _,err = buffer.ReadFrom(file)
    if err != nil {
      fmt.Println(err)
      return notFound(req)
    }
    headers := make(map[string]string,0)
    headers["Content-Type"] = "application/octect-stream"
    content := buffer.Bytes()
    headers["Content-Length"] = strconv.Itoa(len(content))
    responseStruct := response.NewResponse(200,"OK",content,headers)
    return responseStruct.Marshal()
}

func notFound(req *request.Request) []byte{
    fancy := []byte("<html>Not found 404</html>")
    headers := make(map[string]string,0)
    headers["Content-Length"] = strconv.Itoa(len(fancy))
    headers["Content-Type"] = "text/html"
    return response.NewResponse(404,"Not found",fancy,headers).Marshal()
}

