package main

import (
	"flag"
	"net/http"
	"net/smtp"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

const recipientsSize = 8

// sliceFlag defines a new type to hold multiple recipients
type sliceFlag []string

// String implements the stringer interface
func (s *sliceFlag) String() string {
	return ("Implementation of the String interface")
}

// Set adds new entries to the slice
func (r *sliceFlag) Set(value string) error {
	*r = append(*r, value)
	return nil
}

// smtpServer configuration struct
type smtpServer struct {
	host string
	port string
}

// serverName URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

// mailSender sends email notification messages in case of errors
func mailSender(from string, password string, to []string, message []byte) {

	smtpServerInstance := &smtpServer{host: "smtp.gmail.com", port: "587"}

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpServerInstance.host)

	// Sending email
	err := smtp.SendMail(smtpServerInstance.Address(), auth, from, to, message)
	if err != nil {
		log.Errorf("Error sending email notification", err)
		return
	}

	log.Info("Notification e-mail sent.")
}

// siteChecker verifies the status of the site
func siteChecker(url string) (bool, string) {
	r, err := regexp.Compile(`^200 OK$`)
	if err != nil {
		log.Fatal("Unable to compile the regexp.")
	}
	client := &http.Client{}
	rsp, err := client.Get(url)
	if err != nil {
		log.Errorf("%v\n", err)
	}

	return r.MatchString(rsp.Status), rsp.Status
}

func main() {
	var recipientsList sliceFlag
	flag.Var(&recipientsList, "to", "Recipient e-mail")
	siteUrl := flag.String("url", "http://www.example.com", "Url name")
	senderAccount := flag.String("from", "sender@gmail.com", "Sender account")
	senderPassword := flag.String("password", "mypassword", "Sender password")

	flag.Parse()

	for {
		ok, status := siteChecker(*siteUrl)

		mailMessage := []byte("Subject: Site status notification.\r\n" +
			"\r\n" +
			"This is an alert message for " + *siteUrl + ". Status is down with error " + status + "\r\n")

		if ok {
			log.Infof("%s is up. Status: %s\n", *siteUrl, status)
			//mailMessage := fmt.Sprintf("ALERT: The site" + *siteUrl + "is down!")
		} else {
			log.Warnf("%s is down. Status: %s\n", *siteUrl, status)
			//mailMessage := fmt.Sprintf("ALERT: The site" + *siteUrl + "is down!")
			mailSender(*senderAccount, *senderPassword, recipientsList, mailMessage)
		}
		time.Sleep(300 * time.Second)
	}
}
