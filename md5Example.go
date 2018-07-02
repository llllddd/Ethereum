package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//将初始信息分组为512bit每组，每32bit经过一次计算，共有4轮，每轮16次计算，一轮处理512bit
const (
	init0 = 0x67452301
	init1 = 0xEFCDAB89
	init2 = 0x98BADCFE
	init3 = 0x10325476
)
const (
	S11 = 7
	S12 = 12
	S13 = 17
	S14 = 22
	S21 = 5
	S22 = 9
	S23 = 14
	S24 = 20
	S31 = 4
	S32 = 11
	S33 = 16
	S34 = 23
	S41 = 6
	S42 = 10
	S43 = 15
	S44 = 21
)

var m_Buffer = make([]byte, 64)
var m_Result = make([]byte, 16)
var m_ChainingVariable = make([]uint32, 4)
var m_Count = make([]uint32, 2)

//定义四个非线性函数
func F(x, y, z uint32) uint32 {
	return (x & y) | ((^x) & z)
}

func G(x, y, z uint32) uint32 {
	return (x & z) | (y & (^z))
}

func H(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func I(x, y, z uint32) uint32 {
	return y ^ (x | (^z))
}

//初始化分组变量

func Initialize() {
	m_Count[0], m_Count[1] = 0, 0
	m_ChainingVariable[0] = init0
	m_ChainingVariable[1] = init1
	m_ChainingVariable[2] = init2
	m_ChainingVariable[3] = init3
}

//循环左移
func LeftRotate(opNumber uint32, opBit uint32) uint32 {
	left := opNumber
	right := opNumber
	return (left << opBit) | (right >> (32 - opBit))
}

//每次计算循环移位的位数以及加的常数都不同总体可以分为四种
func FF(a *uint32, b, c, d uint32, Mj, s, Ti uint32) {
	temp := *a + F(b, c, d) + Mj + Ti
	*a = b + LeftRotate(temp, s)
}

func GG(a *uint32, b, c, d uint32, Mj, s, Ti uint32) {
	temp := *a + G(b, c, d) + Mj + Ti
	*a = b + LeftRotate(temp, s)
}

func HH(a *uint32, b, c, d uint32, Mj, s, Ti uint32) {
	temp := *a + H(b, c, d) + Mj + Ti
	*a = b + LeftRotate(temp, s)
}

func II(a *uint32, b, c, d uint32, Mj, s, Ti uint32) {
	temp := *a + I(b, c, d) + Mj + Ti
	*a = b + LeftRotate(temp, s)
}

//4byte转1uint32
func byteToInt(source []byte) []uint32 {
	var destination []uint32
	var temp uint32
	for i := 0; i < len(source)/4; i++ {

		bytesBuffer := bytes.NewBuffer(source[i*4 : (i*4 + 4)])
		binary.Read(bytesBuffer, binary.LittleEndian, &temp)
		destination = append(destination, temp)
	}
	return destination
}

//1uint32转byte
func intToByte(source []uint32) []byte {
	var result [][]byte
	for i := 0; i < len(source); i++ {
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.LittleEndian, source[i])
		result = append(result, bytesBuffer.Bytes())
	}

	return bytes.Join(result, []byte(""))
}

