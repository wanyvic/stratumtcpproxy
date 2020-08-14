/*
 * @Date: 2020-08-14 14:44:38
 * @LastEditors: liuruihao@huobi.com
 * @LastEditTime: 2020-08-14 17:02:27
 * @FilePath: /tcpproxy/channel.go
 */
package main

import (
	"bufio"
	"encoding/json"
	"net"

	logging "github.com/ipfs/go-log"
)

type JsonRpcRequest struct {
	Id     interface{}       `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}
type Channel struct {
	inBind      *net.TCPConn
	outBind     *net.TCPConn
	outBindAddr *net.TCPAddr
	log         *logging.ZapEventLogger
	reader      bufio.Reader
}

func NewChannel(outBindAddr *net.TCPAddr, inBind *net.TCPConn) *Channel {
	c := &Channel{
		inBind:      inBind,
		outBindAddr: outBindAddr,
	}
	c.log = logging.Logger(inBind.RemoteAddr().String())
	return c
}

func (c *Channel) Run() {
	defer c.Close()
	var err error
	if c.outBind, err = net.DialTCP("tcp", nil, c.outBindAddr); err != nil {
		c.log.Error("outBind connection failed: ", err)
		return
	}
	go c.client2Server()
	c.server2Client()
}

func (c *Channel) client2Server() {
	scanner := bufio.NewScanner(c.inBind)
	msg := JsonRpcRequest{}
	for scanner.Scan() {
		if json.Unmarshal(scanner.Bytes(), &msg) != nil {
			if msg.Method == "mining.authorize" && len(msg.Params) > 0 {
				c.log = logging.Logger(c.inBind.RemoteAddr().String() + string(msg.Params[0]))
			}
		}
		c.log.Debug("s<--: ", scanner.Text())
		if _, err := c.outBind.Write(scanner.Bytes()); err != nil {
			c.log.Error("write to outbind failed: ", err)
		}
	}
	c.log.Info("inbind conn close")
}
func (c *Channel) server2Client() {
	scanner := bufio.NewScanner(c.outBind)
	for scanner.Scan() {
		c.log.Debug("s-->: ", scanner.Text())
	}
	c.log.Info("outbind conn close")
}

func (c *Channel) Close() error {
	if c.inBind != nil {
		c.inBind.Close()
	}
	if c.outBind != nil {
		c.outBind.Close()
	}
	return nil
}
