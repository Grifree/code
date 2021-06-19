package onlyRedis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/grifree/code/friend/m"
	md "github.com/grifree/code/friend/memory_datastorage"
	"strconv"
)

type Biz struct {
	RDB *redis.Client
}

func (dep Biz) AddFriend(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(ok bool, err error){
	friendKey := md.UserFriend{}.KeyName(userID)
	targetFriendKey := md.UserFriend{}.KeyName(targetUserID)
	replyOK := "reply:ok"
	script := `
		local friendKey = KEYS[1]
		local targetFriendKey = KEYS[2]
		local userID = ARGV[1]
		local targetUserID = ARGV[2]
		local replyOK = ARGV[3]
		-- 互加好友 (TODO 判断都是好友了就返回false)
		redis.call('SADD', friendKey, targetUserID) --key不存在时自动生成一个集合
		redis.call('SADD', targetFriendKey, userID)
		return replyOK
	`
	keys := []string{
		/* 1 */friendKey,
		/* 2 */targetFriendKey,
	}
	argv := []interface{}{
		/* 1 */strconv.FormatUint(userID.Uint64(), 10),
		/* 2 */strconv.FormatUint(targetUserID.Uint64(), 10),
		/* 3 */replyOK,
	}
	var result string
	result, err = dep.RDB.Eval(ctx, script, keys, argv...).Text();if err != nil {
		return
	}
	switch result {
	case replyOK:
		return true, nil
	default:
		return false, errors.New("unexpected")
	}
}
func (dep Biz) FriendList(ctx context.Context, userID m.IDUser)(list []m.IDUser, err error) {
	friendKey := md.UserFriend{}.KeyName(userID)
	var result []string
	result, err = dep.RDB.SMembers(ctx, friendKey).Result();if err != nil {
		return
	}
	for _, idString := range result {
		var id uint64
		id, err = strconv.ParseUint(idString, 10, 64);if err != nil {
			return
		}
		list = append(list, m.NewIDUser(id))
	}
	return
}

func (dep Biz) IsFriend(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(isFriend bool, err error){
	// 暂不考虑单方好友的错误情况，都是好友时才返回true
	friendKey := md.UserFriend{}.KeyName(userID)
	var is bool
	is, err = dep.RDB.SIsMember(ctx, friendKey, targetUserID.Uint64()).Result();if err != nil {
		return
	}
	if is == false {
		return false, nil
	}
	targetFriendKey := md.UserFriend{}.KeyName(targetUserID)
	is, err = dep.RDB.SIsMember(ctx, targetFriendKey, userID.Uint64()).Result();if err != nil {
		return
	}
	if is == false {
		return false, nil
	}
	return true,nil
}

func (dep Biz) DeleteFriend(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(ok bool, err error){
	friendKey := md.UserFriend{}.KeyName(userID)
	targetFriendKey := md.UserFriend{}.KeyName(targetUserID)
	replyOK := "reply:ok"
	script := `
		local friendKey = KEYS[1]
		local targetFriendKey = KEYS[2]
		local userID = ARGV[1]
		local targetUserID = ARGV[2]
		local replyOK = ARGV[3]
		--互删好友
		redis.call('SREM', friendKey, targetUserID)
		redis.call('SREM', targetFriendKey, userID)
		return replyOK
	`
	keys := []string{
		/* 1 */friendKey,
		/* 2 */targetFriendKey,
	}
	argv := []interface{}{
		/* 1 */strconv.FormatUint(userID.Uint64(), 10),
		/* 2 */strconv.FormatUint(targetUserID.Uint64(), 10),
		/* 3 */replyOK,
	}
	var result string
	result, err = dep.RDB.Eval(ctx, script, keys, argv...).Text();if err != nil {
		if errors.Is(err, redis.Nil) {
			// 脚本出错，遗漏返回值
		}
		return
	}
	switch result {
	case replyOK:
		return true, nil
	default:
		return false, errors.New("unexpected")
	}
}

func (dep Biz) CommonFriendList(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(list []m.IDUser, err error){
	friendKey := md.UserFriend{}.KeyName(userID)
	targetFriendKey := md.UserFriend{}.KeyName(targetUserID)
	var result []string
	result, err = dep.RDB.SInter(ctx, friendKey, targetFriendKey).Result();if err != nil {
		return
	}
	for _, idString := range result {
		var id uint64
		id, err = strconv.ParseUint(idString, 10, 64);if err != nil {
			return
		}
		list = append(list, m.NewIDUser(id))
	}
	return
}

func (dep Biz) CommonFriendCount(ctx context.Context, userID m.IDUser, targetUserID m.IDUser)(count uint64, err error){
	friendKey := md.UserFriend{}.KeyName(userID)
	targetFriendKey := md.UserFriend{}.KeyName(targetUserID)
	return dep.RDB.SInterStore(ctx, "destination", friendKey, targetFriendKey).Uint64()
}