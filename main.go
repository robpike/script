// Copyright 2020 Rob Pike. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Script is an interactive driver for a line-at-a-time command language such as
// the shell.
//
// It reads each line of the script file, waiting for a newline on standard input
// to proceed. After receiving a newline, it prints the next line of of the script
// and also feeds it to a single running command instance, which prints the output
// of the command to standard output. Thus script serves as a way to control the
// sequential input to the command and thus supervise its activities.
//
// The argument file is in whatever format the command normally accepts. Comments,
// in whatever single-line form the command accepts, will be passed to the command
// but can serve as documentation during execution. There is no facility to send
// more than one line of input at a time.
//
// Normally the user types an empty line to proceed, advancing to the next line of
// the script. But a non-empty line is passed to the command without advancing the
// script, allowing the user to inject extra commands.
//
package main // import "robpike.io/cmd/script"
import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 3 {
		log.Fatal("Usage: script script-file command -args...\n")
	}
	log.SetPrefix("script: ")
	text, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	script := bytes.Split(text, []byte("\n"))
	if len(script) == 0 {
		log.Fatalf("%q is empty file", os.Args[1])
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	input, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	scan := bufio.NewScanner(os.Stdin)
	// Loop over lines of user input.
	for lineNo := 0; lineNo < len(script) && scan.Scan(); {
		if len(scan.Bytes()) == 0 {
			// User just hit newline. Send the current script line.
			line := script[lineNo]
			lineNo++
			fmt.Printf("%s\n", line)
			_, err = input.Write(append(line, '\n'))
		} else {
			// User typed some text. Send that instead.
			_, err = input.Write(scan.Bytes())
			input.Write([]byte("\n")) // Error ignored; that's OK.
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := scan.Err(); err != nil {
		log.Fatal(err)
	}
}
