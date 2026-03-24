package algorithm

import (
	"fmt"
	"log"

	"github.com/tyler-smith/go-bip32"
)

func BIP32Master(seed []byte) *bip32.Key {

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Master (xprv):", masterKey.String())
	return masterKey
}

func BIP32Child(k *bip32.Key, child uint32) *bip32.Key {
	ck, err := k.NewChildKey(child)
	if err != nil {
		log.Fatalf("derive child %d: %v", child, err)
	}
	return ck
}
