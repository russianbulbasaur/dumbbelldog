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
    "bufio"
    "errors"
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
    req.Compressor.Write([]byte(echo))
    req.Compressor.Close()
    res := req.CompressedBuffer.Bytes()
    headers := make(map[string]string,0)
    headers["Content-Type"] = "text/plain"
    headers["Content-Encoding"] = "gzip"
    headers["Content-Length"] = strconv.Itoa(len(res))
    fmt.Println(res)
    responseStruct := response.NewResponse(200,"OK",res,headers)
    return responseStruct.Marshal()
}


func userAgent(req *request.Request) []byte{
    output := []byte(req.Headers["User-Agent"].Val)
    req.Compressor.Write(output)
    headers := make(map[string]string,0)
    headers["Content-Type"] = "text/plain"
    headers["Content-Length"] = strconv.Itoa(len(output))
    responseStruct := response.NewResponse(200,"OK",req.CompressedBuffer.Bytes(),headers)
    return responseStruct.Marshal()
}


func files(req *request.Request) []byte {
    target := req.RLine.Target
    fileName := strings.TrimPrefix(target,"/files/")
    var buffer *bytes.Buffer = bytes.NewBuffer(make([]byte,0))
    file,err := os.Open(fileName)
    if err != nil {
      fmt.Println(err)
      if errors.Is(err,os.ErrNotExist) && req.RLine.Method=="POST"{
        result,err := newFile(req,fileName)
        if err != nil {
          fmt.Println(err)
          return notFound(req)
        }
        return result
      }
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
    req.Compressor.Write(content)
    headers["Content-Length"] = strconv.Itoa(len(content))
    responseStruct := response.NewResponse(200,"OK",req.CompressedBuffer.Bytes(),headers)
    return responseStruct.Marshal()
}


func newFile(req *request.Request,fileName string) ([]byte,error) {
  fmt.Println(string(req.Body))
  file,err := os.OpenFile(fileName,os.O_CREATE|os.O_RDWR,0444)
  if err != nil {
    return nil,err
  }
  bufWriter := bufio.NewWriter(file)
  bufWriter.Write(req.Body)
  bufWriter.Flush()
  file.Close()
  headers := make(map[string]string,0)
  headers["Content-Type"] = "text/plain"
  headers["Content-Length"] = "0"
  responseStruct := response.NewResponse(201,"CREATED",[]byte{},headers)
  return responseStruct.Marshal(),nil
}

func notFound(req *request.Request) []byte{
    fancy := []byte("<html>Not found 404</html>")
    headers := make(map[string]string,0)
    headers["Content-Length"] = strconv.Itoa(len(fancy))
    headers["Content-Type"] = "text/html"
    return response.NewResponse(404,"Not found",fancy,headers).Marshal()
}

