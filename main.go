package main

import (
	"fmt"
	"github.com/4zv4l/chacha"
	"os"
	"strings"
)

func usage() {
	fmt.Printf("usage : %s [encrypt/decrypt] <file> <key>\n", os.Args[0])
	fmt.Println("Careful, decrypt a non-encrypted file may corrupt it")
}

// handleError "handle" (hum) errors :)
func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	args := os.Args
	// check if there is enough arguments
	if len(args) < 4 {
		usage()
		return
	}
	in, err := os.OpenFile(args[2], os.O_RDWR, 0600)
	handleError(err)
	defer in.Close()
	// default nonce
	nonce := "startHere"
	switch {
	case args[1] == "encrypt":
		outName := fmt.Sprintf("%s.cha", in.Name())
		out, err := os.OpenFile(outName, os.O_CREATE|os.O_RDWR, 0600)
		handleError(err)
		defer out.Close()
		err = chacha.EncryptFile(in, out, args[3], nonce)
		handleError(err)
		//os.Remove(args[2])
		break
	case args[1] == "decrypt":
		outName := strings.TrimSuffix(in.Name(), ".cha") + ".dec"
		out, err := os.OpenFile(outName, os.O_CREATE|os.O_RDWR, 0600)
		handleError(err)
		defer out.Close()
		err = chacha.DecryptFile(in, out, args[3], nonce)
		handleError(err)
		//os.Remove(args[2])
		break
	default:
		usage()
		return
	}
}
