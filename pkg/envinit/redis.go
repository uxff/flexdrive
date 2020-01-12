package envinit

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/uxff/flexdrive/pkg/log"
)

// 备注: 标准库扩展,不记录任何日志

var redisPool map[string]*redis.Pool

func init() {
	redisPool = make(map[string]*redis.Pool)
}

// InitRedis 初始化 redis main 函数利解决
func InitRedis(key, server, password string, maxConn int) {
	redisPool[key] = &redis.Pool{
		MaxIdle:     500,
		MaxActive:   maxConn,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			tStart := time.Now()
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if time.Now().Sub(tStart).Seconds() > 1 {
				log.Warnf("Dial Redis too long lookup??", time.Now().Sub(tStart))
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// RdsConn 所有操作绑定的这个链接上
type RdsConn struct {
	Conn redis.Conn
}

// GetRdsConn 获取一个新链接, 需要检查是否正确
func GetRdsConn(key string) (*RdsConn, bool) {
	r := &RdsConn{}
	pool, ok := redisPool[key]
	if !ok {
		return nil, false
	}
	r.Conn = pool.Get()
	return r, r.Conn != nil
}

// Close 链接用完需要释放
func (r *RdsConn) Close() {
	r.Conn.Close()
}

// Key 操作

// GetKeys KEYS
func (r *RdsConn) GetKeys(tag string) ([]string, error) {
	return redis.Strings(r.Conn.Do("keys", tag))
}

// Expire 设置key的过期时间
func (r *RdsConn) Expire(key string, t int) error {
	_, err := redis.Int(r.Conn.Do("EXPIRE", key, t))
	return err
}

// Exists 验证 key 是否存在
func (r *RdsConn) Exists(key string) (bool, error) {
	return redis.Bool(r.Conn.Do("EXISTS", key))
}

// Del 删除键
func (r *RdsConn) Del(key ...interface{}) (int, error) {
	return redis.Int(r.Conn.Do("DEL", key...))
}

// String 操作

// SetNx 如果key不存在,就添加
func (r *RdsConn) SetNx(key string, val interface{}) (int, error) {
	return redis.Int(r.Conn.Do("SETNX", key, val))
}

// Set 不管key存不存在,都添加
func (r *RdsConn) Set(key string, val interface{}) (interface{}, error) {
	return r.Conn.Do("SET", key, val)
}

// SetEx 不管key存不存在,都添加
func (r *RdsConn) SetEx(key string, expire int64, val interface{}) (interface{}, error) {
	return r.Conn.Do("SETEX", key, expire, val)
}

func (r *RdsConn) GetInt(key string) (int64, error) {
	return redis.Int64(r.Conn.Do("GET", key))
}

func (r *RdsConn) GetString(key string) (string, error) {
	return redis.String(r.Conn.Do("GET", key))
}

// Hash 操作

//Incr 计数器加1
func (r *RdsConn) Incr(key string) (int, error) {
	return redis.Int(r.Conn.Do("INCR", key))
}

//Decr 计数器减1
func (r *RdsConn) Decr(key string) (int, error) {
	return redis.Int(r.Conn.Do("DECR", key))
}

//Get 获取key值
func (r *RdsConn) Get(key string) (int, error) {
	return redis.Int(r.Conn.Do("GET", key))
}

// HIncrBy +-
func (r *RdsConn) HIncrBy(key, field string, n int) (int, error) {
	return redis.Int(r.Conn.Do("HINCRBY", key, field, n))
}

// HGetInt 获取一个int值
func (r *RdsConn) HGetInt(key, field string) (int, error) {
	return redis.Int(r.Conn.Do("HGET", key, field))
}

// HGetInt 获取一个int值
func (r *RdsConn) HGetInt64(key, field string) (int64, error) {
	return redis.Int64(r.Conn.Do("HGET", key, field))
}

// HGetString 获取一个hash值
func (r *RdsConn) HGetString(key string, field interface{}) (string, error) {
	s, err := redis.String(r.Conn.Do("HGET", key, field))
	if err == redis.ErrNil {
		return "", nil
	}
	return s, err
}

// HMGetString 获取部分key值
func (r *RdsConn) HMGetString(key string, fields []interface{}) ([]string, error) {
	var all []interface{}
	all = append(all, key)
	all = append(all, fields...)
	return redis.Strings(r.Conn.Do("HMGET", all...))
}

// HSet HSet
func (r *RdsConn) HSet(key string, field string, val interface{}) (int, error) {
	return redis.Int(r.Conn.Do("HSET", key, field, val))
}

// HMSet
// value 可以为slice []interface{}    {"k1","v1","k2","v2"}
//
// value 可以为map  map[interface{}]interface{}
// {
// 		"k1":"v1",
// 		"k2":"v2",
// }
// 使用详情见test
func (r *RdsConn) HMSet(key string, value interface{}) (interface{}, error) {
	return r.Conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(value)...)
}

// HSetNx HSet不存在就设置值
func (r *RdsConn) HSetNx(key string, field string, val interface{}) (int, error) {
	return redis.Int(r.Conn.Do("HSETNX", key, field, val))
}

//HDel HDel
func (r *RdsConn) HDel(key string, field string) (int, error) {
	return redis.Int(r.Conn.Do("HDEL", key, field))
}

// HGetAllInt 获取所有
func (r *RdsConn) HGetAllInt(key string) (map[string]int, error) {
	return redis.IntMap(r.Conn.Do("HGETALL", key))
}

// HGetAllString 获取所有
func (r *RdsConn) HGetAllString(key string) (map[string]string, error) {
	return redis.StringMap(r.Conn.Do("HGETALL", key))
}

// HLen ...
func (r *RdsConn) HLen(key string) (int, error) {
	return redis.Int(r.Conn.Do("HLEN", key))
}

// HKeys ...
func (r *RdsConn) HKeys(key string) ([]string, error) {
	return redis.Strings(r.Conn.Do("HKEYS", key))
}

////////////////////////////////////////////////////////////////
// List 操作

// RPush 插入队列尾部
func (r *RdsConn) RPush(key string, value string) (int, error) {
	return redis.Int(r.Conn.Do("RPUSH", key, value))
}

// LPop 从队列头部取出数据
func (r *RdsConn) LPop(key string) (string, error) {
	str, err := redis.String(r.Conn.Do("LPOP", key))
	if err == redis.ErrNil {
		return "", nil
	}
	return str, err
}

// Set 操作

// SAdd 集合添加元素 1: 新元素添加成功 0: 元素已存在
func (r *RdsConn) SAdd(key string, val interface{}) (int, error) {
	return redis.Int(r.Conn.Do("SADD", key, val))
}

// SisMember 检查集合中是否有成员val
func (r *RdsConn) SisMember(key string, val interface{}) (int, error) {
	return redis.Int(r.Conn.Do("SISMEMBER", key, val))
}

//SPop 随机弹出一个数据
func (r *RdsConn) SPop(key string) (string, error) {
	str, err := redis.String(r.Conn.Do("SPOP", key))
	if err == redis.ErrNil {
		return "", nil
	}
	return str, err
}

// SortedSet 操作

//ZRevRange ...
func (r *RdsConn) ZRevRange(key string, start, end int64) ([]string, error) {
	return redis.Strings(r.Conn.Do("ZREVRANGE", key, start, end))
}

// ZRange ...
func (r *RdsConn) ZRange(key string, start, end int64) ([]string, error) {
	return redis.Strings(r.Conn.Do("ZRANGE", key, start, end))
}

// ZRangeWithScore ...
func (r *RdsConn) ZRangeWithScore(key string, start, end int64) (map[string]int64, error) {
	return redis.Int64Map(r.Conn.Do("ZRANGE", key, start, end, "withscore"))
}

// ZAdd 添加 1 新创建 0 老key
func (r *RdsConn) ZAdd(key string, member string, v int64) (int64, error) {
	return redis.Int64(r.Conn.Do("ZADD", key, v, member))
}

// ZIncrBy 增量
func (r *RdsConn) ZIncrBy(key string, member string, v int64) (int64, error) {
	return redis.Int64(r.Conn.Do("ZINCRBY", key, v, member))
}

// Pub/Sub 操作

// Publish ...
func (r *RdsConn) Publish(key string, value string) (int64, error) {
	return redis.Int64(r.Conn.Do("PUBLISH", key, value))
}
