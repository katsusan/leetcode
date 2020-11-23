refer：
    https://github.com/farmerjohngit/myblog/issues/6
    https://eli.thegreenplace.net/2018/basics-of-futexes/

# 1. futex的简单理解性定义

```
	 //uaddr指向一个地址，val代表这个地址期待的值，当*uaddr==val时，才会进行wait
	 int futex_wait(int *uaddr, int val);
	 //唤醒n个在uaddr指向的锁变量上挂起等待的进程
	 int futex_wake(int *uaddr, int n)
```

# 2. 设计前身

```
void lock() {
    while (!trylock(lockval)) {
        wait()  //add current thread to the waiting queue and yield
    }

    //get lock...
}

bool trylock(int val) {
    int count = 0;
    while(!CAS(v, 0, 1)) {
        if (++i > 10) {
            return false;   //spin for 10 times and give up temporarily
        }
    }
    return true;
}

void unlock() {
    CAS(lockval, 1, 0);
    notify();   //wake up thread in the waiting queue 
}
```

上面的设计有一个缺陷，就在于trylock失败的时候，如果trylock与wait之间有其它线程unlock了，
那么当前线程将不会被再被锁的notify所唤起。

# 3. futex的wait与wake的解释
    - int futex_wait(int *uaddr, int val);  
        //检查*uaddr与val是否相等，不相等的话立刻返回让用户态继续trylock，
        //相等的话将当前线程插入到一个队列等待唤起
    - int futex_wake(int *uaddr, int n);
        //唤醒n个等待在uaddr上的线程

# 4. 用法

```C++
#include <sys/syscall.h>
#include <sys/time.h>

// libc没有futex系统调用对应的wrapper，需要自己实现
int futex(int* uaddr, int futex_op, int val, const struct timespec* timeout,
          int* uaddr2, int val3) {
  return syscall(SYS_futex, uaddr, futex_op, val, timeout, uaddr2, val3);
}

// Waits for the futex at futex_addr to have the value val, ignoring spurious
// wakeups. This function only returns when the condition is fulfilled; the only
// other way out is aborting with an error.
void wait_on_futex_value(int* futex_addr, int val) {
  while (1) {
    int futex_rc = futex(futex_addr, FUTEX_WAIT_PRIVATE, val, NULL, NULL, 0);
    if (futex_rc == -1) {
      if (errno != EAGAIN) {
        perror("futex");
        exit(1);
      }
    } else if (futex_rc == 0) {
      if (*futex_addr == val) {
        // This is a real wakeup.
        return;
      }
    } else {
      abort();
    }
  }
}

// A blocking wrapper for waking a futex. Only returns when a waiter has been
// woken up.
void wake_futex_blocking(int* futex_addr) {
  while (1) {
    int futex_rc = futex(futex_addr, FUTEX_WAKE_PRIVATE, 1, NULL, NULL, 0);
    if (futex_rc == -1) {
      perror("futex wake");
      exit(1);
    } else if (futex_rc > 0) {
      return;
    }
  }
}

void* thd_func(void *p) {
    int* shared_data = (int*)p;
    printf("child waiting for 0xA\n");
    wait_on_futex_value(shared_data, 0xA);

    printf("child writing 0xB\n");
    *shared_data = 0xB;
    wake_futex_blocking(shared_data);
    return NULL;
}

int main(int argc, char* argv[]) {
    int v = 0;
    int* shared_v = &v;

    pthread_t childt;
    pthread_create(&childt, NULL, thd_func, (void*)shared_v);

    // Parent thread
    printf("parent writing 0xA\n");
    *shared_v = 0xA;
    wake_futex_blocking(shared_v);

    printf("parent waiting for 0xB\n");
    wait_on_futex_value(shared_v, 0xB);

    pthread_join(childt, NULL);
    return 0;
}
```

