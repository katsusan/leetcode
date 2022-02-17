refer:
    https://segmentfault.com/a/1190000003063859
epoll implementation:
    https://programming.vip/docs/linux-kernel-notes-epoll-implementation-principle.html

# 1. IO过程
    linux下每次IO访问，数据会先拷贝到内核缓冲区，然后才会从内核的缓冲区拷贝到用户地址空间
    比如read操作有以下两个阶段：
    - 等待数据准备
    - 将数据从内核拷贝到用户空间

    基于以上阶段，分为以下五个网络模式：
    - 阻塞I/O
    - 非阻塞I/O
    - I/O多路复用
    - 信号驱动I/O
    - 异步I/O


# 2. 阻塞I/O
    默认情况linux下所有socket都是阻塞的。当调用recv/recvfrom/recvmsg时，在数据被拷贝到内核缓冲区，内核会将其拷贝
    到用户空间并解除用户进程阻塞。

# 3. 非阻塞I/O
    linux下可以用fcntl将socket设置为非阻塞I/O。当调用recv/recvfrom/recvmsg时，当内核缓冲区中数据没有准备好，则会
    立刻返回-1，并将error设为EAGAIN或EWOULDBLOCK。
    此种模式下用户进程需不断询问内核数据是否准备就绪。

# 4. I/O多路复用
    即常说的select/poll/epoll。不断轮询所有监视的socket，有数据就绪就通知用户进程。
    
    select会一直阻塞直到以下3种情况之一发生：
    - 某个文件描述符准备就绪
    - 此次调用被某个信号处理中断
    - timeout超时
    
    int select(int nfds, fd_set *readfds, fd_set *writefds, fd_set *exceptfds, struct timeval *timeout);
    
    example：

    ```C
    #include <stdio.h>
    #include <stdlib.h>
    #include <sys/time.h>
    #include <sys/types.h>
    #include <unistd.h>

    int
    main(void)
    {
        fd_set rfds;
        struct timeval tv;
        int retval;

        /* Watch stdin (fd 0) to see when it has input. */
        FD_ZERO(&rfds);
        FD_SET(0, &rfds);

        /* Wait up to five seconds. */
        tv.tv_sec = 5;
        tv.tv_usec = 0;

        retval = select(1, &rfds, NULL, NULL, &tv);
        /* Don't rely on the value of tv now! */

        if (retval == -1)
            perror("select()");
        else if (retval)
            printf("Data is available now.\n");
            /* FD_ISSET(0, &rfds) will be true. */
        else
            printf("No data within five seconds.\n");

        exit(EXIT_SUCCESS);
    }
    ```

    poll与select相同，也会阻塞直至那3个条件中某个成立。
    区别在于如果你想监视fd 1, 4, 6, 8, 13，select会循环遍历0-13的文件描述符来确定数据就绪情况，
    而poll只会监视用pollfd指定的这几个，以及select会告诉你"某个描述符可读/可写",而poll会更详细
    地提供描述符的具体情况(POLLIN/POLLOUT/POLLWRBAND/POLLPRI...)等等。

    int poll(struct pollfd *fds, nfds_t nfds, int timeout);

    epoll为select/poll的增强改进版本，
    它对描述符的操作有两种模式，水平触发(level trigger)和边缘触发(edge trigger)。
    默认采用LT模式。

    - The performance of epoll scales much better than select() and poll() when monitoring large numbers of file descriptors.
    - The epoll API permits either level-triggered or edge-triggered notification. By contrast, select() and poll() provide only level-triggered notification, 
        and signal-driven I/O provides only edge-triggered notification.

    LT：支持block/nonblock，对通知已就绪的描述符不作任何操作的话还会继续提示
    ET: 支持nonblock，对通知过的描述符状态即使不作操作也不再通知

    例：
    1. 我们已经把一个用来从管道中读取数据的文件句柄(RFD)添加到epoll描述符
    2. 这个时候从管道的另一端被写入了2KB的数据
    3. 调用epoll_wait(2)，并且它会返回RFD，说明它已经准备好读取操作
    4. 然后我们读取了1KB的数据
    5. 调用epoll_wait(2)......

    如果是LT模式，那么在第5步调用epoll_wait(2)之后，仍然能受到通知。
    如果是ET模式，那么第5步epoll_wait可能会挂起，只有其它事件发生时才会再次返回，此时缓冲区内剩余数据可能会被放弃等待。
    因此recv的时候如果返回的大小等于buffer大小需要再次尝试recv一次。


# 4.1 epoll API  

步骤分为三步：
- create a context in the kernel using epoll_create     // This file descriptor is not used for I/O, it is a handle for kernel data structures.
    + recording a list of file descriptors that this process has declared an interest in monitoring — the interest list
    + maintaining a list of file descriptors that are ready for I/O — the ready list.

- add and remove file descriptors to/from the context using epoll_ctl

- wait for events in the context using epoll_wait

int epoll_create(int size)；//创建一个epoll的句柄，size用来告诉内核这个监听的数目一共有多大，返回epfd
int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event)；//注册要监听的事件类型
int epoll_wait(int epfd, struct epoll_event * events, int maxevents, int timeout);



