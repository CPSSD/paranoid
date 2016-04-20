package main

import (
	"encoding/gob"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"os"
	"path"
	"sync"
	"time"
)

const unlockQueryInterval time.Duration = time.Second * 30
const unlockTimeout time.Duration = time.Minute * 10
const lockWaitDuration time.Duration = time.Minute * 1

type keyResponse struct {
	uuid  string
	piece *keyman.KeyPiece
}

func requestKeyPiece(uuid string, recievedPieceChan chan keyResponse) {
	piece, err := pnetclient.RequestKeyPiece(uuid)
	if err != nil {
		log.Errorf("Error requesting key piece from node %s: %s", uuid, err)
		return
	}
	recievedPieceChan <- keyResponse{
		uuid:  uuid,
		piece: piece,
	}
}

//Attempt to unlock the state machine
func Unlock() {

	timer := time.NewTimer(unlockQueryInterval)
	defer timer.Stop()
	timeout := time.After(unlockTimeout)

	generation := keyman.StateMachine.GetCurrentGeneration()
	if generation == -1 {
		log.Fatal("Failed to unlock system, not part of a generation")
	}

	peers, err := keyman.StateMachine.GetNodes(generation)
	if err != nil {
		log.Fatal("Failed to unlock system:", err)
	}
	for i := 0; i < len(peers); i++ {
		if peers[i] == globals.ThisNode.UUID {
			peers = append(peers[:i], peers[i+1:]...)
			break
		}
	}

	var pieces []*keyman.KeyPiece
	pieces = append(pieces, globals.HeldKeyPieces.GetPiece(generation, globals.ThisNode.UUID))

	recievedPieceChan := make(chan keyResponse, len(peers))
	var keyRequestWait sync.WaitGroup

	for {
		select {
		case <-timeout:
			log.Fatal("Failed to unlock system before timeout")
		case <-timer.C:
			if len(peers) == 0 {
				log.Fatal("No peers to request peers from")
			}
			for i := 0; i < len(peers); i++ {
				keyRequestWait.Add(1)
				go func() {
					defer keyRequestWait.Done()
					requestKeyPiece(peers[i], recievedPieceChan)
				}()
			}
			timer.Reset(unlockQueryInterval)
		case keyData := <-recievedPieceChan:
			for i := 0; i < len(peers); i++ {
				if peers[i] == keyData.uuid {
					pieces = append(pieces, keyData.piece)
					peers = append(peers[:i], peers[i+1:]...)
					key, err := keyman.RebuildKey(pieces)
					if err != nil {
						log.Warn("Could not rebuild key:", err)
						break
					}
					globals.EncryptionKey = key
					globals.SystemLocked = false

					done := make(chan bool, 1)
					go func() {
						keyRequestWait.Wait()
						done <- true
					}()
					for {
						select {
						case <-recievedPieceChan:
						case <-done:
							log.Info("Successfully unlocked system.")
							return
						}
					}
				}
			}
		}
	}
}

func LoadPieces() {
	if _, err := os.Stat(path.Join(globals.ParanoidDir, "meta", "pieces")); os.IsNotExist(err) {
		log.Info("Filesystem not locked. Will not attepmt to load KeyPieces.")
		return
	}
	globals.SystemLocked = true
	piecePath := path.Join(globals.ParanoidDir, "meta", "pieces")
	file, err := os.Open(piecePath)
	if err != nil {
		// If the file doesn't exist, ignore it, because it could just be the first run.
		if os.IsNotExist(err) {
			log.Debugf("KeyPiece GOB file %s does not exist.", piecePath)
			return
		}
		log.Fatalf("Unable to open %s for reading pieces: %s", piecePath, file)
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	err = dec.Decode(&globals.HeldKeyPieces)
	if err != nil {
		log.Fatal("Failed decoding GOB KeyPiece data:", err)
	}
}
