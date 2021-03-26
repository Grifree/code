package pd

type IDUser string
type IDInviteCode string
func (id IDInviteCode) String() string {return string(id)}
type InviteCode struct {
	ID IDInviteCode
	UserID IDUser
	IsUse bool
	CreatedDate string
}