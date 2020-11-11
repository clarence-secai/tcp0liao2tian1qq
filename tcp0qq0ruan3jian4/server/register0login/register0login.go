package register0login

import (
	"encoding/json"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/common"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/dao"
	"go7tcp0liao2tian1app/20200619tcp0qq0ruan3jian4/server/utils"
)

func DealLogin(body string) string {
	var user common.User
	err := json.Unmarshal([]byte(body), &user)
	utils.ErrExit("40行", err)
	return dao.Login(user.UserName, user.PassWord)
}

func DealRegist(body string) string {
	var user common.User
	err := json.Unmarshal([]byte(body), &user)
	utils.ErrExit("40行", err)
	return dao.Regist(user.UserName, user.PassWord, user.Friends, user.Groups)
}