package main

import "fmt"

// If an extension wordlist was not supplied, generate default extensions.

func default_extensions() []string {
	output := []string{}
	for i := 0; i < 2000; i++ {
		j := fmt.Sprintf("%d", i)
		output = append(output, j)
	}
	for i := 2000; i <= 20000; i += 100 {
		j := fmt.Sprintf("%d", i)
		output = append(output, j)
	}
	for i := 2001; i <= 20001; i += 100 {
		j := fmt.Sprintf("%d", i)
		output = append(output, j)
	}
	return output
}

// If a password-cracking wordlist was not supplied, generate default numeric passwords.

func default_passwords() []string {
	output := []string{}
	for i := 0; i < 9999; i++ {
		j := fmt.Sprintf("%04d", i)
		output = append(output, j)
	}
	return output
}
