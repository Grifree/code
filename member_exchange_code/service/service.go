package service

import (
	"context"
)

type Service struct {

}
/* 并发场景 需求:
实现一个页面，访问这个页面有可能出现兑换码(uuid), 前100个访问的，能拿到兑换码。
不需要用户标识来源
限定每天只需新增100名会员, 抢的码是会员兑换资格码
每天9点开放100个
*/
func (dep Service) getExchangeCode_demo1(ctx context.Content) (vipInviteCode string, reject error) {
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
}
/* demo1:有redis做数量判断后,sql的数量判断就重复了,且延时性问题加剧.此外sql事务和redis减数量不是原子性操作,会导致部分兑换码没有发出去.*/

func (dep Service) getExchangeCode_demo2(ctx context.Content) (vipInviteCode string, reject error) {
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
}
/* demo2: 需要明确点在于
1兑换码是需要存储的
2存储于redis和sql都是可以的,且存在于redis更便于操作,利于解决并发问题
3redis尽可能不存储持久化的数据的问题好解决,可以每日同步到sql中
另个需求考虑点,避免一个人领取多个兑换码,可以加入用户体系,加入每个userid兑换的限制
*/