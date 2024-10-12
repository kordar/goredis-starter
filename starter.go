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
	name string
	load func(moduleName string, itemId string, item map[string]string)
}

func NewRedisModule(name string, load func(moduleName string, itemId string, item map[string]string)) *RedisModule {
	return &RedisModule{name, load}
}

func (m RedisModule) Name() string {
	return m.name
}

func (m RedisModule) _load(id string, cfg map[string]string) {
	if err := goframeworkgoredis.AddRedisInstanceArgs(id, cfg, _tlsConfig, _dialer, _onConnect); err != nil {
		logger.Fatalf("[%s] initializing goredis: %v", m.Name(), err)
		return
	}

	if m.load != nil {
		m.load(m.name, id, cfg)
		logger.Debugf("[%s] triggering custom loader completion", m.Name())
	}

	logger.Infof("[%s] loading module '%s' successfully", m.Name(), id)
}

func (m RedisModule) Load(value interface{}) {
	items := cast.ToStringMap(value)
	for key, item := range items {
		section := cast.ToStringMapString(item)
		m._load(key, section)
	}
}

func (m RedisModule) Close() {
}
