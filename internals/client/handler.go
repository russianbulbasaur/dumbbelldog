package client
import "dumbbelldog/internals/response"
import "dumbbelldog/internals/request"
import "strings"
import "strconv"
func (client *Client) parseRequest(requestBytes []byte) []byte {
  request := request.ParseRequest(requestBytes)
  var responseStruct *response.Response;
  if request.RLine.Target == "/index" {
     responseStruct = response.NewResponse(200,"OK",[]byte{},
   make(map[string]string,0))
  }else if strings.HasPrefix(request.RLine.Target,"/echo"){
    echo := strings.TrimPrefix(request.RLine.Target,"/echo/")
    headers := make(map[string]string,0)
    headers["Content-Type"] = "text/plain"
    headers["Content-Length"] = strconv.Itoa(len(echo))
    responseStruct = response.NewResponse(200,"OK",[]byte(echo),headers)
  } else {
    responseStruct = response.NewResponse(404,"Not found",[]byte{},make(map[string]string,0))
  }
  return responseStruct.Marshal()
}

