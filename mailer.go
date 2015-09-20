package main

import (
    "log"
    "net/smtp"
    )

func main() {
  // Set up authentication information.
  auth := smtp.PlainAuth(
    "",
    "syrohei@gmail.com", // foo@gmail.com
    "prgjsckcxspfcwhz",
    "smtp.gmail.com",
    )
        // Connect to the server, authenticate, set the sender and recipient,
        // and send the email all in one step.

  subject := "BTC/USD rate Will be Increase!"
  msg := "18bsT6FEXbfgT18Ask3gV2BTEq6k8GeUdx"   
  body := "Subject:" + subject + "\n" + msg
  err := smtp.SendMail(
    "smtp.gmail.com:587",
    auth,
    "syrohei@gmail.com", //foo@gmail.com
    []string{"syrohei@gmail.com"},
    []byte(body),
    )
  if err != nil {
    log.Fatal(err)
  }
}
