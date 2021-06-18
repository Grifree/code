
// Generate by https://tools.goclub.vip
package m
import (
	sq "github.com/goclub/sql"
)

type TableUserFriend struct {
	sq.WithoutSoftDelete
}
// 给 TableName 加上指针 * 能避免 db.InsertModel(user) 这种错误， 应当使用 db.InsertModel(&user) 或
func (*TableUserFriend) TableName() string { return "user_friend" }
type UserFriend struct {
	UserID       IDUser  `db:"user_id"sq:"ignoreUpdate" `
	FriendUserID IDUser  `db:"friend_user_id"sq:"ignoreUpdate" `
	TableUserFriend
	sq.DefaultLifeCycle
}
func (v UserFriend) PrimaryKey() []sq.Condition {
	return sq.And(v.Column().UserID, sq.Equal(v.UserID)).
		And(v.Column().FriendUserID, sq.Equal(v.FriendUserID))
}

func (v TableUserFriend) Column() (col struct{
	UserID        sq.Column
	FriendUserID  sq.Column

}) {
	col.UserID        = "user_id"
	col.FriendUserID  = "friend_user_id"

	return
}
