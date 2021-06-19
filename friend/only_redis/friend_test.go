package onlyRedis_test

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/grifree/code/friend/m"
	md "github.com/grifree/code/friend/memory_datastorage"
	onlyRedis "github.com/grifree/code/friend/only_redis"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var rdb *redis.Client
func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
	})
}

func TestOnlyRedis(t *testing.T) {
	biz := onlyRedis.Biz{
		RDB:rdb,
	}
	ctx := context.Background()
	//log.Print(biz.AddFriend(ctx, m.NewIDUser(1), m.NewIDUser(2)))
	//log.Print(biz.FriendList(ctx, m.NewIDUser(1)))
	//log.Print(biz.IsFriend(ctx, m.NewIDUser(1), m.NewIDUser(2)))
	//log.Print(biz.DeleteFriend(ctx, m.NewIDUser(1), m.NewIDUser(2)))
	//log.Print(biz.CommonFriendList(ctx, m.NewIDUser(1), m.NewIDUser(2)))
	//log.Print(biz.CommonFriendCount(ctx, m.NewIDUser(1), m.NewIDUser(2)))

	// 清空数据
	{
		friendKey := md.UserFriend{}.KeyName(m.NewIDUser(1))
		targetFriendKey := md.UserFriend{}.KeyName(m.NewIDUser(2))
		delCount, err := rdb.Del(ctx, friendKey, targetFriendKey).Uint64();if err != nil {
			panic(err)
		}
		log.Print("delCount ",delCount)
	}
	// 测试 1 2 的好友关系
	//is(1,2) // false
	//is(2,1) // false
	{
		isFriend, err := biz.IsFriend(ctx, m.NewIDUser(1), m.NewIDUser(2))
		assert.NoError(t, err)
		assert.Equal(t, false, isFriend)

		isFriend, err = biz.IsFriend(ctx, m.NewIDUser(2), m.NewIDUser(1))
		assert.NoError(t, err)
		assert.Equal(t, false, isFriend)
	}
	//add(1,2) // ok
	//add(2,1) // repeat
	{
		ok, err := biz.AddFriend(ctx, m.NewIDUser(1), m.NewIDUser(2))
		assert.NoError(t, err)
		assert.Equal(t, true, ok)

		ok, err = biz.AddFriend(ctx, m.NewIDUser(2), m.NewIDUser(1))
		assert.NoError(t, err)
		assert.Equal(t, true, ok)
	}
	//list(1) // [2]
	//list(2) // [1]
	{
		list, err := biz.FriendList(ctx, m.NewIDUser(1))
		assert.NoError(t, err)
		assert.Equal(t, []m.IDUser{2}, list)

		list, err = biz.FriendList(ctx, m.NewIDUser(2))
		assert.NoError(t, err)
		assert.Equal(t, []m.IDUser{1}, list)
	}
	//is(1,2) // true
	//is(2,1) // true
	{
		isFriend, err := biz.IsFriend(ctx, m.NewIDUser(1), m.NewIDUser(2))
		assert.NoError(t, err)
		assert.Equal(t, true, isFriend)

		isFriend, err = biz.IsFriend(ctx, m.NewIDUser(2), m.NewIDUser(1))
		assert.NoError(t, err)
		assert.Equal(t, true, isFriend)
	}
	// 测试一个人有多个好友
	//add(1,3) // ok
	//list(1) // [2, 3]
	//list(3) // [1]
	{
		ok, err := biz.AddFriend(ctx, m.NewIDUser(1), m.NewIDUser(3))
		assert.NoError(t, err)
		assert.Equal(t, true, ok)

		list, err := biz.FriendList(ctx, m.NewIDUser(1))
		assert.NoError(t, err)
		assert.Equal(t, []m.IDUser{2,3}, list)

		list, err = biz.FriendList(ctx, m.NewIDUser(3))
		assert.NoError(t, err)
		assert.Equal(t, []m.IDUser{1}, list)
	}
	// 测试删除
	//delete(1, 2) // ok
	//is(1,2) // false
	//list(1) // [3]
	{
		ok, err := biz.DeleteFriend(ctx, m.NewIDUser(1), m.NewIDUser(2))
		assert.NoError(t, err)
		assert.Equal(t, true, ok)

		isFriend, err := biz.IsFriend(ctx, m.NewIDUser(1), m.NewIDUser(2))
		assert.NoError(t, err)
		assert.Equal(t, false, isFriend)

		list, err := biz.FriendList(ctx, m.NewIDUser(1))
		assert.NoError(t, err)
		assert.Equal(t, []m.IDUser{3}, list)
	}
}
