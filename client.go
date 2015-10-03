// client.go
package main

import (
"fmt"
"log"
"net"
"net/rpc/jsonrpc"
"os"
"strconv"
)


	type buyingStocks struct {
		stockSymbolAndPercentage string
		budget float64
	}
	type responseBuying struct {
		tradeid int
		stocks []string
		unvestedAmount float64
	}

	type ResponseResult struct {
		Message string
	}

	func main() {

		arg := len(os.Args[1:])


		if arg == 2 {
			fmt.Println("buy stocks ")

			

			client, err := net.Dial("tcp", "127.0.0.1:1234")
			if err != nil {
				log.Fatal("dialing:", err)
			}

			f, err := strconv.ParseFloat(os.Args[2], 64)
			if err != nil {
				log.Fatal("dialing:", err)
			}
			args := &buyingStocks{os.Args[1],f}
			var reply ResponseResult
			c := jsonrpc.NewClient(client)
			err = c.Call("Stocktradingsystem.PurchStocks", args, &reply)
			if err != nil {
				log.Fatal("arith error:", err)
			}
			fmt.Println("\n",reply.Message)
		} else if arg == 1 {
			client, err := net.Dial("tcp", "127.0.0.1:1234")
			fmt.Println("second section")
			//fmt.Println("Arguement 1 =" + os.Args[1])

			tradeid1,err := strconv.Atoi(os.Args[1])
		
			var reply ResponseResult
			c := jsonrpc.NewClient(client)
			err = c.Call("TradingSystem.CheckingPortfolio", tradeid1, &reply)
			if err != nil {
				log.Fatal("arith error:", err)
			}
			fmt.Println("\n",reply.Message)

		}else{
			fmt.Println("Please check your input")
			fmt.Println("Syntax to run this system : go run rpc_client.go  GOOG:50%,YHOO:30%,GOOG:10% 10000")
		}




	}
