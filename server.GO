	// server.go
	package main
	import (
"log"
"net"
"net/rpc"
"net/rpc/jsonrpc"
	//"log"
"net/http"
//"os"
"strings"
"fmt"
"strconv"
"net/url"
"math"
ejson "encoding/json"
)
	var NUMB  int
	var responseBuying_obj [100] responseBuying 	
	type Args struct {
		X, Y int
	}
	type buyingStocks struct {
		stockSymbolAndPercentage string
		budget float64
	}

	type responseBuying struct {
		tradeid int
		stocks []string
		unvestedAmount float64
	}
	type Stocks struct{
		stock_name string
		shareCount int
		stock_price string
	}
	type ResponseResult struct {
		Message string
	}
	type TradingSystem struct {}


	func BuyShareLive (budget float64, brand string) (string,float64){
		var Url *url.URL
		Url, err := url.Parse("https://query.yahooapis.com")
		if err != nil {
			panic("Error Panic")
		}
		Url.Path += "/v1/public/yql"
		parameters := url.Values{}
		parameters.Add("q", "select * from yahoo.finance.quote where symbol in ('"+brand+"')")
		Url.RawQuery = parameters.Encode()
		Url.RawQuery += "&format=json&diagnostics=true"
		Url.RawQuery += "&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
	//fmt.Printf("Encoded URL is %q\n", Url.String())
		res, err := http.Get(Url.String())
									//
		if err != nil {
			panic("Error Panic")
		}
		defer res.Body.Close()
		var v map[string]interface{}
		dec:= ejson.NewDecoder(res.Body);
		if err := dec.Decode(&v); err != nil {
			fmt.Println("ERROR: " + err.Error())
		}   
	   // person := new(Person)
	   // json.Unmarshal([]byte(res.String()), person)
		symbol := v["query"].(map[string]interface{})["results"].(map[string]interface{})["quote"].(map[string]interface{})["symbol"]
		closingValue := v["query"].(map[string]interface{})["results"].(map[string]interface{})["quote"].(map[string]interface{})["prevprice"]
		fmt.Println("The symbol is ",symbol,"prevprice=",closingValue)
		f, err := strconv.ParseFloat(closingValue.(string), 64)
		shareCount :=  int(budget/f)
	//shareCount_float:= float64(shareCount)
		unvestedAmount := budget - float64(shareCount)*f
		share := brand +":"+strconv.Itoa(shareCount)+":$"+closingValue.(string)
		//fmt.Println("Share String is "+ share)
		return share,unvestedAmount
	}


	func (t *TradingSystem) CheckingPortfolio(tradeid int, reply *ResponseResult) error {

		fmt.Println("\nSection1--- PortFolio\n\n")
	//fmt.Println("Track ID =",tradeid)
		reply.Message = "\nstocks -->"
		var currentMarketValue  float64
		fmt.Println(responseBuying_obj[tradeid])
		for _,each_brand := range responseBuying_obj[tradeid].stocks{
		//fmt.Println(each_brand)
			values := strings.Split(each_brand,":")
		//fmt.Println(values[2])
			brand := values[0]
			shareCount := values[1]
			stock_value,err := strconv.ParseFloat(strings.Trim(values[2], "$"),64)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(stock_value)
			var Url *url.URL
			Url, err = url.Parse("https://query.yahooapis.com")
			if err != nil {
				panic("Error Panic")
			}
			Url.Path += "/v1/public/yql"
			parameters := url.Values{}
			 parameters.Add("q", "select * from yahoo.finance.quote where symbol in ('"+brand+"')")//
			 Url.RawQuery = parameters.Encode()
			 Url.RawQuery += "&format=json&diagnostics=true"
			 Url.RawQuery += "&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback="
			 res, err := http.Get(Url.String())
									//
			 if err != nil {
			 	panic("Error Panic")
			 }
			 defer res.Body.Close()
			 var v map[string]interface{}
			 dec:= ejson.NewDecoder(res.Body);
			 if err := dec.Decode(&v); err != nil {
			 	fmt.Println("ERROR: " + err.Error())
			 }   
			 symbol := v["query"].(map[string]interface{})["results"].(map[string]interface{})["quote"].(map[string]interface{})["symbol"]
			 closingValue := v["query"].(map[string]interface{})["results"].(map[string]interface{})["quote"].(map[string]interface{})["prevprice"]
			 fmt.Println("The symbol is ",symbol,"prevprice =",closingValue)
			 float_TradePrice, err := strconv.ParseFloat(closingValue.(string), 64)
			 if err != nil {
			 	panic("Error Panic")
			 }
			 var FinalMsg string
			 	//float_TradePrice = float_TradePrice - 1.0
			 fmt.Println("stock_value =",stock_value," float_TradePrice =",float_TradePrice)
			 if stock_value < float_TradePrice{
			 	FinalMsg = brand +":"+ shareCount +":"+ "+$"+closingValue.(string) + " "
			 }else if stock_value > float_TradePrice{
			 	FinalMsg = brand +":"+ shareCount +":"+"-$"+closingValue.(string) + " "
			 }else if stock_value == float_TradePrice{
			 	FinalMsg = brand +":"+ shareCount +":"+"=$"+closingValue.(string)+" "
			 }
			 reply.Message += FinalMsg
			 //MarketValue
			 currentMarketValue_this,err := strconv.ParseFloat(shareCount,64)
			 if err != nil {
			 	panic("Error Panic")
			 }
			 currentMarketValue += currentMarketValue_this*float_TradePrice
			}
			reply.Message += "\ncurrentMarketValue -->"+strconv.FormatFloat(currentMarketValue,'f',6,64)
			reply.Message += "\nunvestedAmount -->"+strconv.FormatFloat(responseBuying_obj[tradeid].unvestedAmount,'f',6,64)

			
			return nil
		}


		func (t *TradingSystem) PurchStocks(args *buyingStocks, reply *ResponseResult) error {

			fmt.Println("\nSection2,,,,,, trading \n\n")
			result := strings.Split(args.stockSymbolAndPercentage, ",")
			responseBuying_obj[NUMB].tradeid = NUMB
			for _,brand := range result {
				
				stock_data := strings.Split(brand,":");
				
				stock_name := stock_data[0]
				percentages := stock_data[1]
				allocated_budget,err := CallParsing(percentages,64) 
				if err != nil {
					fmt.Println(err)
				}
				availAmnt := allocated_budget * args.budget
				fmt.Println("allocated_budget = " + FloatToString(allocated_budget))
				fmt.Println("availAmnt = " + FloatToString(availAmnt))
				shareCount,unvestedAmount := BuyShareLive(availAmnt, stock_name)
				responseBuying_obj[NUMB].unvestedAmount += unvestedAmount
				fmt.Println("shareCount = ",shareCount ,"unvestedAmount= ",responseBuying_obj[NUMB].unvestedAmount)
				responseBuying_obj[NUMB].stocks = append(responseBuying_obj[NUMB].stocks,shareCount)
			}
			stocks:=queueTogether(responseBuying_obj[NUMB].stocks)
			reply.Message = "\ntradeid ="+ strconv.Itoa(responseBuying_obj[NUMB].tradeid) + "\nstocks =" + stocks + "\nunvestedAmount =" + strconv.FormatFloat(responseBuying_obj[NUMB].unvestedAmount,'f',6,64)
			NUMB = NUMB + 1
			
			fmt.Println("\n--------------\n\n")
			return nil
		}


		func queueTogether(input []string) string{
			var finResult string
			
			for _,stringval := range input{
				stringval = " " + stringval
				finResult += stringval
			}
			   return finResult
		}
		func FloatToString(input_num float64) string {
	    
			return strconv.FormatFloat(input_num, 'f', 6, 64)
		}


		func CallParsing(s string, bitSize int) (f float64, err error) {
			i := strings.Index(s, "%")
			if i < 0 {
				return 0, fmt.Errorf("no percentage sign available")
			}
			f, err = strconv.ParseFloat(s[:i], bitSize)
			if err != nil {
				return 0, err
			}
			return f / 100, nil
		}


		func DigitFunc(f float64, places int) (float64) {
			shift := math.Pow(10, float64(places))
			return digitMath(f * shift) / shift;    
		}

		func digitMath(f float64) float64 {
			return math.Floor(f + .5)
		}

		func main() {

			fmt.Println("\n--------------Servers are running now---------\n\n")
			NUMB = 1
			cal := new(TradingSystem)
			server := rpc.NewServer()
			server.Register(cal)
			server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
			listener, e := net.Listen("tcp", ":1234")
			if e != nil {
				log.Fatal("listen problem:", e)
			}
			for {
				if conn, err := listener.Accept(); err != nil {
					log.Fatal("problem: " + err.Error())
				} else {
					log.Printf("A new conection has been formed \n\n\n\n")
					go server.ServeCodec(jsonrpc.NewServerCodec(conn))
				}
			}
		}
