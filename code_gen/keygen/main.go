package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/lestrrat-go/jwx/v3/jwk"
)

func CreateKey() (jwk.Key, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	jwkKey, err := jwk.Import(privateKey)
	if err != nil {
		return nil, err
	}
	if err = jwk.AssignKeyID(jwkKey); err != nil {
		return nil, err
	}
	if err = jwkKey.Set(jwk.KeyUsageKey, "sig"); err != nil {
		return nil, err
	}
	if err = jwkKey.Set(jwk.AlgorithmKey, "RS256"); err != nil {
		return nil, err
	}
	return jwkKey, nil
}

func main() {
	outputPath := flag.String("output", "jwk.json", "Where to safe the new JWK-Key")
	flag.Parse()

	f, err := os.Create(*outputPath)
	if err != nil {
		log.Fatal("Output path could not be opened ", err)
	}
	defer f.Close()

	key, err := CreateKey()
	if err != nil {
		log.Fatal("key could not be created ", err)
	}
	if err = json.NewEncoder(f).Encode(key); err != nil {
		log.Fatal("key could not be written ", err)
	}
	log.Println("New Key written to", *outputPath)
}
