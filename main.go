package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"os"
)

const (
	FROM     = "adheeshgarg0611@gmail.com"
	FROMNAME = "Adheesh Garg"
	SUBJECT  = "Test Email with MIME"
	SMTPHOST = "smtp.gmail.com"
	SMTPPORT = "587"
)

var PASSWORD string

func real() string {
	if PASSWORD != "" {
		return PASSWORD
	}
	data, err := os.ReadFile("pass.txt")
	if err != nil {
		log.Fatal(err)
	}
	PASSWORD = string(data)
	return PASSWORD
}

func main() {
	// Assign the password
	real()
	TO := "adheeshgarg2005@gmail.com"
	// Create the email body
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Write the headers
	fmt.Fprintf(&body, "From: %s <%s>\r\n", FROMNAME, FROM)
	fmt.Fprintf(&body, "To: %s\r\n", TO)
	fmt.Fprintf(&body, "Subject: %s\r\n", SUBJECT)
	fmt.Fprintf(&body, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&body, "Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary())
	fmt.Fprintf(&body, "\r\n")

	// Write the plain text part
	textPart, err := writer.CreatePart(map[string][]string{
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Content-Transfer-Encoding": {"7bit"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(textPart, "This is the plain text part of the message.\r\n")

	// Write the HTML part
	htmlPart, err := writer.CreatePart(map[string][]string{
		"Content-Type":              {"text/html; charset=UTF-8"},
		"Content-Transfer-Encoding": {"7bit"},
	})
	if err != nil {
		log.Fatal(err)
	}
    fmt.Fprintf(htmlPart, `<html><body><p>This is the HTML part of the message.</p><img src="cid:image1"></body></html>`)
	fmt.Fprintf(htmlPart, "\r\n")

    // Write the image part
    imagePart, err := writer.CreatePart(map[string][]string{
        "Content-Type":              {"image/png"},
        "Content-Transfer-Encoding": {"base64"},
        "Content-Disposition":       {"inline; filename=\"image.png\""},
        "Content-ID":                {"<image1>"},
    })
    if err != nil {
        log.Fatal(err)
    }
	imageContent, err := os.ReadFile("image_base64.txt")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Fprintf(imagePart, "%s\r\n", imageContent)

	// Close the writer to finalize the MIME message
	writer.Close()

	// Set up authentication information.
	auth := smtp.PlainAuth("", FROM, PASSWORD, SMTPHOST)

	// Send the email
	err = smtp.SendMail(SMTPHOST+":"+SMTPPORT, auth, FROM, []string{TO}, body.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully!")
}
