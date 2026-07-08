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
	// Ensure URL has a scheme
	if !strings.HasPrefix(curURL, "http://") && !strings.HasPrefix(curURL, "https://") && !strings.HasPrefix(curURL, "file://") {
		curURL = "https://" + curURL
	}
	for {
		links, forms, base := browse(curURL)
		if links == nil && forms == nil {
			break
		}
		fmt.Print("Enter link number, 'f' for form, 'b' back, 'r' refresh, or 'q' to quit: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "q" {
			break
		}
		if input == "b" {
			fmt.Println("No previous page (history not implemented)")
			continue
		}
		if input == "r" {
			continue
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
				baseURL, err := url.Parse(base)
				if err != nil {
					fmt.Fprintln(os.Stderr, "interwebz: invalid base URL:", err)
					continue
				}
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
	// Handle file:// URLs
	if strings.HasPrefix(pageURL, "file://") {
		return browseFile(pageURL[7:])
	}

	resp, err := http.Get(pageURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "interwebz: error fetching URL:", err)
		return nil, nil, pageURL
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "interwebz: HTTP error: %d %s\n", resp.StatusCode, resp.Status)
		return nil, nil, pageURL
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "interwebz: error reading body:", err)
		return nil, nil, pageURL
	}
	html := string(body)

	// Determine base URL from response or original URL
	baseURL := pageURL
	baseRe := regexp.MustCompile(`(?i)<base\s+[^>]*href\s*=\s*["']([^"']+)["']`)
	if m := baseRe.FindStringSubmatch(html); len(m) > 1 {
		baseURL = m[1]
	}

	// Strip script/style/comments
	scriptRe := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	html = scriptRe.ReplaceAllString(html, "")
	styleRe := regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	html = styleRe.ReplaceAllString(html, "")
	commentRe := regexp.MustCompile(`(?is)<!--.*?-->`)
	html = commentRe.ReplaceAllString(html, "")

	// Find and list links
	linkRe := regexp.MustCompile(`<a\s+[^>]*href\s*=\s*["']([^"']+)["'][^>]*>(.*?)</a>`)
	matches := linkRe.FindAllStringSubmatch(html, -1)
	links := []string{}
	linkTexts := []string{}
	for _, m := range matches {
		href := m[1]
		text := stripTags(m[2])
		text = strings.TrimSpace(text)
		if text == "" {
			text = href
		}
		links = append(links, href)
		linkTexts = append(linkTexts, text)
	}

	// Find forms
	forms := parseForms(html)

	// Strip HTML tags for text view
	re := regexp.MustCompile("<[^>]+>")
	text := re.ReplaceAllString(html, "")
	// Decode common HTML entities
	text = decodeHTMLEntities(text)
	// Collapse multiple blank lines
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	fmt.Println("\n--- Page Content ---")
	fmt.Println(text)
	if len(links) > 0 {
		fmt.Println("\n--- Links ---")
		for i, l := range links {
			display := l
			if len(display) > 80 {
				display = display[:77] + "..."
			}
			fmt.Printf("[%d] %s (%s)\n", i+1, linkTexts[i], display)
		}
	}
	if len(forms) > 0 {
		fmt.Println("\n--- Forms ---")
		for i, f := range forms {
			fmt.Printf("[f] Form %d: method=%s action=%s fields=%v\n", i+1, f.method, f.action, f.fields)
		}
	}
	return links, forms, baseURL
}

// browseFile reads and displays a local HTML file
func browseFile(path string) ([]string, []form, string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "interwebz: error reading file:", err)
		return nil, nil, "file://" + path
	}
	html := string(data)
	baseURL := "file://" + path

	scriptRe := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	html = scriptRe.ReplaceAllString(html, "")
	styleRe := regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	html = styleRe.ReplaceAllString(html, "")
	commentRe := regexp.MustCompile(`(?is)<!--.*?-->`)
	html = commentRe.ReplaceAllString(html, "")

	linkRe := regexp.MustCompile(`<a\s+[^>]*href\s*=\s*["']([^"']+)["'][^>]*>(.*?)</a>`)
	matches := linkRe.FindAllStringSubmatch(html, -1)
	links := []string{}
	linkTexts := []string{}
	for _, m := range matches {
		href := m[1]
		text := stripTags(m[2])
		text = strings.TrimSpace(text)
		if text == "" {
			text = href
		}
		links = append(links, href)
		linkTexts = append(linkTexts, text)
	}

	forms := parseForms(html)

	re := regexp.MustCompile("<[^>]+>")
	text := re.ReplaceAllString(html, "")
	text = decodeHTMLEntities(text)
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	fmt.Println("\n--- Page Content ---")
	fmt.Println(text)
	if len(links) > 0 {
		fmt.Println("\n--- Links ---")
		for i, l := range links {
			display := l
			if len(display) > 80 {
				display = display[:77] + "..."
			}
			fmt.Printf("[%d] %s (%s)\n", i+1, linkTexts[i], display)
		}
	}
	if len(forms) > 0 {
		fmt.Println("\n--- Forms ---")
		for i, f := range forms {
			fmt.Printf("[f] Form %d: method=%s action=%s fields=%v\n", i+1, f.method, f.action, f.fields)
		}
	}
	return links, forms, baseURL
}

