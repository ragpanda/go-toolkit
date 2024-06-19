package gin_server

import (
	"context"

	"github.com/gin-gonic/gin"
)

const BizDataKey = "_bizdata"

type BizData struct {
	UserID string         `json:"user_id"`
	LogID  string         `json:"log_id"`
	Custom map[string]any `json:"custom"`
}

func NewBizData() *BizData {
	return &BizData{
		Custom: make(map[string]any),
	}
}

func (b *BizData) SetKey(k string, v any) {
	b.Custom[k] = v
}

func (b *BizData) GetKey(k string) any {
	return b.Custom[k]
}

func (b *BizData) DeepCopy() *BizData {
	newData := NewBizData()
	newData.UserID = b.UserID
	newData.LogID = b.LogID
	for k, v := range b.Custom {
		newData.Custom[k] = v
	}
	return newData
}

func BizDataMw(c *gin.Context) {
	d := c.Value(BizDataKey)
	if d == nil {
		SetBizDataToGinCtx(c, NewBizData())
	}

	c.Next()
}

func GetBizData(c context.Context) *BizData {
	d := c.Value(BizDataKey)
	if d == nil {
		return NewBizData()
	}
	return d.(*BizData)
}

func SetBizData(c context.Context, d *BizData) context.Context {
	c = context.WithValue(c, BizDataKey, d)
	return c
}

func SetBizDataToGinCtx(c *gin.Context, d *BizData) {
	c.Set(BizDataKey, d)
}
