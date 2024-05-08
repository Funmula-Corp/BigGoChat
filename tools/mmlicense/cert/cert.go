package cert

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"encoding/pem"
	"log"
)

//go:embed private.pem
var PrivateKey []byte

//go:embed public.pem
var PublicKey []byte

func ParsePrivateKey() (pKey *rsa.PrivateKey) {
	privateKey, rest := pem.Decode(PrivateKey)
	if len(rest) > 0 {
		log.Fatalln("failed to decode private key buffer", rest)
	}

	var (
		err error
		key any
	)
	if key, err = x509.ParsePKCS8PrivateKey(privateKey.Bytes); err != nil {
		log.Fatalln("failed to parse private key: ", err)
	}

	switch key := key.(type) {
	case *rsa.PrivateKey:
		pKey = key
	default:
		log.Fatalf("invalid private key type: %T", key)
	}
	return
}

func SignBuffer(data []byte) (signature []byte) {
	hash := sha512.New()
	hash.Write(data)

	var err error
	if signature, err = rsa.SignPKCS1v15(rand.Reader, ParsePrivateKey(), crypto.SHA512, hash.Sum(nil)); err != nil {
		log.Fatalln("error generating signature: ", err)
	}
	return
}

func ValidateLicense(signed []byte) {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(signed)))

	_, err := base64.StdEncoding.Decode(decoded, signed)
	if err != nil {
		log.Fatalf("encountered error decoding license: %s\r\n", err)
	}

	// remove null terminator
	for len(decoded) > 0 && decoded[len(decoded)-1] == byte(0) {
		decoded = decoded[:len(decoded)-1]
	}

	if len(decoded) <= 256 {
		log.Fatalln("Signed license not long enough")
	}

	plaintext := decoded[:len(decoded)-256]
	signature := decoded[len(decoded)-256:]

	block, _ := pem.Decode(PublicKey)

	public, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("Encountered error signing license: %s\r\n", err)
	}

	rsaPublic := public.(*rsa.PublicKey)

	h := sha512.New()
	h.Write(plaintext)
	d := h.Sum(nil)

	err = rsa.VerifyPKCS1v15(rsaPublic, crypto.SHA512, d, signature)
	if err != nil {
		log.Fatalf("Invalid signature: %s\r\n", err)
	}
}
