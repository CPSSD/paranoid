package libpfs

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
)

////////////////////// ENCRYPTION ////////////////////////

// CipherBlock Stores the cipher block
var CipherBlock cipher.Block

//TODO: Remove this temporary function
// init is the temporary function until we have proper way of using keys
//func init() {
// Generated using "openssl enc -aes-128-cbc -nosalt -P -md sha1"
// Password: a
// IV: 377667B81B7FEEA3771EADAB99061711
// encryptionKey := []byte("86F7E437FAA5A7FCE15D1DDCB9EAEAEA")
// _ = GenerateAESCipherBlock(encryptionKey)
//}

// GenerateAESCipherBlock generates the cipher used for encryption and decryption of data
// It takes in a byte array key and returns an error if the key is not
// of size 16, 24 or 32 or when the cipher failed to initialize.
func GenerateAESCipherBlock(key []byte) (cipherBlock cipher.Block, err error) {
	switch len(key) {
	case 16, 24, 32:
		break
	default:
		return nil, fmt.Errorf("bad key length (%d)", len(key))
	}
	cipherBlock, err = aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cannot create cipher: %s", err)
	}
	return cipherBlock, nil
}

// CheckIsCipher checks the "status" of the cipher
func CheckIsCipher() bool {
	return CipherBlock != nil
}

// SetCipher sets a cipher to be used for data encyption and decryption
func SetCipher(cb cipher.Block) {
	CipherBlock = cb
}

// Encrypt encrypts the data and returns a bytes.Buffer with the results
func Encrypt(data []byte) (encryptedData bytes.Buffer, length int) {
	i := 0

	cipherBlockSize := CipherBlock.BlockSize()
	encBuf := make([]byte, cipherBlockSize)

	for ; i+cipherBlockSize+1 < len(data); i += cipherBlockSize {
		CipherBlock.Encrypt(encBuf, data[i:i+cipherBlockSize])
		encryptedData.Write(encBuf)
	}
	if i < len(data) {
		raw := make([]byte, cipherBlockSize)
		length = copy(raw, data[i:])
		CipherBlock.Encrypt(encBuf, raw)
		encryptedData.Write(encBuf)
	}

	return encryptedData, length
}

// Decrypt decrypts the data and returns a bytes.Buffer with the results
func Decrypt(data []byte, length int) (dec bytes.Buffer) {
	i := 0

	cipherBlockSize := CipherBlock.BlockSize()
	decBuf := make([]byte, cipherBlockSize)

	for ; i+cipherBlockSize < len(data); i += cipherBlockSize {
		CipherBlock.Decrypt(decBuf, data[i:i+cipherBlockSize])
		dec.Write(decBuf)
	}
	if i < len(data) {
		CipherBlock.Decrypt(decBuf, data[i:i+cipherBlockSize])
		dec.Write(decBuf[:length])
	}

	return dec
}

// LastBlockSize reads the size of the last block from the beginning of the file
func LastBlockSize(r *os.File) (size int, err error) {
	buf := []byte{byte(0)}
	_, err = r.Read(buf)
	if err != nil {
		return 0, err
	}
	return int(buf[0]), nil
}

// GetLastBlock gets the last block of the file
func GetLastBlock(r *os.File) (data []byte, err error) {
	buf := make([]byte, CipherBlock.BlockSize())
	stats, err := r.Stat()
	if err != nil {
		return buf, err
	}
	size := stats.Size()
	_, err = r.ReadAt(buf, size-int64(len(buf)))
	if err == io.EOF {
		return buf, nil
	}
	return buf, err
}
