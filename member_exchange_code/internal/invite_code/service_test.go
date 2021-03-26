package inviteCodeS_test

import (
	"context"
	xtime "github.com/goclub/time"
	"github.com/google/uuid"
	"log"
	"strconv"
	"sync"
	"time"
)

type Service struct {

}
/* 并发场景 需求:
实现一个页面，访问这个页面有可能出现兑换码(uuid), 前100个访问的，能拿到兑换码。
不需要用户标识来源
限定每天只需新增100名会员, 抢的码是会员兑换资格码
每天9点开放100个
*/
/*func (dep Service) getExchangeCode_demo1(ctx context.Content) (vipInviteCode string, reject error) {
	// 开始时间校验
	{
		if now < date(`${today} 09:00:00`){
			reject = "兑换时间未开始";return
		}
	}
	// 缓存 数量校验
	{
		count := redis.Increase(`vipInviteCodeCount:${today}`, 1, 25hour)
		if count > 100 {
			redis.subtract(`vipInviteCodeCount:${today}`, 1)
			reject = "今日兑换码已抢完";return
		}
	}
	// 数据 数量校验&创建数据
	{
		isRollback,reject := beginTX(ctx, func() {
			count, err := sql.count(`select count(id) from exchane_code where created_date = ?`, today);if reject != nil {
				return tx.rollback(err)
			}
			if count > 100 {
				return tx.rollback("今日兑换码已抢完")
			}
			vipInviteCode, err := sql.create(`insert exchane_code cloumn(id,is_use,created_date) value(?,?,?)`, "@uuid", false, today);if reject != nil {
				return tx.rollback(err)
			}
			tx.commit()
		});if reject != nil {
			redis.subtract(`vipInviteCodeCount:${today}`, 1)
			return
		}
	}
	return
}*/
/* demo1:有redis做数量判断后,sql的数量判断就重复了,且延时性问题加剧.此外sql事务和redis减数量不是原子性操作,会导致部分兑换码没有发出去.*/

/*func (dep Service) getExchangeCode_demo2(ctx context.Content) (vipInviteCode string, reject error) {
	// redis lua 脚本
	{
		count := redis.Count("hash",`vipInviteCode:${today}`, 25hour)
		if count < 100 {
			vipInviteCode := "@uuid"
			redis.all("hash",`vipInviteCode:${today}`,{
				vipInviteCode:`{"is_use":false,"created_at":now}`
			})
		}
	}
	return vipInviteCode
}*/
/* demo2: 需要明确点在于
1兑换码是需要存储的
2存储于redis和sql都是可以的,且存在于redis更便于操作,利于解决并发问题
3redis尽可能不存储持久化的数据的问题好解决,可以每日同步到sql中
另个需求考虑点,避免一个人领取多个兑换码,可以加入用户体系,加入每个userid兑换的限制
*/

package main

import (
"context"
"errors"
red "github.com/goclub/redis"
xtime "github.com/goclub/time"
"github.com/google/uuid"
"github.com/mediocregopher/radix/v3"
radix3 "github.com/redis-driver/mediocregopher-radix-v3"
redScript "github.com/redis-driver/script"
"log"
"strconv"
"sync"
"time"
)

func main () {
	ctx := context.Background()
	_, err := red.DEL{
		Key:  CDKeyHashKey(),
	}.Do(ctx, client) ; if err != nil {
		log.Print(err)
	}
	var i uint64
	wg := sync.WaitGroup{}
	for i=0;i<5;i++ {
		wg.Add(1)
		go func(userID uint64) {
			log.Print(GetCDKey(ctx, userID))
			wg.Done()
		}(i)
	}
	for i=0;i<15;i++ {
		wg.Add(1)
		go func(userID uint64) {
			log.Print(GetCDKey(ctx, userID))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
func CDKeyHashKey() string {
	return "cdkey:date:" + xtime.FormatChinaDate(time.Now())
}
func GetCDKey(ctx context.Context, userID uint64) (string, error) {
	chinaNow := xtime.NewChinaTime(time.Now())
	startTime := time.Date(chinaNow.Year(), chinaNow.Month(), chinaNow.Day(), 9,0,0,0,xtime.LocationChina)
	if chinaNow.Before(startTime) {
		return "", errors.New("尚未开始")
	}
	var result string
	hashKey := CDKeyHashKey()
	maxLimit := "10"
	cdKey := uuid.New().String()
	keys := []string{
		/* 1 */ hashKey,
		/* 2 */ maxLimit,
		/* 3 */ cdKey,
	}
	argv := []string{
		/* 1 */ strconv.FormatUint(userID, 10),
	}
	script := `
local count = redis.call("HLEN", KEYS[1])
if (count < tonumber(KEYS[2]))
then
	local reply = redis.call("HSETNX", KEYS[1], ARGV[1], KEYS[3])
	if (reply == 0)
	then
		return "got"
	else
		return "ok"
	end
else
	return "outOfStock"
end
`
	err := client.RedisScript(ctx, redScript.Script{
		ValuePtr: &result,
		Script:   script,
		Keys:     keys,
		Args:     argv,
	}) ; if err != nil {
		return "", err
	}
	switch result {
	case "got":
		return "", errors.New("你已经抽过了,userID:" + strconv.FormatUint(userID, 10))
	case "ok":
		return cdKey, nil
	case "outOfStock":
		return "", errors.New("未抽到，明天再来吧")
	default:
		return "", errors.New("script result not match")
	}
}

var client red.Client

func init() {
	pool, err := radix.NewPool("tcp", "127.0.0.1:6379", 10); if err != nil {
		panic(err)
	}
	client = radix3.Client{Core: pool}
}