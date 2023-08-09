// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/smtp"
// )

// func main() {
// 	// Sender's email details
// 	from := "avlashabhi336@gmail.com"
// 	password := "izpdbnaivxzwpucv"

// 	// Recipient's email address
// 	to := "abhilash.a@techversantinfo.com"

// 	// SMTP server configuration
// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := 587

// 	// Message to be sent
// 	message := []byte("Subject: Email Notification\n\nThis is the body of the email.")

// 	// Authentication
// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	// Sending email
// 	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, from, []string{to}, message)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Email sent successfully!")
// }

// package main

// func main() {

// 	msg := gomail.NewMessage()
// 	msg.SetHeader("From", "avlashabhi336@gmail.com")
// 	msg.SetHeader("To", "abhilash.a@techversantinfo.com")
// 	msg.SetHeader("Subject", "email check")
// 	msg.SetBody("text/html", "<b>Thi i sthe email sent by avlashabhi</b>")
// 	msg.Attach("./5.jpg")

// 	n := gomail.NewDialer("smtp.gmail.com", 587, "avlashabhi336@gmail.com", "izpdbnaivxzwpucv")

// 	// Send the email
// 	if err := n.DialAndSend(msg); err != nil {
// 		panic(err)
// 	}

// }

package main

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

func main() {

	// Sender data.
	from := "avlashabhi336@gmail.com"
	password := "izpdbnaivxzwpucv"

	// Receiver email address.
	to := []string{
		"assim.s@techversantinfo.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("index.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "IronMan",
		Message: "Where is Jarvis",
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}
