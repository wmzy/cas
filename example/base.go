package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/wmzy/go-cas"
)

type handler struct {
}

func (self *handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	cas2 := cas.NewV2("https://your-cas-server.com")
	if req.URL.Path == "/login" {
		http.Redirect(res, req, cas2.GetLoginURL("http://localhost:8080/admin"), http.StatusFound)
	}
	if req.URL.Path == "/admin" {
		ticket := req.URL.Query().Get("ticket")
		user, err := cas2.GetUser(ticket, "http://localhost:8080/admin")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(
			res,
			user.User,
			"\n",
			user.Attributes.Get("email"),
			"\n",
			user.Attributes.Get("name"),
			"\n",
			user.Attributes.Get("displayName"),
			"\n",
			user.Attributes.Get("departmentName"),
			"\n",
			user.Attributes.Get("sessionTicket"),
			user,
		)
	}
}

func main() {
	var h = &handler{}

	var server = &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	fmt.Println("server run at: http://localhost:8080/login")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("run error: %v", err)
		os.Exit(1)
	}
}
