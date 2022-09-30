package conn

//
//import (
//	"fmt"
//	"net"
//)
//
//var (
//	Req_REGISTER  byte = 1 // 1 --- c register cid
//	Res_REGISTER  byte = 2 // 2 --- s response
//	Req_HEARTBEAT byte = 3 // 3 --- s send heartbeat req
//	Res_HEARTBEAT byte = 4 // 4 --- c send heartbeat res
//	Req           byte = 5 // 5 --- cs send data
//	Res           byte = 6 // 6 --- cs send ack
//)
//
//var Dch chan bool
//var Rch chan []byte
//var Wch chan []byte
//
//func main() {
//	Dch = make(chan bool)
//	Rch = make(chan []byte)
//	Wch = make(chan []byte)
//	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:6666")
//	conn, err := net.DialTCP("tcp", nil, addr)
//	//	conn,err := net.Dial("tcp","127.0.0.1:6666")
//	if err != nil {
//		fmt.Println("连接服务端失败:", err.Error())
//		return
//	}
//	fmt.Println("已连接服务器")
//	defer conn.Close()
//	go Handler(conn)
//	select {
//	case <-Dch:
//		fmt.Println("关闭连接")
//	}
//}
//
//func Handler(conn *net.TCPConn) {
//	// 直到register ok
//	data := make([]byte, 128)
//	for {
//		conn.Write([]byte{Req_REGISTER, '2'})
//		conn.Read(data)
//		//		fmt.Println(String(data))
//		if data[0] == Res_REGISTER {
//			break
//		}
//	}
//	//	fmt.Println("i'm register")
//	go RHandler(conn)
//	go WHandler(conn)
//	go Work()
//}
//
//func RHandler(conn *net.TCPConn) {
//
//	for {
//		// 心跳包,回复ack
//		data := make([]byte, 128)
//		i, _ := conn.Read(data)
//		if i == 0 {
//			Dch <- true
//			return
//		}
//		if data[0] == Req_HEARTBEAT {
//			fmt.Println("recv ht pack")
//			conn.Write([]byte{Res_REGISTER, 'h'})
//			fmt.Println("send ht pack ack")
//		} else if data[0] == Req {
//			fmt.Println("recv data pack")
//			fmt.Printf("%v\n", string(data[2:]))
//			Rch <- data[2:]
//			conn.Write([]byte{Res, '#'})
//		}
//	}
//}
//
//func WHandler(conn net.Conn) {
//	for {
//		select {
//		case msg := <-Wch:
//			fmt.Println((msg[0]))
//			fmt.Println("send data after: " + string(msg[1:]))
//			conn.Write(msg)
//		}
//	}
//
//}
//
//func Work() {
//	for {
//		select {
//		case msg := <-Rch:
//			fmt.Println("work recv " + string(msg))
//			Wch <- []byte{Req, 'x', 'x'}
//		}
//	}
//}
