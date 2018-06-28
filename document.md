## 以太坊账户地址

以太坊有两种不同的账户类型:合约账户和EOA(Externally Owned Accounts)外部账户。账户的所有权是由密钥来验证的。密钥由用户存储在文件或数据库中也就是钱包。用户钱包中的数字密钥完全独立于以太坊协议，并且可以由用户的钱包软件生成和管理而并不需要依照区块链或接入互联网。而以太坊的交易需要一个有效的数字签名被包含在区块链中，它只能由一个账户的私钥生成；因此，任何一个具有该密钥副本的人都有对账户的控制。交易中的数字签名证明了资金的真正所有者。

账户中的密钥包含两部分：公钥和私钥,公钥可以暴露给其他人，而私钥只有账户的所有者才能知道这些密钥通常由钱包管理。

在以太坊的交易支付部分中，预期的接收方由一个以太坊地址表示，大多数情况下，地址有账户的公钥生成，对应于公钥，但是并不是所有的以太坊地址都表示公钥，也可以表示合约。以太坊地址是用户接触到的密钥的唯一表示。

## 椭圆曲线加密

公钥加密方案基于困难的数学难题，比如著名的RSA算法是基于大整数的分解难题。以太坊中的公钥加密方案使用的是椭圆曲线加密方案，它是基于椭圆曲线上的离散对数问题设计的一种公钥加密方案。具体来看一下椭圆曲线加密方案。

首先椭圆曲线上的所有点做成一个有限域，在曲线上进行乘法运算是简单的，但是做除法是困难的。
因此在进行加密是我们必须有一个选好的椭圆曲线，以太坊中的椭圆曲线和比特币中的一样为secp256k1，实际上以太坊中调用的就是比特币的加密库。曲线为<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;y^2=x^3&plus;7" title="y^2=x^3+7" />，
椭圆曲线上的计算实际上是在有限域上进行计算，对于一个有限域<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;F_{p}" title="F_{p}" />，其上的计算结果都要模素数<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;p=2^{256}-2^{32}-2^{9}-2^{8}-2^{7}-2^{6}-2^{4}-1" title="p=2^{256}-2^{32}-2^{9}-2^{8}-2^{7}-2^{6}-2^{4}-1" /> 。并且所有的运算结果同样会在曲线上。椭圆曲线上的所有点由一个基点<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;G" title="G" />得到，我们假设无穷远处的点为零点<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;O" title="O" /> 则基点的阶为<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;n" title="n" />满足<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;n*G=O" title="n*G=O" />我们同时定义椭圆曲线上的加法和乘法

加法：连接曲线上两点<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;p_{1}" title="p_{1}" />，<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;p_{2}" title="p_{2}" /> 交曲线于 <img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;p_{3}^{'}(x,y)" title="p_{3}^{'}(x,y)" /> ，则<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;p_{1}&plus;p_{2}=p_{3}(x,y)" title="p_{1}+p_{2}=p_{3}(x,y)" />

乘法:<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;k*P=P&plus;P&plus;\cdots&plus;P" title="k*P=P+P+\cdots+P" />

### 密钥的生成
以太坊账户私钥的生成，即产生一个256bit的随机数，通常是由一个hash算法得到比如Hash256或Keaccak256。以太坊的公钥实际上是椭圆曲线上的一个点<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;(x,y)" title="(x,y)" />满足椭圆曲线方程。公钥是对私钥做椭圆曲线乘法得到的：<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;PK=sk*G" title="PK=sk*G" />， <img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;PK" title="PK" />和<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;G" title="G" />都是椭圆曲线上的点，<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;G"/>是椭圆曲线的基点。公钥可以由私钥得到的意味着仅通过私钥我们就可以计算出公钥，但是通过想要公钥得到私钥是困难的因为椭圆曲线上的除法难于计算。反过来得到<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;sk"/>的除法计算我们称为椭圆曲线上的离散对数计算，这是一个困难问题。在以太坊中，可以看到公钥表示为66个十六进制的字符即3
3bytes，以太坊中只使用了未压缩的公钥因此与公钥相关的前缀为0x04.

