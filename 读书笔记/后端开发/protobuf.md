refer： https://developers.google.com/protocol-buffers/docs/gotutorial

protoc  -ID:\Projects\protoc-3.6.1-win32\include --go_out=dstdir my.proto

# 1. 示例入门

```go
// [START declaration]
syntax = "proto3";
package tutorial;

import "google/protobuf/timestamp.proto";
// [END declaration]

// [START java_declaration]
option java_package = "com.example.tutorial";
option java_outer_classname = "AddressBookProtos";
// [END java_declaration]

// [START csharp_declaration]
option csharp_namespace = "Google.Protobuf.Examples.AddressBook";
// [END csharp_declaration]

// [START golang_declaration]
option go_package = "github.com/protocolbuffers/protobuf/examples/addressbook";
// [END golang_declaration]

message Person {
  string name = 1;
  int32 id = 2;  // Unique ID number for this person.
  string email = 3;

  enum PhoneType {
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
  }

  message PhoneNumber {
    string number = 1;
    PhoneType type = 2;
  }

  repeated PhoneNumber phones = 4;

  google.protobuf.Timestamp last_updated = 5;
}

// Our address book file is just one of these.
message AddressBook {
  repeated Person people = 1;
}
```

# 2 编码原理
refer：https://weibo.com/ttarticle/p/show?id=2309404659498626187307

Protobuf简单来说属于变种的base128，探索之前可以先了解base64的原理。

# 2.1 base64
base64将数据按6bit分组，剩下少于6bit的低位补0，然后在每组6bit的高位补2个0,   
假设你要传输"base"这4个字符，则按如下流程：

`b→0x62, a→0x61, s→0x73, e→0x65`

`01100010 01100001 01110011 01100101`   
  ↓ 按6bit分组   
`011000 100110 000101 110011 011001 010000`   
  ↓ 高位补0   
`00011000 00100110 00000101 00110011 00011001 00010000`  
  ↓      
`24 38 5 51 25 16`   
  ↓ 对应base64 alphabet   
`Y m F z Z Q`

