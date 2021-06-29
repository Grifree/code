package union

import (
	"context"
	"errors"
	md "github.com/Grifree/code/friend/internal/memory_datastorage"
	"github.com/go-redis/redis/v8"
	sq "github.com/goclub/sql"
	"github.com/grifree/code/friend/m"
	"strconv"
	"time"
)

type Biz struct {
	DB *sq.Database
	RDB *redis.Client
}

func sortID(userID m.IDUser, targetUserID m.IDUser) (m.IDUser, m.IDUser){
	if(userID < targetUserID){
		return userID, targetUserID
	}else if(targetUserID < userID){
		return targetUserID, userID
	}
	// 示例代码 不考虑出错情况
	return userID, targetUserID
}
func (dep Biz) Add(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(ok bool, err error){
	// todo 写占用锁
	firstID, secondID := sortID(userID, targetUserID)
	result, err := dep.DB.Insert(ctx, sq.QB{
		Raw: sq.Raw{
			Query: "INSERT IGNORE user_friend (user_id, friend_user_id) VALUES (?, ?,)",
			Values: []interface{}{firstID, secondID},
		},
	}) ; if err != nil {
		return
	}
	row, err := result.RowsAffected();if err != nil {
		return
	}
	switch row {
	case 0:
		return false, errors.New("已经是好友了")
	case 1:
		// sql成功 row ==1 此时断了, 没有清除redis, 可以不考虑 短暂现象
		dep.RDB.Del(ctx,
			md.UserFriend{}.KeyName(userID),
			md.NoUserFriend{}.KeyName(userID),
			md.UserFriend{}.KeyName(targetUserID),
			md.NoUserFriend{}.KeyName(targetUserID),
		)
		return true, nil
	default:
		// todo record biz error
		return false, errors.New("unexpected")
	}
}
func (dep Biz) syncData(ctx context.Context, userID m.IDUser)(noFriend bool, err error){
	noFriendKey := md.NoUserFriend{}.KeyName(userID)
	friendKey :=  md.UserFriend{}.KeyName(userID)
	result, err := dep.RDB.Get(ctx, noFriendKey).Result();if err != nil {
		return 
	}
	if result == "1" {
		return true, nil
	}
	exist, err := dep.RDB.Exists(ctx, noFriendKey).Uint64();if err != nil {
		return 
	}
	if exist == 1 {
		return false, nil
	}
	var list []m.IDUser 
	err = dep.DB.QuerySlice(ctx, &list, sq.QB{
		Raw:                 sq.Raw{
			Query:  `SELECT user_id as id FROM user_friend WHERE friend_user_id=?
			UNION
			SELECT friend_user_id as id FROM user_friend WHERE user_id=?`,
			Values: []interface{}{userID, userID},
		},
		Review:              "",
	});if err != nil {
		return
	}
	if len(list) > 0 {
		Keys := []string{
			/*1*/friendKey,
		}
		Argv := []interface{}{
			/*1*/list,
			/*2*/time.Minute*2,
		}
		Script := `
			local friendKey = KEYS[1]
			local members = ARGV[1]
			local expireTime = ARGV[2]
			reids.call('sadd', friendKey, members)
			redis.call('expires', friendKey, expireTime)
		`
		dep.RDB.Eval(ctx, Script, Keys, Argv...)
		return false, nil
	}else{
		dep.RDB.SetEX(ctx, friendKey, 1, time.Minute*2)
		return true, nil
	}
}
func (dep Biz) List(ctx context.Context, userID m.IDUser)(list []m.IDUser, err error){
	noFriend, err := dep.syncData(ctx, userID);if err != nil {
		return 
	}
	if noFriend {
		return []m.IDUser{}, nil
	}
	idList, err := dep.RDB.SMembers(ctx, md.UserFriend{}.KeyName(userID)).Result();if err != nil {
		return 
	}
	for _, idString := range idList {
		id, err := strconv.ParseUint(idString, 10, 64);if err != nil {
			return 
		}
		list = append(list, m.NewIDUser(id))
	}
	return list, nil
}
func (dep Biz) Mutual(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(list []m.IDUser, err error){
	userNoFriend, err := dep.syncData(ctx, userID);if err != nil {
		return
	}
	if userNoFriend {
		return []m.IDUser{}, nil
	}
	targetNoFriend, err := dep.syncData(ctx, targetUserID);if err != nil {
		return
	}
	if targetNoFriend {
		return []m.IDUser{}, nil
	}
	idList, err := dep.RDB.SInter(ctx, md.UserFriend{}.KeyName(userID), md.UserFriend{}.KeyName(targetUserID)).Result();if err != nil {
		return
	}
	for _, idString := range idList {
		id, err := strconv.ParseUint(idString, 10, 64);if err != nil {
			return
		}
		list = append(list, m.NewIDUser(id))
	}
	return list, nil
}
func (dep Biz) Is(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(is bool, err error){
	userNoFriend, err := dep.syncData(ctx, userID);if err != nil {
		return
	}
	if userNoFriend {
		return false, nil
	}
	return dep.RDB.SIsMember(ctx, md.UserFriend{}.KeyName(userID), targetUserID).Result()
}
func (dep Biz) Delete(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(ok bool, err error){
	// todo 写占用锁
	firstID, secondID := sortID(userID, targetUserID)
	result, err := dep.DB.HardDelete(ctx, sq.QB{
		Raw:                 sq.Raw{
			Query:  `DELETE FROM user_friend WHERE user_id=? AND friend_user_id=?`,
			Values: []interface{}{firstID, secondID},
		},
		Review:              "todo",
	});if err != nil {
		return
	}
	row, err := result.RowsAffected();if err != nil {
		return
	}
	switch row {
	case 0:
		return false, errors.New("不存在好友关系")
	case 1:
		dep.RDB.Del(ctx,
			md.UserFriend{}.KeyName(userID),
			md.NoUserFriend{}.KeyName(userID),
			md.UserFriend{}.KeyName(targetUserID),
			md.NoUserFriend{}.KeyName(targetUserID),
		)
		return true, nil
	default:
		// todo record biz error
		return false, errors.New("unexpected")
	}
}