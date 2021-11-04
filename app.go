package antsgo

import (
	"fmt"

	"github.com/antsgo/antsgo/conf"
	"github.com/antsgo/antsgo/db"
	"github.com/antsgo/antsgo/log"
	redis_app "github.com/antsgo/antsgo/redis"
	"github.com/antsgo/antsgo/server"

	"github.com/go-redis/redis"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func App() *app {
	return local
}

var local *app

func NewApp() *app {
	app := &app{}
	app.conf = conf.NewConf()
	app.logger = log.NewLogger(app.conf)
	var err error
	if !app.conf.Db.IsEmpty() {
		app.db, err = db.NewDB(app.conf.Db)
		if err != nil {
			app.Logger("").Panicf("数据库连接失败:%+v", err)
		}
	}
	if !app.conf.Redis.IsEmpty() {
		app.redis, err = redis_app.NewRedis(app.conf.Redis)
		if err != nil {
			app.Logger("").Panicf("Redis连接失败:%+v", err)
		}
	}
	local = app
	return local
}

type app struct {
	conf                conf.Conf
	logger              map[string]*logrus.Logger
	db                  *gorm.DB
	redis               *redis.Client
	registerGrpcHandler func(*grpc.Server)
	registerHttpHandler func(*runtime.ServeMux, *grpc.ClientConn)
}

func (a *app) GetConf() conf.Conf {
	return a.conf
}

func (a *app) Logger(name string) *logrus.Logger {
	if name == "" {
		name = "default"
	}
	value, ok := a.logger[name]
	if ok {
		return value
	}
	if name != "default" {
		err := errors.New(fmt.Sprintf("无效的logger类型：%s", name))
		a.Logger("").Errorf("%+v", err)
		return a.Logger("")
	}
	return nil
}

func (a *app) GetDB() *gorm.DB {
	return a.db
}

func (a *app) GetRedis() *redis.Client {
	return a.redis
}

func (a *app) SetRegisterGrpcHandler(callback func(*grpc.Server)) {
	a.registerGrpcHandler = callback
}

func (a *app) SetRegisterHttpHandler(callback func(*runtime.ServeMux, *grpc.ClientConn)) {
	a.registerHttpHandler = callback
}

func (a *app) RunGrpc() {
	if a.registerGrpcHandler != nil {
		server.NewGrpc(a.conf, a.Logger(""), a.registerGrpcHandler)
	}
	c := make(chan bool)
	<-c
}

func (a *app) RunGrpcGateway() {
	if a.registerHttpHandler != nil {
		server.NewGrpcGateway(a.conf, a.Logger(""), a.registerHttpHandler)
	}
	c := make(chan bool)
	<-c
}

func (a *app) Run() {
	go a.RunGrpc()
	go a.RunGrpcGateway()
	c := make(chan bool)
	<-c
}
