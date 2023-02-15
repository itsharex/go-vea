package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
)

var (
	Logger     *logrus.Logger
	Redis      *redis.Client
	IpSearcher *xdb.Searcher
)
