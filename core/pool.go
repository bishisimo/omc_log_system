/*
@author '彼时思默'
@time 2020/5/11 下午4:07
@describe:
*/
package core

import (
	"github.com/bishisimo/store"
	"github.com/go-resty/resty/v2"
	"sync"
)

//var MsgChanPoolStruct=sync.Pool{New: func() interface{}{return make(MsgChan, ChanSize)}}
//var MsgStrChanPoolStruct=sync.Pool{New: func() interface{}{return make(MsgStrChan, ChanSize)}}

var MsgChanPool *MsgChanPoolStruct
var MsgStrChanPool *MsgStrChanPoolStruct
var RequestPool *RequestPoolStruct
var S3Pool *S3PoolStruct
var CosPool *CosPoolStruct

func init() {
	MsgChanPool = NewMsgChanPool()
	MsgStrChanPool = NewMsgStrChanPool()
	RequestPool = NewRequestPool()
	S3Pool = NewS3Pool()
	CosPool = NewCosPool()
}

func NewMsgChanPool() *MsgChanPoolStruct {
	return &MsgChanPoolStruct{
		Body: sync.Pool{New: func() interface{} { return make(MsgChan, ChanSize) }},
	}
}

func NewMsgStrChanPool() *MsgStrChanPoolStruct {
	return &MsgStrChanPoolStruct{
		Body: sync.Pool{New: func() interface{} { return make(MsgStrChan, ChanSize) }},
	}
}

type MsgChanPoolStruct struct {
	Body sync.Pool
}

func (p *MsgChanPoolStruct) Put(x interface{}) {
	p.Body.Put(x)
}
func (p *MsgChanPoolStruct) Get() MsgChan {
	return p.Body.Get().(MsgChan)
}

type MsgStrChanPoolStruct struct {
	Body sync.Pool
}

func (p *MsgStrChanPoolStruct) Put(x MsgStrChan) {
	p.Body.Put(x)
}
func (p *MsgStrChanPoolStruct) Get() MsgStrChan {
	return p.Body.Get().(MsgStrChan)
}

func NewRequestPool() *RequestPoolStruct {
	client := resty.New()
	return &RequestPoolStruct{
		Body: sync.Pool{New: func() interface{} { return client.NewRequest() }},
	}
}

type RequestPoolStruct struct {
	Body sync.Pool
}

func (p *RequestPoolStruct) Put(x *resty.Request) {
	p.Body.Put(x)
}
func (p *RequestPoolStruct) Get() *resty.Request {
	return p.Body.Get().(*resty.Request)
}

func NewS3Pool() *S3PoolStruct {
	return &S3PoolStruct{
		Body: sync.Pool{New: func() interface{} { return store.NewS3Connector() }},
	}
}

type S3PoolStruct struct {
	Body sync.Pool
}

func (p *S3PoolStruct) Put(x store.StoreFace) {
	p.Body.Put(x)
}
func (p *S3PoolStruct) Get() store.StoreFace {
	return p.Body.Get().(store.StoreFace)
}

func NewCosPool() *CosPoolStruct {
	return &CosPoolStruct{
		Body: sync.Pool{New: func() interface{} { return store.NewCosConnector() }},
	}
}

type CosPoolStruct struct {
	Body sync.Pool
}

func (p *CosPoolStruct) Put(x store.StoreFace) {
	p.Body.Put(x)
}
func (p *CosPoolStruct) Get() store.StoreFace {
	return p.Body.Get().(store.StoreFace)
}