### 椭圆曲线密码加解密步骤：

1. 生成私钥<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;sk\leftarrow&space;[1,n-1]" title="sk\leftarrow [1,n-1]" />，公钥为<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;PK=sk*G">

2. 加密：选择随机数<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;r\leftarrow&space;[1,n-1]" title="r\leftarrow [1,n-1]" />，明文信息<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;M" title="M" />

3. 通过公钥计算得到密文<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;C_x=r*G,C_y=M&plus;r*PK" title="C_x=r*G,C_y=M+r*PK" /> 密文为 <img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;C=(C_x,C_y)" title="C=(C_x,C_y)" />

4. 解密：由密文通过私钥<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;sk" title="sk" />可以计算得到明文<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;M=C_{y}-sk*C_{x}=M+r*PK-sk*(rG)">


## 椭圆曲线数字签名算法
以太坊中使用的数字签名算法为椭圆曲线数字签名算法(Elliptic Curve Digital Signature Algorithm,ECDSA),椭圆曲线签名算法同样依赖于生成的密钥对<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;(sk,PK)">.以太坊中的数字签名有三个目的：首先，签名证明了私钥的所有者，它是隐含的帐户的所有者，表明了以太币支出，或合同的执行是经过授权的。第二，证明了授权是不可抵赖的。第三，防止交易数据在交易完成后被篡改。同时，数字签名的数学方案有两部分组成
### 计算签名：
设私钥为<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;sk">，公钥为<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;PK=sk*G">，椭圆曲线参数 <img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;(Curve,G,n),n*G=0">

使用私钥进行签名

1. 计算，<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;e=HASH(m)">通过hash算法计算消息m的hash值，例如sha2，以太坊中keccak256s算法
2. 取hash值得最高有效位起的<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;\log{n}">bits，得到<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;z">
3. 选取随机数<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;k\leftarrow[1,n-1]">，计算<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;M=k*G">，为椭圆曲线上的一点<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;(x,y)">.
4. 计算<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;r=x\mod{n}">若<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;r=0">返回第三步
5. 计算<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;s=k^{-1}(z+r*sk)\mod{n}">，若<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;s=0">,返回第三步
6. 得到签名<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;(r,s)">

### 公钥验证签名    
验证算法为上述签名算法的逆过程，使用签名<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;(r,s)">和公钥<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;PK">计算得到曲线上的一个点<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;Q">即验证通过。

1. 判断签名<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;(r,s)\in[1,n-1]">，计算得到消息的hash值<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;e=HASH(m)">取hash值得最高有效位起的<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;\log{n}">bits，得到<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;z">
2. 计算参数：<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;w=s^{-1}\mod{n}">
3. 得到若点<img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;Q=z*s^{-1}+r*s^{-1}*PK">在椭圆曲线上则签名验证成功
4. <img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;r=Q_{x}">则签名有效

## 哈希算法
哈希函数的输入称为预图像或者消息，输出被称为哈希或摘要，它们大多数依据分组密码的思路，将原消息压缩和混淆。它将任意长度的字符串映射为固定长度的bit串，密码学hash函数是一类单向函数，即通过输入我们很容易能得到输出，但是由输出则很难得到输入的信息，比如同余方程。哈希函数由以下五种主要的性质：

* 确定性：相同的消息总是得到相同的哈希值
* 高效性：任意给定的消息，哈希算法可以快速得到其哈希值
* 不可逆：有限时间类，由哈希值计算得到原消息是不可能的
* 雪崩性：原消息微小的改动即会导值所得到的哈希值得巨大差异
* 低碰撞：两个不同的消息得到同样的哈希值得概率是非常低的

上述性质使得哈希函数在以太坊中有着广泛的应用包括数据签名，消息完整性验证，工作量证明，随机数生成等等。
目前通用的哈希算法有四代标准，SHA1有名的MD5，SHA2的 SHA256，SHA512算法，及以太坊中使用的SHA3标准的Keacck256算法。

