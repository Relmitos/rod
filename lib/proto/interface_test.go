package proto_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
)

type Client struct {
	sessionID  string
	methodName string
	params     interface{}
	err        error
	ret        interface{}
}

var _ proto.Client = &Client{}

func (c *Client) Call(ctx context.Context, sessionID, methodName string, params interface{}) (res []byte, err error) {
	c.sessionID = sessionID
	c.methodName = methodName
	c.params = params
	return utils.MustToJSONBytes(c.ret), c.err
}

func (c *Client) GetTargetSessionID() proto.TargetSessionID { return "" }

func (c *Client) GetContext() context.Context { return nil }

func (t T) CallErr() {
	client := &Client{err: errors.New("err")}
	t.Eq(proto.PageEnable{}.Call(client).Error(), "err")
}

func (t T) ParseMethodName() {
	d, n := proto.ParseMethodName("Page.enable")
	t.Eq("Page", d)
	t.Eq("enable", n)
}

func (t T) GetType() {
	method := proto.GetType("Page.enable")
	t.Eq(reflect.TypeOf(proto.PageEnable{}), method)
}

func (t T) TimeCodec() {
	raw := []byte("123.123")
	var duration proto.MonotonicTime
	t.E(json.Unmarshal(raw, &duration))

	t.Eq(123123, duration.Milliseconds())

	data, err := json.Marshal(duration)
	t.E(err)
	t.Eq(raw, data)

	raw = []byte("123")
	var datetime proto.TimeSinceEpoch
	t.E(json.Unmarshal(raw, &datetime))

	t.Eq(123, datetime.Unix())

	data, err = json.Marshal(datetime)
	t.E(err)
	t.Eq(raw, data)
}

func (t T) NormalizeInputDispatchMouseEvent() {
	e := proto.InputDispatchMouseEvent{
		Type: proto.InputDispatchMouseEventTypeMouseWheel,
	}

	data, err := json.Marshal(e)
	t.E(err)

	t.Eq(`{"type":"mouseWheel","x":0,"y":0,"deltaX":0,"deltaY":0}`, string(data))

	ee := proto.InputDispatchMouseEvent{
		Type: proto.InputDispatchMouseEventTypeMouseMoved,
	}

	data, err = json.Marshal(ee)
	t.E(err)

	t.Eq(`{"type":"mouseMoved","x":0,"y":0}`, string(data))
}

func (t T) Rect() {
	rect := proto.DOMQuad{
		336, 382, 361, 382, 361, 421, 336, 412,
	}

	t.Eq(348.5, rect.Center().X)
	t.Eq(399.25, rect.Center().Y)

	res := &proto.DOMGetContentQuadsResult{}
	t.Nil(res.OnePointInside())

	res = &proto.DOMGetContentQuadsResult{Quads: []proto.DOMQuad{rect}}
	pt := res.OnePointInside()
	t.Eq(348.5, pt.X)
	t.Eq(399.25, pt.Y)
}

func (t T) InputTouchPointMoveTo() {
	p := &proto.InputTouchPoint{}
	p.MoveTo(1, 2)

	t.Eq(1, p.X)
	t.Eq(2, p.Y)
}

func (t T) CookiesToParams() {
	list := proto.CookiesToParams([]*proto.NetworkCookie{{
		Name:  "name",
		Value: "val",
	}})

	t.Eq(list[0].Name, "name")
	t.Eq(list[0].Value, "val")
}
