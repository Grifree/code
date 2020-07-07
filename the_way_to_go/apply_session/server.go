package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	cookie, _ := r.Cookie("username")
	fmt.Fprint(w, cookie)

}