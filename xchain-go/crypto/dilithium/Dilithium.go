package dilithium

/*
#cgo CFLAGS: -I./
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include "sign.h"
#include "ext.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

var (
	CRYPTO_BYTES = C.ulonglong(1487)

	ErrGenerateKeyFailed = errors.New("生成私钥错误")
	ErrInvalidMsgLen     = errors.New("无效的信息长度,需要32字节")
	ErrInvalidKey        = errors.New("无效的私钥")
	ErrSignFailed        = errors.New("签名失败")
	ErrGetPubicKeyFailed = errors.New("生成公钥失败")
)

type PublicKey []byte
type PrivateKey []byte

func GenerateSk() ([]byte, []byte, error) {
	var pubkey = make([]byte, 896)
	var seckey = make([]byte, 2096)
	pubkeydata := (*C.uchar)(unsafe.Pointer(&pubkey[0]))
	seckeydata := (*C.uchar)(unsafe.Pointer(&seckey[0]))

	if C.genkey(pubkeydata, seckeydata) != 0 {
		return nil, nil, ErrGenerateKeyFailed
	}
	return pubkey, seckey, nil
}

func Sign(msg []byte, seckey []byte) ([]byte, error) {
	if len(msg) != 32 {
		return nil, ErrInvalidMsgLen
	}
	if len(seckey) != 2096 {
		return nil, ErrInvalidKey
	}
	mlen := C.ulonglong(len(msg))
	var smlent uint64
	//smlen := (*C.ulonglong)(unsafe.Pointer(&smlent))
	smlen := C.ulonglong(smlent)
	sm := make([]byte, mlen+CRYPTO_BYTES)
	var smdata = (*C.uchar)(unsafe.Pointer(&sm[0]))
	var msgdata = (*C.uchar)(unsafe.Pointer(&msg[0]))
	var seckeydata = (*C.uchar)(unsafe.Pointer(&seckey[0]))

	if C.chain_sign(smdata, smlen, msgdata, mlen, seckeydata) != 0 {
		return nil, ErrSignFailed
	}
	return sm, nil
}

func VerifySignature(pubkey, msg, signature []byte) bool {
	m1len, m1 := verifySignature(pubkey, msg, signature)
	if int(m1len) != len(msg) {
		return false
	}
	for i := 0; i < len(msg); i++ {
		if msg[i] != m1[i] {
			return false
		}
	}

	return true
}

func verifySignature(pubkey, msg, signature []byte) (uint64, []byte) {
	mlen := C.ulonglong(len(msg))
	smlen := C.ulonglong(len(signature))
	sigdata := (*C.uchar)(unsafe.Pointer(&signature[0]))
	keydata := (*C.uchar)(unsafe.Pointer(&pubkey[0]))
	m1 := make([]byte, mlen+CRYPTO_BYTES)
	var ml = new(uint64)
	mlen1 := (*C.ulonglong)(unsafe.Pointer(ml))
	m1data := (*C.uchar)(unsafe.Pointer(&m1[0]))

	if C.chain_verify_sign(m1data, mlen1, sigdata, smlen, keydata) != 0 {
		return 0, nil
	}

	return *ml, m1
}

func GetPk(sk []byte) ([]byte, error) {
	if len(sk) != 2096 {
		return nil, ErrInvalidKey
	}
	seckeydata := (*C.uchar)(unsafe.Pointer(&sk[0]))

	var pubkey = make([]byte, 896)
	pubkeydata := (*C.uchar)(unsafe.Pointer(&pubkey[0]))

	if C.getpk(seckeydata, pubkeydata) != 0 {
		return nil, ErrGetPubicKeyFailed
	}
	return pubkey, nil
}
