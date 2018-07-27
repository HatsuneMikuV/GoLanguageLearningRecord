package main

import (
	"20180408/bank"
	"20180408/memo"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

/*
	第九章　基于共享变量的并发
*/


//一， 竞争条件
//1.在一个线性(就是说只有一个goroutine的)的程序中，程序的执行顺序只由程序的逻辑来决定
//2.一个特定类型的一些方法和操作函数，对于某个类型来说，如果其所有可访问的方法和操作都是并发安全的话，那么类型便是并发安全的
//3.在一个程序中有非并发安全的类型的情况下，我们依然可以使这个程序并发安全
//4.导出包级别的函数一般情况下都是并发安全的
//5.竞争条件指的是程序在多个goroutine交叉执行操作时，没有给出正确的结果
//6.避免数据竞争：第一种方法是不要去写变量
//7.避免数据竞争：避免从多个goroutine访问变量
//8.避免数据竞争：允许很多goroutine去访问变量，但是在同一个时刻最多只有一个goroutine在访问

func test_competition()  {

	//-----------

	bank.Init()

	var wg sync.WaitGroup
	wg.Add(1)
	// Alice:
	go func() {
		defer wg.Done()
		bank.Deposit(200)                // A1
	}()
	wg.Add(1)
	// Bob:
	go func(){
		defer wg.Done()
		bank.Deposit(100)
	}()

	wg.Add(1)
	go func(){
		defer wg.Done()
		res := bank.Withdraw(200)
		if res{
			fmt.Println("取款成功")
		}else {
			fmt.Println("取款失败")
		}
	}()
	wg.Wait()

	res := bank.Balance()
	fmt.Printf("剩余:%d" , res)
}

//二，sync.Mutex互斥锁
//1.一个只能为1和0的信号量叫做二元信号量(binary semaphore)
//2.这种互斥很实用，而且被sync包里的Mutex类型直接支持
//3.Lock和Unlock的调用是在所有路径中都严格配对的
//4.一个deferred Unlock即使在临界区发生panic时依然会执行
//5.defer调用只会比显式地调用Unlock成本高那么一点点，不过却在很大程度上保证了代码的整洁性
func test_sync_Mutex()  {

}

//三，sync.RWMutex读写锁
//1.一种特殊类型的锁，其允许多个只读操作并行执行，但写操作会完全互斥
//2.这种锁叫作“多读单写”锁(multiple readers, single writer lock)，Go语言提供的这样的锁是sync.RWMutex
//3.RLock只能在临界区共享变量没有任何写入操作时可用
//4.RWMutex只有当获得锁的大部分goroutine都是读操作，而锁在竞争条件下，
// 也就是说，goroutine们必须等待才能获取到锁的时候，RWMutex才是最能带来好处的
func test_sync_RWMutex()  {

}

//四，内存同步
//1.在一个独立的goroutine中，每一个语句的执行顺序是可以被保证的
//2.所有并发的问题都可以用一致的、简单的既定的模式来规避
//3.多个goroutine都需要访问的变量，使用互斥条件来访问
func test_memory()  {
	var x, y int
	go func() {
		x = 1 // A1
		fmt.Print("y:", y, " ") // A2
	}()
	go func() {
		y = 1                   // B1
		fmt.Print("x:", x, " ") // B2
	}()
}

//五，sync.Once初始化
//1.初始化延迟到需要的时候再去做就是一个比较好的选择--懒加载
//2.所有并发的问题都可以用一致的、简单的既定的模式来规避
//3.多个goroutine都需要访问的变量，使用互斥条件来访问

//Icon用到了懒初始化(lazy initialization)
var icons map[string]image.Image
var loadIconsOnce sync.Once

func test_sync_Once()  {

}
func loadIcons() {
	icons = make(map[string]image.Image)
	//icons["spades.png"] = loadIcon("spades.png")
	//icons["hearts.png"] = loadIcon("hearts.png")
	//icons["diamonds.png"] = loadIcon("diamonds.png")
	//icons["clubs.png"] = loadIcon("clubs.png")
}
func Icon(name string) image.Image {
	loadIconsOnce.Do(loadIcons)
	return icons[name]
}

//六，竞争条件检测
//1.Go的runtime和工具链为我们装备了动态分析工具--竞争检查器(the race detector)

//七，示例: 并发的非阻塞缓存
//1.duplicate suppression(重复抑制/避免)
func test_memoizing()  {
	m := memo.New(httpGetBody)
	for _, url := range incomingURLs() {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("%s, %s, %d bytes\n",
			url, time.Since(start), len(value.([]byte)))
	}
}
func test_memoizing_two()  {
	m := memo.New(httpGetBody)
	var n sync.WaitGroup
	for _, url := range incomingURLs() {
		n.Add(1)
		go func(url string) {
			start := time.Now()
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s, %d bytes\n",
				url, time.Since(start), len(value.([]byte)))
			n.Done()
		}(url)
	}
	n.Wait()
}
func incomingURLs() []string  {
	URLs := []string{"https://golang.org","https://godoc.org","https://play.golang.org"}
	return  URLs
}
func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//八，Goroutines和线程
//1.每一个OS线程都有一个固定大小的内存块(一般会是2MB)来做栈，
// 这个栈会用来存储当前正在被调用或挂起(指在调用其它函数时)的函数的内部变量
//2.goroutine会以一个很小的栈开始其生命周期，一般只需要2KB
//3.和OS线程不太一样的是一个goroutine的栈大小并不是固定的
//4.goroutine的栈的最大值有1GB
//5.Go的运行时包含了其自己的调度器
//6.Go调度器的工作和内核的调度是相似的，但是这个调度器只关注单独的Go程序中的goroutine
//7.Go调度器并不是用一个硬件定时器而是被Go语言"建筑"本身进行调度的
//8.Go的调度方式不需要进入内核的上下文，所以重新调度一个goroutine比调度一个线程代价要低得多
//9.Go的调度器使用GOMAXPROCS的变量来决定会有多少个操作系统的线程同时执行Go的代码,默认的值是CPU的核心数
//10.goroutine没有可以被程序员获取到的身份(id)的概念。这一点是设计上故意而为之，由于thread-local storage总是会被滥用
//11.Go鼓励更为简单的模式，这种模式下参数对函数的影响都是显式的
func test_Goroutines()  {

}
func main() {

	//一， 竞争条件
	//test_competition()

	//二，sync.Mutex互斥锁
	//test_sync_Mutex()

	//三，sync.RWMutex读写锁
	//test_sync_RWMutex()

	//四，内存同步
	//test_memory()


	//七，示例: 并发的非阻塞缓存
	//test_memoizing()
	test_memoizing_two()
}

