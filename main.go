package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	/*
		peerChan := initMDNS(node, "aikotoba")

		for {
			peer := <-peerChan
			fmt.Println("Peer: ", peer, "が見つかりました。接続します。")

			if err := node.Connect(context.Background(), peer); err != nil {
				fmt.Println("接続失敗、続行")
				continue
			}

			stream, err := node.NewStream(context.Background(), peer.ID, "chat/1.1.0")
			if err != nil {
				panic(err)
			}

			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			go streamReader(rw)
			go streamWriter(rw)
		}
	*/
	if len(os.Args) > 2 {
		addr, err := multiaddr.NewMultiaddr(os.Args[2])
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

		stream, err := node.NewStream(context.Background(), info.ID, "chat/1.1.0")
		if err != nil {
			panic(err)
		}

		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		go streamReader(rw)
		go streamWriter(rw)
	} else {
		node.SetStreamHandler("chat/1.1.0", handleStream)
	}

	//プロセスの停止まで待つ
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	fmt.Println("shut down: ", s)
}

func handleStream(stream network.Stream) {
	fmt.Println("new Stream open")
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	go streamReader(rw)
	go streamWriter(rw)
}

func streamWriter(rw *bufio.ReadWriter) {
	scanner := bufio.NewScanner(os.Stdin)
	w := rw.Writer
	for {
		scanner.Scan()
		w.Write(scanner.Bytes())
	}
}

func streamReader(rw *bufio.ReadWriter) {
	r := rw.Reader
	buf := make([]byte, 128)
	for {
		_, err := r.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(buf))
	}
}
