package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func main() {

	pricesPtr := flag.String("set", "CH1,3.11,AP1,6.00,CF1,11.23,MK1,4.75,OM1,3.69", "A string which specifies a list of product and price paires.")

	itemsPtr := flag.String("items", "AP1,AP1,OM1,AP1", "A string which specifies a list of shopping items.")

	flag.Parse()

	names, nums, ok := validatePricelist(*pricesPtr)
	if !ok {
		return
	}

	//prices is a map which records the pair of product and its price
	prices := setPrices(names, nums)

	items, ok := validateShoppingItems(*itemsPtr, prices)
	if !ok {
		return
	}

	//strs := []string{"AP1", "AP1", "CH1", "AP1"}

	for {
		isPrint := true
		calculate(items, prices, isPrint)

		s, ok := userInput()
		if !ok {
			return
		}

		_, ok = validateUserInput(s, prices)
		if !ok {
			fmt.Println("")
			fmt.Println("Enter Y/y for continue, N/n for quit")
			var str string
			fmt.Scanln(&str)
			if strings.ToLower(str) == "n" {
				return
			}
			continue
		}

		items = append(items, s)
	}
}

func validatePricelist(str string) ([]string, []float32, bool) {

	nns := strings.Split(str, ",")

	if len(nns)%2 != 0 {
		fmt.Println("ERROR: the input string is not valid, it should be a list of pairs of product and price.")
		return nil, nil, false
	}

	names := make([]string, len(nns)/2)
	nums := make([]float32, len(nns)/2)

	for i, j := 0, 0; i < len(nns); i, j = i+2, j+1 {

		_, err := strconv.ParseFloat(nns[i], 32)
		if err == nil {
			fmt.Printf("ERROR: <%s> is not valid product name\n", nns[i])
			return nil, nil, false

		}

		v, err := strconv.ParseFloat(nns[i+1], 32)
		if err != nil {
			fmt.Printf("ERROR: <%s> is not a float number\n", nns[i+1])
			return nil, nil, false
		}

		names[j] = strings.TrimSpace(nns[i])
		nums[j] = float32(v)
	}

	return names, nums, true
}

func setPrices(names []string, nums []float32) map[string]float32 {

	res := make(map[string]float32)
	for i := 0; i < len(names); i++ {
		res[names[i]] = nums[i]
	}

	return res
}

func validateShoppingItems(str string, prices map[string]float32) ([]string, bool) {

	if len(str) == 0 {
		fmt.Printf("ERROR: the shopping list is empty\n")
		return nil, false
	}

	items := strings.Split(str, ",")

	res := make([]string, len(items))

	for i, v := range items {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		if _, ok := prices[v]; !ok {
			fmt.Printf("ERROR: ***%s*** did not existed in the product list\n", v)
			return nil, false
		}
		res[i] = v
	}
	return res, true
}

func userInput() (string, bool) {

	fmt.Println("")
	fmt.Println("Enter Y/y for continue, N/n for quit")
	var str string
	fmt.Scanln(&str)
	if strings.ToLower(str) == "n" {
		return "", false
	}
	fmt.Println("Enter next item name (one item only):")
	fmt.Scanln(&str)
	str = strings.TrimSpace(str)
	return strings.ToUpper(str), true
}

func validateUserInput(str string, prices map[string]float32) (string, bool) {

	if _, ok := prices[str]; !ok {
		fmt.Printf("ERROR: ***%s*** was not a valid shopping item\n", str)
		return "", false
	}
	return str, true
}

func calculate(strs []string, prices map[string]float32, isPrint bool) float32 {

	cnts_map, onOff_map := setDiscount(strs, prices)

	res := float32(0.0)

	printHead(isPrint)

	for _, v := range strs {

		printNormal(isPrint, v, prices[v])
		res += prices[v]

		discount := getDiscount(v, cnts_map, onOff_map, prices, isPrint)
		res += discount
	}

	printEnd(isPrint, res)

	return res
}

//Return two maps: cnts_map, onOff_map which will be used to calculate the discount value.
func setDiscount(products []string, prices map[string]float32) (map[string]int, map[string]bool) {

	cnts_map := make(map[string]int)
	onOff_map := make(map[string]bool)
	for k, _ := range prices {
		cnts_map[k] = 0
		onOff_map[k] = false
	}

	for _, s := range products {
		cnts_map[s] += 1
	}

	//CHMK policy
	if cnts_map["CH1"] > 0 {
		//onOff_ch1 = true
		onOff_map["CH1"] = true
	}

	return cnts_map, onOff_map
}

//Calculate the discount value according to different discount policy.
func getDiscount(str string, cnts_map map[string]int, onOff_map map[string]bool, prices map[string]float32, isPrint bool) float32 {

	var res float32

	switch str {
	case "CF1":
		{
			if onOff_map["CF1"] {
				printDiscount(isPrint, "BOGO", -prices["CF1"])
				res = -prices["CF1"]
			}
			onOff_map["CF1"] = !onOff_map["CF1"]
			return res
		}
	case "MK1":
		{
			if onOff_map["CH1"] {
				printDiscount(isPrint, "CHMK", -prices["MK1"])
				res = -prices["MK1"]
				onOff_map["CH1"] = false
			}
			return res
		}
	case "AP1":
		{

			//APOM policy, prefer to use APOM discount, as it will give the largest discount than APPL policy
			if cnts_map["OM1"] > 0 {
				diff := -prices["AP1"] / 2
				printDiscount(isPrint, "APOM", diff)
				res = diff
				cnts_map["OM1"] -= 1

				//return here, as one product only can have one discount
				return res
			}

			//APPL policy
			if cnts_map["AP1"] >= 3 {
				res = 4.50 - prices["AP1"]
				printDiscount(isPrint, "APPL", res)
			}
			return res
		}
	default:
		return res

	}
}

func printNormal(ok bool, str string, v float32) {

	if !ok {
		return
	}

	fmt.Printf("%s\t\t\t\t%8.2f\n", str, v)
}

func printDiscount(ok bool, code string, v float32) {
	if !ok {
		return
	}

	fmt.Printf("\t\t%s\t\t%8.2f\n", code, v)
}

func printHead(ok bool) {

	if !ok {
		return
	}

	fmt.Println("")
	fmt.Println("Item\t\t\t\t   Price")
	fmt.Println("----\t\t\t\t   -----")
}

func printEnd(ok bool, v float32) {

	if !ok {
		return
	}

	fmt.Printf("----------------------------------------\n")
	fmt.Printf("%40.2f\n", v)
}
