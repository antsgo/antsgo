package conf

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

func NewConf() Conf {
	var flagconf string
	flag.StringVar(&flagconf, "conf", "./configs/config.yaml", "config path, eg: -conf config.yaml")
	flag.Parse()

	v := viper.New()
	v.SetConfigFile(flagconf)
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var conf Conf
	if err := v.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("Fatal error config: %s \n", err))
	}
	return conf
}

type Conf struct {
	Logger       map[string]ConfLogger
	ServerGrpc   ConfigServer `json:"grpc" mapstructure:"grpc"`
	ServerGrpcGw ConfigServer `json:"grpc_gw" mapstructure:"grpc_gw"`
	Db           ConfigDB
	Redis        ConfigRedis
}

type ConfLogger struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	MaxSize int    `json:"max_size" mapstructure:"max_size"`
	DayNum  int    `json:"day_num" mapstructure:"day_num"`
	FileNum int    `json:"file_num" mapstructure:"file_num"`
}

type ConfigServer struct {
	Port int64
}

func (c ConfigServer) IsEmpty() bool {
	return reflect.DeepEqual(c, ConfigServer{})
}

type ConfigDB struct {
	Driver      string `json:"driver"`
	Addr        string `json:"addr"`
	DbName      string `json:"db_name" mapstructure:"db_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Config      string `json:"config"`
	Prefix      string `json:"prefix"`
	MaxIdleConn int    `json:"max_idle_conn" mapstructtrue:"max_idle_conn"`
	MaxOpenConn int    `json:"max_open_conn" mapstructure:"max_open_conn"`
	LogMode     bool   `json:"log_mode"`
}

func (c ConfigDB) IsEmpty() bool {
	return reflect.DeepEqual(c, ConfigDB{})
}

type ConfigRedis struct {
	Addr              string `json:"addr"`
	Password          string `json:"password"`
	Db                int    `json:"db"`
	ExpireUserService int64  `json:"expire_user_service"`
}

func (c ConfigRedis) IsEmpty() bool {
	return reflect.DeepEqual(c, ConfigRedis{})
}
