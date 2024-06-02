package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
)

type Email struct {
	Date    string   `json:"date"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

const indexName = "enron-emails"
const defaultZincUrl = "http://localhost:4080"

func main() {

	skipProfilingParam := os.Getenv("SKIP_PROFILING")
	skipProfiling := false

	if skipProfilingParam == "TRUE" {
		skipProfiling = true
	}

	if !skipProfiling {

		// Setup profiling
		cpuFile, err := os.Create("profiling/cpu.prof")
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		defer cpuFile.Close()

		if err := pprof.StartCPUProfile(cpuFile); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()

		// Setup memory profiling
		memFile, err := os.Create("profiling/mem.prof")
		if err != nil {
			log.Fatal("Could not create memory profile: ", err)
		}
		defer memFile.Close()

		defer func() {
			runtime.GC() // Collect garbage before writing the memory profile
			if err := pprof.WriteHeapProfile(memFile); err != nil {
				log.Fatal("Could not write memory profile: ", err)
			}
		}()
	}

	rootPath := "./enron_mail_20110402/maildir"
	emails := []Email{}
	batchSize := 1000

	zincSearchUrl := os.Getenv("ZINCSEARCH_URL")

	if zincSearchUrl == "" {
		zincSearchUrl = defaultZincUrl
	}

	indexCreated, err := createZincIndex(zincSearchUrl)

	if !indexCreated {
		message := ""
		if err != nil {
			message = err.Error()
		}
		fmt.Println("Index not created, skipping document upload. ", message)
		// return
	}

	fmt.Println("Downloading dataset...")
	err = downloadEnronDb("enron_mail_20110402.tgz")
	if err != nil {
		log.Fatal("Error downloading the file: ", err)
		return
	}

	fmt.Println("Unpacking dataset...")
	err = untarGz("enron_mail_20110402.tgz", "./")
	if err != nil {
		log.Fatal("Error unpacking the file: ", err)
	}

	fmt.Println("Starting ingest of emails")
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// fmt.Println(path)
		if !info.IsDir() {
			email, err := parseEmail(path)
			if err != nil {
				fmt.Println("Error parsing email:", err)
				return nil // Ignore, let's process the rest of the files
			}
			emails = append(emails, email)

			if len(emails) >= batchSize {
				jsonData := createBulkJSON(emails)
				postToZinc(jsonData, zincSearchUrl)
				emails = []Email{}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through directories:", err)
		return
	}

	if len(emails) > 0 {
		jsonData := createBulkJSON(emails)
		postToZinc(jsonData, zincSearchUrl)
	}

	os.Remove("./enron_mail_20110402.tgz")
	os.RemoveAll("./enron_mail_20110402")
}

func parseEmail(filePath string) (Email, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Email{}, err
	}
	return extractEmailFields(string(content)), nil
}

func extractEmailFields(content string) Email {
	lines := strings.Split(content, "\n")
	email := Email{}
	readingBody := false
	bodyBuilder := strings.Builder{}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			readingBody = true // Start reading the body after the headers
			continue
		}
		if readingBody {
			bodyBuilder.WriteString(line + "\n")
			continue
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) < 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		switch key {
		case "Date":
			email.Date = value
		case "From":
			email.From = value
		case "To":
			email.To = strings.Split(value, ", ")
		case "Cc":
			email.Cc = strings.Split(value, ", ")
		case "Bcc":
			email.Bcc = strings.Split(value, ", ")
		case "Subject":
			email.Subject = value
		}
	}
	email.Body = bodyBuilder.String()
	return email
}

func createBulkJSON(emails []Email) string {
	records := make([]map[string]interface{}, 0, len(emails))
	for _, email := range emails {
		record := map[string]interface{}{
			"Date":    email.Date,
			"From":    email.From,
			"To":      email.To,
			"Cc":      email.Cc,
			"Bcc":     email.Bcc,
			"Subject": email.Subject,
			"Body":    email.Body,
		}
		records = append(records, record)
	}
	bulkData := map[string]interface{}{
		"index":   indexName,
		"records": records,
	}
	bulkJSON, _ := json.MarshalIndent(bulkData, "", "  ")
	return string(bulkJSON)
}
