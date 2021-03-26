package inviteCodeS

import (
	"context"
	"encoding/json"
	"fmt"
	xtime "github.com/goclub/time"
	md "github.com/grifree/code/member_exchange_code/internal/memory_data"
	pd "github.com/grifree/code/member_exchange_code/internal/persistence_data"
	"time"
	sq "github.com/goclub/sql"
)

type Service struct {

}
func (dep Service) InviteCode(ctx context.Context, userID pd.IDUser) (invitedCode pd.IDInviteCode, reject error) {
	// 开始时间校验
	{
		todayDate := xtime.FormatChinaDate(time.Now())
		var startTime time.Time
		startTime, reject = xtime.ParseChinaTime(todayDate + " 09:00:00");if reject != nil {
			return
		}
		after := startTime.After(time.Now())
		if after == false {
			reject = fmt.Errorf("兑换时间未开始");return
		}
	}
	// 数据缓存
	redisKey := md.InviteCodeHash{}.RedisKey()
	countLimit := 100
	IDInvitedCode := pd.IDInviteCode(sq.UUID())
	data := md.InviteCode{
		ID:IDInvitedCode,
		UserID:userID,
		CreatedAt:xtime.NewChinaTime(time.Now()),
	}
	var dataByte []byte
	dataByte, reject = json.Marshal(data);if reject != nil {
		return
	}
	ok, count := redis.NewScript(`
		local ok = false
		local count = redis.call("HLEN", `+redisKey+`)
		if(count < `+string(countLimit)+`) then
			local res = redis.call("HSETNX", `+redisKey+`,`+IDInvitedCode.String()+`,`+string(dataByte)+`)
			if res == 1 then 
				count = count + 1
				ok = true
			end
		end
		return ok, count
	`)
	if ok == false {
		if count >= 100 {
			reject = fmt.Errorf("今日已兑换完");return
		}
		reject = fmt.Errorf("请重试");return
	}
	invitedCode = IDInvitedCode
	return
}
func (dep Service) SyncInviteCodeToPD(ctx context.Context, date string) (reject error) {

}