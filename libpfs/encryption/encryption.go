package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"os"
)

var (
	Encrypted   bool
	cipherBlock cipher.Block
)

// GenerateAEScipherBlock generates the cipher used for encryption and decryption of data
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

// SetCipher sets a cipher to be used for data encyption and decryption
func SetCipher(cb cipher.Block) {
	cipherBlock = cb
}

func GetCipherSize() int {
	return cipherBlock.BlockSize()
}

// Encrypt encrypts the given data in place
func Encrypt(data []byte) error {
	cipherBlockSize := cipherBlock.BlockSize()
	if len(data)%cipherBlockSize != 0 {
		return errors.New("can not encrypt data not of size n * blocksize")
	}

	for i := 0; i < len(data); i += cipherBlockSize {
		cipherBlock.Encrypt(data[i:i+cipherBlockSize], data[i:i+cipherBlockSize])
	}
	return nil
}

// Decrypt decrypts the given data in place
func Decrypt(data []byte) error {
	cipherBlockSize := cipherBlock.BlockSize()
	if len(data)%cipherBlockSize != 0 {
		return errors.New("can not decrypt data not of size n*blocksize")
	}

	for i := 0; i < len(data); i += cipherBlockSize {
		cipherBlock.Decrypt(data[i:i+cipherBlockSize], data[i:i+cipherBlockSize])
	}
	return nil
}

// LastBlockSize reads the size of the last block from the beginning of the file
func LastBlockSize(r *os.File) (size int, err error) {
	buf := []byte{byte(255)}
	n, err := r.ReadAt(buf, 0)
	if err != nil {
		return 0, fmt.Errorf("error getting last block: %s", err)
	}
	if n != 1 {
		return 0, errors.New("error getting last block")
	}
	return int(buf[0]), nil
}
