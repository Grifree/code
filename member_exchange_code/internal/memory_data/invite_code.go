package md

import (
	xtime "github.com/goclub/time"
	pd "github.com/grifree/code/member_exchange_code/internal/persistence_data"
	"strings"
	"time"
)

type InviteCode struct {
	ID pd.IDInviteCode `json:"id"`
	UserID pd.IDUser `json:"userID"`
	CreatedAt xtime.ChinaTime `json:"createdAt"`
}
type InviteCodeHash map[pd.IDInviteCode]InviteCode
func (InviteCodeHash) RedisKey() string {
	return strings.Join([]string{"inviteCode", xtime.FormatChinaDate(time.Now())}, ":")
}