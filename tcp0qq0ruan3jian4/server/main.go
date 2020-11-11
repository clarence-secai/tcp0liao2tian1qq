package main

import (
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/ioMes"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/origin"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/utils"
)

func main() {
	listener := origin.Start()
	defer listener.Close() //一般不需要

	//循环接受随时可能的新增客户端拨来的专属于该用户的连接
	for {
		conn, err := listener.Accept()
		utils.ErrExit("12行", err)

		//针对每一个用户都开一个协程，来与该客户端维护一个专属于该用户的长连接保持通话
		go ioMes.IoWithConn(conn)
	}
}


