package main

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"time"
)

const unlockQueryInterval time.Duration = time.Second * 30
const lockWaitDuration time.Duration = time.Minute * 1

// Chunks key and sends the pieces to other nodes on the network.
func Lock() error {
	numPieces := int64(len(globals.Nodes.GetAll()) + 1)
	requiredPieces := numPieces/2 + 1
	log.Info("Generating pieces.")
	pieces, err := keyman.GeneratePieces(globals.EncryptionKey, numPieces, requiredPieces)
	if err != nil {
		log.Error("Could not chunk key:", err)
		return fmt.Errorf("could not chunk key: %s", err)
	}
	// We always keep the first piece and distribute the rest
	globals.HeldKeyPieces[globals.ThisNode] = pieces[0]

	for i := 1; i < len(pieces); i++ {
		pnetclient.SendKeyPiece(pieces[i])
	}
	// Delete our copy of the full key
	globals.EncryptionKey = nil
	globals.SystemLocked = true

	return nil
}

// Get our keypieces from the other nodes and rebuild our key
func Unlock() error {
	log.Info("Attempting to unlock.")
	pieces, err := pnetclient.RequestKeyPieces()
	if err != nil {
		log.Error("Could not get key pieces.")
		return fmt.Errorf("could not get key pieces: %s", err)
	}
	// Add our own KeyPiece in with the others.
	pieces = append(pieces, globals.HeldKeyPieces[globals.ThisNode])
	key, err := keyman.RebuildKey(pieces)
	if err != nil {
		log.Warn("Could not rebuild key:", err)
		return fmt.Errorf("could not rebuild key: %s", err)
	}
	globals.EncryptionKey = key
	globals.SystemLocked = false
	log.Info("Successfully unlocked system.")

	return nil
}

// UnlockWorker is run in a goroutine and periodically collects our KeyPieces
// from the other nodes and attempts to unlock our filesystem. It will run continuously
// until PFSD terminates or until the filesystem is unlocked.
func UnlockWorker() {
	defer globals.Wait.Done()

	timer := time.NewTimer(unlockQueryInterval)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-globals.Quit:
			if !ok {
				return
			}
		case <-timer.C:
			err := Unlock()
			if err != nil {
				log.Info("Failed to unlock system.")
				timer.Reset(unlockQueryInterval)
			} else {
				return
			}
		}
	}
}
