package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Swirl: a simple HTTP client that can follow redirects and save the response to a file (poor man's curl - from scratch. Forked from the original "swirl" source code I made.)
func main() {
	// Define flags
	urlPtr := flag.String("url", "", "URL of the page to fetch (required)")
	methodPtr := flag.String("X", "GET", "HTTP method for the request (default: GET)")
	outputPtr := flag.String("o", "", "Output file path (default: print to console)")
	followRedirectsPtr := flag.Bool("L", false, "Follow HTTP redirects (default: false)")

	flag.Parse()

	// Check for required flag (-url)
	if *urlPtr == "" {
		fmt.Println("Error: -url flag is required")
		flag.PrintDefaults()
		return
	}

	// Create an HTTP client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if *followRedirectsPtr { // Use the boolean variable
				return nil
			}
			return fmt.Errorf("too many redirects") // Or another error
		},
	}

	// Create a new request
	req, err := http.NewRequest(*methodPtr, *urlPtr, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching content:", err)
		return
	}

	defer resp.Body.Close() // Ensure body is closed

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:", resp.StatusCode)
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Handle output
	if *outputPtr != "" {
		// Write content to file
		err = os.WriteFile(*outputPtr, body, 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Println("Content saved to:", *outputPtr)
	} else {
		// Print content to console
		fmt.Println(string(body))
	}
}
