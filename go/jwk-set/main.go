package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"gopkg.in/square/go-jose.v2"
)

func main() {
	pk1, pem, err := examplePKSerialize()
	if err != nil {
		panic(err)
	}

	pk2, err := examplePKDeserialize(pem)
	if err != nil {
		panic(err)
	}

	fmt.Println("******************* pk1 *******************")

	err = exampleJWKSet(pk1)
	if err != nil {
		panic(err)
	}

	err = exampleJWS(pk1)
	if err != nil {
		panic(err)
	}

	fmt.Println("******************* pk2 *******************")

	err = exampleJWKSet(pk2)
	if err != nil {
		panic(err)
	}

	err = exampleJWS(pk2)
	if err != nil {
		panic(err)
	}
}

func examplePKSerialize() (*rsa.PrivateKey, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, "", err
	}

	var pemKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	var buf bytes.Buffer
	err = pem.Encode(&buf, pemKey)
	if err != nil {
		return nil, "", err
	}

	fmt.Println(buf.String())

	return privateKey, buf.String(), nil
}

func examplePKDeserialize(pemData string) (*rsa.PrivateKey, error) {
	pemBlock, _ := pem.Decode([]byte(pemData))
	if pemBlock.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected type: %s", pemBlock.Type)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func exampleJWKSet(privateKey *rsa.PrivateKey) error {
	jwkSet := &jose.JSONWebKeySet{}
	{
		jwk := jose.JSONWebKey{}
		jwk.Key = privateKey
		jwk.KeyID = "test"
		jwk.Algorithm = "RS256"
		jwk.Use = "sig"

		jwkSet.Keys = append(jwkSet.Keys, jwk)

		jwkSet.Keys = append(jwkSet.Keys, jwk.Public())
	}

	b, err := json.MarshalIndent(jwkSet, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}

func exampleJWS(privateKey *rsa.PrivateKey) error {
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: privateKey}, nil)
	if err != nil {
		return err
	}

	var payload = []byte("Lorem ipsum dolor sit amet")
	object, err := signer.Sign(payload)
	if err != nil {
		return err
	}

	{
		serialized, err := object.CompactSerialize()
		if err != nil {
			return err
		}

		fmt.Println(serialized)

		object, err = jose.ParseSigned(serialized)
		if err != nil {
			panic(err)
		}

		output, err := object.Verify(&privateKey.PublicKey)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(output))
	}
	{
		serialized := object.FullSerialize()

		fmt.Println(serialized)

		object, err = jose.ParseSigned(serialized)
		if err != nil {
			panic(err)
		}

		output, err := object.Verify(&privateKey.PublicKey)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(output))
	}

	return nil
}
