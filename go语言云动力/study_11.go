package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

/* 标准包 */


//超链接包
//接着study_10，用于联调调用

func test_eight()  {
	r, err := http.Get("http://localhost:1234/图灵")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer r.Body.Close()
	fmt.Println("Hedaer:")

	for k, v := range r.Header{
		fmt.Printf("%s :%s\n", k, v)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Body: \n%s\n", body)
}
//get
func test_nine()  {
	host := "http://localhost:1234/"

	user := url.Values{
		"Name":{"Megan fox"},
		"Sex":{"female"},
	}
	q := host + "?" + user.Encode()
	fmt.Println("\n" + q + "\n")
	r, err := http.Get(q)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("body：%s", body)
}

//post
type User struct {
	Name string
	Date time.Time
}

var users []User

var homePage = template.Must(template.New("home").Parse(
	`<html><body>{{range.}}
{{.Name}} dates on{{.Date.Format "3:04pm, 2 Jan"}}<br/>
{{end}}

<form action="/post"method="post">
姓名：<input type='text'name='Name'/><br/>
<input type='submit'value='提交'/>
</form></body></html>
`))

func home(w http.ResponseWriter, _ *http.Request)  {
	if err := homePage.Execute(w, users); err != nil {
		log.Printf("%v", err)
	}
}

func post(w http.ResponseWriter, r *http.Request)  {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v", err)
		return
	}
	users = append(users, User{
		r.FormValue("Name"),
		time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func test_ten()  {
	http.HandleFunc("/", home)
	http.HandleFunc("/post", post)
	log.Println("http://localhost:1234")
	log.Println(http.ListenAndServe(":1234", nil))
}


func home_cookie(w http.ResponseWriter, r *http.Request)  {
	var us []User
	cookie, err := r.Cookie("user")
	fmt.Println(cookie)
	if err == http.ErrNoCookie {
		v := time.Now().Format(time.RFC3339)
		str := []string{}
		c := http.Cookie{
			"user",
			v,
			"/",
			"",
			time.Time{},
			"",
			5,
			false,
			true,
			"",
			str,
		}
		http.SetCookie(w, &c)
		log.Printf("New user: %s", v)
	}else {
		//us = users[cookie.Value]
		log.Printf("Return user: %s", cookie.Value)
	}
	if err := homePage.Execute(w, us); err != nil {
		log.Printf("%v", err)
	}
}

func post_cookie(w http.ResponseWriter, r *http.Request)  {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v", err)
		return
	}
	users = append(users, User{
		r.FormValue("Name"),
		time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func test_cookie()  {
	http.HandleFunc("/", home_cookie)
	http.HandleFunc("/post", post_cookie)
	log.Println("http://localhost:1234")
	log.Println(http.ListenAndServe(":1234", nil))
}

func main() {

	//test_eight()

	//get
	//test_nine()

	//post  测试此处时请注释前面两个函数的调用，并且停止运行study_10
	test_ten()
	//测试cookie  需要注释test_ten()
	//test_cookie()
}
