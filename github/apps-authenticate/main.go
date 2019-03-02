package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	// https://developer.github.com/apps/building-github-apps/authenticating-with-github-apps/
	//require 'openssl'
	//require 'jwt'  # https://rubygems.org/gems/jwt
	//
	//# Private key contents
	//private_pem = File.read(YOUR_PATH_TO_PEM)
	//private_key = OpenSSL::PKey::RSA.new(private_pem)
	//
	//# Generate the JWT
	//payload = {
	//	# issued at time
	//iat: Time.now.to_i,
	//	# JWT expiration time (10 minute maximum)
	//exp: Time.now.to_i + (10 * 60),
	//	# GitHub App's identifier
	//iss: YOUR_APP_ID
	//}
	//
	//jwt = JWT.encode(payload, private_key, "RS256")
	//puts jwt

	err := exec()
	if err != nil {
		panic(err)
	}
}

func exec() error {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return err
	}
	var pemFile string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".private-key.pem") {
			pemFile = file.Name()
			break
		}
	}
	if pemFile == "" {
		return errors.New("can't find *.private-key.pem")
	}
	pemData, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return err
	}

	pk, err := jwt.ParseRSAPrivateKeyFromPEM(pemData)
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.StandardClaims{
			Issuer:    os.Getenv("GITHUB_APP_ID"),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	)
	// fmt.Printf("%#v\n", token)

	accessToken, err := token.SignedString(pk)
	if err != nil {
		return err
	}

	fmt.Printf(`curl -i -H "Authorization: Bearer %s" -H "Accept: application/vnd.github.machine-man-preview+json" https://api.github.com/app`, accessToken)

	return nil
}
