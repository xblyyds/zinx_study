package main

import "go-study/zinx/znet"

func main() {
	s := znet.NewServer("[zinx v0.1]")
	s.Serve()
}
