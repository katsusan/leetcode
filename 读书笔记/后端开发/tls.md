
### 1. RSA

https://en.wikipedia.org/wiki/RSA_(cryptosystem)

RSA includes four steps:
- key generation
- key distribution
- encryption
- decryption

basic principle:
find three very largin positive integers: e, d, n,
(m^e)^d ≡ m (mod n).

The public key can be known by everyone and is used for encrypting messages. 

The public key is represented by the integers n and e, and the private key by the integer d.
(n is also used during the decryption process, so it might be considered to be a part of the private key too)

m represents the message.

#### 1.1 Key Generation

1. Choose two distinct prime numbers `p` and `q`.
  - `p` and `q` should be chosen at random.
  - `p` and `q` are kept secret.

2. Compute `n` = p * q.
  - `n` is used as modulus for both public/private keys. it's length is usually key length.
  - `n` is released as part of public key.

3. Computes `λ(n)`, 
  - [bonus]: λ(n) is Carmichael's totient function. 
    n = pq → λ(n)=lcm(λ(p), λ(q)) [lcm:最小公倍数] → p,q are prime → λ(p)=p-1, λ(q)=q-1 → λ(n)=lcm(p-1, q-1)

4. Choose `e` such that `1 < e < λ(n)` and `gcd(e, λ(n)) = 1`, that is, gcd(e, λ(n)) = 1.[e和λ(n)的最大公约数为1]
  - `e` should has short bit-length and small Hamming weight, most commonly e = 2^16+1 = 65537.
  - `e` is released as part of public key.

5. Determine `d` as `d ≡ e^-1 (mod λ(n))`, that is, de ≡ 1 (mod λ(n))  // d,e互为模逆元
  - `d` is kept secret as the private key exponent.

#### 1.2 Key Distribution

Bob wants to send message to Alice use RSA.

Alice tranmits her public key(n, e) to Bob, private key d is never distributed.

#### 1.3 Encryption

Bob wants to send message M to Alice.
- turn M into m, so that 0 < m < n using an agreed-upon padding scheme.
- computes cyphertext c ≡ m^e (mod n) using Alice's public key e.
  + [bonus]: m^e (mod n), the modular exponentiation, can be computed using `(a*b) mod n = ((a mod n) * (b mod n)) mod n`
- Bob transmits c to Alice.

#### 1.4 Decryption

Alice use her private key exponent d to recover the original message m from c.

c^d ≡ (m^e)^d ≡ m (mod n) 

#### 1.5 Example

- 选定2个质数p=61, q=53
- n = pq = 3233
- λ(n) = lcm(p-1, q-1) = lcm(60, 52) = 780
- e = 1~780 and e is coprime to λ(n) → e can be 17
- de ≡ 1 mod (λ(n)) → 17d ≡ 1 (mod 780) → 扩展欧几里得算法(Extended Euclidean algorithm)→ d = 413

- public key = (n, e) = (3233, 17)
- for padding plaintext message m, encryption funtion is: 
  c(m) = m^e (mod n), E.g. 加密m=65, 则c(m) = 65^17 (mod 3233) = 2790

- private key = (n, d) = (3233, 413)
- for an encrypted ciphertext c, decryption function is: 
  m(c) = c^d (mod n), E.g. 解密c=2790, 则m(c) = 2790^413 (mod 3233) = 65

#### 1.6 Signing messages

// background: anyone can use Bob's public key to send him an encrypted   
// message, so Bob cannot verify the message is from Alice.   
Alice wants to send a signed message to Bob.   
  ↓   
// h = hash(m), take signature as h^d (mod n)   
Alice hashes the message and raises it to the power of d modulos n.   
  ↓   
Bob receives the message then use the same hash algorithm, `h2 = hash(m)`,   
and raises the signature `h^d (mod n)` to the power of e modulos n.
  ↓   
Bob verifies that h2 ≡ (h^d)^e (mod n).   

Thus the keys may be swapped without loss of generality, that is,    
a private key of a key pair may be used either to:   
// 私钥有两种用途,一种是用于非对称加密,来解密接收方收到的密文;   
// 另一种是用于数字签名, 用于验证发送方是否是消息的发送者.

- Decrypt a message only intended for the recipient, which may be encrypted 
  by anyone having the public key (asymmetric encrypted transport).
- Encrypt a message which may be decrypted by anyone, but which can only be
  encrypted by one person; this provides a digital signature.


#### 1.7 Proof using Fermat's little theorem

We want to prove `m^de ≡ m (mod pq)` is valid(p,q are both prime):

```
de ≡ 1 (mod λ(n)) → de = kλ(n) + 1   

  ↓   

m^de = m^(kλ(n)+1) = m * (m^λ(n))^k 

  ↓

λ(n) = lcm(p-1, q-1) => λ(n) = x(p-1) = y(q-1)

  ↓

m * (m^λ(n))^k ≡ m * (m^x(p-1))^k ≡ m * (m ^ (p-1))^xk ≡ m * 1^xk ≡ m mod p (Fermat's little theorem/费马小定理)

  ↓

m^de ≡ m mod p   
  &   
m^de ≡ m mod q   

  ↓

m^de ≡ m mod pq
```

### 2. https
refer: https://coolshell.cn/articles/21708.html   

RSA这种非对称加密算法在存在中间人攻击的情况下会产生安全问题.
即Bob无法确认'Alice的公钥'是否来源于Alice.

因此引入了数字证书(CA, Certificate Authority)这一技术.
如下:
![](https://coolshell.cn/wp-content/uploads/2022/01/certificate-768x399.png)


1. Alice把自己的相关信息和公钥打包成一个CSR文件发给CA.
2. CA在物理世界对Alice验证通过之后,用自己的私钥对CSR文件进行签名.
3. Bob拿到2)中Alice的证书之后,用CA的公钥对该证书解密,得到Alice的公钥.
4. 3)中得到的公钥确定为Alice的公钥,可以拿来解密Alice的消息.

上述的证书机构CA还可以给下级CA颁发证书, 最顶级的称为根证书.   
客户端在该情况下需要验证证书链上的所有证书.   

用openssl生成证书的example:

- 生成CA的证书（公钥）ca.crt 和私钥 ca.key.
  openssl req -newkey rsa:2048 -new -nodes -x509 -days 365 -out ca.crt -keyout ca.key \   
    -subj "/C=SO/ST=Earth/L=Mountain/O=CoolShell/OU=HQ/CN=localhost"
- 生成 alice 的私钥
  openssl genrsa -out alice.key 2048
- 生成 Alice 的 CSR – Certificate Signing Request
  openssl req -new -key alice.key 365 -out alice.csr \   
    -subj "/C=CN/ST=Beijing/L=Haidian/O=CoolShell/OU=Test/CN=localhost.alice"
- 使用 CA 给 Alice 签名证书
  openssl x509  -req -in alice.csr \
    -extfile <(printf "subjectAltName=DNS:localhost.alice") \ 
    -CA ca.crt -CAkey ca.key  \
    -days 365 -sha256 -CAcreateserial \
    -out alice.crt

#### 2.2 双向认证mtls

传统的 TLS 认证过程是:
- 客户端连接到服务器
- 服务器提供其 TLS 证书
- 客户端验证服务器的证书
- 客户端和服务器通过加密的 TLS 连接交换信息

在mtls中(加粗的为额外步骤):
- 客户端连接到服务器
- 服务器提供其 TLS 证书
- 客户端验证服务器的证书
- **客户端出示其 TLS 证书**
- **服务器验证客户端的证书**
- **服务器授予访问权限**
- 客户端和服务器通过加密的 TLS 连接交换信息







