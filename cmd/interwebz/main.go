package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// interwebz: text-based web browser with navigation, relative links, form support, and JS stripping
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "interwebz: usage: interwebz URL")
		os.Exit(1)
	}
	curURL := os.Args[1]
	for {
		links, forms, base := browse(curURL)
		if len(links) == 0 && len(forms) == 0 {
			break
		}
		fmt.Print("Enter link number, 'f' for form, or 'q' to quit: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "q" {
			break
		}
		if input == "f" && len(forms) > 0 {
			curURL = handleForm(forms, base)
			continue
		}
		num := -1
		fmt.Sscanf(input, "%d", &num)
		if num >= 1 && num <= len(links) {
			link := links[num-1]
			u, err := url.Parse(link)
			if err != nil || !u.IsAbs() {
				baseURL, _ := url.Parse(base)
				link = baseURL.ResolveReference(u).String()
			}
			curURL = link
		} else {
			fmt.Println("Invalid selection.")
		}
	}
}

// browse fetches a URL, prints text, and returns a list of links, forms, and the base URL
func browse(pageURL string) ([]string, []form, string) {
	resp, err := http.Get(pageURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "interwebz: error fetching URL:", err)
		return nil, nil, pageURL
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "interwebz: HTTP error: %d\n", resp.StatusCode)
		return nil, nil, pageURL
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "interwebz: error reading body:", err)
		return nil, nil, pageURL
	}
	html := string(body)
	// Strip <script>...</script> blocks (JavaScript)
	scriptRe := regexp.MustCompile(`(?is)<script.*?>.*?</script>`) // (?is) = case-insensitive, dot matches newline
	html = scriptRe.ReplaceAllString(html, "")
	// Find and list links
	linkRe := regexp.MustCompile(`<a [^>]*href=["']([^"']+)["'][^>]*>(.*?)</a>`)
	matches := linkRe.FindAllStringSubmatch(html, -1)
	links := []string{}
	for _, m := range matches {
		links = append(links, m[1])
	}
	// Find forms
	forms := parseForms(html)
	// Strip HTML tags for text view
	re := regexp.MustCompile("<[^>]+>")
	text := re.ReplaceAllString(html, "")
	fmt.Println("\n--- Page Content ---")
	fmt.Println(text)
	if len(links) > 0 {
		fmt.Println("\n--- Links ---")
		for i, l := range links {
			fmt.Printf("[%d] %s\n", i+1, l)
		}
	}
	if len(forms) > 0 {
		fmt.Println("\n--- Forms ---")
		for i, f := range forms {
			fmt.Printf("[f] Form %d: method=%s action=%s fields=%v\n", i+1, f.method, f.action, f.fields)
		}
	}
	return links, forms, pageURL
}

type form struct {
	method string
	action string
	fields []string
}

// parseForms extracts forms and their input fields from HTML
func parseForms(html string) []form {
	formRe := regexp.MustCompile(`(?is)<form[^>]*method=["']?(get|post)["']?[^>]*action=["']?([^"'> ]+)["']?[^>]*>(.*?)</form>`)
	inputRe := regexp.MustCompile(`(?is)<input[^>]*name=["']?([^"'> ]+)["']?[^>]*>`)
	forms := []form{}
	formMatches := formRe.FindAllStringSubmatch(html, -1)
	for _, fm := range formMatches {
		inputs := inputRe.FindAllStringSubmatch(fm[3], -1)
		fields := []string{}
		for _, in := range inputs {
			fields = append(fields, in[1])
		}
		forms = append(forms, form{method: strings.ToUpper(fm[1]), action: fm[2], fields: fields})
	}
	return forms
}

// handleForm prompts user for input and submits the form
func handleForm(forms []form, base string) string {
	fmt.Printf("Select form number (1-%d): ", len(forms))
	scanner := bufio.NewScanner(os.Stdin)
	num := 1
	if scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d", &num)
	}
	if num < 1 || num > len(forms) {
		fmt.Println("Invalid form number.")
		return base
	}
	f := forms[num-1]
	values := url.Values{}
	for _, field := range f.fields {
		fmt.Printf("Enter value for '%s': ", field)
		scanner.Scan()
		values.Set(field, scanner.Text())
	}
	formURL := f.action
	u, err := url.Parse(formURL)
	if err != nil || !u.IsAbs() {
		baseURL, _ := url.Parse(base)
		formURL = baseURL.ResolveReference(u).String()
	}
	if f.method == "GET" {
		return formURL + "?" + values.Encode()
	} else {
		resp, err := http.PostForm(formURL, values)
		if err != nil {
			fmt.Fprintln(os.Stderr, "interwebz: form POST error:", err)
			return base
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Fprintf(os.Stderr, "interwebz: HTTP error: %d\n", resp.StatusCode)
			return base
		}
		body, _ := io.ReadAll(resp.Body)
		// Save response to a temp file and return file:// URL
		tmp := "interwebz_form_response.html"
		os.WriteFile(tmp, body, 0644)
		return "file://" + tmp
	}
}
