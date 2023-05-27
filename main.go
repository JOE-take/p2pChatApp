package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	//引数で指定されたポート番号のすべてのIP
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", string(os.Args[1])))

	//新しいノードの作成で使用する秘密鍵を生成
	r := rand.Reader
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	//sourceMultiAddrをListenし、prvKeyを秘密鍵に持つノードを作成
	node, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey))
	if err != nil {
		panic(err)
	}
	defer node.Close()

	node.SetStreamHandler("chat/1.2.0", handleStream)

	//ノードのAddrInfoを作成
	peerInfo := peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}

	//IPアドレスからP2Pアドレスにする
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("libp2p node address:", addrs[1])

	//mDNSから見つかったノード受け取るためのチャネル
	peerChan := initMDNS(node, "aikotoba")

	for {
		//見つかったノードをpeerで受け取る
		peer := <-peerChan
		fmt.Println("peer found: ", peer, "connecting")

		//peerに接続
		if err := node.Connect(context.Background(), peer); err != nil {
			fmt.Println("failed to connect, continue")
			continue
		}

		//ストリームを開始
		stream, err := node.NewStream(context.Background(), peer.ID, "chat/1.2.0")
		if err != nil {
			panic(err)
		} else {
			//正しくストリームが開ければ、データの読み書きを開始
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			go streamRead(rw)
			go streamWrite(rw)
		}
	}
}

// ストリームの開始が要求された時に呼ばれる
func handleStream(stream network.Stream) {
	fmt.Println("new Stream open")
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	go streamRead(rw)
	go streamWrite(rw)
}

func streamWrite(rw *bufio.ReadWriter) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		rw.WriteString(scanner.Text())
		rw.Flush()
	}
}

func streamRead(rw *bufio.ReadWriter) {
	buf := make([]byte, 128)
	for {
		read, err := rw.Read(buf)
		if err != nil {
			panic(err)
		}
		str := string(buf[:read])
		fmt.Printf("\x1b[32m%s\x1b[0m\n", str)
	}
}
