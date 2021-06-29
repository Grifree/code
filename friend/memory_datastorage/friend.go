package md

import (
	"github.com/grifree/code/friend/m"
	"strconv"
)

type UserFriend struct {

}
func (UserFriend) KeyName(userID m.IDUser) string{
	return "friend:"+strconv.FormatUint(userID.Uint64(), 10)
}