![base64 alphabet](https://wx2.sinaimg.cn/large/008aq1Aply4gshfogehmwj30f008qglu.jpg)

由于原始数据按6bit分组，每3字节的数据要用4字节来表示，所以要按4Byte对齐，在后面补"=",
也就是YmFzZQ==，与命令实验的结果一致。

```sh
➜  / echo -n base | base64
YmFzZQ==
```

# 2.2 Base 128

Base64因其高位填充了2个0，因此有效编码的利用率最高只有75%,而顺其扩展可以想到按7bit分组，高位补0/1，
这就是Base128的实现思路。但由于Ascii只有128个字符，其中还有一些不可打印字符，因此Base64目前未被
编码效率更高的Base128取代。

虽然不能完美扩展到Base128，但现有的方法都是基于Base64编码扩展的变种，比如LEB128、Base85以及下面要讲的
Base 128 Varints。

# 2.3 Base 128 Varints
refer：https://developers.google.com/protocol-buffers/docs/encoding

# 2.3.1 Varints

Varints are a method of serializing integers using one or more bytes.
Smaller numbers take a smaller number of bytes. 

varint中除了最后一个字节，其余字节都有一个msb(most significant bit),msb代表是否后面还有后续字节。
比如整数1编码后为(0)0000001，用1个字节就能表示，msb为0。
而整数300的二进制表示为100101100，编码后为(1)0101100 (0)0000010。
解码的步骤如下：
- 去除高位的msb   →   0101100 0000010
- 将字节流逆序(msb为0的字节存储在高位，小端模式)  →   0000010 0101100
- 拼接剩下来的字节      →   100101100

反过来编码的步骤如下：
- 将字节按7bit分组  →   0000010 0101100  
- 将分组逆序    →   0101100 0000010
- 添加msb      →   (1)0101100 (0)0000010

注：protobuf的varints最多编码8字节的数据，因现代计算机大多数最高支持处理64位整型。

# 2.3.2 其它类型

protobuf除了整数外还支持以下数据类型(wire type):

![](https://wx3.sinaimg.cn/large/008aq1Aply4gshfsor7vqj30oh083q30.jpg)

上节讲的"整数"指的不是编程语言概念里一般的整数，而是图里的varint。

protobuf实际上对数据编码要经过下面两个步骤：
- 将编程语言的数据结构转化为wire type。
- 根据不同的wire type进行编码，前面提到的Base 128 Varints用来编码Varint类型，   
  剩下的有其它编码规则：   
    - uint32  →   wire type 0   →   varint
    - bool    →   wirt type 0   →   varint
    - string  →   wire type 2   →   length-delimited   
    - ...

下面举一些其它类型到wire type的转码规则：
- 有符号整型：采用ZigZag编码将signed int32/int64转化为wire type 0
  ZigZag的编码规则如下：
  - (n << 1) ^ (n >> 31)    // for 32-bit signed integer
  - (n << 1) ^ (n >> 63)    // for 64-bit signed integer   
  数学上来讲就是：2n(n>=0); -2n-1(n<0)

  如果不用ZigZag编码，那么负的uint64位由于高位为1，所以Base 128 Varints编码后   
  需要10bytes，使用ZigZag编码后可以减小绝对值较小的负数的占用空间。   
  ZigZag虽然节省了部分负数情况下的字节大小，但使得正数的占用空间也变大，需要取舍，  
  业务上可以尽可能用uint型数据。

- 定长数据(64bit): 直接用小端模式，不转换
- 字符串：用length + data 格式来表示，length用Base 128 Varints编码，data为utf8字节流。


# 3 消息结构

Protobuf采用proto3作为DSL来描述其支持的消息结构;

```protobuf
syntax = "proto3";

message SearchRequest {
  string query = 1;
  int32 page_num = 2;
  int32 result_per_page = 3;
}
```

`string query = 1`后面"="指定的是field number，一个message内各个字段的field number   
必须不同。

field number在protobuf的源码中被称为tag，tag由field number和type组成。
- field number左移3bits
- 在低位3bits里写入wire type

比如对于`string query = 1`而言：   
field_number = 1,  wire type = 2 (length-delimited)
- 00000001
- 00001000  // 左移3位
- 00001010  // 末尾3位填上wire type(2)

// protobuf源码里tag的编码  
```Go
// https://github.com/protocolbuffers/protobuf-go/blob/v1.27.1/encoding/protowire/wire.go#L505
// EncodeTag encodes the field Number and wire Type into its unified form.
func EncodeTag(num Number, typ Type) uint64 {
	return uint64(num)<<3 | uint64(typ&7)
}
```

虽然返回的tag是uint64，但不是uint64范围的所有数都可以用，有一部分属于保留：
- 0   →   protobuf规定tag必须为正整数
- 19000~19999   →   protobuf供内部使用的保留位

换而言之：
- field number不必从1开始，可以从合法范围内任意开始
- 不同字段的field number不必连续，只需合法且不同

但实际上由于较小的数字经过Base 128 Varints编码后占用的空间小，大多数还是选择从1开始。

修改proto文件时的注意点：
- field number一旦被分配就不应该被修改，除非能协同到所有的接收方/发送方。
- tag中不携带field name信息，更改field name不会影响消息的结构，   
  只需要字段的field number和wire type对上即可。

由于tag中携带的是wire type，而一个wire type可以被解码成多种类型，具体解码成哪一种由其
proto文件决定，因此修改proto文件中的类型可能会导致错误。

编码example：   

![](https://wx1.sinaimg.cn/large/008aq1Aply4gshfujftvoj60ta0n975802.jpg)


# 4 嵌套消息

嵌套消息在protobuf里按照wire type 2 (length-delimited), 将被嵌套的消息按照
规则编码后前面加入Base 128 Varints编码的长度即可。

编码example：
![](https://wx3.sinaimg.cn/large/008aq1Aply4gshfujfulgj30qt0okwfb.jpg)


# 5 重复消息

proto3中对于repeated字段会使用packed encoding，而对于non-repeated字段，
如果parser接受到多份同样的field，会遵从如下规则：
- numeric type和string  →   取最后收到的值
- embedded message field  →   递归merge
- repeated field    →   append到已有字段

Packed Repeated Fields：  
Version 2.1.0 introduced packed repeated fields, which in proto2 are declared like repeated fields but with the special [packed=true] option.   
In proto3, repeated fields of scalar numeric types are packed by default.

// 这个功能虽然有点像repeated fields，但编码方式不同   
These function like repeated fields, but are encoded differently.

- A packed repeated field containing zero elements does not appear in the encoded message. // 0个元素的repeated field不会编码
- Otherwise, all of the elements of the field are packed into a single key-value pair with wire type 2 (length-delimited).
  // 其它情况下所有元素会被打包到一个key-value键值对，wire type指定为2
- Each element is encoded the same way it would be normally, except without a key preceding it.
  // 每个元素按照正常规则编码，但前面不会加上key

```proto
message Test4 {
  repeated int32 d = 4 [packed=true];
}
```
以消息Test4为例，假设为字段d赋值3,270,86942，那么编码后的结果如下：

```
22        // key (field number 4, wire type 2)
06        // payload size (6 bytes)
03        // first element (varint 3)
8E 02     // second element (varint 270)
9E A7 05  // third element (varint 86942)
```

// 只有原生数值类型(使用varint、32bit或64bit的wire type)的repeated字段可以声明为packed
Only repeated fields of primitive numeric types (types which use the varint, 32-bit, or 64-bit wire types)
can be declared "packed".

# 6 字段顺序

// field number可以以任何顺序在proto文件中使用，同时这个顺序不影响序列化的顺序   
Field numbers may be used in any order in a .proto file. 
The order chosen has no effect on how the messages are serialized.
- Do not assume the byte output of a serialized message is stable.
- By default, repeated invocations of serialization methods on the same protocol buffer message instance 
  may not return the same byte output.
  // 默认情况下重复对同样的protobuf消息调用同样序列化方法可能返回不同的字节流，但相同版本的protobuf对一段消息的解码结果一致。

