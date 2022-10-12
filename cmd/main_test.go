package main

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"io"
	"strings"
	"testing"

	"filippo.io/age"
	"github.com/gopherlearning/gophkeeper/internal/bech32"
	"github.com/rs/zerolog/log"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/curve25519"
)

func toAgePrivate(priv ed25519.PrivateKey) string {

	s, err := bech32.Encode("AGE-SECRET-KEY-", priv)
	if err != nil {
		log.Info().Msg(err.Error())
	}

	return strings.ToUpper(s)
}
func toAgePublic(pub ed25519.PublicKey) string {
	s, err := bech32.Encode("age", pub)
	if err != nil {
		log.Info().Msg(err.Error())
	}

	return s
}

func TestMain(t *testing.T) {
	// Generate a mnemonic for memorization or user-friendly seeds
	// entropy, _ := bip39.NewEntropy(256)
	// mnemonic, _ := bip39.NewMnemonic(entropy)

	// // Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	// seed := bip39.NewSeed(mnemonic, "")

	// masterKey, _ := bip32.NewMasterKey(seed)
	// publicKey := masterKey.PublicKey()

	// // Display mnemonic and keys
	// fmt.Println("Mnemonic: ", mnemonic)
	// fmt.Println("Master private key: ", masterKey)
	// fmt.Println("Master public key: ", publicKey)

	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return
	}

	// seed := bip39.NewSeed(mnemonic, "")
	fmt.Println(bip39.IsMnemonicValid(mnemonic))
	fmt.Println(mnemonic)

	mnemonic1 := "require solid liar million never coyote century quality uncover soft job agent"
	seed1 := bip39.NewSeed(mnemonic1, "")

	priv1, err := curve25519.X25519(seed1[:32], curve25519.Basepoint)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	pub1, err := curve25519.X25519(priv1, curve25519.Basepoint)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	mnemonic2 := "erupt pulse green fashion disagree index prepare own expect exercise wife fresh"
	seed2 := bip39.NewSeed(mnemonic2, "")
	priv2, err := curve25519.X25519(seed2[:32], curve25519.Basepoint)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	pub2, err := curve25519.X25519(priv2, curve25519.Basepoint)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	mnemonic3 := "prize sniff practice series pause state clock guilt ribbon dolphin cancel isolate"
	seed3 := bip39.NewSeed(mnemonic3, "")
	// pub3, priv3, err := ed25519.GenerateKey(bytes.NewReader(seed3[:32]))
	priv3, err := curve25519.X25519(seed3[:32], curve25519.Basepoint)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	pub3, err := curve25519.X25519(priv3, curve25519.Basepoint)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	fmt.Printf("1:\nAge: %s\nAge: %s\n2:\nAge: %s\nAge: %s\n3:\nAge: %s\nAge: %s\n", toAgePrivate(priv1), toAgePublic(pub1), toAgePrivate(priv2), toAgePublic(pub2), toAgePrivate(priv3), toAgePublic(pub3))

	rec1, err := age.ParseX25519Recipient(toAgePublic(pub1))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	rec2, err := age.ParseX25519Recipient(toAgePublic(pub2))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	out := &bytes.Buffer{}

	w, err := age.Encrypt(out, rec1, rec2)
	if err != nil {
		log.Fatal().AnErr("Failed to create encrypted file: %v", err)
	}
	if _, err := io.WriteString(w, "Black lives matter."); err != nil {
		log.Fatal().AnErr("Failed to write to encrypted file: %v", err)
	}
	if err := w.Close(); err != nil {
		log.Fatal().AnErr("Failed to close encrypted file: %v", err)
	}

	fmt.Printf("Encrypted file size: %d\n", out.Len())

	encrypted := make([]byte, out.Len())
	copy(encrypted, out.Bytes())

	identity, err := age.ParseX25519Identity(toAgePrivate(priv1))
	if err != nil {
		log.Fatal().AnErr("Failed to parse private key: %v", err)
	}

	r, err := age.Decrypt(bytes.NewReader(encrypted), identity)
	if err != nil {
		log.Fatal().AnErr("Failed to open encrypted file: %v", err)
		return
	}
	out = &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		log.Error().AnErr("Failed to read encrypted file: %v", err)
	}

	fmt.Printf("File contents: %q\n", out.Bytes())

	identity, err = age.ParseX25519Identity(toAgePrivate(priv2))
	if err != nil {
		log.Error().AnErr("Failed to parse private key: %v", err)
	}

	r, err = age.Decrypt(bytes.NewReader(encrypted), identity)
	if err != nil {
		log.Error().AnErr("Failed to open encrypted file: %v", err)
		return
	}
	out = &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		log.Error().AnErr("Failed to read encrypted file: %v", err)
	}

	fmt.Printf("File contents: %q\n", out.Bytes())

	identity, err = age.ParseX25519Identity(toAgePrivate(priv3))
	if err != nil {
		log.Error().AnErr("Failed to parse private key: %v", err)
	}

	r, err = age.Decrypt(bytes.NewReader(encrypted), identity)
	if err != nil {
		log.Err(err).Msg("")
		fmt.Println(err)
		return
	}
	out = &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		log.Error().AnErr("Failed to read encrypted file: %v", err)
	}

	fmt.Printf("File contents: %q\n", out.Bytes())
	// age1eaqdpyzkcapytca6vz5dzrn3k9z9wpgzr89zxyzhg6vw6uatvqxqns4xjz
	// age1rnvcjt3a0t6tfpsgrthx3dlm4wcs9rnumktdqtcdrn0gac8xuvzsdp6l3c
	// masterKey, _ := bip32.NewMasterKey(seed)
	// publicKey := masterKey.PublicKey()
	// // crypto.GenerateKey()
	// prv1, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// prv2, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// message := []byte("Hello, world.")
	// ct, err := ecies.Encrypt(rand.Reader, &prv2.PublicKey, message, nil, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Encrypted: %s\n", string(ct))
	// pt, err := prv2.Decrypt(ct, nil, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Decrypted: %s\n", string(pt))

	// if !bytes.Equal(pt, message) {
	// 	log.Fatal(err)
	// }

	// _, err = prv1.Decrypt(ct, nil, nil)
	// if err == nil {
	// 	log.Fatal(err)
	// }

}
