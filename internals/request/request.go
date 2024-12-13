package request
import "dumbbelldog/constants"
import "bytes"
import "fmt"
type requestLine struct{
  method string
  Target string
  version string
}

func parseRequestLine(requestLineBytes []byte) requestLine {
  parts := bytes.Split(requestLineBytes,[]byte{' '})
  if len(parts) != 3  {
    fmt.Println("Requet line less than 3")
    return requestLine{}
  }
  return requestLine{
    method: string(parts[0]),
    Target: string(parts[1]),
    version: string(parts[2]),
  }
} 

type header struct{
  key string
  val string
}

func parseHeaders(headersBytes [][]byte) []header {
  var headers []header;
  for _,headerBytes := range headersBytes {
    parts := bytes.Split(headerBytes,[]byte{':',' '})
    headers = append(headers,header{key:string(parts[0]),val:string(parts[1])})
  }
  return headers;
}

type Request struct {
  RLine requestLine
  headers []header
}

func ParseRequest(request []byte) *Request {
  parts := bytes.Split(request,[]byte{constants.CR,constants.LF})
  rLine := parseRequestLine(parts[0])
  headers := parseHeaders(parts[1:len(parts)-2])
  //body := parseBody(parts[len(parts)-1])
  return &Request{
    RLine: rLine,
    headers: headers,
  }
}
