# 1.内核tcp/ip协议栈中的tcp包的处理流程 

```C
void ip_protocol_deliver_rcu(struct net *net, struct sk_buff *skb, int protocol)
{
    ipprot = rcu_dereference(inet_protos[protocol]);
	if (ipprot) {
		if (!ipprot->no_policy) {
			if (!xfrm4_policy_check(NULL, XFRM_POLICY_IN, skb)) {
				kfree_skb(skb);
				return;
			}
			nf_reset(skb);
		}
		ret = ipprot->handler(skb);     // old way
        ret = INDIRECT_CALL_2(ipprot->handler, tcp_v4_rcv, udp_rcv, skb)    //https://elixir.bootlin.com/linux/latest/source/net/ipv4/ip_input.c#L204


int tcp_v4_rcv(struct sk_buff *skb) //
{
	if (sk->sk_state == TCP_LISTEN) {
		ret = tcp_v4_do_rcv(sk, skb);
		goto put_and_return;
	}    
}

↓

int tcp_v4_do_rcv(struct sock *sk, struct sk_buff *skb)
{
    if (tcp_rcv_state_process(sk, skb)) {
		rsk = sk;
		goto reset;
	}
}

↓

int tcp_rcv_state_process(struct sock *sk, struct sk_buff *skb)
{
    switch (sk->sk_state) {
        case TCP_ESTABLISHED:
		tcp_data_queue(sk, skb);
		queued = 1;
		break;        
}

↓

static void tcp_data_queue(struct sock *sk, struct sk_buff *skb)    // https://elixir.bootlin.com/linux/latest/source/net/ipv4/tcp_input.c#L4850
{
	struct tcp_sock *tp = tcp_sk(sk);
    ...
tcp_data_queue:
    ...
		if (!sock_flag(sk, SOCK_DEAD))
			tcp_data_ready(sk);
		return;
}

↓

void tcp_data_ready(struct sock *sk)        // https://elixir.bootlin.com/linux/latest/source/net/ipv4/tcp_input.c#L4837
{
	const struct tcp_sock *tp = tcp_sk(sk);
	int avail = tp->rcv_nxt - tp->copied_seq;

	if (avail < sk->sk_rcvlowat && !sock_flag(sk, SOCK_DONE))
		return;

	sk->sk_data_ready(sk);
}

↓
sk->sk_data_ready	=	sock_def_readable;

static void sock_def_readable(struct sock *sk)      // net/core/sock.c
{
	struct socket_wq *wq;

	rcu_read_lock();
	wq = rcu_dereference(sk->sk_wq);
	if (skwq_has_sleeper(wq))
		wake_up_interruptible_sync_poll(&wq->wait, EPOLLIN | EPOLLPRI |
						EPOLLRDNORM | EPOLLRDBAND);
	sk_wake_async(sk, SOCK_WAKE_WAITD, POLL_IN);
	rcu_read_unlock();
}


rfc1122 4.2.2.13 // https://tools.ietf.org/html/rfc1122#page-87
	If such a host issues a CLOSE call while received data is 
	still pending in TCP, or if new data is received after CLOSE is
	called, its TCP SHOULD send a RST to show that data was lost.

rfc2525 2.17 // https://tools.ietf.org/html/rfc2525#page-50

实现位于 https://elixir.bootlin.com/linux/latest/source/net/ipv4/tcp.c#L2460

