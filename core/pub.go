/*
@author '彼时思默'
@time 2020/5/9 上午9:23
@describe: 发布者
*/
package core

import "context"

type Pub struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	MsgChan
}

func NewPub() *Pub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pub{
		Ctx:     ctx,
		Cancel:  cancel,
		MsgChan: MsgChanPool.Get(),
	}
}
