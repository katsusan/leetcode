package godemo

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
)

func testRedis() {
	redisclnt := redis.NewClient(&redis.Options{
		Addr: "192.168.1.50:6379",
	})
	defer redisclnt.Close()
	log.Println(redisclnt.Ping())
}

//aquireLockWithTimeout will try to acquire specified 'lock' key and failed if time exceeded.
func aquireLockWithTimeout(clnt *redis.Client, lock string, aqTimeout time.Duration, lockTimeout time.Duration) {
	identifier, err := uuid.NewV4()
	if err != nil {
		log.Println("failed to generate UUID,", err)
		return
	}

	lockname := "lock:" + lock
	end := time.Now().Add(aqTimeout)

	for time.Now().Unix() < end.Unix() {
		//succeed to acquire the lock
		if clnt.SetNX(lockname, identifier, 0).Val() == true {
			log.Println("acquire lock OK")
			clnt.Expire(lockname, lockTimeout)
			return
		} else if clnt.TTL(lockname).Val() == 0 {
			//set expire time to prevent infinite wait(dead lock)
			clnt.Expire(lockname, lockTimeout)
		}

		time.Sleep(1 * time.Millisecond)
	}
	log.Println("failed to acquire lock until timeout arrived")
	return
}
