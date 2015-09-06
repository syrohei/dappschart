
package main

import (
     "fmt"
     "time"
     mgo "gopkg.in/mgo.v2"
     "gopkg.in/mgo.v2/bson"
)

type Class struct {
     Date string 
     Ticker struct {
      Buy string
      Sell string
      High string
      Low string
      Last string
      Vol string
     }
}

type Data struct {
  Date  string  `bson:"date"`
  Buy string `bson:"buy"`
  Sell string `bson:"sell"`
  High string `bson:"high"`
  Low string `bson:"low"`
  Last string `bson:"last"`
  Vol string `bson:"vol"`
  }



func (class Class) String() string {
  return fmt.Sprintf("Id=%d Name=%d", class.Date, class.Ticker.Buy)

}


func main() {

     session, err := mgo.Dial("172.31.42.49")
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
       panic(err)
       }
     fmt.Printf("%+v\n",p)


     time.Sleep(10 * time.Second)
     }

}
