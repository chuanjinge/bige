package messages

import (
	"encoding/binary"
	"fmt"

	"github.com/buguang01/util"
	"github.com/gogo/protobuf/proto"
)

type GateProtoMessageHandle struct {
	msgHead   uint32                 //消息头
	msgMaxLen uint16                 //消息最大长度
	routelist map[uint32]interface{} //消息路由列表
}

func GateProtoMessageSetMsgHead(msghead uint32) options {
	return func(msghandle IMessageHandle) {
		msghandle.(*GateProtoMessageHandle).msgHead = msghead
	}
}

func GateProtoMessageSetMsgMaxLen(max uint16) options {
	return func(msghandle IMessageHandle) {
		msghandle.(*GateProtoMessageHandle).msgMaxLen = max
	}
}

func GateProtoMessageHandleNew(opts ...options) (msghandle *GateProtoMessageHandle) {
	msghandle = &GateProtoMessageHandle{
		routelist: make(map[uint32]interface{}),
		msgHead:   uint32(0x12340000),
		msgMaxLen: ^uint16(0),
	}
	for _, opt := range opts {
		opt(msghandle)
	}
	return msghandle
}

func (msghandle *GateProtoMessageHandle) GateMarshal(gate *GateMessage, data interface{}) ([]byte, error) {
	return nil, nil
}

//编码
func (msghandle *GateProtoMessageHandle) Marshal(msgid uint32, data interface{}) ([]byte, error) {
	return nil, nil
	// buff := &bytes.Buffer{}
	// in_data, err := json.Marshal(data)
	// tmpbuf := make([]byte, 4)
	// pklen := uint32(len(in_data)+8) | msghandle.msgHead
	// binary.BigEndian.PutUint32(tmpbuf, pklen)
	// buff.Write(tmpbuf)
	// binary.BigEndian.PutUint32(tmpbuf, msgid)
	// buff.Write(tmpbuf)
	// buff.Write(in_data)
	// return buff.Bytes(), err
}

//解码
func (msghandle *GateProtoMessageHandle) Unmarshal(buff []byte) (data interface{}, err error) {
	pklen := binary.BigEndian.Uint32(buff[:4])
	pklen = pklen ^ msghandle.msgHead
	if pklen != uint32(len(buff)) {
		return nil, fmt.Errorf("MsgLen Error:%d.", pklen)
	}
	buff = buff[4:]
	msgid := binary.BigEndian.Uint32(buff[:4])
	msget, err := msghandle.GetRoute(msgid)
	if err != nil {
		return nil, err
	}
	buff = buff[4:]
	err = proto.Unmarshal(buff, msget.(proto.Message))
	return msget, err
}

//设置消息路由
func (msghandle *GateProtoMessageHandle) SetRoute(msgid uint32, msg interface{}) {
	msghandle.routelist[msgid] = msg
}

//按消息拿出消息处理实例
func (msghandle *GateProtoMessageHandle) GetRoute(msgid uint32) (result interface{}, err error) {
	if msget, ok := msghandle.routelist[msgid]; ok {
		return util.ReflectNew(msget)
	}
	return nil, fmt.Errorf("Not exist MsgID:%d.", msgid)
}

//一个消息是否收完了
//返回这个消息应该的长度，和是否收完的信息
func (msghandle *GateProtoMessageHandle) CheckMaxLenVaild(buff []byte) (msglen uint32, ok bool) {
	pklen := binary.BigEndian.Uint32(buff[:4])
	pklen = pklen ^ msghandle.msgHead
	if pklen > uint32(msghandle.msgMaxLen) {
		return 0, false
	}
	if pklen > uint32(len(buff)) {
		return pklen, false
	}
	return pklen, true
}