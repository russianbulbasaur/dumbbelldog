package client
import "dumbbelldog/internals/response"
import "dumbbelldog/internals/request"
import "strconv"
import "strings"
import "fmt"
func (client *Client) parseRequest(requestBytes []byte) []byte {
  request := request.ParseRequest(requestBytes)
  target := request.RLine.Target
  fmt.Println(target)
  target = strings.Split(target,"/")[1]
  if f,ok := client.serverPaths[target];ok {
    return f(request)
  }else {
    return client.serverPaths["not_found"](request)
  }
}


func notFound() []byte{
    fancy := []byte("<html>Not found 404</html>")
    headers := make(map[string]string,0)
    headers["Content-Length"] = strconv.Itoa(len(fancy))
    headers["Content-Type"] = "text/html"
    return response.NewResponse(404,"Not found",fancy,headers).Marshal()
}
