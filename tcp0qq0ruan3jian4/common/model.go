package common

import "net"

type User struct {
	UserName string
	PassWord int64
	Friends  []string //好友姓名列表
	Groups   []string //所加入群的若干群名
	Status   string   //用户当前状态，如在线、忙碌等，本项目中暂未开发该部分功能
}

func NewUser(userName string, passWord int64) *User {
	var u User
	u.UserName = userName
	u.PassWord = passWord
	u.Friends = []string{"Tencent"}
	u.Groups = []string{}
	return &u
}

type OnlineUser struct {
	Name     string //用户姓名
	Conn     net.Conn
	Addr     string
	Location string //当前所在群，用于在某一群时不接受其他群的群消息
}
type Group struct {
	GroupName string
	Owner     *User         //实际只用上了用户的姓名，其他的密码、好友列表切片、所加群列表切片用不上
	Members   []*OnlineUser //群成员姓名组成的切片
}

func NewGroup(groupName string, ownerName string, conn net.Conn) *Group {
	//var menbers []*OnlineUser
	//menbers = append(menbers,&OnlineUser{Name:ownerName,Conn: conn})
	group := Group{
		GroupName: groupName,
		Owner:     &User{UserName: ownerName},
		//Members: menbers,        //todo:这样没法通过thisUser.Location = intoThisGroup来统一标记
	}
	return &group
}
