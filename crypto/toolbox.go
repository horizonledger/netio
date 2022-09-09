package crypto

//general crypto toolbox
//higher level functions building on btcec

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	//"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type Keypair struct {
	PrivKey btcec.PrivateKey `json:"privkey"`
	PubKey  btcec.PublicKey  `json:"pubkey"`
}

type KeypairH struct {
	PrivKey string `json:"privkey"`
	PubKey  string `json:"pubkey"`
	Address string `json:"address"`
}

type KeypairAll struct {
	PrivKey btcec.PrivateKey `json:"privkey"`
	PubKey  btcec.PublicKey  `json:"pubkey"`
	Address string           `json:"address"`
}

type KeypairPub struct {
	PubKey  btcec.PublicKey `json:"pubkey"`
	Address string          `json:"address"`
}

// TODO only from pubkey type
func Address(pubkey string) string {
	return "P" + GetSHAHash(pubkey)[:12]
}

func PairFromHex(hexString string) Keypair {
	pkBytes, _ := hex.DecodeString(hexString)
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	kp := Keypair{PrivKey: *privKey, PubKey: *pubKey}
	return kp
}

func PairFromSecret(secret string) Keypair {
	hasher := sha256.New()
	hasher.Write([]byte(secret))
	hashedsecret := hex.EncodeToString(hasher.Sum(nil))

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), []byte(hashedsecret))
	kp := Keypair{PrivKey: *privKey, PubKey: *pubKey}
	return kp
}

func KeypairToHex(k Keypair) KeypairH {
	return KeypairH{PrivKey: PrivKeyToHex(k.PrivKey), PubKey: PubKeyToHex(k.PubKey), Address: Address(PubKeyToHex(k.PubKey))}
}

func KeypairFromHex(k KeypairH) Keypair {
	return Keypair{PrivKey: PrivKeyFromHex(k.PrivKey), PubKey: PubKeyFromHex(k.PubKey)}
}

func PrivKeyToHex(privkey btcec.PrivateKey) string {
	return hex.EncodeToString(privkey.Serialize())
}

func PrivKeyFromHex(hexString string) btcec.PrivateKey {
	//TODO handle errors
	pkBytes, _ := hex.DecodeString(hexString)
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	return *privKey
}

func PubKeyToHex(pubkey btcec.PublicKey) string {
	return string(hex.EncodeToString(pubkey.SerializeCompressed()))
}

func PubKeyFromHex(hexString string) btcec.PublicKey {
	pubKeyBytes, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println("err ", err)
		//return
	}
	pubKey, _ := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	return *pubKey
}

// Decode hex-encoded serialized signature
func SignatureFromHex(hexString string) btcec.Signature {
	//TODO handle errors
	sigBytes, err := hex.DecodeString(hexString)

	if err != nil {
		fmt.Println("err ", err)
		//return
	}
	signature, err := btcec.ParseSignature(sigBytes, btcec.S256())
	if err != nil {
		fmt.Println(err)
		//return
	}
	return *signature
}

func SignatureToHex(sig btcec.Signature) string {
	//TODO handle errors
	x := sig.Serialize()
	return string(hex.EncodeToString(x))

}

func VerifyMessageSignPub(signature btcec.Signature, pubkey btcec.PublicKey, message string) bool {

	messageHash := MsgHash(message)
	verified := signature.Verify(messageHash, &pubkey)
	return verified
}

func VerifyMessageSign(signature btcec.Signature, keypair Keypair, message string) bool {

	messageHash := MsgHash(message)
	verified := signature.Verify(messageHash, &keypair.PubKey)
	return verified
}

func SignMsgHash(privkey btcec.PrivateKey, message string) btcec.Signature {
	messageHash := chainhash.DoubleHashB([]byte(message))
	//fmt.Println("messageHash ", messageHash)
	signature, err := privkey.Sign(messageHash)
	if err != nil {
		fmt.Println("err ", err)
		//return
	}
	return *signature
}

// //TODO
// func FaucetTx() {

// }

// func CreateSignedTx(tx block.Tx, kp Keypair) block.Tx {

// 	txjson, _ := json.Marshal(tx)
// 	signature := SignMsgHash(kp.PrivKey, string(txjson))
// 	pubkey_string := PubKeyToHex(kp.PubKey)
// 	tx.SenderPubkey = pubkey_string
// 	sighex := hex.EncodeToString(signature.Serialize())
// 	tx.Signature = sighex
// 	signedtx := tx
// 	return signedtx
// }