## 账户地址的生成

以太坊通过椭圆曲线加密算法来实现对交易的签名与验证，路径github.com/ethereum/go-ethereum/crypto/下的代码包负责所有与加密相关的操作。
 
 以太坊中的地址是使用单向hash函数Keacck-256从公钥或合约中派生出来的唯一标识符.
 首先公钥是从椭圆曲线加密算法中得到64bytes字符串，以太坊中使用的是未压缩的公钥包含椭圆曲线上的点的所有信息，公钥的前缀是04，最终的公钥为：
 04 + X-coordinate (32 bytes/64 hex) + Y coordinate (32 bytes/64 hex)，公钥由私钥导出.

 由唯一的私钥得到与之对应的公钥后，使用Keacck-256算法计算公钥的hash，最后只保留最后的20bytes(大端表示)作为以太坊的地址，通常地址是由16进制表示因此有前缀0x。

### 相关代码

go语言包中自带的crypto/ecdsa包中包含了关于椭圆曲线的结构体声明和操作的函数，以太坊也是通过调用它来生成账户的私钥并产生公钥的.
ECDSA的公钥结构体，通过一个elliptic.Curve接口的实现体来提供椭圆曲线的所有属性和相关操作；

公钥的成员(X,Y)即为生成的未压缩的公钥.
私钥是以太坊账户中存储的唯一可以用来验证账户身份的信息，但实际上它也包含有公钥的结构体，D是算法生成的私钥，根据不同的用途可以使用结构体PrivateKey或PrivateKey.D

```go
type PublicKey struct {
    elliptic.Curve
    X, Y *big.Int
}

type PrivateKey struct {
    PublicKey
	D *big.Int
}
```

GenerateKey函数生成PrivateKey类型的私钥，其中也包含了用来生成账户地址的的公钥信息。

```go
func GenerateKey(c elliptic.Curve, rand io.Reader) (*PrivateKey, error) 
```

在生成地址之前需要将publicKey字符串类型和ecdsa.PublicKey{}类型格式转换，在github.com/ethereum/go-ethereum/crypto/crypto.go中定义了相关的转换函数.
代码实际上是完成了big.int 类型到 []byte的转换，在实际调用时要注意返回的[]byte类型字符串是由三部分组成，S256()返回的是基于secp256k1椭圆曲线参数的接口的实现类.

```go
func ToECDSAPub(pub []byte) *ecdsa.PublicKey {  
    x, y := elliptic.Unmarshall(S256(), pub)  
    return &ecdsa.PublicKey{Curve:S256(), X:x, Y:y}  
} 

func FromECDSAPub(pub *ecdsa.PublicKey) []byte {  
    return elliptic.Marshall(S256(), pub.X, pub.Y)  
} 
```

在生成地址之前需要得到公钥字符串的hash值，以太坊中使用两个自定义类型来表示32bytes的hash值和20bytes的地址定义在github.com/ethereum/go-ethereum/common/types.go中

```go
const (  
    HashLength = 32  
    AddressLength = 20  
)  
type Hash [HashLength]byte  
type Address [AddressLength]byte 
```

在github.com/ethereum/go-ethereum/crypto/crypto.go中的Keccak256和Keccak256Hash定义了得到字符串hash值的算法。需要注意的是函数中的参数是公钥的字符串而通过前面格式转换函数得到的公钥字符串中pubkey[0]并不是公钥所以只需要传入pubkey[1:]来得到对应的hash值

```go
func Keccak512(data ...[]byte) []byte 
func Keccak256Hash(data ...[]byte) (h common.Hash) 
```

最后将大端表示的公钥的hash值取最后的20bytes作为账户的地址，这两个函数实际上是用来验证字符串是否满足账户地址的格式，函数位置在github.com/ethereum/go-ethereum/common
```go
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}
```
## 交易的签名

以太坊中的每个交易在被放进区块中时都经过椭圆曲线签名算法进行数字签名

