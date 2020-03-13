package main

// If an extension wordlist was not supplied, generate default extensions.

func default_extensions() []int {
	output := []int{}
	for i := 0; i < 2000; i++ {
		output = append(output, i)
	}
	for i := 2000; i <= 20000; i += 100 {
		output = append(output, i)
	}
	for i := 2001; i <= 20001; i += 100 {
		output = append(output, i)
	}
	return output
}
