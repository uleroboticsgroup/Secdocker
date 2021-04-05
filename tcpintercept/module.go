package tcpintercept

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/hex"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/uleroboticsgroup/Secdocker/httpserver"
)

//Data is a thin wrapper that provides metadata that may be useful when mangling bytes on the network.
type Data struct {
	FromClient bool        //FromClient is true is the data sent is coming from the client (the device you are proxying)
	Bytes      []byte      //Bytes is a byte slice that contians the TCP data
	TLSConfig  *tls.Config //TLSConfig is a TLS server config that contains Trudy's TLS server certficiate.
	ServerAddr net.Addr    //ServerAddr is net.Addr of the server
	ClientAddr net.Addr    //ClientAddr is the net.Addr of the client (the device you are proxying)
	Forbbiden  bool
}

//DoMangle will return true if Data needs to be sent to the Mangle function.
func (input Data) DoMangle() bool {
	return true
}

//Mangle can modify/replace the Bytes values within the Data struct. This can
//be empty if no programmatic mangling needs to be done.
func (input *Data) Mangle() {
	if len(input.Bytes) > 0 {
		r := bytes.NewReader(input.Bytes)
		reader := bufio.NewReader(r)
		req, err := http.ReadRequest(reader)

		if err == nil {
			if strings.Contains(req.URL.Path, "/containers/create") {
				data := httpserver.ProcessCreateContainer(req)

				// If data is empty, it has forbidden options
				if len(data) > 0 {
					stringReader := bytes.NewReader(data)
					stringReadCloser := ioutil.NopCloser(stringReader)
					req.Body = stringReadCloser
					req.ContentLength = int64(len(data))

					buff := bytes.NewBuffer([]byte{})
					req.Write(buff)
					input.Bytes = buff.Bytes()
				} else {
					message := `Option forbidden`
					stringReader := bytes.NewReader([]byte(message))
					stringReadCloser := ioutil.NopCloser(stringReader)
					response := http.Response{}
					response.Status = "403 Forbidden"
					response.StatusCode = 403
					response.Proto = "HTTP/1.1"
					response.Body = stringReadCloser
					response.Close = true
					response.ContentLength = int64(len(message))

					buff := bytes.NewBuffer([]byte{})
					response.Write(buff)
					input.Bytes = buff.Bytes()

					input.Forbbiden = true
				}
			}
		}
	}
}

//Drop will return true if the Data needs to be dropped before going through
//the pipe.
func (input Data) Drop() bool {
	return false
}

//PrettyPrint returns the string representation of the data. This string will
//be the value that is logged to the console.
func (input Data) PrettyPrint() string {
	return hex.Dump(input.Bytes)
}

//DoPrint will return true if the PrettyPrinted version of the Data struct
//needs to be logged to the console.
func (input Data) DoPrint() bool {
	return true
}

//DoIntercept returns true if data should be sent to the Trudy interceptor.
func (input Data) DoIntercept() bool {
	return false
}

//Deserialize should replace the Data struct's Bytes with a deserialized bytes.
//For example, unpacking a HTTP/2 frame would be deserialization.
func (input *Data) Deserialize() {

}

//Serialize should replace the Data struct's Bytes with the serialized form of
//the bytes. The serialized bytes will be sent over the wire.
func (input *Data) Serialize() {

}

//BeforeWriteToClient is a function that will be called before data is sent to
//a client.
func (input *Data) BeforeWriteToClient(p Pipe) {

}

//AfterWriteToClient is a function that will be called after data is sent to
//a client.
func (input *Data) AfterWriteToClient(p Pipe) {

}

//BeforeWriteToServer is a function that will be called before data is sent to
//a server.
func (input *Data) BeforeWriteToServer(p Pipe) {

}

//AfterWriteToServer is a function that will be called after data is sent to
//a server.
func (input *Data) AfterWriteToServer(p Pipe) {

}