* 使用私钥（签名秘钥）从交易信息中创建签名 
* 允许任何人使用交易信息和公钥验证签名   
	
在对以太坊中的交易进行签名时首先要得到交易的hash值，这里使用的和生成账户地址时是同样的Keacck256函数，不同的是将交易进行序列化时用的是
	以太坊自己定义的RLP编码github.com/ethereum/go-ethereum/rlp下定义了所有与rlp编码的有关操作。
	对一个以太坊中的交易进行签名包含以下步骤：
	
1. 创建交易，完整的交易数据结构包含九个部分: nonce, gasPrice, startGas, to, value, data, v, r, s
	
2. 生成交易的RLP编码的序列化信息
	
3. 计算Keaccak256 hash值
	
4. 使用私钥对hash值签名 <img src="https://latex.codecogs.com/png.latex?\inline&space;\dpi{100}&space;Sig=F_{sig}(F_{keccak256}(m),sk), Sig = (R,S)">

### 相关代码

在github.com/ethereum/go-ethereum/core/types中定义了交易的结构体类型，并且提供了新建交易的函数接口
一个完整的交易必须包含有转入方地址，转账金额，以及每个交易的独立的gasprice和gaslimit，签名[R||S||V]初始为0.

```go
type Transaction struct {
	data txdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

type txdata struct {
	AccountNonce uint64          
	Price        *big.Int        
	GasLimit     uint64          
	Recipient    *common.Address 
	Amount       *big.Int        
	Payload      []byte         
	// 签名值
	V *big.Int
	R *big.Int
	S *big.Int
	Hash *common.Hash
}

func NewTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction 
```
以太坊中对需要序列化的信息设计了rlp编码的规则，github.com/ethereum/go-ethereum/rlp中定义了所有与rlp编码有关的函数，在对交易序列化的过程中
实际上是完成types.Transaction{}结构体到[]byte类型的转换.
```go
//types/encode.go
func EncodeToBytes(val interface{}) ([]byte, error) {
	eb := encbufPool.Get().(*encbuf)
	defer encbufPool.Put(eb)
	eb.reset()
	if err := eb.encode(val); err != nil {
		return nil, err
	}
	return eb.toBytes(), nil
}
```
得到序列化的交易信息的hash值后，再调用椭圆曲线交易签名算法。对交易签名，需要账户的私钥以及交易的hash值。
以太坊对bitcoin的secp256k1 C库进行了封装，代码在github.com/ethereum/go-ethereum/crypto/secp256k1,并且定了签名函数。
以太坊中的数字签名在计算过程所生成的签名, 是一个长度为65bytes的字节数组，
它被截成三段放进交易中，前32bytes赋值给成员变量R, 中间32bytes赋值给S，最后1byte赋给V，由于R、S、V声明的类型都是*big.Int, 上述赋值存在[]byte 到big.Int的类型转换
V表示的字节是用来再恢复公钥时加速运算的，因为恢复公钥时会得到两种不同的结果通过V的奇偶性来判断哪一个才是正确的结果，提高了运算效率

```go
func Sign(msg []byte, seckey []byte) ([]byte, error) {
	if len(msg) != 32 {
		return nil, ErrInvalidMsgLen
	}
	if len(seckey) != 32 {
		return nil, ErrInvalidKey
	}
	seckeydata := (*C.uchar)(unsafe.Pointer(&seckey[0]))
	if C.secp256k1_ec_seckey_verify(context, seckeydata) != 1 {
		return nil, ErrInvalidKey
	}

	var (
		msgdata   = (*C.uchar)(unsafe.Pointer(&msg[0]))
		noncefunc = C.secp256k1_nonce_function_rfc6979
		sigstruct C.secp256k1_ecdsa_recoverable_signature
	)
	if C.secp256k1_ecdsa_sign_recoverable(context, &sigstruct, msgdata, seckeydata, noncefunc, nil) == 0 {
		return nil, ErrSignFailed
	}

	var (
		sig     = make([]byte, 65)
		sigdata = (*C.uchar)(unsafe.Pointer(&sig[0]))
		recid   C.int
	)
	C.secp256k1_ecdsa_recoverable_signature_serialize_compact(context, sigdata, &recid, &sigstruct)
	sig[64] = byte(recid) // add back recid to get 65 bytes sig
	return sig, nil
}
```

