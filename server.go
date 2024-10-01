package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var key = []byte("mysecretkey1234567890123456")

const (
	port          = ":7070"
	acceptTimeout = 10 * time.Second
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	userAns := make(chan string)

	go func() {
		var response string
		fmt.Println("Incoming file , want to accept? (yes/no)")
		fmt.Scanln(&response)
		userAns <- response
	}()

	// Timeout handling
	select {
	case response := <-userAns:
		if response != "yes" {
			fmt.Println("file transfer rejected!")
			return
		}
		fmt.Println("File transfer accepted,Recieving file...")
		err := recieveFiles(conn)
		if err != nil {
			fmt.Println("Error recieving files:", err)
			return
		}
	case <-time.After(acceptTimeout):
		fmt.Println("Connection timed out , no response, aborting conn.....")
		return
	}

}

func recieveFiles(conn net.Conn) error {
	// recieve encrypted files
	files, err := os.Create("encrypted_file.txt")
	if err != nil {
		return err
	}
	defer files.Close()

	_, err = io.Copy(files, conn)
	if err != nil {
		return err
	}

	fmt.Println("File recieved succesfully!")

	err = decryptFile("encrypted_file.txt", "decrypted_file.txt")
	if err != nil {
		return err
	}

	fmt.Println("File decrypted and saved as 'decrypted_file.txt'")
	return nil

}

func decryptFile(src, dst string) error {

	// read from src file
	inFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// destination file after decryption
	outFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// decryption logic : compare the hashes of both files based on secretKey (key here)

	// i don't understand a f*ck of it
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	bytedata := make([]byte, aes.BlockSize)

	if _, err := io.ReadFull(inFile, bytedata); err != nil {
		return err
	}
	stream := cipher.NewCFBDecrypter(block, bytedata)
	reader := &cipher.StreamReader{S: stream, R: inFile}

	_, err = io.Copy(outFile, reader)
	return err

}

func main() {
	cn, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("error starting server:", err)
	}

	defer cn.Close()

	fmt.Println("starting server on port: ", port)

	for {
		conn, err := cn.Accept()
		if err != nil {
			fmt.Println("Error Accepting connections:", err)
			continue
		}
		go handleConnection(conn)
	}
}
