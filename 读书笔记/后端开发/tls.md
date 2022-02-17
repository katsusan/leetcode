
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

5. Determine `d` as `d ≡ e^-1 (mod λ(n))`, that is, de ≡ 1 (mod λ(n))
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



### 2. https流程


