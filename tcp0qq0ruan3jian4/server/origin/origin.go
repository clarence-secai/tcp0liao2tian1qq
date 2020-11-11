package origin

import (
	"fmt"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/common"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/utils"
	"net"
)

//由于conn这一数据类型不好存储进数据库，需直接在服务端维护,原因也可以是如果把好友
//的conn也存进数据库，则要求不仅有所有用户的数据表，而且每个用户有一个好友姓名、conn数
//据表，这会导致同一个人,出现在多个人的好友数据表中，比较重复，不如统一维护在服务端。同时，
//因需要频繁的使用conn，每次都从数据库拿取会降低效率。且用户一旦下线，下次上线又会是新的conn
var OnlineUsers = make(map[string]*common.OnlineUser) //姓名为键的所有在线用户列表
var OnlineGroup = make(map[string]*common.Group)      //群名为键的所有群列表
var LeaveWordTime int						//一对一好友聊天留言条数

func Start()net.Listener{
	listener, err := net.Listen("tcp", "localhost:8888")
	utils.ErrExit("10行", err)

	fmt.Println("服务端开始监听")
	return listener
}