// stripTags removes HTML tags from a string
func stripTags(s string) string {
	re := regexp.MustCompile("<[^>]+>")
	return re.ReplaceAllString(s, "")
}

// decodeHTMLEntities decodes common HTML entities
func decodeHTMLEntities(s string) string {
	// Build entity strings from character codes to avoid Go source interpretation
	amp := mkEnt("amp")
	lt := mkEnt("lt")
	gt := mkEnt("gt")
	quot := mkEnt("quot")
	apos39 := mkEnt("#39")
	apos27 := mkEnt("#x27")
	nbsp := mkEnt("nbsp")
	copyEnt := mkEnt("copy")
	reg := mkEnt("reg")

	repl := []struct{ old, new string }{
		{amp, "&"},
		{lt, "<"},
		{gt, ">"},
		{quot, "\""},
		{apos39, "'"},
		{apos27, "'"},
		{nbsp, " "},
		{copyEnt, "(c)"},
		{reg, "(r)"},
	}
	for _, r := range repl {
		s = strings.ReplaceAll(s, r.old, r.new)
	}

	// Decode numeric entities
	numRe := regexp.MustCompile(`&#(\d+);`)
	s = numRe.ReplaceAllStringFunc(s, func(m string) string {
		m2 := numRe.FindStringSubmatch(m)
		if len(m2) > 1 {
			var code int
			fmt.Sscanf(m2[1], "%d", &code)
			if code >= 32 && code <= 126 {
				return string(rune(code))
			}
		}
		return m
	})
	return s
}

// mkEnt builds an HTML entity string like "&name;" from the name part
func mkEnt(name string) string {
	return string([]byte{'&'}) + name + string([]byte{';'})
}

type form struct {
	method string
	action string
	fields []string
}

// parseForms extracts forms and their input fields from HTML
func parseForms(html string) []form {
	formRe := regexp.MustCompile(`(?is)<form[^>]*>(.*?)</form>`)
	inputRe := regexp.MustCompile(`(?is)<input[^>]*name\s*=\s*["']?([^"'>\s]+)["']?[^>]*>`)
	actionRe := regexp.MustCompile(`(?i)action\s*=\s*["']([^"']+)["']`)
	methodRe := regexp.MustCompile(`(?i)method\s*=\s*["']([^"']+)["']`)
	forms := []form{}
	formMatches := formRe.FindAllStringSubmatch(html, -1)
	for _, fm := range formMatches {
		formHTML := fm[1]
		action := ""
		if m := actionRe.FindStringSubmatch(fm[0]); len(m) > 1 {
			action = m[1]
		}
		method := "GET"
		if m := methodRe.FindStringSubmatch(fm[0]); len(m) > 1 {
			method = strings.ToUpper(m[1])
		}
		inputs := inputRe.FindAllStringSubmatch(formHTML, -1)
		fields := []string{}
		for _, in := range inputs {
			fields = append(fields, in[1])
		}
		forms = append(forms, form{method: method, action: action, fields: fields})
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
	if formURL == "" {
		formURL = base
	}
	u, err := url.Parse(formURL)
	if err != nil || !u.IsAbs() {
		baseURL, err := url.Parse(base)
		if err != nil {
			fmt.Fprintln(os.Stderr, "interwebz: invalid base URL:", err)
			return base
		}
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
		tmp := "interwebz_form_response.html"
		os.WriteFile(tmp, body, 0644)
		return "file://" + tmp
	}
}
