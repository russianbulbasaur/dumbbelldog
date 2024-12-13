package response
import "fmt"
import "bytes"
import "dumbbelldog/constants"
import "strconv"
type statusLine struct {
  version string
  statusCode string
  reason string
}

type header struct {
  key string
  val string
}

type Response struct {
  status statusLine
  headers []header
  body []byte
}

func NewResponse(statusCode int,reason string,body []byte,headers map[string]string) *Response {
  headersList := make([]header,0)
  for key,val := range headers {
    headersList = append(headersList,header{
      key:key,
      val:val,
    })
  }
  return &Response {
    status : statusLine {
      version: "HTTP/1.1",
      statusCode: strconv.Itoa(statusCode),
      reason: reason,
    },
    headers: headersList,
    body: body,
  }
}


func (response *Response) Marshal() []byte {
  var buffer *bytes.Buffer = bytes.NewBuffer(make([]byte,0))
  //Status line marshal
  buffer.WriteString(fmt.Sprintf("%s %s %s",response.status.version,response.status.statusCode,response.status.reason))
  buffer.WriteByte(constants.CR)
  buffer.WriteByte(constants.LF)
  for _,header := range response.headers {
    buffer.WriteString(fmt.Sprintf("%s: %s",header.key,header.val))
    buffer.Write([]byte{constants.CR,constants.LF})
  }
  buffer.Write([]byte{constants.CR,constants.LF})
  buffer.Write(response.body)
  return buffer.Bytes()
}




