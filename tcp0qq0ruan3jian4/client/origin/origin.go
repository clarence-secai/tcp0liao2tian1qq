package origin

import (
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/client/utils"
	"net"
)

//初始化一个可以和server端进行消息往来的连接
var OneConn net.Conn
var err0 error
var ChanCommon = make(chan string, 2) //最多只需装两个字符串
func Origin() {
	OneConn, err0 = net.Dial("tcp", "localhost:8888")
	utils.ErrExit("38行", err0)
}