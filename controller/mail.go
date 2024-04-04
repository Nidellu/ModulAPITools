package controller

import "gopkg.in/gomail.v2"

func SendEmail(pesan string) {

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	sender := "vincentnatanaelchandra@gmail.com"
	password := "mtls uegh gqup qmfd"

	recipient := "brianzefanya71@gmail.com"

	message := gomail.NewMessage()
	message.SetHeader("From", sender)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", "Test Subject")
	message.SetBody("text/plain", pesan)

	dialer := gomail.NewDialer(smtpHost, smtpPort, sender, password)

	if err := dialer.DialAndSend(message); err != nil {
		panic(err)
	}

	println("Email sent successfully!")
}
