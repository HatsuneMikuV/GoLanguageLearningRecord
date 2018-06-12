package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/cmplx"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/*
	Goroutines和Channels
	1.并发程序指同时进行多个任务的程序
	2.goroutine
	3.channel
*/

//一，Goroutines
//1.在Go语言中，每一个并发的执行单元叫作一个goroutine
//2.新的goroutine会用go语句来创建,go语句会使其语句中的函数在一个新创建的goroutine中运行
//3.主函数返回时，所有的goroutine都会被直接打断，程序退出
func test_goroutine()  {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n) // slow
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}
func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}

//二，示例: 并发的Clock服务
//1.网络编程是并发大显身手的一个领域
//2.go语言的net包，提供编写一个网络客户端或者服务器程序的基本组件
func test_clock()  {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleConn(conn)
	}
}
func handleConn(c net.Conn)  {
	defer c.Close()

	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(time.Second * 1.0)
	}
}

//三，示例: 并发的Echo服务
//1.让服务使用并发不只是处理多个客户端的请求，甚至在处理单个连接时也可能会用到

func test_echo()  {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go echo_handleConn(conn)
	}
}
func echo_handleConn(c net.Conn)  {
	input := bufio.NewScanner(c)
	for input.Scan() {
		str := input.Text()
		if str == "EOF" {
			break
		}
		if len(str) > 0 {
			go echo(c, str, 1*time.Second)
		}
	}
	c.Close()
}
func echo(c net.Conn, shut string, delay time.Duration)  {
	fmt.Fprintf(c, "\t", strings.ToUpper(shut))
	time.Sleep(delay)
	fmt.Fprintf(c, "\t", shut)
	time.Sleep(delay)
	fmt.Fprintf(c, "\t", strings.ToLower(shut))
}

//四，Channels
//1.goroutine是Go语言程序的并发体的话，channels则是它们之间的通信机制
//2.一个channel是一个通信机制,它可以让一个goroutine通过它给另一个goroutine发送值信息
//3.channel的零值也是nil
//4.两个相同类型的channel可以使用==运算符比较，同样可以和nil比较
//5.channel有发送和接受两个主要操作，都是通信行为
//6.<-  channel在左侧是发送，在右侧是接收，不使用接收数据操作符也是合法的
//7.如果channel被关闭，则任何发送操作都会导致panic异常，
// 	但是依然能接收到成功发送的数据，如果没有数据，则产生零值数据
//8.channel的make创建第二个参数，对应其容量

//9.无缓存Channels的发送和接收操作将导致两个goroutine做一次同步操作，因此被称为同步Channels
//10.happens before
//11.一个Channel的输出作为下一个Channel的输入, 这种串联的Channels就是所谓的管道（pipeline）
func test_Channels()  {

	//演示在netcat的channel_84
	//test_echo()

	//串行channel
	//test_pipeline()

	//使用单向channel实现串行
	//单向操作的channel中的只接受方是不需要关闭的，发送者关闭即可
	//任何双向channel向单向channel变量的赋值操作都将导致该隐式转换
	//naturals := make(chan int)
	//squares := make(chan int)
	//go counter_chan(naturals)
	//go squarer(squares, naturals)
	//printer(squares)

	//带缓存channel
	test_channel_cache()
}

//12.串行的channel是无法知道另一个channel是否关闭的，
// 但是关闭的channel会产生第二个值，为bool类型，ture代表接收到值，false表示channel已被关闭没有可接收的值
func test_pipeline()  {
	naturals := make(chan int)
	squares := make(chan int)

	max := 100

	// Counter
	go func() {
		for x := 0; ; x++  {
			if x > max {
				close(naturals)
				break
			}
			naturals <- x
		}
	}()

	// Squarer
	go func() {
		for {
			x, ok := <- naturals
			if !ok {
				fmt.Println("naturals close")
				break
			}
			squares <- x * x
		}
		close(squares)
	}()

	// Printer (in main goroutine)
	for {
		x, ok := <-squares
		if !ok {
			fmt.Println("squares close")
			break
		}
		fmt.Println(x)
	}
}
func counter_chan(out chan<- int) {
	for x := 0; x < 100; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}
