package redis_test

import (
	"context"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	"log"
	"testing"
	"time"
)

// 在线实验 https://try.redis.io/
// 命令参考 http://doc.redisfans.com/
/* 教程
	https://www.jianshu.com/p/e02a1028973e
	https://studygolang.com/articles/25522?fr=sidebar
	https://github.com/go-redis/redis
	https://www.jianshu.com/p/e02a1028973e
 */

var redisDB *redis.Client
var ctx = context.Background()
//连接
func initRedisDB() {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize:   1000,
		ReadTimeout: time.Millisecond * time.Duration(100),
		WriteTimeout: time.Millisecond * time.Duration(100),
		IdleTimeout: time.Second * time.Duration(60),
	})

	pong, err := redisDB.Ping(ctx).Result()
	if err != nil {
		panic("init redis error")
	} else {
		fmt.Println("init redis ok, pong:", pong)
	}
}
// 设值
func setString(key string, val string, expTime int32) {
	redisDB.Set(ctx, key, val, time.Duration(expTime) * time.Second)
}
// 取值
func getString(key string) (string, bool) {
	r, err := redisDB.Get(ctx,key).Result()
	if err != nil {
		fmt.Println("get fail, err:", err)
		return "", false
	}
	return r, true
}


/*func TestRedisString(t *testing.T) {
	initRedisDB()
	// 设值 字符串
	setString("name", "grifree", 10)
	// 取值
	name, ok := getString("name")
	fmt.Println("name:", name, ", ok:", ok)
	// 获取过期时间
	tm, err := redisDB.TTL(ctx, "name").Result()
	fmt.Println("tm:", tm, ", err:", err)
	// 不存在才设置 过期时间 nx ex
	value, err := redisDB.SetNX(ctx, "name", "grifree2", 10*time.Second).Result()
	fmt.Println("setnx \nvalue:", value, ", err:", err)
}*/

// https://www.pianshen.com/article/503765436/
/*func TestRedisList(t *testing.T) {
	initRedisDB()

	// 不存在时, 创建并插值
	//len, err := redisDB.RPush(ctx, "list","a","b","c").Result()
	//fmt.Println(len, err)
	//len, err = redisDB.RPush(ctx, "list","d").Result()
	//fmt.Println(len, err)

	// 获取长度
	//len,err = redisDB.LLen(ctx, "list").Result()
	//fmt.Println(len, err)

	// 获取值
	array,err := redisDB.LRange(ctx, "list",0,100).Result()
	fmt.Println(array, err)

	// 修剪
	//trimOK, err := redisDB.LTrim(ctx, "list",3,7).Result()
	//fmt.Println(trimOK, err)

	// 返回位置index的值
	value, err := redisDB.LIndex(ctx, "list", 0).Result()
	fmt.Println(value, err)

	// 给index位置的元素赋值
	//lsetOK, err := redisDB.LSet(ctx, "list", 1, "hello").Result()
	//fmt.Println(lsetOK, err)

	// 删除并返回 首元素
	//value,err = redisDB.LPop(ctx, "list").Result()
	//fmt.Println(value, err)

	// 删除并返回 尾元素
	//value,err = redisDB.RPop(ctx, "list").Result()
	//fmt.Println(value, err)
}*/

func TestRedisHash(t *testing.T){
	initRedisDB()

	// 添加
	datas := map[string]interface{}{
		"name": "LI LEI",
		"sex":  "male",
		"age":  28,
		"tel":  123445578,
	}
	setOK, err := redisDB.HMSet(ctx, "hash_test", datas).Result()
	fmt.Println("setOK:",setOK, err)

	// 获取
	rets, err := redisDB.HMGet(ctx, "hash_test", "name", "sex").Result()
	fmt.Println("rets:", rets, err)

	// 获取全部
	retAll, err := redisDB.HGetAll(ctx, "hash_test").Result()
	fmt.Println("retAll:", retAll, err)

	// 存在
	bExist, err := redisDB.HExists(ctx, "hash_test", "tel").Result()
	fmt.Println("bExist:", bExist, err)

	bRet, err := redisDB.HSetNX(ctx, "hash_test", "id", 100).Result()
	log.Println("bRet:", bRet, err)
}