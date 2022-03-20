package main

import (
	"fmt"
	"github.com/4zv4l/chacha"
	"os"
	"strings"
	"sync"
)

// usage shows how to use the program
func usage() {
	fmt.Printf("usage : %s [encrypt/decrypt] <key> <file> <file> ...\n", os.Args[0])
}

// handleError "handle" (hum) errors :)
func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	// check if there is enough arguments
	args := os.Args
	if len(args) < 4 {
		usage()
		return
	}
	
	// set and initialise the needed variables
	var (
		wg sync.WaitGroup
		isEncrypt bool
		nonce string = "helloThere"
		key string = args[2]
		files []string = args[3:]
	)

	// check if encrypt or decrypt
	if args[1] == "encrypt" {
		isEncrypt = true 
	} else if args[1] == "decrypt" { 
		isEncrypt = false
	} else { usage(); return }

	// for each file given => process them in a goroutine
	for _, file := range files {
		wg.Add(1)
		go func(filename string, isEncrypt bool, wg *sync.WaitGroup){
			defer wg.Done()

			// open the input file in read-only mode
			in, err := os.OpenFile(filename, os.O_RDONLY, 0600)
			handleError(err)
			defer in.Close()

			switch {
			
			// encrypt mode
			case isEncrypt:
				// setup the output file
				outName := fmt.Sprintf("%s.cha", in.Name())
				out, err := os.OpenFile(outName, os.O_CREATE|os.O_WRONLY, 0600)
				handleError(err)
				defer out.Close()
				err = chacha.EncryptFile(in, out, key, nonce)
				handleError(err)
			
			// decrypt mode
			case !isEncrypt:
				// setup the output file
				outName := strings.TrimSuffix(in.Name(), ".cha") + ".dec"
				out, err := os.OpenFile(outName, os.O_CREATE|os.O_WRONLY, 0600)
				handleError(err)
				defer out.Close()
				err = chacha.DecryptFile(in, out, key, nonce)
				handleError(err)

			}
		}(file, isEncrypt, &wg)
	}

	// wait for each file to be processed
	wg.Wait()
}