//13.带缓存的channel内部有一个元素队列，容量在创建时通过第二个参数指定
//14.带缓存channel会因为队列满导致发送等待，队列空接收等待
//15.发送的内容添加到队列尾部，接收内容是从队列头部取出并删除
func test_channel_cache()  {
	chan_cache := make(chan string, 3)
	//channel缓存容量
	fmt.Println(cap(chan_cache))

	chan_cache <- "A"
	chan_cache <- "B"
	chan_cache <- "C"
	//有效元素个数
	fmt.Println(len(chan_cache))

	fmt.Println(<-chan_cache) // "A"
	//有效元素个数
	fmt.Println(len(chan_cache))

	chan_cache <- "AA"
	//有效元素个数
	fmt.Println(len(chan_cache))

	fmt.Println(<-chan_cache) // "B"
	fmt.Println(<-chan_cache) // "C"
	//有效元素个数
	fmt.Println(len(chan_cache))

	result := mirroredQuery()
	fmt.Println(result)
}
func mirroredQuery() string {
	responses := make(chan string, 3)
	go func() { responses <- request("asia.gopl.io") }()
	go func() { responses <- request("europe.gopl.io") }()
	go func() { responses <- request("americas.gopl.io") }()
	/*
		练习 8.11： 紧接着8.4.4中的mirroredQuery流程，
		实现一个并发请求url的fetch的变种。
		当第一个请求返回时，直接取消其它的请求。
		----添加一个select，来终止goroutine
	*/
	select {
	case <-done:
		return "终止成功"
	case <-responses:
		return <-responses
	}
}
func request(hostname string) (response string) {
	if hostname == "asia.gopl.io" {
		time.Sleep(1 * time.Second)
		return "111"
	}else if hostname == "europe.gopl.io" {
		time.Sleep(2 * time.Second)
		return "222"
	}else if hostname == "americas.gopl.io" {
		time.Sleep(500 * time.Millisecond)
		return "333"
	}
	return "000"
}


//五，并发的循环
func test_concurrent(num int64)  {

	/*
		练习 8.4： 修改reverb2服务器，
		在每一个连接中使用sync.WaitGroup来计数活跃的echo goroutine。
		当计数减为零时，关闭TCP连接的写入，像练习8.3中一样。
		验证一下你的修改版netcat3客户端会一直等待所有的并发“喊叫”完成，
		即使是在标准输入流已经关闭的情况下。
	 */
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		if num == 87 {
			go handleConn_concurrent_87(conn)
		}else  {
			go handleConn_concurrent(conn)
		}
	}


	test_thr3_concurrent()
}
func handleConn_concurrent(c net.Conn)  {
	input := bufio.NewScanner(c)

	var wg sync.WaitGroup // number of working goroutines

	for input.Scan() {
		str := input.Text()
		if str == "EOF" {
			break
		}
		if len(str) > 0 {
			wg.Add(1)
			go func(c net.Conn, shut string, delay time.Duration) {
				defer wg.Done()
				fmt.Fprintf(c, "\t", strings.ToUpper(shut))
				time.Sleep(delay)
				fmt.Fprintf(c, "\t", shut)
				time.Sleep(delay)
				fmt.Fprintf(c, "\t", strings.ToLower(shut))
			}(c , str, 1*time.Second)
		}
	}
	wg.Wait()
	c.Close()
}

func test_thr3_concurrent()  {
	/*
			练习 8.5： 使用一个已有的CPU绑定的顺序程序，
			比如在3.3节中我们写的Mandelbrot程序计算程序，
			并将他们的主循环改为并发形式，使用channel来进行通信。
			在多核计算机上这个程序得到了多少速度上的改进？
			使用多少个goroutine是最合适的呢？
	*/

	color_chan := make(chan color.Color)

	var wg sync.WaitGroup // number of working goroutines

	fmt.Println(time.Now())
	fileName := "mandelbrot_concurrent.png"
	//绘制Mandelbrot图像
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height = 1024, 1024
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for py := 0; py < height; py++ {
		yyy := float64(py) / height * (ymax - ymin) + ymin
		for px := 0; px < width; px++ {
			xxx := float64(px) / width * (xmax - xmin) + xmin
			z := complex(xxx, yyy)
			wg.Add(1)
			go func(zz complex128) {
				defer wg.Done()
				color_chan <- mandelbrot_concurrent(zz)
			}(z)

			colo := <- color_chan
			img.Set(px, py, colo)
		}
	}
	fmt.Println(time.Now())

	wg.Wait()
	close(color_chan)

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	png.Encode(file, img)
}
func mandelbrot_concurrent(z complex128) color.Color  {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v * v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast * n}
		}
	}
	return color.Black
}


