package cache

import (
	"sync"
	redigo "github.com/maxwell92/gokits/redigo"
	redis "github.com/garyburd/redigo/redis"
	mylog "github.com/maxwell92/gokits/log"

)

var log = mylog.Log

type RedisCache struct {
	pool *redis.Pool
}

var instance *RedisCache
var once sync.Once

func RedisCacheInstance() *RedisCache {
	return instance
}

func NewRedisCache() *RedisCache {
	once.Do(func() {
		instance = new(RedisCache)
		instance.pool = redigo.NewRedisClient()
	})
	return instance
}
/*
func (rc *RedisCache) Update(key, value string) error {
	/*
	fns := []func(string, string)(bool, error){rc.set, rc.lpush}
	if ok, err := rc.transaction(key, value, fns); !ok {
		log.Errorf("RedisCache Update Error: error=%s")
		return err
	}

	if ok, err := rc.set(key, value); !ok {
		log.Errorf("RedisCache Update Error: error=%s", err)
		return err
	}
	log.Infof("RedisCache Update Success")
	return nil
}

func (rc *RedisCache) Delete(key string) error {
	if ok, err := rc.Delete(key); !ok {
		log.Errorf("")
		return err
	}
	log.Infof("RedisCache Delete Key Success")
	return nil
}

func (rc *RedisCache) Find(key string) ([]string, error) {
	len := rc.llen(key)
	start := "0"
	end := strconv.Itoa(int(len - 1))
	results, err := rc.lrange(key, start, end)
	if err != nil {
		log.Errorf("RedisCache Find Key Error: error=%s", err)
		return nil, err
	}

	if results != nil {
		log.Infof("RedisCache Find Key Success: total %d", len(results))
		return results, nil
	}

	log.Errorf("RedisCache Find Key Error")
	return nil, nil
}
*/

func (rc *RedisCache) Get(key string) string {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Get Connection Error: Nil")
		return ""
	}
	defer conn.Close()

	ok, err := rc.Exist(key)
	if err != nil {
		log.Errorf("RedisCache Get Exist key %s Error: error=%s", key, err)
		return ""
	}
	if !ok {
		log.Warnf("RedisCache Get key %s Not Exist", key)
		return ""
	}

	result, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Errorf("RedisCache Get Key Error: error=%s", err)
		return ""
	}
	return result
}

func (rc *RedisCache) Set(key, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Set Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	result, err := conn.Do("SET", key, value)
	if err != nil {
		log.Errorf("RedisCache Set Key Error: error=%s", err)
		return false, err
	}
	log.Tracef("RedisCache Set Key Success: result=%v", result)
	return true, nil
}

func (rc *RedisCache) Mget(keys []string) interface{} {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache MGet Connection Error: Nil")
		return nil
	}
	defer conn.Close()

	for _, k := range keys {
		if exists, err := rc.Exist(k); !exists || err != nil {
			log.Errorf("RedisCache MGet Non-exist key Error")
			return nil
		}
	}

	// TODO: it should be keys...
	result, err := redis.Strings(conn.Do("MGET", keys[0]))
	if err != nil {
		log.Errorf("RedisCache MGet Key Error: error=%s", err)
		return nil
	}
	return result
}

func (rc *RedisCache) Delete(key string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Delete Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	result, err := conn.Do("DEL", key)
	if err != nil {
		log.Errorf("RedisCache Delete Key Error: error=%s", err)
		return false, err
	}
	log.Tracef("RedisCache Delete Key Success: result=%v", result)
	return true, nil
}

func (rc *RedisCache) Lpush(key, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache LPush Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("LPUSH", key, value)
	if err != nil {
		log.Errorf("RedisCache LPush Key Error: error=%s", err)
		return false, err
	}
	return true, nil
}

func (rc *RedisCache) Lrange(key, start, end string) ([]string, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache LRange Connection Error: Nil")
		return nil, nil
	}
	defer conn.Close()

	result, err := redis.Strings(conn.Do("LRANGE", key, start, end))
	if err != nil {
		log.Errorf("RedisCache LRange Key Error: error=%s", err)
		return nil, err
	}
	return result, nil
}

func (rc *RedisCache) Llen(key string) int32 {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache LLen Connection Error: Nil")
		return 0
	}
	defer conn.Close()

	result, err := redis.Int(conn.Do("LLEN", key))
	if err != nil {
		log.Errorf("RedisCache LLen Key Error: error=%s", err)
		return 0
	}
	return int32(result)
}

func (rc *RedisCache) Lrem(key, count, value string) (bool, error){
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache LLen Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("LREM", key, count, value)
	if err != nil {
		log.Errorf("RedisCache LLen Key Error: error=%s", err)
		return false, err
	}
	return true, nil
}

func (rc *RedisCache) Exist(key string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Exist Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	results, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil{
		log.Errorf("RedisCache Exist Key Error: error=%s",  err)
		return false, err
	}

	if results == 0 {
		return false, nil
	}
	return true, nil
}

func (rc *RedisCache) Sadd(key, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Sadd Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	results, err := redis.Int(conn.Do("SADD", key, value))
	if err != nil{
		log.Errorf("RedisCache Sadd Key Error: error=%s",  err)
		return false, err
	}

	if results == 0 {
		return false, nil
	}
	return true, nil
}

func (rc *RedisCache) Smember(key string) ([]string, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Smember Connection Error: Nil")
		return nil, nil
	}
	defer conn.Close()

	results, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil{
		log.Errorf("RedisCache Smember Key Error: error=%s",  err)
		return nil, err
	}

	return results, nil
}

func (rc *RedisCache) Scard(key string) (int32, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Scard Connection Error: Nil")
		return 0, nil
	}
	defer conn.Close()

	results, err := redis.Int(conn.Do("SCARD", key))
	if err != nil{
		log.Errorf("RedisCache Scard Key Error: error=%s",  err)
		return 0, err
	}

	return int32(results), nil
}

func (rc *RedisCache) Srem(key, value string) (bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Srem Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()

	results, err := conn.Do("SREM", key, value)
	if err != nil{
		log.Errorf("RedisCache Srem Key Error: error=%s",  err)
		return false, err
	}

	if results == 0 {
		return false, nil
	}
	return true, nil
}

/*
func (rc *RedisCache) Transaction(key, value string, fns func(string, string))(bool, error) {
	conn := rc.pool.Get()
	if conn == nil {
		log.Errorf("RedisCache Transaction Connection Error: Nil")
		return false, nil
	}
	defer conn.Close()
	sync.Mutex.Lock()
	_, err := conn.Do("MULTI")
	if err != nil {
		log.Errorf("RedisCache Transaction Key Error: error=%s", err)
		return false, err
	}

	for _, f := range fns {
		f(key, value)
	}

	_, err1 := conn.Do("EXEC")
	if err1 != nil {
		log.Errorf("RedisCache Transaction Key Error: error=%s", err)
		return false, err1
	}
	sync.Mutex.Unlock()
	return true, nil
}


func (rc *RedisCache) Search(keyword string) (*map[string]string, error) {return nil, nil}
func (rc *RedisCache) Watch(name string) error {return nil}
func (rc *RedisCache) Create(db string) error {return nil}
func (rc *RedisCache) Index(columns []string) error {return nil}
*/
