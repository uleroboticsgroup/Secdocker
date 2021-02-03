package tcpintercept

import (
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var connectionCount uint
var websocketConn *websocket.Conn
var websocketMutex *sync.Mutex
var tlsConfig *tls.Config

func ServeConnection() {
	var tcpport string
	//var tlsport string

	//var x509 string
	//var key string

	var showConnectionAttempts bool

	flag.StringVar(&tcpport, "tcp", "8999", "Listening port for non-TLS connections.")
	//flag.StringVar(&tlsport, "tls", "6443", "Listening port for TLS connections.")
	//flag.StringVar(&x509, "x509", "./certificate/trudy.cer", "Path to x509 certificate that will be presented for TLS connection.")
	//flag.StringVar(&key, "key", "./certificate/trudy.key", "Path to the corresponding private key for the specified x509 certificate")
	flag.BoolVar(&showConnectionAttempts, "show", true, "Show connection open and close messages")

	flag.Parse()

	tcpport = ":" + tcpport
	//tlsport = ":" + tlsport
	setup(tcpport, "", "", "", showConnectionAttempts)
}

func setup(tcpport, tlsport, x509, key string, show bool) {

	//Setup non-TLS TCP listener!
	tcpAddr, err := net.ResolveTCPAddr("tcp", tcpport)
	if err != nil {
		log.Printf("There appears to be an error with the TCP port you specified. See error below.\n%v\n", err.Error())
		return
	}
	tcpListener := new(TCPListener)
	tcpListener.Listen("tcp", tcpAddr, &tls.Config{})

	log.Printf("[INFO] Listening for all TCP connections on port %s\n", tcpport)

	go websocketHandler()
	connectionDispatcher(tcpListener, "TCP", show)

}

func connectionDispatcher(listener TrudyListener, name string, show bool) {
	defer listener.Close()
	for {
		fd, conn, err := listener.Accept()
		if err != nil {
			continue
		}

		p := new(TrudyPipe)
		if name == "TLS" {
			err = p.New(connectionCount, fd, conn, true)
		} else {
			err = p.New(connectionCount, fd, conn, false)
		}

		if err != nil {
			log.Println("[ERR] Error creating new pipe.")
			continue
		}
		if show {
			log.Printf("[INFO] ( %v ) %v Connection accepted!\n", connectionCount, name)
		}
		go clientHandler(p, show)
		go serverHandler(p)
		connectionCount++
	}
}

func errHandler(err error) {
	if err != nil {
		panic(err)
	}
}

//clientHandler manages data that is sent from the client to the server.
func clientHandler(pipe Pipe, show bool) {
	if show {
		defer log.Printf("[INFO] ( %v ) Closing TCP connection.\n", pipe.Id())
	}

	buffer := make([]byte, 65535)

	for {
		bytesRead, clientReadErr := pipe.ReadFromClient(buffer)

		if clientReadErr != io.EOF && clientReadErr != nil {
			break
		}

		if clientReadErr != io.EOF && bytesRead == 0 {
			continue
		}

		data := Data{FromClient: true,
			Bytes:      buffer[:bytesRead],
			TLSConfig:  tlsConfig,
			ServerAddr: pipe.ServerInfo(),
			ClientAddr: pipe.ClientInfo()}

		data.Deserialize()

		if data.Drop() {
			continue
		}

		if data.DoMangle() {
			data.Mangle()
			bytesRead = len(data.Bytes)
		}

		if data.Forbbiden {
			fmt.Println("DROPPED")
			bytesRead = len(data.Bytes)

			_, clientWriteErr := pipe.WriteToClient(data.Bytes[:bytesRead])
			if clientWriteErr != nil || clientReadErr == io.EOF {
				break
			}
			continue
		}

		if data.DoIntercept() {
			if websocketConn == nil {
				log.Printf("[ERR] Websocket Connection has not been setup yet! Cannot intercept.")
				continue
			}
			websocketMutex.Lock()
			bs := fmt.Sprintf("% x", data.Bytes)
			if err := websocketConn.WriteMessage(websocket.TextMessage, []byte(bs)); err != nil {
				log.Printf("[ERR] Failed to write to websocket: %v\n", err)
				websocketMutex.Unlock()
				continue
			}
			_, moddedBytes, err := websocketConn.ReadMessage()
			websocketMutex.Unlock()
			if err != nil {
				log.Printf("[ERR] Failed to read from websocket: %v\n", err)
				continue
			}
			str := string(moddedBytes)
			str = strings.Replace(str, " ", "", -1)
			moddedBytes, err = hex.DecodeString(str)
			if err != nil {
				log.Printf("[ERR] Failed to decode hexedited data.")
				continue
			}
			data.Bytes = moddedBytes
			bytesRead = len(moddedBytes)
		}

		if data.DoPrint() {
			log.Printf("%v -> %v\n%v\n", data.ClientAddr.String(), data.ServerAddr.String(), data.PrettyPrint())
		}

		data.Serialize()

		data.BeforeWriteToServer(pipe)
		bytesRead = len(data.Bytes)
		_, serverWriteErr := pipe.WriteToServer(data.Bytes[:bytesRead])
		if serverWriteErr != nil || clientReadErr == io.EOF {
			break
		}

		data.AfterWriteToServer(pipe)
	}
}

//serverHandler manages data that is sent from the server to the client.
func serverHandler(pipe Pipe) {
	buffer := make([]byte, 65535)

	defer pipe.Close()

	for {
		bytesRead, serverReadErr := pipe.ReadFromServer(buffer)

		if serverReadErr != io.EOF && serverReadErr != nil {
			break
		}

		if serverReadErr != io.EOF && bytesRead == 0 {
			continue
		}

		data := Data{FromClient: false,
			Bytes:      buffer[:bytesRead],
			TLSConfig:  tlsConfig,
			ClientAddr: pipe.ClientInfo(),
			ServerAddr: pipe.ServerInfo()}

		data.Deserialize()

		if data.Drop() {
			continue
		}

		if data.DoMangle() {
			//data.Mangle()
			bytesRead = len(data.Bytes)
		}

		if data.DoIntercept() {
			if websocketConn == nil {
				log.Printf("[ERR] Websocket Connection has not been setup yet! Cannot intercept.")
				continue
			}
			websocketMutex.Lock()
			bs := fmt.Sprintf("% x", data.Bytes)
			if err := websocketConn.WriteMessage(websocket.TextMessage, []byte(bs)); err != nil {
				log.Printf("[ERR] Failed to write to websocket: %v\n", err)
				websocketMutex.Unlock()
				continue
			}
			_, moddedBytes, err := websocketConn.ReadMessage()
			websocketMutex.Unlock()
			if err != nil {
				log.Printf("[ERR] Failed to read from websocket: %v\n", err)
				continue
			}
			str := string(moddedBytes)
			str = strings.Replace(str, " ", "", -1)
			moddedBytes, err = hex.DecodeString(str)
			if err != nil {
				log.Printf("[ERR] Failed to decode hexedited data.")
				continue
			}
			data.Bytes = moddedBytes
			bytesRead = len(moddedBytes)
		}

		if data.DoPrint() {
			log.Printf("%v -> %v\n%v\n", data.ServerAddr.String(), data.ClientAddr.String(), data.PrettyPrint())
		}

		data.Serialize()

		data.BeforeWriteToClient(pipe)
		bytesRead = len(data.Bytes)

		_, clientWriteErr := pipe.WriteToClient(data.Bytes[:bytesRead])
		if clientWriteErr != nil || serverReadErr == io.EOF {
			break
		}

		data.AfterWriteToClient(pipe)
	}
}

func websocketHandler() {
	websocketMutex = &sync.Mutex{}
	upgrader := websocket.Upgrader{ReadBufferSize: 65535, WriteBufferSize: 65535}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, editor)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var err error
		websocketConn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[ERR] Could not upgrade websocket connection.")
			return
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

const editor string = `<!-- this wonderful page was found here: https://github.com/xem/hex -->
<body onload='
// Reset the textarea value
m.value="00";
// Init the top cell content
for(i=0;i<16;i++)
  t.innerHTML+=(0+i.toString(16)).slice(-2)+" ";
'>
<!-- TRUDY SPECIFIC CODE ADDED FOR THIS PROJECT -->
<h1> ~ Trudy Intercept ~ </h1>
<script>
    var url = window.location.href
    var arr = url.split("/");
    var ws_url = "ws://" + arr[2] + "/ws"
    var socket = new WebSocket(ws_url)
    socket.onmessage = function (event) {
	document.getElementById('m').value = event.data
	document.getElementById('m').oninput()
	document.getElementById('send').disabled = false
    }
    var sender = function() {
        socket.send(document.getElementById('m').value)
	document.getElementById('send').disabled = true
        document.getElementById('m').value = "00"
        document.getElementById('m').oninput()
    }
</script>
<button onclick="sender()" id='send' disabled=true>send</button>
<!-- END TRUDY SPECIFIC CODE -->
</body>
<table border><td><pre><td id=t><tr><td id=l width=80>00000000<td><textarea spellcheck=false id=m oninput='
// On input, store the length of clean hex before the textarea caret in b
b=value
.substr(0,selectionStart)
.replace(/[^0-9A-F]/ig,"")
.replace(/(..)/g,"$1 ")
.length;
// Clean the textarea value
value=value
.replace(/[^0-9A-F]/ig,"")
.replace(/(..)/g,"$1 ")
.replace(/ $/,"")
.toUpperCase();
// Set the height of the textarea according to its length
style.height=(1.5+value.length/47)+"em";
// Reset h
h="";
// Loop on textarea lines
for(i=0;i<value.length/48;i++)
  
  // Add line number to h
  h+=(1E7+(16*i).toString(16)).slice(-8)+" ";
// Write h on the left column
l.innerHTML=h;
// Reset h
h="";
// Loop on the hex values
for(i=0;i<value.length;i+=3)
  
  // Convert them in numbers
  c=parseInt(value.substr(i,2),16),
  
  // Convert in chars (if the charCode is in [64-126] (maybe more later)) or ".".
  h=63<c&&127>c?h+String.fromCharCode(c):h+".";
  
// Write h in the right column (with line breaks every 16 chars)
r.innerHTML=h.replace(/(.{16})/g,"$1 ");
// If the caret position is after a space or a line break, place it at the previous index so we can use backspace to erase hex code
if(value[b]==" ")
  b--;
// Put the textarea caret at the right place
setSelectionRange(b,b)'
cols=48></textarea><td width=160 id=r>.</td>
</table>
<style>
*{margin:0;padding:0;vertical-align:top;font:1em/1em courier}
#m{height:1.5em;resize:none;overflow:hidden}
#t{padding:0 2px}
#w{position:absolute;opacity:.001}
</style>
`