//六，并发的web爬虫
//用bfs(广度优先)算法来抓取整个网站
//每一个彼此独立的抓取命令可以并行进行IO，最大化利用网络资源
func test_web_crawler()  {

	//初步11111
	//crawl_one()

	//优化并发数量22222
	//第二个问题是这个程序永远都不会终止，即使它已经爬到了所有初始链接衍生出的链接
	//crawl_one()

	//优化并发程序能够终止33333
	crawl_one()
}
/*
	练习 8.6： 为并发爬虫增加深度限制。
	也就是说，如果用户设置了depth=3，
	那么只有从首页跳转三次以内能够跳到的页面才能被抓取到。
*/
var depths int64 = 3
var depthFirst int64 = 0

var tokens = make(chan struct{}, 20)
func web_crawl_two(url string) []string {
	fmt.Println(url)
	if depthFirst >= depths {
		return nil
	}
	depthFirst++
	tokens <- struct{}{} // acquire a token
	list, err := web_Extract(url)
	defer func() { <-tokens }() // release the token
	if err != nil {
		log.Print(err)
	}
	return list
}
func crawl_one()  {
	worklist := make(chan []string)

	//33333333
	var n int // number of pending sends to worklist
	// Start with the command-line arguments.
	n++
	//这个版本中，计算器n对worklist的发送操作数量进行了限制

	// Start with the command-line arguments.
	go func() {
		list := []string{
			"http://gopl.io/",
			"http://gopl.io/",
			"https://golang.org/help/",
			"https://golang.org/doc/",
			"https://golang.org/blog/"}

		worklist <- list
	}()


	// Crawl the web concurrently.
	seen := make(map[string]bool)

	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					//11111
					//worklist <- web_crawl_one(link)

					//22222
					worklist <- web_crawl_two(link)
				}(link)
			}
		}
	}
}
func web_crawl_one(url string) []string {
	fmt.Println(url)
	list, err := web_Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}
func web_Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	web_forEachNode(doc, visitNode, nil)
	return links, nil
}
func web_forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		web_forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

//七，基于select的多路复用
//1.和switch语句稍微有点相似，也会有几个case和最后的default选择支
//2.每一个case代表一个通信操作(在某个channel上进行发送或者接收)并且会包含一些语句组成的一个语句块
//3.一个接收表达式可能只包含接收表达式自身,或者包含在一个简短的变量声明中
//4.select会等待case中有能够执行的case时去执行,执行后，其他通信是不会执行
//5.没有任何case的select会永远等待下去，写作select{}
//6.time.Tick所建goroutine依然运行，但没有其他channel接收其值，导致goroutine泄露，因此使用代码所示方式
//7.nil channel的作用：1>对一个nil的channel发送和接收操作会永远阻塞
// 					  2>在select语句中操作nil的channel永远都不会被select到

