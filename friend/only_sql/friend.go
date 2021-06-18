package onlySQL

import (
	"context"
	"errors"
	sq "github.com/goclub/sql"
	"github.com/grifree/code/friend/m"
)

type Biz struct {
	DB *sq.Database
}
func sortUserID(userID m.IDUser, friendUserID m.IDUser)(m.IDUser,m.IDUser){
	if userID < friendUserID {
		return userID , friendUserID
	}
	if friendUserID < userID {
		return friendUserID, userID
	}
	// 实际上这里应该报错了,但是练习就暂不考虑出错的情况
	return userID , friendUserID
}
func (dep Biz) AddFriend(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(ok bool, err error){
	if userID == targetUserID {
		return false, errors.New("不能加自己")
	}
	if targetUserID == 0 {
		return false, errors.New("好友不能空")
	}
	// 校验userID targetUserID 是否存在
	{
		has, err := dep.DB.Has(ctx, sq.QB{
			Table: &m.User{},
			Where: sq.And(m.User{}.Column().ID, sq.Equal(targetUserID)),
			Review: "SELECT 1 FROM `user` WHERE `id` = ? LIMIT ?",
		});if err != nil {
			return false, err
		}
		if has == false {
			return false, errors.New("好友不存在")
		}
	}
	id, friendID := sortUserID(userID,targetUserID)
	// sq(`Insert ignore INTO user_friend (user_id, friend_user_id) VALUES (?,?)`, id, friendID)
	result, err := dep.DB.Core.ExecContext(ctx,"INSERT IGNORE INTO `user_friend` (`user_id`,`friend_user_id`) VALUES (?,?)",id,friendID);if err != nil {
		return false,err
	}
	rowAffect, err := result.RowsAffected();if err != nil {
		return false,err
	}
	switch rowAffect {
	case 0:
		return false, errors.New("已经是好友了")
	case 1:
		return true, nil
	default:
		return false, errors.New("unexpected")
	}
}
type FriendListReply struct {
	UserID m.IDUser
	Name string
}
func (dep Biz) FriendList(ctx context.Context,userID m.IDUser) (list []FriendListReply, err error) {
	row, err := dep.DB.Core.QueryContext(ctx,`
		select id,name
		from USER
		where id in (
			select user_id
			from user_friend
			where friend_user_id = ?
		)
		union
		select id,name
		from USER
		where id in (
			select friend_user_id
			from user_friend
			where user_id = ?
		)
	`,userID,userID);if err != nil {
		return
	}
	for row.Next() {
		item := FriendListReply{}
		err = row.Scan(&item.UserID, &item.Name);if err != nil {
			break
		}
		list = append(list, item)
	};if err != nil {
		return
	}
	return
}
func (dep Biz) IsFriend(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(isFriend bool, err error){
	if userID == targetUserID {
		return false, errors.New("不能加自己")
	}
	if targetUserID == 0 {
		return false, errors.New("好友不能空")
	}
	// 校验userID targetUserID 是否存在
	{
		has, err := dep.DB.Has(ctx, sq.QB{
			Table: &m.User{},
			Where: sq.And(m.User{}.Column().ID, sq.Equal(targetUserID)),
			Review: "SELECT 1 FROM `user` WHERE `id` = ? LIMIT ?",
		});if err != nil {
		return false, err
	}
		if has == false {
			return false, errors.New("好友不存在")
		}
	}

	/*
		sq(`select 1 from user_friend where user_id = ? && friend_user_id = ? limit 1`, id, friendID)
	*/
	id, friendID := sortUserID(userID,targetUserID)
	return dep.DB.Has(ctx, sq.QB{
		Table:               &m.TableUserFriend{},
		Where:               sq.And(m.TableUserFriend{}.Column().UserID, sq.Equal(id)).
			And(m.TableUserFriend{}.Column().FriendUserID, sq.Equal(friendID)),
		Review:              "SELECT 1 FROM `user_friend` WHERE `user_id` = ? AND `friend_user_id` = ? LIMIT ?",
	})
}
func (dep Biz) DeleteFriend(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(ok bool, err error){
	id, friendID := sortUserID(userID,targetUserID)
	/*
	 sq(`DELETE FROM user_friend WHERE user_id = ? && friend_user_id = ? limit 1`, id, friendID)
	*/
	result, err := dep.DB.HardDelete(ctx,sq.QB{
		Table:               &m.TableUserFriend{},
		Where:               sq.And(m.TableUserFriend{}.Column().UserID, sq.Equal(id)).
			And(m.TableUserFriend{}.Column().FriendUserID, sq.Equal(friendID)),
		Limit:               1,
		Review:              "DELETE FROM `user_friend` WHERE `user_id` = ? AND `friend_user_id` = ? LIMIT ?",
	});if err != nil {
		return
	}
	rowAffect, err := result.RowsAffected();if err != nil {
		return
	}
	switch rowAffect {
	case 0:
		return false, errors.New("不存在好友关系")
	case 1:
		return true, nil
	default:
		return false, errors.New("unexpected")
	}
}
