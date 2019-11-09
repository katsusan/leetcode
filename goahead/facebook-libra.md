最近fb推出了号称跨主权货币天秤币libra，抱着尝鲜的态度试了下，感觉前景不错。

根据developers.libra.org的说法，只要5步就可以简单启动本地的libra测试网。

##
第四步需要编译libra，会需要大量内存，不够的话可能需要扩展swap空间，下面是扩展步骤。

touch /var/swapfile     //创建swap文件
dd if=/dev/zero of=/var/swapfile bs=1024 count=2000000  //用零填充2G空间
mkswap /var/swapfile    //设定为swap区
swapon /var/swapfile    //开启swap交换，关闭用swapoff
##


1. clone到本地
    git clone https://github.com/libra/libra.git && cd libra

2. 切换到testnet分支
    git checkout testnet

3. 安装必要依赖
    ./scripts/dev_setup.sh

4. 运行测试网
    ./scripts/cli/start_cli_testnet.sh

5. 尝试基本交易指令
    libra% account create   //创建账户A
    >> Creating/retrieving next account from wallet
    Created/retrieved account #0 address 38099a5b17c2db81babca08bbcf3399299f6ced16d23c28be77dbc68a0e4dba0

    libra% account create   //创建账户B
    >> Creating/retrieving next account from wallet
    Created/retrieved account #1 address 45d4f95771caa7f66d50ebf3101eeffbda9a8e166b07dacefa6befca0d21cb30

    libra% account list     //列出所有账户
    User account index: 0, address: 38099a5b17c2db81babca08bbcf3399299f6ced16d23c28be77dbc68a0e4dba0, sequence number: 0, status: Local
    User account index: 1, address: 45d4f95771caa7f66d50ebf3101eeffbda9a8e166b07dacefa6befca0d21cb30, sequence number: 0, status: Local


    libra% account mint 0 200   //账户A加200
    >> Minting coins
    Mint request submitted
    

    libra% query balance 0      //查询余额
    Balance is: 200.000000


    libra% query sequence 0     //查询交易序列号，转出方+1，被转入账户不变
    >> Getting current sequence number
    Sequence number is: 0
    libra% query sequence 1
    >> Getting current sequence number
    Sequence number is: 0


    libra% transfer 0 1 50      //A向BKB转50
    >> Transferring
    Transaction submitted to validator
    To query for transaction status, run: query txn_acc_seq 0 0 <fetch_events=true|false>

    libra% query balance 0      //查询余额发现转账成功
    Balance is: 150.000000
    libra% query balance 1
    Balance is: 50.000000

    libra% query txn_acc_seq 0 0 true       //交易细节
    >> Getting committed transaction by account and sequence number
    Committed transaction: SignedTransaction { 
    raw_txn: RawTransaction { 
        sender: 38099a5b17c2db81babca08bbcf3399299f6ced16d23c28be77dbc68a0e4dba0, 
        sequence_number: 0, 
        payload: {, 
            transaction: peer_to_peer_transaction, 
            args: [ 
                {ADDRESS: 45d4f95771caa7f66d50ebf3101eeffbda9a8e166b07dacefa6befca0d21cb30},
                {U64: 50000000}, 
            ]
        }, 
        max_gas_amount: 140000, 
        gas_unit_price: 0, 
        expiration_time: 1573286980s, 
    }, 
    public_key: Ed25519PublicKey(
        PublicKey(CompressedEdwardsY: [165, 12, 220, 83, 188, 243, 141, 126, 250, 126, 8, 97, 60, 70, 86, 22, 221, 101, 249, 96, 212, 142, 182, 39, 95, 141, 226, 243, 74, 254, 205, 122]), EdwardsPoint{
            X: FieldElement51([1997212141317716, 1981575496736290, 711217020023999, 241306037824169, 1668563508638375]),
            Y: FieldElement51([1675365069884581, 2124398465077201, 1096863584966936, 1881899434207792, 2160402451349032]),
            Z: FieldElement51([1, 0, 0, 0, 0]),
            T: FieldElement51([881527528983752, 1886896662238058, 23518343360366, 134029746237099, 649488372298501])
        }),
    ), 
    signature: Ed25519Signature(
        Signature( R: CompressedEdwardsY: [242, 238, 89, 236, 29, 145, 40, 81, 65, 250, 196, 10, 125, 216, 15, 220, 125, 180, 200, 116, 215, 67, 133, 109, 232, 8, 213, 2, 129, 223, 105, 112], s: Scalar{
            bytes: [177, 210, 218, 247, 23, 50, 117, 245, 46, 42, 50, 175, 211, 31, 0, 216, 17, 100, 89, 164, 181, 172, 127, 133, 241, 206, 93, 247, 197, 91, 144, 1],
        } ),
    ), 
    }
    Events: 
    ContractEvent { key: 0a9b081a7ac15892ccfb8db6420d9617104d0f9105f511f421b904390b676250, index: 0, type: Struct(StructTag { address: 0000000000000000000000000000000000000000000000000000000000000000, module: Identifier("LibraAccount"), name: Identifier("SentPaymentEvent"), type_params: [] }), event_data: AccountEvent { amount: 50000000, account: 45d4f95771caa7f66d50ebf3101eeffbda9a8e166b07dacefa6befca0d21cb30 } }
    ContractEvent { key: 8e697c610a892d9ede5eda15ff72c5402a5b24987c0654a8ae8a29daf97d7add, index: 0, type: Struct(StructTag { address: 0000000000000000000000000000000000000000000000000000000000000000, module: Identifier("LibraAccount"), name: Identifier("ReceivedPaymentEvent"), type_params: [] }), event_data: AccountEvent { amount: 50000000, account: 38099a5b17c2db81babca08bbcf3399299f6ced16d23c28be77dbc68a0e4dba0 } }

