// Implementation of a KeyValueServer.

package p0

import (
	"bufio"
	"bytes"
	// "fmt"
	"io"
	"net"
	"strconv"
)

const MAX_MESSAGE_QUEUE_LENGTH = 500

// Stores a connection and corresponding message queue.
type client struct {
	connection       net.Conn
	messageQueue     chan []byte
	quitSignal_Read  chan int
	quitSignal_Write chan int
}

// Used to specify DBRequests
type db struct {
	isGet bool
	key   string
	value []byte
}

// Implements KeyValueServer.
//NOTE 接口实现的问题  ==  多态
type keyValueServer struct {
	listener          net.Listener
	currentClients    []*client
	newMessage        chan []byte
	newConnection     chan net.Conn			//channel用来保证线程安全更改cur_clients
	deadClient        chan *client
	dbQuery           chan *db
	dbResponse        chan *db
	countClients      chan int
	clientCount       chan int
	quitSignal_Main   chan int
	quitSignal_Accept chan int
}

// Initializes a new KeyValueServer.
func New() KeyValueServer {
	return &keyValueServer{
		nil,
		nil,
		make(chan []byte),
		make(chan net.Conn),
		make(chan *client),
		make(chan *db),
		make(chan *db),
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int)}
}

// Implementation of Start for keyValueServer.
func (kvs *keyValueServer) Start(port int) error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	kvs.listener = ln
	init_db()

	go runServer(kvs)
	go acceptRoutine(kvs)

	return nil
}

// Implementation of Close for keyValueServer.
func (kvs *keyValueServer) Close() {
	kvs.listener.Close()
	kvs.quitSignal_Main <- 0
	kvs.quitSignal_Accept <- 0
}

// Implementation of Count.
func (kvs *keyValueServer) Count() int {
	kvs.countClients <- 0
	return <-kvs.clientCount
}

// Main server routine.
func runServer(kvs *keyValueServer) {
	// defer fmt.Println("\"runServer\" ended.")

	for {
		select {
		// Send the message to each client's queue.
		//NOTE: 题目要求 get_response返回给所有的connected_client
		case newMessage := <-kvs.newMessage:
			for _, c := range kvs.currentClients {
				// If the queue is full, drop the oldest message.
				if len(c.messageQueue) == MAX_MESSAGE_QUEUE_LENGTH {
					<-c.messageQueue
				}
				c.messageQueue <- newMessage
			}
			// Add a new client to the client list.
		case newConnection := <-kvs.newConnection:
			c := &client{
				newConnection,
				make(chan []byte, MAX_MESSAGE_QUEUE_LENGTH),
				make(chan int),
				make(chan int)}
			kvs.currentClients = append(kvs.currentClients, c)
			go readRoutine(kvs, c)
			go writeRoutine(c)

			// Remove the dead client.
		case deadClient := <-kvs.deadClient:
			for i, c := range kvs.currentClients {
				if c == deadClient {
					kvs.currentClients =
						append(kvs.currentClients[:i], kvs.currentClients[i+1:]...)
					break
				}
			}

			// Run a query on the DB
		case request := <-kvs.dbQuery:
			// response required for GET query
			if request.isGet {
				v := get(request.key)
				kvs.dbResponse <- &db{
					value: v,
				}
			} else {
				put(request.key, request.value)
			}

			// Get the number of clients.
		case <-kvs.countClients:
			kvs.clientCount <- len(kvs.currentClients)

			// End each client routine.
		case <-kvs.quitSignal_Main:
			for _, c := range kvs.currentClients {
				c.connection.Close()
				c.quitSignal_Write <- 0
				c.quitSignal_Read <- 0
			}
			return
		} 
	}
}

// One running instance; accepts new clients and sends them to the server.
func acceptRoutine(kvs *keyValueServer) {
	// defer fmt.Println("\"acceptRoutine\" ended.")

	for {
		select {
		case <-kvs.quitSignal_Accept:
			return
		default:
			//NOTE: server = Listen + Accept  client=Dial获取conn
			conn, err := kvs.listener.Accept()
			if err == nil {
				kvs.newConnection <- conn
			}
		}
	}
}

// One running instance for each client; reads in
// new  messages and sends them to the server.
func readRoutine(kvs *keyValueServer, c *client) {
	// defer fmt.Println("\"readRoutine\" ended.")

	clientReader := bufio.NewReader(c.connection)

	// Read in messages.
	for {
		select {
		case <-c.quitSignal_Read:
			return
		default:
			message, err := clientReader.ReadBytes('\n')
			//NOTE: 这里使用io.EOF 判断tcp的结束情况   思考为什么接收到EOF == tcp的FIN信号
			//NOTE:  dead_client read_routine会自动结束吗？==> err!=nil时退出
			//NOTE:  dead_client write_routine会自动结束吗？ ==>等待 GC？？
			//NOTE:  特别注意TCP链接是全双工
			if err == io.EOF {
				kvs.deadClient <- c
			} else if err != nil {
				return
			} else {
				tokens := bytes.Split(message, []byte(","))
				if string(tokens[0]) == "put" {
					key := string(tokens[1][:])

					// do a "put" query
					kvs.dbQuery <- &db{
						isGet: false,
						key:   key,
						value: tokens[2],
					}
				} else {
					// remove trailing \n from get,key\n request
					keyBin := tokens[1][:len(tokens[1])-1]
					key := string(keyBin[:])

					// do a "get" query
					kvs.dbQuery <- &db{
						isGet: true,
						key:   key,
					}

					// NOTE:  这里channel dbresponse 阻塞/同步 作用
					response := <-kvs.dbResponse
					// NOTE: keyBin+','+response.value
					// NOTE: append slice类似于byte流的相加？？？
					kvs.newMessage <- append(append(keyBin, ","...), response.value...)
				}
			}
		}
	}
}

// One running instance for each client; writes messages
// from the message queue to the client.
func writeRoutine(c *client) {
	// defer fmt.Println("\"writeRoutine\" ended.")

	for {
		select {
		case <-c.quitSignal_Write:
			return
		case message := <-c.messageQueue:
			c.connection.Write(message)
		}
	}
}
