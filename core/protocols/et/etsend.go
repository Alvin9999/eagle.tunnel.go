/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-03-19 20:08:49
 * @LastEditTime: 2019-08-24 10:50:24
 */

package et

import (
	"errors"

	"github.com/eaglexiang/eagle.tunnel.go/core/protocols/et/comm"
	"github.com/eaglexiang/go-logger"
	mynet "github.com/eaglexiang/go-net"
	"github.com/eaglexiang/go-tunnel"
)

// Send 发送ET请求
func (et *ET) Send(e *mynet.Arg) error {
	// 选择Sender
	newE, err := comm.ParseNetArg(e)
	if err != nil {
		return err
	}
	sender, ok := comm.SubSenders[comm.CMDType(newE.TheType)]
	if !ok {
		logger.Error("no tcp sender")
		return errors.New("no tcp sender")
	}
	// 进入子协议业务
	return sender.Send(newE)
}

func (et *ET) checkVersionOfRelayer(t *tunnel.Tunnel) error {
	req := comm.DefaultArg.Head
	_, err := t.WriteRight([]byte(req))
	if err != nil {
		return err
	}
	reply, err := t.ReadRightStr()
	if reply != "valid valid valid" {
		logger.Warning("invalid reply for et version check: ",
			reply)
		return errors.New("invalid reply")
	}
	return nil
}

func (et *ET) checkLocalUser(t *tunnel.Tunnel) (err error) {
	if comm.DefaultArg.LocalUser.ID == "null" {
		return nil // no need to check
	}
	user := comm.DefaultArg.LocalUser.ToString()
	_, err = t.WriteRight([]byte(user))
	if err != nil {
		return err
	}
	reply, err := t.ReadRightStr()
	if err != nil {
		return err
	}
	if reply != "valid" {
		logger.Error("invalid reply for check local user: ", reply)
		return errors.New("invalid reply")
	}
	t.SetSpeedLimiter(comm.DefaultArg.LocalUser.SpeedLimiter())
	return nil
}