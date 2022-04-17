package utils

import (
	"crypto/sha256"
	"fmt"
	"github.com/itchyny/base58-go"
	"log"
	"math/big"
)

func Sha256Of(initialString string) []byte {
	encoder := sha256.New()
	encoder.Write([]byte(initialString))
	return encoder.Sum(nil)
}

func Base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return string(encoded)
}

func GenerateShortUrl(url string) string {
	urlHashBytes := Sha256Of(url)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := Base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:10]
}
