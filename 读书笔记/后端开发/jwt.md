from https://jwt.io/introduction/

1. JWT全称为JSON Web Tokens，一般由以下三部分组成。
    <Header>.<Payload>.<Signature>

    - <header>: 包括两部分，type和algorithm, type即"JWT", algorithm指定签名算法，如HMAC SHA256，RSA
        形式-> {
            "alg": "HS256",
            "typ": "JWT"
        }
        然后将这个json用Base64Url编码就形成了JWT的第一部分
    
    - <payload>: Payload包含了某个实体(通常是user)以及一些附加信息，总共由三种claim
        + Registered claims: 这是一些预先定义好且推荐但不强制的claim，如 iss (issuer), exp (expiration time), sub (subject), aud (audience)
            **claim名必须为3位以便兼容**
        + Public claims：可以任意定义但必须避免与IANA定义的字段(https://www.iana.org/assignments/jwt/jwt.xhtml)冲突
        + Private claims：自定义用来在各方之间共享信息的claim
        一个Payload的例子：
        {
            "sub": "1234567890",
            "name": "John Doe",
            "admin": true
        }
    - <signature>: 对上述的header和payload以及指定的密钥用header中的加密算法加密。
        HMACSHA256(
            base64UrlEncode(header) + "." +
            base64UrlEncode(payload),
            your-secret-key)

2. details

implementation in Golang:
 
```Golang
func main() {
	header := `{"alg":"HS256","typ":"JWT"}`
	payload := `{"sub":"1234567890","name":"John Doe","iat":1516239022}`
	secret := `abc`
	var b bytes.Buffer
	headerBase64 := base64.RawURLEncoding.EncodeToString([]byte(header))
	fmt.Printf("header base64: %s\n", headerBase64)     // eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9

	payloadBase64 := base64.RawURLEncoding.EncodeToString([]byte(payload))
	fmt.Printf("payload base64: %s\n", payloadBase64)   // eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ

	b.WriteString(headerBase64)
	b.WriteString(".")
	b.WriteString(payloadBase64)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(headerBase64 + "." + payloadBase64))
	sum := h.Sum(nil)
	sig := base64.RawURLEncoding.EncodeToString(sum)
	fmt.Printf("sum: %x\n", sum)
	fmt.Printf("sig base64: %s\n", sig)     // nZ86hUWPdG43W6HVSGFy6DJnDVOZhx8a73LhQ3gIxY8
}
```

finally token   = base64urlEncoding(header) + '.' + base64urlEncoding(payload) + '.' + base64urlEncoding(signature)
                = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.nZ86hUWPdG43W6HVSGFy6DJnDVOZhx8a73LhQ3gIxY8

3. Use

see https://jwt.io/introduction and https://en.wikipedia.org/wiki/JSON_Web_Token.

typically in HTTP header `Authorization: Bearer eyJhbGci...<snip>...yu5CSpyHI`.

