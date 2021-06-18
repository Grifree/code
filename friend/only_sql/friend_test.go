package onlySQL_test

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	sq "github.com/goclub/sql"
	"github.com/grifree/code/friend/m"
	onlySQL "github.com/grifree/code/friend/only_sql"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var db *sq.Database
func init () {
	var err error
	var dbClose func() error
	db, dbClose, err = sq.Open("mysql", sq.MysqlDataSource{
		// 生产环境请使用环境变量或者配置中心配置数据库地址，不要硬编码在代码中
		User:     "root",
		Password: "somepass",
		Host:     "127.0.0.1",
		Port:     "3306",
		DB:       "code_demo",
		Query: map[string]string{
			"charset": "utf8",
			"parseTime": "True",
			"loc": "Local",
		},
	}.FormatDSN()) ; if err != nil {
		// 大部分创建数据库连接失败应该panic
		panic(err)
	}
	// 使用 init 方式连接数据库则无需 close ，依赖注入场景下才需要 close
	_ = dbClose()
}

func TestOnlyFriend(t *testing.T) {
	biz := onlySQL.Biz{
		DB: db,
	}
	ctx := context.Background()
	// 清空表
	{
		result,err := biz.DB.Core.ExecContext(ctx, "truncate table user_friend");if err != nil {
			panic(err)
		}
		log.Print(result.RowsAffected())
	}
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

		ok, err = biz.AddFriend(ctx, m.NewIDUser(1), m.NewIDUser(2))
		assert.Equal(t, false, ok)
		assert.Error(t, err, "已经是好友了")
	}
	//list(1) // [2]
	//list(2) // [1]
	{
		list, err := biz.FriendList(ctx, m.NewIDUser(1))
		log.Print("list(1): ",list)
		assert.NoError(t, err)
		assert.Equal(t, []onlySQL.FriendListReply{
			{m.NewIDUser(2),"b"},
		}, list)

		list, err = biz.FriendList(ctx, m.NewIDUser(2))
		log.Print("list(2): ",list)
		assert.NoError(t, err)
		assert.Equal(t, []onlySQL.FriendListReply{
			{m.NewIDUser(1),"a"},
		}, list)
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
		log.Print("list(1): ",list)
		assert.NoError(t, err)
		assert.Equal(t, []onlySQL.FriendListReply{
			{m.NewIDUser(2),"b"},
			{m.NewIDUser(3),"c"},
		}, list)

		list, err = biz.FriendList(ctx, m.NewIDUser(3))
		log.Print("list(3): ",list)
		assert.NoError(t, err)
		assert.Equal(t, []onlySQL.FriendListReply{
			{m.NewIDUser(1),"a"},
		}, list)
	}
	// 测试删除
	//delete(1, 2) // ok
	//delete(1, 2) // false
	//is(1,2) // false
	//list(1) // [3]
	{
		ok, err := biz.DeleteFriend(ctx, m.NewIDUser(1), m.NewIDUser(2))
		assert.NoError(t, err)
		assert.Equal(t, true, ok)

		ok, err = biz.DeleteFriend(ctx, m.NewIDUser(1), m.NewIDUser(2))
		assert.Error(t, err, "不存在好友关系")
		assert.Equal(t, false, ok)

		isFriend, err := biz.IsFriend(ctx, m.NewIDUser(1), m.NewIDUser(2))
		assert.NoError(t, err)
		assert.Equal(t, false, isFriend)

		list, err := biz.FriendList(ctx, m.NewIDUser(1))
		log.Print("list(1): ",list)
		assert.NoError(t, err)
		assert.Equal(t, []onlySQL.FriendListReply{
			{m.NewIDUser(3),"c"},
		}, list)
	}
}
