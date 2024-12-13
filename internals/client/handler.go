package client
import "dumbbelldog/internals/response"
import "dumbbelldog/internals/request"
import "strings"
import "strconv"
import "fmt"
import "os"
import "bytes"
func (client *Client) parseRequest(requestBytes []byte) []byte {
  request := request.ParseRequest(requestBytes)
  target := request.RLine.Target
  var responseStruct *response.Response;
  if target == "/index" {
     responseStruct = response.NewResponse(200,"OK",[]byte{},
   make(map[string]string,0))
  }else if strings.HasPrefix(target,"/echo"){
    echo := strings.TrimPrefix(target,"/echo/")
    headers := make(map[string]string,0)
    headers["Content-Type"] = "text/plain"
    headers["Content-Length"] = strconv.Itoa(len(echo))
    responseStruct = response.NewResponse(200,"OK",[]byte(echo),headers)
  } else if target=="/user-agent"{
    output := []byte(request.Headers["User-Agent"].Val)
    fmt.Println(string(output))
    headers := make(map[string]string,0)
    headers["Content-Type"] = "text/plain"
    headers["Content-Length"] = strconv.Itoa(len(output))
    responseStruct = response.NewResponse(200,"OK",output,headers)
  }else if strings.HasPrefix(target,"/files"){
    fileName := strings.TrimPrefix(target,"/files/")
    var buffer *bytes.Buffer = bytes.NewBuffer(make([]byte,0))
    file,err := os.Open(fileName)
    if err != nil {
      fmt.Println(err)
      return notFound()
    }
    _,err = buffer.ReadFrom(file)
    if err != nil {
      fmt.Println(err)
      return notFound()
    }
    headers := make(map[string]string,0)
    headers["Content-Type"] = "application/octect-stream"
    content := buffer.Bytes()
    headers["Content-Length"] = strconv.Itoa(len(content))
    responseStruct = response.NewResponse(200,"OK",content,headers)
  }else {
    return notFound()
  }
  return responseStruct.Marshal()
}


func notFound() []byte{
    fancy := []byte("<html>Not found 404</html>")
    headers := make(map[string]string,0)
    headers["Content-Length"] = strconv.Itoa(len(fancy))
    headers["Content-Type"] = "text/html"
    return response.NewResponse(404,"Not found",fancy,headers).Marshal()
}
