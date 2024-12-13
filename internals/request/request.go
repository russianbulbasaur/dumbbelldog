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

type Header struct{
  Key string
  Val string
}

func parseHeaders(headersBytes [][]byte) map[string]Header {
  headers := make(map[string]Header,0)
  for _,headerBytes := range headersBytes {
    parts := bytes.Split(headerBytes,[]byte{':',' '})
    headers[string(parts[0])] = Header{Key:string(parts[0]),Val:string(parts[1])}
  }
  return headers;
}

type Request struct {
  RLine requestLine
  Headers map[string]Header
}

func ParseRequest(request []byte) *Request {
  parts := bytes.Split(request,[]byte{constants.CR,constants.LF})
  rLine := parseRequestLine(parts[0])
  headers := parseHeaders(parts[1:len(parts)-2])
  //body := parseBody(parts[len(parts)-1])
  return &Request{
    RLine: rLine,
    Headers: headers,
  }
}
