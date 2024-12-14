package request
import "dumbbelldog/constants"
import "bytes"
import "fmt"
import "io"
import "compress/gzip"
import "strings"
type requestLine struct{
  Method string
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
    Method: string(parts[0]),
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
  Body []byte
  Compressor io.WriteCloser
  CompressedBuffer *bytes.Buffer
}

func ParseRequest(request []byte) *Request {
  parts := bytes.Split(request,[]byte{constants.CR,constants.LF})
  rLine := parseRequestLine(parts[0])
  headers := parseHeaders(parts[1:len(parts)-2])
  var compressedBuffer *bytes.Buffer = bytes.NewBuffer(make([]byte,0))
  var compressor io.WriteCloser = getDefaultCompressor(compressedBuffer)
  if compression,ok := headers["Accept-Encoding"]; ok {
    encodings := strings.Split(compression.Val,",")
    for _,encoding := range encodings {
      if strings.TrimSpace(encoding) == "gzip" {
        fmt.Println("using gzip")
        compressor = gzip.NewWriter(compressedBuffer)
        //heddaders["Accept-Encoding"].Val = "gzip"
        break
      }
    }
  }
  return &Request{
    RLine: rLine,
    Headers: headers,
    Body: parts[len(parts)-1],
    Compressor: compressor,
    CompressedBuffer: compressedBuffer,
  }
}

func getDefaultCompressor(writer io.Writer) *defaultCompressor {
  return &defaultCompressor{
    dst: writer,
  }
}

type defaultCompressor struct{
  dst io.Writer
}

func (compressor *defaultCompressor) Write(p []byte) (int,error) {
  n,err := compressor.Write(p)
  return n,err
}

func (compressor *defaultCompressor) Close() error {
  return nil
}
