/*
 * @Date: 2020-08-14 14:44:38
 * @LastEditors: liuruihao@huobi.com
 * @LastEditTime: 2020-08-14 17:38:22
 * @FilePath: /tcpproxy/channel.go
 */
package main

import (
	"bufio"
	"encoding/json"
	"net"

	logging "github.com/ipfs/go-log/v2"
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
	defer c.Close()
	scanner := bufio.NewScanner(c.inBind)
	msg := JsonRpcRequest{}
	for scanner.Scan() {
		b := scanner.Bytes()
		if err := json.Unmarshal(b, &msg); err == nil {
			if msg.Method == "mining.authorize" && len(msg.Params) > 0 {
				user := ""
				if err := json.Unmarshal(msg.Params[0], &user); err != nil {
					c.log.Error("json Unmarshal failed: ", err)
				} else {
					c.log = logging.Logger(c.inBind.RemoteAddr().String() + " " + user)
				}
			}
		} else {
			c.log.Error("json Unmarshal failed: ", err)
		}
		c.log.Debug("s<--: ", scanner.Text())
		b = append(b, '\n')
		if _, err := c.outBind.Write(b); err != nil {
			c.log.Error("write to outbind failed: ", err)
		}
	}
	c.log.Info("inbind conn close")
}
func (c *Channel) server2Client() {
	defer c.Close()
	scanner := bufio.NewScanner(c.outBind)
	for scanner.Scan() {
		c.log.Debug("s-->: ", scanner.Text())
		if _, err := c.inBind.Write([]byte(scanner.Text() + "\n")); err != nil {
			c.log.Error("write to inbind failed: ", err)
		}
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