func test_select_more()  {

	/*
		time.Tick->goroutine泄露,处理方式
		ticker := time.NewTicker(1 * time.Second)
		<-ticker.C    // receive from the ticker's channel
		ticker.Stop() // cause the ticker's goroutine to terminate
	*/

	//-------11111
	//test_countdown_one()
	//test_countdown_two()
	//test_countdown_thr()


	//-------22222

	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch:
			fmt.Println(x) // "0" "2" "4" "6" "8"
		case ch <- i:
		}
	}

	//
	test_concurrent(87)
}
/*
	练习 8.8： 使用select来改造8.3节中的echo服务器，为其增加超时，
	这样服务器可以在客户端10秒中没有任何喊话时自动断开连接。
*/
func handleConn_concurrent_87(c net.Conn)  {
	input := bufio.NewScanner(c)

	var wg sync.WaitGroup // number of working goroutines
	//abort := make(chan struct{})
	abort := make(chan string)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for  {
			select {
			case <-time.After(10 * time.Second):
				c.Close()
			case str := <-abort:
				wg.Add(1)
				go func(c net.Conn, shut string, delay time.Duration) {
					defer wg.Done()
					fmt.Fprintf(c, "\t", strings.ToUpper(shut))
					time.Sleep(delay)
					fmt.Fprintf(c, "\t", shut)
					time.Sleep(delay)
					fmt.Fprintf(c, "\t", strings.ToLower(shut))
				}(c , str, 1*time.Second)
			}
		}
	}()
	for input.Scan() {
		str := input.Text()
		if str == "exit" {
			break
		}
		if len(str) > 0 {
			abort <- str
		}

	}
	wg.Wait()
	c.Close()
}

func test_countdown_one()  {
	fmt.Println("Commencing countdown.")
	ticker := time.NewTicker(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-ticker.C
	}
	ticker.Stop()
	fmt.Println("launch()")
}
func test_countdown_two()  {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	fmt.Println("Commencing countdown.  Press return to abort.")
	select {
	case <-time.After(10 * time.Second):
		// Do nothing.
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	fmt.Println("launch()")
}
func test_countdown_thr()  {
	// ...create abort channel...
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown.  Press return to abort.")
	ticker := time.NewTicker(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-ticker.C:
			// Do nothing.
		case <-abort:
			fmt.Println("Launch aborted!")
			return
		}
	}
	ticker.Stop()
	fmt.Println("launch()")
}


