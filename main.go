package main

import (
    "log"
    "net/smtp"
)

func main() {
    from := "adheeshgarg0611@gmail.com"
    password := real()
    to := "adheeshgarg0611@gmail.com"
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    subject := "Subject: Test Email\n"
    body := "This is a test email."
    message := []byte(subject + "\n" + body)

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Email sent successfully")
}