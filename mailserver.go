
package main

import (
     "fmt"
     "log"
     "time"
     "net/smtp"
     mgo "gopkg.in/mgo.v2"
     "gopkg.in/mgo.v2/bson"
)

type Data struct {
  Date  string  `bson:"date"`
  Buy string `bson:"buy"`
  Sell string `bson:"sell"`
  High string `bson:"high"`
  Low string `bson:"low"`
  Last string `bson:"last"`
  Vol string `bson:"vol"`
  Ema1 string `bson:"ema1"`
  Ema1Ave string `bson:"ema1ave"`
  DiffAve string `bson:"diffave"`
  Diff string `bson:"diff"`
  Up string `bson:"up"`
  }


func main() {

     //session, err := mgo.Dial("172.31.8.24")
     session, err := mgo.Dial("localhost")
     if err != nil {
      panic(err)
      }


     defer session.Close()

     db := session.DB("test")

     
     for {

     var p []Data

     query := db.C("okcoin_btc_cny")
     err = query.Find(bson.M{}).Limit(1).Sort("-date").All(&p)
     if err != nil {
        log.Fatal(err)
        }
     fmt.Printf("%+v\n",p)

       auth := smtp.PlainAuth(
           "",
           "syrohei@gmail.com", // foo@gmail.com
           "prgjsckcxspfcwhz",
           "smtp.gmail.com",
           )
       // Connect to the server, authenticate, set the sender and recipient,
       // and send the email all in one step.

     if p[0].Up == "UPoverZERO" {

               subject := "BTC/USD rate Will be Increase!"
               rate := p[0].Last
               msg := "18bsT6FEXbfgT18Ask3gV2BTEq6k8GeUdx"   
               body := "Subject:" + subject + "\n" + msg + "\n" + rate
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
     if p[0].Up == "DOWNoverZERO" {

               subject := "BTC/USD rate Will be Decrease!"
               rate := p[0].Last
               msg := "1JTF1QpJ6yNhtF6fRUEM14x6AxBL8F9TyE"   
               body := "Subject:" + subject + "\n" + msg + "\n" + rate
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

     time.Sleep(10 * time.Second)
     }
  

}
