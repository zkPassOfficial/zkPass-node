package node

import (
	"log"
	"net/http"
	"zkpass-node/app/middleware"
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
	// g.Use(gin.Logger())

	//todo auth
	g.Use(middleware.Auth())

	//todo
	//1. get sid from http header
	//2. session :=sm.GetSession(sid)
	//3.bind session to http context

	g.GET("/login/:sid", func(c *gin.Context) {
		sid := c.Param("sid")
		if sm.Has(sid) {
			log.Println("Error: Session already exists ", sid)
			// TODO: for Session with new rotated of key, needs to recreate it
		}
		s := new(session.Session)
		s.Id = sid
		// c.destroyChan = cm.destroyChan

		sm.Add(s)

		c.JSON(http.StatusOK, gin.H{
			"mes": "sid:" + sid,
		})
	})

	api.Register(g)

	g.Run(":3333")

	go sm.CheckSessions()
	go sm.CheckKickoffSignal()

	return n, nil
}
