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
  Ema1Ave string `bson:"ema1ave"`
  DiffAve string `bson:"diffave"`
  Diff string `bson:"diff"`
  Up string `bson:"up"`
  }



func (class Class) String() string {
  return fmt.Sprintf("Id=%d Name=%d", class.Date, class.Ticker.Buy)

}

func Sum(a *[]Data) (sumEma1 float64) {
   for _, v := range *a {

      Ema1, _:=strconv.ParseFloat(v.Ema1,64)
      sumEma1 += Ema1

   }

   return sumEma1 
}
func isCount(a *[]Data) (count string) {
   for _, v := range *a {

      Up := v.Up
      if Up == "DOWNoverZERO" || Up == "UPoverZERO" || Up == "Mount" {
        count = "Mount" 
      }

   }

   return count 
}



func main() {

     session, err := mgo.Dial("172.31.8.24")
     //session, err := mgo.Dial("localhost")
     if err != nil {
      panic(err)
      }


     defer session.Close()

     db := session.DB("test")


     
     target := "https://www.okcoin.com/api/v1/ticker.do"

     client := &http.Client{}

     req, err := http.NewRequest("GET", target, nil)
     if err != nil {
          log.Fatal(err)
     }

     values := url.Values{} // url.Valuesオブジェクト生成
     values.Add("symbol", "btc") // key-valueを追加
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

     value := 200

     var lastData []Data
     query := db.C("okcoin_btc_cny")
     err = query.Find(bson.M{}).Limit(value-1).Sort("-date").All(&lastData)
     if err != nil {

     } else { 

     sumEma1 := Sum(&lastData)

     lastema1, _:=strconv.ParseFloat(lastData[0].Ema1,64)
     currentIndex, _:=strconv.ParseFloat(classes.Ticker.Last,64)
     var ema = lastema1*(1.0-2.0/13.0)+currentIndex*2.0/13.0
     ema1ave := (ema + sumEma1)/float64(value)
     currentdiff, _:=strconv.ParseFloat(lastData[0].Ema1,64)
     currentdiff =  ema - ema1ave 
     lastdiff, _:= strconv.ParseFloat(lastData[0].DiffAve,64)
     //lastdiff = lastdiff - currentdiff
     countbuff := lastData[:20]
     upcount := isCount(&countbuff)

     if currentdiff  >0.0 {
       if lastdiff < 0.0{
         if upcount != "Mount" {
         upcount = "UPoverZERO"
         }
       } else {
       upcount = "UP"
       }
     } else { 
       if lastdiff > 0.0000000{
         if upcount != "Mount" {
         upcount = "DOWNoverZERO"
 	 }
       } else {
	 upcount = "DOWN" 
       }
     }
     
     data := &Data {
      Date: classes.Date,
      Buy: classes.Ticker.Buy,
      Sell: classes.Ticker.Sell,
      High: classes.Ticker.High,
      Low: classes.Ticker.Low,
      Last: classes.Ticker.Last,
      Vol: classes.Ticker.Vol,
      Ema1: strconv.FormatFloat(ema, 'f', 6,64),
      Ema1Ave: strconv.FormatFloat(ema1ave, 'f', 6,64),
      DiffAve: strconv.FormatFloat(currentdiff, 'f', 6,64),
      Up: upcount,
      }
     err = query.Insert(data)
     if err != nil {
      panic(err)
      }

     p := new(Data)
     current := db.C("okcoin_btc_cny").Find(bson.M{"date": classes.Date})
     current.One(&p)
     fmt.Printf("%+v\n",p)

     time.Sleep(10 * time.Second)
     }
  }

}
