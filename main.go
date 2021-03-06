package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
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
func mailSender(smtpHost string, smtpPort string, from string, password string, to []string, message []byte) error {

	smtpServerInstance := &smtpServer{host: smtpHost, port: smtpPort}

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpServerInstance.host)

	// Sending email
	err := smtp.SendMail(smtpServerInstance.Address(), auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
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

// healthCheck returns a status string for the webcheckd service
func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status OK.\n")
}

func main() {
	// Set logrus custom formatter
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	// Set commandline flags
	var recipientsList sliceFlag
	var urlList sliceFlag
	flag.Var(&recipientsList, "to", "Recipient e-mail")
	flag.Var(&urlList, "url", "Url list")
	senderAccount := flag.String("from", "sender@gmail.com", "Sender account")
	senderPassword := flag.String("password", "mypassword", "Sender password")
	smtpHost := flag.String("host", "smtp.gmail.com", "SMTP Server")
	smtpPort := flag.String("port", "587", "SMTP Port")
	intervalSecs := flag.Int64("interval", 300, "Interval in seconds")
	flag.Parse()

	if len(urlList) == 0 {
		log.Fatal("URL list cannot be empty")
	}

	// Start a os.Signal channel to accept signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-sigs
		log.Info("Shutting down webcheckd service\n")
		os.Exit(0)
	}()

	// Start embedded web server for service health probes
	go func() {
		http.HandleFunc("/healthz", healthCheck)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Start a deferred closure for site check
	defer func() {
		for {
			// Loop over the full URL list
			for _, siteUrl := range urlList {

				ok, status := siteChecker(siteUrl)

				mailMessage := []byte("Subject: Site alert notification.\r\n" +
					"\r\n" +
					"This is an alert message for " + siteUrl +
					".\nWebsite status is down with the following error: " + status +
					".\nPlease take action immediately." + "\r\n")

				if ok {
					log.Infof("%s is up. Status: %s\n", siteUrl, status)
				} else {
					log.Warnf("%s is down. Status: %s\n", siteUrl, status)
					//mailMessage := fmt.Sprintf("ALERT: The site" + *siteUrl + "is down!")
					err := mailSender(*smtpHost, *smtpPort, *senderAccount, *senderPassword, recipientsList, mailMessage)
					if err != nil {
						log.Errorf("Error sending email notification", err)
					}
					log.Info("Notification e-mail sent.")
				}
			}
			time.Sleep(time.Duration(*intervalSecs) * time.Second)
		}
	}()

	log.Infof("Webcheckd loop started for %s.\n", strings.Join(urlList, ", "))

}