//md5计算过程
func ProcessOfMd5(groups []byte) {
	a := m_ChainingVariable[0]
	b := m_ChainingVariable[1]
	c := m_ChainingVariable[2]
	d := m_ChainingVariable[3]
	M := byteToInt(groups)

	FF(&a, b, c, d, M[0], S11, 0xd76aa478)
	FF(&d, a, b, c, M[1], S12, 0xe8c7b756)
	FF(&c, d, a, b, M[2], S13, 0x242070db)
	FF(&b, c, d, a, M[3], S14, 0xc1bdceee)
	FF(&a, b, c, d, M[4], S11, 0xf57c0faf)
	FF(&d, a, b, c, M[5], S12, 0x4787c62a)
	FF(&c, d, a, b, M[6], S13, 0xa8304613)
	FF(&b, c, d, a, M[7], S14, 0xfd469501)
	FF(&a, b, c, d, M[8], S11, 0x698098d8)
	FF(&d, a, b, c, M[9], S12, 0x8b44f7af)
	FF(&c, d, a, b, M[10], S13, 0xffff5bb1)
	FF(&b, c, d, a, M[11], S14, 0x895cd7be)
	FF(&a, b, c, d, M[12], S11, 0x6b901122)
	FF(&d, a, b, c, M[13], S12, 0xfd987193)
	FF(&c, d, a, b, M[14], S13, 0xa679438e)
	FF(&b, c, d, a, M[15], S14, 0x49b40821)

	GG(&a, b, c, d, M[1], S21, 0xf61e2562)
	GG(&d, a, b, c, M[6], S22, 0xc040b340)
	GG(&c, d, a, b, M[11], S23, 0x265e5a51)
	GG(&b, c, d, a, M[0], S24, 0xe9b6c7aa)
	GG(&a, b, c, d, M[5], S21, 0xd62f105d)
	GG(&d, a, b, c, M[10], S22, 0x2441453)
	GG(&c, d, a, b, M[15], S23, 0xd8a1e681)
	GG(&b, c, d, a, M[4], S24, 0xe7d3fbc8)
	GG(&a, b, c, d, M[9], S21, 0x21e1cde6)
	GG(&d, a, b, c, M[14], S22, 0xc33707d6)
	GG(&c, d, a, b, M[3], S23, 0xf4d50d87)
	GG(&b, c, d, a, M[8], S24, 0x455a14ed)
	GG(&a, b, c, d, M[13], S21, 0xa9e3e905)
	GG(&d, a, b, c, M[2], S22, 0xfcefa3f8)
	GG(&c, d, a, b, M[7], S23, 0x676f02d9)
	GG(&b, c, d, a, M[12], S24, 0x8d2a4c8a)

	HH(&a, b, c, d, M[5], S31, 0xfffa3942)
	HH(&d, a, b, c, M[8], S32, 0x8771f681)
	HH(&c, d, a, b, M[11], S33, 0x6d9d6122)
	HH(&b, c, d, a, M[14], S34, 0xfde5380c)
	HH(&a, b, c, d, M[1], S31, 0xa4beea44)
	HH(&d, a, b, c, M[4], S32, 0x4bdecfa9)
	HH(&c, d, a, b, M[7], S33, 0xf6bb4b60)
	HH(&b, c, d, a, M[10], S34, 0xbebfbc70)
	HH(&a, b, c, d, M[13], S31, 0x289b7ec6)
	HH(&d, a, b, c, M[0], S32, 0xeaa127fa)
	HH(&c, d, a, b, M[3], S33, 0xd4ef3085)
	HH(&b, c, d, a, M[6], S34, 0x4881d05)
	HH(&a, b, c, d, M[9], S31, 0xd9d4d039)
	HH(&d, a, b, c, M[12], S32, 0xe6db99e5)
	HH(&c, d, a, b, M[15], S33, 0x1fa27cf8)
	HH(&b, c, d, a, M[2], S34, 0xc4ac5665)

	II(&a, b, c, d, M[0], S41, 0xf4292244)
	II(&d, a, b, c, M[7], S42, 0x432aff97)
	II(&c, d, a, b, M[14], S43, 0xab9423a7)
	II(&b, c, d, a, M[5], S44, 0xfc93a039)
	II(&a, b, c, d, M[12], S41, 0x655b59c3)
	II(&d, a, b, c, M[3], S42, 0x8f0ccc92)
	II(&c, d, a, b, M[10], S43, 0xffeff47d)
	II(&b, c, d, a, M[1], S44, 0x85845dd1)
	II(&a, b, c, d, M[8], S41, 0x6fa87e4f)
	II(&d, a, b, c, M[15], S42, 0xfe2ce6e0)
	II(&c, d, a, b, M[6], S43, 0xa3014314)
	II(&b, c, d, a, M[13], S44, 0x4e0811a1)
	II(&a, b, c, d, M[4], S41, 0xf7537e82)
	II(&d, a, b, c, M[11], S42, 0xbd3af235)
	II(&c, d, a, b, M[2], S43, 0x2ad7d2bb)
	II(&b, c, d, a, M[9], S44, 0xeb86d391)

	m_ChainingVariable[0] += a
	m_ChainingVariable[1] += b
	m_ChainingVariable[2] += c
	m_ChainingVariable[3] += d
}

func main() {
	test := []byte{0x01, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f}
//	test := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	fmt.Println("文本长度: ",len(test))
	Initialize()//初始分组
	dest := byteToInt(test)
	fmt.Println("转换为int32: ",dest)

	ProcessOfMd5(test)
	aa := intToByte(m_ChainingVariable)
	//fmt.Println(aa)
	fmt.Printf("128位MD5: %x", aa)
}

