package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net"
	"os"
)

var dkey = []byte("mysecretkey1234567890123456")

func clientMain() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <file-to-send>")
		return
	}

	fileToSend := os.Args[1]
	conn, err := net.Dial("tcp", "192.168.1.5:8080")
	if err != nil {
		fmt.Println("Error sending file:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Sending file...")

	err = sendFile(conn, fileToSend)
	if err != nil {
		fmt.Println("Error sending file:", err)
		return
	}
	fmt.Println("File sent success")

}

func sendFile(conn net.Conn, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	// Encrypt the file
	encryptedPath := "encrypted_" + filePath
	err = encryptFile(filePath, encryptedPath)
	if err != nil {
		return nil
	}
	defer os.Remove(encryptedPath)

	// open the encrypted file to send over tcp
	encryptedFile, err := os.Open(encryptedPath)
	if err != nil {
		return nil
	}

	defer encryptedFile.Close()

	_, err = io.Copy(conn, encryptedFile)

	return err
}

func encryptFile(src, dst string) error {
	inFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(dst)
	if err != nil {
		return nil
	}
	defer outFile.Close()
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inFile, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outFile}

	_, err = io.Copy(writer, inFile)
	return err
}
