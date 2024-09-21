package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
)

var cachedPassword string

func real() string {
    if cachedPassword != "" {
        return cachedPassword
    }
    data, err := os.ReadFile("pass.txt")
    if err != nil {
        log.Fatal(err)
    }
    cachedPassword = string(data)
    return cachedPassword
}

const body = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Subject}}</title>
    <style>
        body {
            font-family: 'PolySansTrial', sans-serif;
            background-color: #262626;
            color: #FCF3C5;
        }
        h1 {
            color: #EF4136;
        }
        h2 {
            color: #FF7A01;
        }
        .socials {
            padding: 3%;
            height: 27px;
        }
        .website {
            background-color: white;
            border: 2px solid white;
            color: black;
            font-weight: 600;
            padding: 7.5px;
        }
    </style>
</head>
<body>
    <div style="text-align: center; padding: 20px;">
        {{.Body}} 
        <br/><br/>
        <div class="socials">
            <a href="https://github.com/ACM-VIT" target="_blank">
                <img src="cid:gh.png" alt="github" height="27px" />
            </a>
            <a href="https://www.instagram.com/acmvit/?hl=en" target="_blank">
                <img src="cid:ig.png" alt="instagram" height="27px" />
            </a>
            <a href="https://www.linkedin.com/company/acmvit/" target="_blank">
                <img src="cid:LI.png" alt="linkedin" height="27px" />
            </a>
            <a href="https://www.youtube.com/@acm_vit" target="_blank">
                <img src="cid:yt.png" alt="youtube" height="27px" />
            </a>
        </div>
        <br/>
    </div>
</body>
</html>`

type EmailData struct {
    Subject string
    Body    string
}

func main() {
    from := "adheeshgarg0611@gmail.com"
    password := real()
    recipients := []string{"adheeshgarg0611@gmail.com", "adheeshgarg2005@gmail.com"} 
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    subject := "Test HTML Email with Image"

    // Create a new multipart writer
    var buf bytes.Buffer
    writer := multipart.NewWriter(&buf)         

    // Write the email headers
    buf.WriteString("From: " + from + "\r\n")
    buf.WriteString("To: " + strings.Join(recipients, ", ") + "\r\n")
    buf.WriteString("Subject: " + subject + "\r\n")
    buf.WriteString("MIME-Version: 1.0\r\n")
    buf.WriteString("Content-Type: multipart/related; boundary=" + writer.Boundary() + "\r\n")
    buf.WriteString("\r\n") // End of headers

    // Parse and execute the HTML template
    t, err := template.New("email").Parse(body)
    if err != nil {
        log.Fatal(err)
    }
    var htmlBody bytes.Buffer
    err = t.Execute(&htmlBody, EmailData{
        Subject: subject,
        Body:    "This is the body of the email",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Write the HTML part
    htmlPartHeaders := textproto.MIMEHeader{}
    htmlPartHeaders.Set("Content-Type", "text/html; charset=UTF-8")
    htmlPart, err := writer.CreatePart(htmlPartHeaders)
    if err != nil {
        log.Fatal(err)
    }
    htmlPart.Write(htmlBody.Bytes())

    // List of images to embed
    images := []struct {
        Path     string
        CID      string
        MimeType string
    }{
        {"images/github-mark-white.png", "gh.png", "image/png"},
        {"images/Instagram_logo_2016.png", "ig.png", "image/png"},
        {"images/LinkedIn_logo_initials.png", "LI.png", "image/png"},
        {"images/YouTube_full-color_icon_(2017).png", "yt.png", "image/png"},
    }

    // Attach each image
    for _, img := range images {
        imageBytes, err := os.ReadFile(img.Path)
        if err != nil {
            log.Fatal(err)
        }
        imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

        imagePartHeaders := textproto.MIMEHeader{}
        imagePartHeaders.Set("Content-Type", img.MimeType)
        imagePartHeaders.Set("Content-Transfer-Encoding", "base64")
        imagePartHeaders.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", img.Path))
        imagePartHeaders.Set("Content-ID", fmt.Sprintf("<%s>", img.CID))

        imagePart, err := writer.CreatePart(imagePartHeaders)
        if err != nil {
            log.Fatal(err)
        }
        imagePart.Write([]byte(imageBase64))
    }

    // Close the multipart writer
    writer.Close()

    // Convert the buffer to a byte slice
    message := buf.Bytes()

    auth := smtp.PlainAuth("", from, password, smtpHost)

    batchSize := 50
    for i := 0; i < len(recipients); i += batchSize {
        end := i + batchSize
        if end > len(recipients) {
            end = len(recipients)
        }
        bcc := recipients[i:end]

        log.Printf("Sending batch %d to %d recipients\n", i/batchSize+1, len(bcc))
        err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, bcc, message)
        if err != nil {
            log.Printf("Error sending batch %d: %v\n", i/batchSize+1, err)
            continue // Continue with the next batch instead of stopping
        }
        log.Printf("Batch %d sent successfully\n", i/batchSize+1)
    }

    log.Println("All sent successfully")
}
