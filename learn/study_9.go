package main

import (
	"fmt"
	"math/rand"
	"time"
)

/* 同步通信 */

//去程
func sheep(i int, logo string)  {
	for ; ; i += 2 {
		fmt.Println(i, "只羊--", logo)
	}
}

func sheep_test()  {
	go sheep(1, "a")
	go sheep(2, "b")
	time.Sleep(time.Millisecond)
}

//程道
var wormhole chan time.Time

func deepspace()  {
	fmt.Println(time.Now())
	wormhole <- time.Now()
}

func chan_test()  {
	wormhole = make(chan time.Time)
	go deepspace()
	fmt.Println(<- wormhole)
}

//遍历与关闭
type hole chan time.Time

func deepspace_one(w hole, h int)  {
	defer close(w)
	for ; h > 0; h-- {
		w <- time.Now()
		time.Sleep(time.Second)
	}
}

func consumer(w hole)  {
	for msg := range w  {
		fmt.Println("consumer", msg)
	}
}

func chan_one()  {
	w := make(hole)
	go deepspace_one(w, 8)
	go consumer(w)
	for {
		msg, ok := <- w
		if !ok {
			break
		}
		fmt.Println("chan_one", msg)
	}
	fmt.Println("Done")
}

//MapReduce
type hole_one chan int

func deepspace_two(w hole_one)  {
	w <- rand.Int()
}

func chan_two()  {
	n := 8
	w := make(hole_one, n)
	//Map
	for i := 0; i < n; i++ {
		go deepspace_two(w)
	}

	//Reduce
	t := 0
	for i := 0; i < n; i++ {
		t += <- w
	}
	fmt.Println("Total:", t)
}

//select语句
func deepspece_select(w hole_one, h int)  {
	defer close(w)
	d := time.Duration(rand.Intn(h)) * time.Second

	for ; h > 0; h-- {
		w <- rand.Int()
		time.Sleep(d)
	}
}

func select_one()  {
	n := 8
	w := make(hole_one)
	t := 0
	maxTime := time.Second

	go deepspece_select(w, n)

	Out :for i := 0; i < n; i++ {
			select {
			case n := <- w:
				t += n
			case <- time.After(maxTime):
				fmt.Println("Time out")
				break Out
			}
		}
		fmt.Println("Total:", t)
}

func select_test()  {
	var c, c1, c2, c3 chan int
	var i1, i2 int

	select {
	case i1 = <- c1:
		fmt.Println("received", i1, "from c1\n")
	case c2 <- i2:
		fmt.Println("sent", i2, "to c2\n")
	case i3, ok := (<- c3)://等价于： i3， ok := <-c3
		if ok {
			fmt.Println("received", i3, "from c3\n")
		}else {
			fmt.Println("c3 is closed\n")
		}
	default:
		fmt.Println("no communication\n")
	}
	for  {//发送随机序列01到c
		select {
		case c <- 0:
			//注意，没有语句、没有dallthrough和重叠分支
		case c <- 1:

		}
	}
	select {}  //永远阻塞
}

//程道值
type Read struct {
	key string
	reply chan <- string
}

type Write struct {
	key string
	val string
}

var hole_two = make(chan interface{})

func deepspace_thi()  {
	m := map[string]string{}

	for {
		switch r := (<-hole_two).(type) {
		case Read:
			r.reply <- m[r.key] + " from Mars."
		case Write:
			m[r.key] = r.val
		}
	}
}

func chan_thi_test()  {
	go deepspace_thi()

	hole_two <- Write{"Name", "Martin"}
	home := make(chan  string)
	hole_two <- Read{"Name", home}
	fmt.Println(<-home)
}


//互斥


func main() {

	//去程
	//sheep_test()

	//程道
	//chan_test()

	//遍历与关闭
	//chan_one()

	//MapReduce
	//chan_two()

	//select语句
	//select_one()
	//select_test()

	//程道值
	chan_thi_test()
}