以太坊也定义了直接生成交易的数字签名的接口 在github.com/ethereum/go-ethereum/core/types/transaction_signing.go中的SignTx函数
其本质上也是调用上述secp256k1包中的函数来来完成椭圆曲线数字签名的生成
```go
func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {  
    h := s.Hash(tx)  
    sig, err := crypto.Sign(h[:], prv)  
    if err != nil {  
        return nil, err  
    }  
    return tx.WithSignature(s, sig)  
} 
// /crypto/signature_cgo.go  
func Sign(hash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {  
    if len(hash) != 32 {  
        return nil, fmt.Errorf(...)
    }  
    seckey := math.PaddedBigBytes(prv.D, n:prv.Params().BitSize/8)  
    defer zeroBytes(seckey)  
    return secp256k1.Sign(hash, seckey)  
} 
```

## 交易签名的验证

为了验证签名，必须有签名（R和S）、序列化的交易信息和公钥（对应于用于创建签名的私钥）。实际上，签名的验证是为了证明只有生成该公钥的私钥的所有者才能在该交易上上产生这个签名。
签名验证算法采用交易的hash值、签名者的公钥和签名（R和S值），如果签名对该交易和公钥有效，则返回true。通常交易的签名还包含第三个值V这是为了在恢复公钥时简化计算，提升运算效率。

### 相关代码

以太坊对bitcoin的secp256k1 C库进行了封装，代码在github.com/ethereum/go-ethereum/crypto/secp256k1,定义了恢复公钥和验证签名的函数
 RecoverPubkey函数返回签名的公钥. msg为32bytes的交易的hash值，sig为65bytes签名[R||S||V]
```go
func RecoverPubkey(msg []byte, sig []byte) ([]byte, error) {
	if len(msg) != 32 {
		return nil, ErrInvalidMsgLen
	}
	if err := checkSignature(sig); err != nil {
		return nil, err
	}

	var (
		pubkey  = make([]byte, 65)
		sigdata = (*C.uchar)(unsafe.Pointer(&sig[0]))
		msgdata = (*C.uchar)(unsafe.Pointer(&msg[0]))
	)
	if C.secp256k1_ext_ecdsa_recover(context, (*C.uchar)(unsafe.Pointer(&pubkey[0])), sigdata, msgdata) == 0 {
		return nil, ErrRecoverFailed
	}
	return pubkey, nil
}
```
VerifySignature用来验证对应公钥的签名，签名为[R || S] 格式.
```go
func VerifySignature(pubkey, msg, signature []byte) bool {
	if len(msg) != 32 || len(signature) != 64 || len(pubkey) == 0 {
		return false
	}
	sigdata := (*C.uchar)(unsafe.Pointer(&signature[0]))
	msgdata := (*C.uchar)(unsafe.Pointer(&msg[0]))
	keydata := (*C.uchar)(unsafe.Pointer(&pubkey[0]))
	return C.secp256k1_ext_ecdsa_verify(context, sigdata, msgdata, keydata, C.size_t(len(pubkey))) != 0
}
```
以太坊中也定义了直接验证数字签名是否有效的函数，在github.com/ethereum/go-ethereum/crypto包中，本质上也是调用上面的两个函数来完成数字签名的验证
```go
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
		return false
	}
	if homestead && s.Cmp(secp256k1halfN) > 0 {
		return false
	}
	return r.Cmp(secp256k1N) < 0 && s.Cmp(secp256k1N) < 0 && (v == 0 || v == 1)
}
```
同样以太坊也定义了数字签名恢复交易地址(公钥)的接口 在github.com/ethereum/go-ethereum/core/types/transaction_signing.go中
```go
func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error)
```
