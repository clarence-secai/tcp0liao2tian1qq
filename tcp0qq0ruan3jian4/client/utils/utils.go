package utils

import (
	"encoding/json"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/common"
	"net"
)

func WritePkg(head string, body interface{}, conn net.Conn) {
	var clientMsg common.Msg //contract.go文件中定义的协议
	{
		clientMsg.Head = head

		bodyByte, err := json.Marshal(&body)
		ErrExit("47行", err)

		clientMsg.Body = string(bodyByte) //用户名密码结构体序列化后的字符串
	}
	clientMsgByte, err := json.Marshal(clientMsg)
	ErrExit("51行", err)
	conn.Write(clientMsgByte)
}
