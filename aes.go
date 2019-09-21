package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

//加密算法(AES+CBC),结果转换成16进制，防止数据库出错
func (*DuduAes) Encrypt(content string, securityKey string) (string,error) {
	key := []byte(securityKey)
	result, err := AesEncrypt([]byte(content), key)
	if err != nil {
		//panic(err)
		return "",err
	}
	return hex.EncodeToString(result),nil
}

//解密算法(AES+CBC)
func (*DuduAes) Decrypt(content string, securityKey string) (string,error) {
	key := []byte(securityKey)

	tmpString,_ := hex.DecodeString(content)
	result,err := AesDecrypt(tmpString, key)
	if err != nil {
		//panic(err)
		return "",err
	}

	return string(result),nil
}

//
//func testAes() {
//	// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
//	key := []byte("hundsun@12345678")
//	result, err := AesEncrypt([]byte("{\"username\":\"北京昌平区\",\"Age\":28,\"Name\":\"liang637210\"}"), key)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(hex.EncodeToString(result))
//	fmt.Println(base64.StdEncoding.EncodeToString(result))
//	origData, err := AesDecrypt(result, key)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(string(origData))
//}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData

	// 长度不能小于aes.Blocksize
	if len(crypted) < aes.BlockSize {
		return nil, errors.New("crypto/cipher: ciphertext too short")
	}

	// 必须为aes.Blocksize的倍数
	if len(crypted)%aes.BlockSize != 0 {
		return nil, errors.New("crypto/cipher: ciphertext is not a multiple of the block size")
	}

	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted

	// 长度不能小于aes.Blocksize
	if len(origData) < aes.BlockSize {
		return nil, errors.New("crypto/cipher: ciphertext too short")
	}

	// 必须为aes.Blocksize的倍数
	if len(origData)%aes.BlockSize != 0 {
		return nil, errors.New("crypto/cipher: ciphertext is not a multiple of the block size")
	}

	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}


func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