//八，示例: 并发的目录遍历
func test_concurrent_directory()  {

	// Determine the initial directories.
	//------33333计算大小
	roots := []string{"/Users"}
	if len(roots) == 0 {
		roots = []string{"."}
	}
	// Traverse the file tree.
	fileSizes := make(chan int64)
	//go func() {
	//	for _, root := range roots {
	//		walkDir(root, fileSizes)
	//	}
	//	close(fileSizes)
	//}()

	//------44444计算大小优化
	/*
		因为磁盘系统并行限制，为了优化
		使用sync.WaitGroup (§8.5)来对仍旧活跃的walkDir调用进行计数，
		另一个goroutine会在计数器减为零的时候将fileSizes这个channel关闭
	*/
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir_group(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// Print the results.
	var nfiles, nbytes int64

	//------11111累加大小打印
	//for size := range fileSizes {
	//	nfiles++
	//	nbytes += size
	//}
	//printDiskUsage(nfiles, nbytes)

	//------22222累加大小打印优化
	/*
		主goroutine现在使用了计时器来每500ms生成事件，
		然后用select语句来等待文件大小的消息来更新总大小数据，
		或者一个计时器的事件来打印当前的总大小数据
	*/
	// Print the results periodically.
	var tick <-chan time.Time
	tick = time.Tick(500 * time.Millisecond)
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
}
/*
	练习 8.9： 编写一个du工具，每隔一段时间将root目录下的目录大小计算并显示出来
*/
func test_exerise89(timeS time.Duration, filep string)  {

	for {
		/*
			这一章，其实已经为我们写好了这个练习，
			只需要吧最优化的部分拿出来即可
		*/
		roots := []string{filep}
		if len(roots) == 0 {
			roots = []string{"."}
		}

		// Traverse the file tree.
		fileSizes := make(chan int64)

		var n sync.WaitGroup
		for _, root := range roots {
			n.Add(1)
			go walkDir_group(root, &n, fileSizes)
		}
		go func() {
			n.Wait()
			close(fileSizes)
		}()

		var nfiles, nbytes int64
		for size := range fileSizes {
			nfiles++
			nbytes += size
		}

		//root目录下的目录大小计算并显示出来
		fmt.Println(filep + "  fileSizes")
		printDiskUsage(nfiles, nbytes)

		//每隔一段时间
		time.Sleep(timeS)
	}
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}
func walkDir_group(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir_group(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}
/*
	由于这个程序在高峰期会创建成百上千的goroutine，
	我们需要修改dirents函数，
	用计数信号量来阻止他同时打开太多的文件
*/
// sema is a counting semaphore for limiting concurrency in dirents.
var sema = make(chan struct{}, 20)
func dirents(dir string) []os.FileInfo {
	//对这个程序的一个简单的性能分析可以揭示瓶颈在dirents函数中获取一个信号量
	//select操作可以这种操作被取消，并且可以将取消时的延迟从几百毫秒降低到几十毫秒
	select {
	case sema <- struct{}{}:
	case <-done:
		return nil
	}
	defer func() { <-sema }() // release token
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}


//九，并发的退出
//1.Go语言并没有提供在一个goroutine中终止另一个goroutine的方法
func test_concurrent_exit()  {

	c := os.Stdin

	go func() {
		input := bufio.NewScanner(c)
		for input.Scan() {
			str := input.Text()
			if len(str) > 0 {
				//有输入，即终止计算
				break
			}
		}
		c.Close()
		close(done)
	}()

	roots := []string{"/Users"}
	if len(roots) == 0 {
		roots = []string{"."}
	}

	fileSizes := make(chan int64)

	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir_group(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	var nfiles, nbytes int64

	var tick <-chan time.Time
	tick = time.Tick(500 * time.Millisecond)
loop:
	for {
		select {
		case <-done:
			fmt.Println("done")
			return
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
}

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func test_exertise_810(url string) (resp *http.Response, err error) {
	select {
	case <-done:
		return nil, nil
	}
	res, err := http.NewRequest("GET", url, nil)
	return res.Response, err
}

//十，示例: 聊天服务
func test_ex_chat()  {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn_ex813(conn)
	}
}
type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
	//练习 8.12： 使broadcaster能够将arrival事件通知当前所有的客户端。
	// 为了达成这个目的，你需要有一个客户端的集合，
	// 并且在entering和leaving的channel中记录客户端的名字。
	clientMap = make(map[string]string)
	oldKey = int64(10001)
)
func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				//练习 8.15： 如果一个客户端没有及时地读取数据可能会导致所有的客户端被阻塞。
				// 修改broadcaster来跳过一条消息，而不是等待这个客户端一直到其准备好写
				if !clients[cli] {
					continue
				}
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
func handleConn_ex(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	//练习 8.14：
	// 修改聊天服务器的网络协议这样每一个客户端就可以在entering时可以提供它们的名字。
	// 将消息前缀由之前的网络地址改为这个名字。
	if len(clientMap[who]) <= 0 {
		str := "num" + fmt.Sprint(oldKey)
		clientMap[who] =  str
		oldKey++
	}
	ch <- "You are " + clientMap[who]
	messages <- clientMap[who] + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- clientMap[who] + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- clientMap[who] + " has left"
	conn.Close()
}
//练习 8.13： 使聊天服务器能够断开空闲的客户端连接，
// 比如最近五分钟之后没有发送任何消息的那些客户端。
// 提示：可以在其它goroutine中调用conn.Close()来解除Read调用，
// 就像input.Scanner()所做的那样
func handleConn_ex813(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)

	abort := make(chan string)

	go func() {
		for  {
			select {
			case <-time.After(5 * time.Minute):
				conn.Close()
			case  str := <-abort:
				messages <- str
			}
		}
	}()

	for input.Scan() {
		str := input.Text()
		if str == "exit" {
			break
		}
		if len(str) > 0 {
			abort <- who + ": " + input.Text()
		}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
func main() {
	//一，Goroutines
	//test_goroutine()

	//二，示例: 并发的Clock服务
	//test_clock()

	//三，示例: 并发的Echo服务
	//test_echo()

	//四，Channels
	//test_Channels()

	//五，并发的循环
	//test_concurrent(0)

	//六，并发的web爬虫
	//test_web_crawler()

	//七，基于select的多路复用
	//test_select_more()

	//八，示例: 并发的目录遍历
	//test_concurrent_directory()

	//练习8.9  隔5秒打印一次桌面文件大小
	//test_exerise89(5 * time.Second, "/Users/angle/Desktop")

	//九，并发的退出
	//test_concurrent_exit()

	//十，示例: 聊天服务
	test_ex_chat()
}

