
package main

import (
     "encoding/json"
     "fmt"
     "io/ioutil"
     "log"
     "net/http"
     "net/url"
     "time"
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
  }



func (class Class) String() string {
  return fmt.Sprintf("Id=%d Name=%d", class.Date, class.Ticker.Buy)

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

     p := new(Data)
     current := db.C("okcoin_btc_cny").Find(bson.M{"date": classes.Date})
     current.One(&p)
     fmt.Printf("%+v\n",p)

     time.Sleep(10 * time.Second)
     }
  

}
