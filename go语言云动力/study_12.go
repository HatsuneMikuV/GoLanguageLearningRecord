package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
)

/* 标准包 */

//编码包
//紧接之前，方便测试，分离


//gob
func gob_test()  {
	name := "test.gob"
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)

	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func gob_test_one()  {
	gdp := map[string]float64 {
		"USA":14.58,
		"China":5.92,
		"Japan":5.45,
	}
	name := "test.gob"
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)

	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}

	enc := gob.NewEncoder(file)
	if err := enc.Encode(gdp); err != nil {
		fmt.Println("gob Cannot encode:", err)
		return
	}
}
var GDP map[string]float64

func gob_open_test()  {
	name := "test.gob"
	file, err := os.Open(name)
	defer file.Close()

	if err != nil {
		fmt.Println("gob Cannot open:", err)
		return
	}
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&GDP); err != nil {
		fmt.Println("gob Cannot decode:", err)
		return
	}
	fmt.Println("gob:", GDP)
}

//json
func json_test()  {
	name := "test.json"
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)

	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func json_test_one()  {
	gdp := map[string]float64 {
		"USA":14.58,
		"China":5.92,
		"Japan":5.45,
	}
	name := "test.json"
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)

	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}

	enc := json.NewEncoder(file)
	if err := enc.Encode(gdp); err != nil {
		fmt.Println("json Cannot encode:", err)
		return
	}
}
var GDP_json map[string]float64

func json_open_test()  {
	name := "test.json"
	file, err := os.Open(name)
	defer file.Close()

	if err != nil {
		fmt.Println("json Cannot open:", err)
		return
	}
	dec := json.NewDecoder(file)
	if err := dec.Decode(&GDP_json); err != nil {
		fmt.Println("json Cannot decode:", err)
		return
	}
	fmt.Println("json:", GDP_json)
}

func main() {

	//gob
	gob_test()

	//encode
	gob_test_one()

	//decode
	gob_open_test()



	//json
	json_test()

	//encode
	json_test_one()

	//decode
	json_open_test()
}
