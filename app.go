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

     session, err := mgo.Dial("localhost")
     if err != nil {
      panic(err)
      }


     defer session.Close()

     db := session.DB("test")


     
     target := "https://www.okcoin.cn/api/v1/ticker.do"

     client := &http.Client{}

     req, err := http.NewRequest("GET", target, nil)
     if err != nil {
          log.Fatal(err)
     }

     values := url.Values{} // url.Valuesオブジェクト生成
     values.Add("symbol", "btc_cny") // key-valueを追加
     req.URL.RawQuery = values.Encode()

     for {
     res, err := client.Do(req)
       if err != nil {
         log.Fatal(err)
       }

     defer res.Body.Close()

     if res.StatusCode != http.StatusOK {
          log.Fatal(res)
     }

     body, err := ioutil.ReadAll(res.Body)
     if err != nil {
          log.Fatal(err)
     }

     //fmt.Println(string(body))

     var classes Class 
     err = json.Unmarshal(body, &classes)
     if err != nil {
          log.Fatal(err)
     }

     data := &Data {
      Date: classes.Date,
      Buy: classes.Ticker.Buy,
      Sell: classes.Ticker.Sell,
      High: classes.Ticker.High,
      Low: classes.Ticker.Low,
      Last: classes.Ticker.Last,
      Vol: classes.Ticker.Vol,
      }
     col := db.C("ticker")
     fmt.Printf("%s",classes.Date)
     err = col.Insert(data)
     if err != nil {
      panic(err)
      }

     p := new(Data)
     query := db.C("ticker").Find(bson.M{})
     query.One(&p)
     fmt.Printf("%+v\n",p)


     fmt.Printf(" : %s \n", classes)

     time.Sleep(1 * time.Second)
     }

}
