// main.go
package main

import (
	"fmt"
	"net"

//	"time"
)

func isTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}

func main() {

	// open broadcast UDP port
	fmt.Println("net.DialUDP")
	BROADCAST_IPv4 := net.IPv4(255, 255, 255, 255)
    
	udpConn, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   BROADCAST_IPv4,
		Port: 4444,
	})
	if err != nil {
		fmt.Println("error open udp socket: ", err)
		return
	}

//~ 	fmt.Println("udpConn.LocalAddr()", udpConn.LocalAddr())
//~ 	udpConn.Close()
//~ 	fmt.Println("udpConn.LocalAddr()", udpConn.LocalAddr())
	
	var connected bool = false

	for !connected {
		// send 'marco'
		marcoMsg := "marco|testMarcoPolo"
		nbBytes, err := udpConn.Write([]byte(marcoMsg))
		if err != nil {
			fmt.Println("error sending marco ", err)
			return
		}
		fmt.Printf("wrote %d bytes\n", nbBytes)
		udpConn.Close()

		// read back answer, try 10x 100ms (total 1s)

		// need to open new udp conn on auto-allocated port to read back answer
		udpConnRead, err := net.ListenUDP("udp4", udpConn.LocalAddr().(*net.UDPAddr))
		if err != nil {
			fmt.Println("error open udp read: ", err)
			return
		}

		data := make([]byte, 1024)

		for i := 0; !connected && i < 10; i++ {
			//deadline := time.Now()
			//deadline = deadline.Add(time.Millisecond * 100)
			//udpConnRead.SetReadDeadline(deadline)

			//nbBytes, udpAddr, err := udpConnRead.ReadFromUDP(data)
			//nbBytes, _, err := udpConnRead.ReadFrom(data)
			nbBytes, err := udpConnRead.Read(data)

			if err != nil {
				if !isTimeout(err) {
					fmt.Println("error reading back answer ", err)
					return
				}
				//fmt.Println("timeout reading back answer")
			} else {
				//fmt.Printf("read %d bytes '%s' from %s\n", nbBytes, data[:nbBytes], udpAddr.String())
				fmt.Printf("read %d bytes '%s'\n", nbBytes, data[:nbBytes])
				connected = true
			}
		}
	}
}
