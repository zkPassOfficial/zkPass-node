package node

import (
	"zkpass-node/pkg/api"
	"zkpass-node/pkg/session"

	"github.com/gin-gonic/gin"
)

type Options struct {
	DataDir        string
	SessionMax     int64
	SessionTimeout int64
	SessionLife    int64
}

type ZkNode struct {
	options *Options
	sm      *session.Manager
}

func New(o *Options) (n *ZkNode, err error) {

	sm := session.New(o.SessionMax, o.SessionTimeout, o.SessionLife)
	n = &ZkNode{
		options: o,
		sm:      sm,
	}

	g := gin.Default()
	api.Route(sm, g)
	g.Run(":3333")

	go sm.CheckSessions()
	go sm.CheckKickoffSignal()

	return n, nil
}
