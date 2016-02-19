// +build !integration

package keyman

import (
	"github.com/cpssd/paranoid/logger"
	"os"
	"testing"
)

func TestDistributionGoodRebuild(t *testing.T) {
	Log = logger.New("keyman", "pfsd", os.DevNull)
	keyBytes := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	key, _ := NewKey(keyBytes)
	pieces, _ := GeneratePieces(key, 5, 3)

	// Since we need at least 3 pieces to reconstruct the key, this should work.
	_, err := RebuildKey(pieces[:3])
	if err != nil {
		t.Log(err.(*FingerMismatchError).ExpectedFingerprint)
		t.Log(err.(*FingerMismatchError).ActualFingerprint)
		t.Fatal(err)
	}
}

func TestDistributionBadRebuild(t *testing.T) {
	keyBytes := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	key, _ := NewKey(keyBytes)
	pieces, _ := GeneratePieces(key, 5, 3)

	// Since we need at least 3 pieces to reconstruct the key, this should not work.
	_, err := RebuildKey(pieces[:2])
	if err == nil {
		t.Fatal("expected fingerprint mismatch")
	}
}
