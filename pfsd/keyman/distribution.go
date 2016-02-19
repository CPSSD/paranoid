// Functions for deconstructing and reconstructing an AES key.
// Written according to http://cs.jhu.edu/~sdoshi/crypto/papers/shamirturing.pdf
//
// See docs/key_management_v0.1.md for more information.

package keyman

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
)

const PrimeSize int = 320 // 40 bytes

type KeyPiece struct {
	Data              []byte
	ParentFingerprint [32]byte // The SHA-256 fingerprint of the key it was generated from.
	Prime             *big.Int // The prime number used in the generation of this KeyPiece.
	Seq               int64    // Where f(Seq) = Data, for some polynomial f
}

type KeyPieceSorter []*KeyPiece

func (s KeyPieceSorter) Len() int {
	return len(s)
}

func (s KeyPieceSorter) Less(i, j int) bool {
	return s[i].Seq < s[j].Seq
}

func (s KeyPieceSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type FingerMismatchError struct {
	ExpectedFingerprint [32]byte
	ActualFingerprint   [32]byte
}

func (e *FingerMismatchError) Error() string {
	return fmt.Sprintf("key fingerprint does not match keypiece fingerprint")
}

// Generates numPieces KeyPieces from key. requiredPieces is the number of KeyPieces
// needed to reconstruct the original key.
func GeneratePieces(key *Key, numPieces, requiredPieces int64) ([]*KeyPiece, error) {
	// Some input validation
	if requiredPieces > numPieces {
		return nil, errors.New("requiredPieces cannot be longer than numPieces")
	}
	if requiredPieces <= 0 || numPieces <= 0 {
		return nil, errors.New("requiredPieces or numPieces cannot be less than or equal to 0")
	}

	var keyBig big.Int
	keyBig.SetBytes(key.GetBytes())
	prime, _ := rand.Prime(rand.Reader, PrimeSize)
	coefficients := make([]*big.Int, requiredPieces)
	coefficients[0] = &keyBig
	// I can't believe I have to cast the number 1 to a 64-bit integer ...
	for i := int64(1); i < requiredPieces; i++ {
		tmp, _ := rand.Int(rand.Reader, prime)
		coefficients[i] = tmp
	}

	pieces := make([]*KeyPiece, numPieces)
	for x := int64(1); x <= numPieces; x++ {
		total := big.NewInt(0)
		for i := int64(0); i < requiredPieces; i++ {
			// This computes (coefficient)(x^i)
			degreeTotal := new(big.Int).Mul(coefficients[i], new(big.Int).Exp(big.NewInt(x), big.NewInt(i), nil))
			total.Add(total, degreeTotal)
		}
		total.Mod(total, prime)
		pieces[x-1] = &KeyPiece{
			Data:              total.Bytes(),
			ParentFingerprint: key.GetFingerprint(),
			Prime:             prime,
			Seq:               x,
		}
	}
	return pieces, nil
}

// Rebuild a Key from a set of KeyPieces. This function will succeed iff
// len(pieces) >= requiredPieces from the Generate function.
func RebuildKey(pieces []*KeyPiece) (*Key, error) {
	fingerprint := pieces[0].ParentFingerprint
	for _, v := range pieces {
		if v.ParentFingerprint != fingerprint {
			return nil, errors.New("not all pieces come from the same key")
		}
	}
	prime := pieces[0].Prime
	inputs := make([]int64, len(pieces))
	outputs := make([]*big.Int, len(pieces))
	for i, v := range pieces {
		inputs[i] = v.Seq
		var tmp big.Int
		tmp.SetBytes(v.Data)
		outputs[i] = &tmp
	}

	// Use Lagrange interpolation to find out the key, which is f(0)
	sum := big.NewInt(0)
	for i := int64(0); i < int64(len(pieces)); i++ {
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)
		for j := int64(0); j < int64(len(pieces)); j++ {
			if j != i {
				// Go doesn't let you chain big.Int operations, so we have to have
				// several separate statements instead.
				numerator.Mul(numerator, big.NewInt(-j-1))
				numerator.Mod(numerator, prime)
				denominator.Mul(denominator, big.NewInt(i-j))
				denominator.Mod(denominator, prime)
			}
		}
		output := outputs[i]
		tmp := new(big.Int).Mul(output, numerator)
		tmp.Mul(tmp, denominator.ModInverse(denominator, prime))
		tmp.Mod(tmp, prime)
		sum.Add(sum, prime)
		sum.Add(sum, tmp)
		sum.Mod(sum, prime)
	}
	keyBytes := sum.Bytes()
	keyFingerprint := sha256.Sum256(keyBytes)
	if keyFingerprint != pieces[0].ParentFingerprint {
		// Even if the key is wrong, we return it, for debugging purposes.
		return &Key{keyBytes, keyFingerprint}, &FingerMismatchError{pieces[0].ParentFingerprint, keyFingerprint}
	}
	return &Key{
		bytes:       keyBytes,
		fingerprint: keyFingerprint,
	}, nil
}
