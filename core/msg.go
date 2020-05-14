/*
@author '彼时思默'
@time 2020/5/11 下午4:17
@describe:
*/
package core

var ChanSize = 10240

type Msg = map[string]interface{}
type MsgChan = chan Msg
type MsgStrChan = chan string
