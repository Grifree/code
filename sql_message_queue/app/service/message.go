package service

import (
	"context"
	"encoding/json"
	"errors"
	xjson "github.com/goclub/json"
	sq "github.com/goclub/sql"
	pd "github.com/grifree/code/sql_message_queue/app/persisence_data"
)

type Service struct {
	DB *sq.Database
}

func (dep Service) PublishMessage(ctx context.Context, message pd.Message) error {
	col := pd.MessageQueue{}.Column()
	data, err := xjson.Marshal(message) ; if err != nil {
		return err
	}
	_, err = dep.DB.Insert(ctx, sq.QB{
		Table:             pd.MessageQueueTable{},
		Insert:            []sq.Insert{
			sq.Value(col.Message, data),
			sq.Value(col.Status, pd.MessageQueueStatus(0).Enum().Pending),
		},
	}) ; if err != nil {
		return err
	}
	return nil
}
func (dep Service) ConsumeMessage(ctx context.Context, ownerKey string) (message pd.Message, hasMessage bool, err error) {
	col := pd.MessageQueue{}.Column()
	result, err := dep.DB.Update(ctx, sq.QB{
		Table:             pd.MessageQueueTable{},
		Where:             sq.And(col.Status, sq.Equal(pd.MessageQueueStatus(0).Enum().Pending)),
		Update: []sq.Update{
			sq.Set(col.Status, pd.MessageQueueStatus(0).Enum().Processing),
			sq.Set(col.OwnerKey, ownerKey),
			sq.Update{
				Raw: sq.Raw{
					Query: string(col.DeliveryCounter) + ` = `+ string(col.DeliveryCounter) +` + 1`,
				},
			},
		},
		OrderBy: []sq.OrderBy{{col.ID, sq.ASC},},
		Limit: 1,
	}) ; if err != nil {
		return
	}
	affected, err := result.RowsAffected() ; if err != nil {
		return
	}
	switch affected {
	case 0:
		return  pd.Message{},false,nil

	case 1:
		mq := pd.MessageQueue{}
		hasMessage, err = dep.DB.QueryStruct(ctx, &mq, sq.QB{
			Where: sq.And(col.OwnerKey, sq.Equal(ownerKey)).
				And(col.Status, sq.Equal(pd.MessageQueueStatus(0).Enum().Processing)),
		}) ; if err != nil {
		return
	}
		if hasMessage == false {
			err = errors.New("意外错误， ownerKey 应该被找到") ; return
		}
		err = json.Unmarshal(mq.Message, &message) ; if err != nil {
		return
	}
		return
	default:
		err = errors.New("affected must be 0 or 1");return
	}
}
func (dep Service) DoneMessage(ctx context.Context, ownerKey string)( err error ){
	col := pd.MessageQueue{}.Column()
	_, err = dep.DB.Update(ctx, sq.QB{
		Table:             pd.MessageQueueTable{},
		Where:             sq.And(col.OwnerKey, sq.Equal(ownerKey)),
		Update: []sq.Update{
			sq.Set(col.Status, pd.MessageQueueStatus(0).Enum().ACK),
		},
		Limit: 1,
	}) ; if err != nil {
		return
	}
	return
}
func (dep Service) CheckMessageRetry(ctx context.Context, datetime string, quantityPer uint64) (err error){
	// 将超时消息改为重试状态
	col := pd.MessageQueue{}.Column()

	_, err = dep.DB.Update(ctx, sq.QB{
		Table:	pd.MessageQueueTable{},
		Where:	sq.And(col.Status, sq.Equal(pd.MessageQueueStatus(0).Enum().Processing)).
			And(col.UpdatedAt, sq.OP{
				Symbol: ">",
				Values: []interface{}{datetime},
			}),
		Limit:int(quantityPer),
		Update: []sq.Update{
			sq.Set(col.Status, pd.MessageQueueStatus(0).Enum().Pending),
			sq.Set(col.OwnerKey,""),
		},
	}); if err != nil {
		return
	}
	return
}
func (dep Service) CheckMessageFail(ctx context.Context, datetime string, quantityPer uint64, retryCount uint64) (err error){
	// 重试次数过多的改为失败 待人工介入检测原因
	col := pd.MessageQueue{}.Column()
	_, err = dep.DB.Update(ctx, sq.QB{
		Table:pd.MessageQueueTable{},
		Update: []sq.Update{
			sq.Set(col.Status, pd.MessageQueueStatus(0).Enum().Fail),
		},
		Where:	sq.And(col.Status, sq.Equal(pd.MessageQueueStatus(0).Enum().Processing)).
			And(col.UpdatedAt, sq.OP{
				Symbol: ">",
				Values: []interface{}{datetime},
			}).And(col.DeliveryCounter, sq.GtOrEqualInt(int(retryCount))),
		Limit:int(quantityPer),
	}); if err != nil {
		return
	}
	return
}
