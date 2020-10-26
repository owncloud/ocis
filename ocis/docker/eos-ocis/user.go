package main

import (
	"fmt"
	"os"
	gouser "os/user"
)

func main() {
	username := os.Args[1]
	user, err := gouser.Lookup(username)
	fmt.Printf("User %#v %#v %#v\n", username, user, err)
}
