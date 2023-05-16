package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)
func main() {
	//ノードの作成
	node, err := libp2p.New()
	if  err != nil {
		panic(err)
	}
	defer node.Close()
	
	

	//作成したノードのpeerIDとアドレスをpeerInfoに入れる
	peerInfo := peer.AddrInfo{
		ID: node.ID(),
		Addrs: node.Addrs(),
	}

	//peerInfoを元に、p2pアドレスを得る。
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address:", addrs[0])	

	if len(os.Args) > 1 {
		addr, err := multiaddr.NewMultiaddr(os.Args[1])
		if err != nil {
			panic(err)
		}

		info, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			panic(err)
		}
		
		if err := node.Connect(context.Background(), *info); err != nil {
			panic(err)
		}

		stream, _ := node.NewStream(context.Background(), info.ID, "chat/1.0.0")
		handleStream(stream)
		
	} else {
		
	node.SetStreamHandler("chat/1.0.0", handleStream)

	//プロセスの停止まで待つ
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	fmt.Println("shut down: ", s)
	}

	
}

func handleStream(stream network.Stream){
	fmt.Println("handle!!!!")
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	go readData(rw)
	go writeData(rw)	
}

func readData(rw *bufio.ReadWriter){
	/*
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if str == "quit" {
			return
		}

		fmt.Printf("\x1b[32m%s\x1b[0m", str)
	}
	*/
	fmt.Println("read")
	for {
			buf := make([]byte, 100)
			rw.Read(buf)
			fmt.Printf("\x1b[32m%s\x1b[0m> ", string(buf))
			time.Sleep(time.Second)
		}
}

func writeData(rw *bufio.ReadWriter) {
	/*
	_, err := rw.WriteString("Hello!")
	if err != nil {
		panic(err)
	}
	*/
	fmt.Println("write")
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}