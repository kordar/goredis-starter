package goredis_starter

import (
	"context"
	"crypto/tls"
	goframeworkgoredis "github.com/kordar/goframework-goredis"
	logger "github.com/kordar/gologger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"net"
)

var (
	_tlsConfig *tls.Config
	_dialer    func(ctx context.Context, network, addr string) (net.Conn, error)
	_onConnect func(ctx context.Context, cn *redis.Conn) error
)

func SetTlsConfig(tlsConfig *tls.Config) {
	_tlsConfig = tlsConfig
}

func SetDialerFn(dialer func(ctx context.Context, network, addr string) (net.Conn, error)) {
	_dialer = dialer
}

func SetOnConnectFn(onConnect func(ctx context.Context, cn *redis.Conn) error) {
	_onConnect = onConnect
}

func HasRedisInstance(db string) bool {
	return goframeworkgoredis.HasRedisInstance(db)
}

func CloseRedisInstance(db string) {
	goframeworkgoredis.RemoveRedisInstance(db)
}

type RedisModule struct {
}

func (m RedisModule) Name() string {
	return "goredis_starter"
}

func (m RedisModule) Load(value interface{}) {
	items := cast.ToStringMap(value)
	for key, val := range items {
		section := cast.ToStringMapString(val)
		if err := goframeworkgoredis.AddRedisInstanceArgs(key, section, _tlsConfig, _dialer, _onConnect); err != nil {
			logger.Errorf("[goredis-starter] 初始化redis异常，err=%v", err)
		}
	}
}

func (m RedisModule) Close() {
}
