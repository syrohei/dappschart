package main

import (
     "encoding/json"
     "strconv"
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
  Ema1 string `bson:"ema1"`
  Diff string `bson:"diff"`
  Posi string `bson:"posi"`
  }



func (class Class) String() string {
  return fmt.Sprintf("Id=%d Name=%d", class.Date, class.Ticker.Buy)

}

func Sum(a *[]Data) (sumEma1 float64, sumDiff float64) {
   for _, v := range *a {

      Ema1, _:=strconv.ParseFloat(v.Last,64)
      Diff, _:=strconv.ParseFloat(v.Diff,64)
      sumEma1 += Ema1
      sumDiff += Diff

   }
   return sumEma1, sumDiff
}


func main() {

     //session, err := mgo.Dial("172.31.42.49")
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

     var lastData []Data
     query := db.C("okcoin_btc_cny")
     err = query.Find(bson.M{}).Limit(10).Sort("-date").All(&lastData)
     if err != nil {

     } else { 

     sumEma1, sumDiff := Sum(&lastData)
     fmt.Println(sumEma1/10, sumDiff/10)

     lastema1, _:=strconv.ParseFloat(lastData[0].Ema1,64)
     currentIndex, _:=strconv.ParseFloat(classes.Ticker.Last,64)
     var ema = lastema1*(1.0-2.0/13.0)+currentIndex*2.0/13.0
     
     
     data := &Data {
      Date: classes.Date,
      Buy: classes.Ticker.Buy,
      Sell: classes.Ticker.Sell,
      High: classes.Ticker.High,
      Low: classes.Ticker.Low,
      Last: classes.Ticker.Last,
      Vol: classes.Ticker.Vol,
      Ema1: strconv.FormatFloat(ema, 'f', 6,64),
      Diff: strconv.FormatFloat(ema-lastema1, 'f', 6,64),
      Posi: strconv.FormatFloat(currentIndex-ema, 'f', 6,64),
      }
     err = query.Insert(data)
     if err != nil {
      panic(err)
      }

     p := new(Data)
     current := db.C("okcoin_btc_cny").Find(bson.M{"date": classes.Date})
     current.One(&p)
     fmt.Printf("%+v\n",p)

     time.Sleep(5 * time.Second)
     }
  }

}
