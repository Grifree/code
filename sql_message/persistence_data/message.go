package pd

type IDMessage string
type Message struct {
	ID IDMessage
	Content string
	Status MessageStatus
	Result string
}
type MessageStatus string
func (m MessageStatus) Enum() (e struct{
	Unhandle MessageStatus
	Processing MessageStatus
	Fail MessageStatus
	Success MessageStatus
}) {
	e.Unhandle = "unhandle"
	e.Processing = "processing"
	e.Success = "success"
	e.Fail = "fail"
	return
}