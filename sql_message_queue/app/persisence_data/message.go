package pd

import sq "github.com/goclub/sql"

type Message struct {
	UserID string `json:"userID"`
	Amount float64 `json:"amount"`
}
type MessageQueueTable struct {
	WithoutSoftDelete
}
type WithoutSoftDelete struct {}
func (WithoutSoftDelete) SoftDeleteWhere() sq.Raw {return sq.Raw{}}
func (WithoutSoftDelete) SoftDeleteSet()   sq.Raw {return sq.Raw{}}

func (MessageQueueTable) TableName() string {
	return "message_queue"
}

type IDMessageQueue uint64
type MessageQueue struct {
	MessageQueueTable
	sq.DefaultLifeCycle
	ID IDMessageQueue `db:"id"`
	Message []byte `db:"message"`
	Status MessageQueueStatus `db:"status"`
	OwnerKey string `db:"owner_key"`
	DeliveryCounter uint64 `db:"delivery_counter"`
	sq.CreatedAtUpdatedAt
}
type MessageQueueStatus uint8
func (MessageQueueStatus) Enum() (e struct {
	Pending MessageQueueStatus
	Processing MessageQueueStatus
	ACK MessageQueueStatus
	Fail MessageQueueStatus
}) {
	e.Pending = 0
	e.Processing = 1
	e.ACK = 2
	e.Fail = 3
	return
}
func (MessageQueue) Column() (col struct {
	ID sq.Column
	Message sq.Column
	OwnerKey sq.Column
	Status sq.Column
	DeliveryCounter sq.Column
	UpdatedAt sq.Column
}) {
	col.ID = "id"
	col.Message = "message"
	col.OwnerKey = "owner_key"
	col.Status = "status"
	col.DeliveryCounter = "delivery_counter"
	col.UpdatedAt = "updated_at"
	return
}