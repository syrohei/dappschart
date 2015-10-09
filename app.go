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


func main() {

     target := "https://www.okcoin.com/api/v1/ticker.do"

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
     fmt.Println(classes)
     time.Sleep(10 * time.Second)
     }
  }

