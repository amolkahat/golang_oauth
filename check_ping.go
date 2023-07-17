package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

)

type Result struct {
	ip   string
	ping bool
	ssh  bool
	vlan bool
}

func main() {
    T := make(chan bool)
	ips := []*Result{{"192.168.0.2", false, false, false}, {"127.0.0.1", false, false, false}}
	for _, i := range ips {
		go i.check_ping(T)
	}
	for _, i := range ips {
        x := <- T
		if x {
			fmt.Printf("Able to connect: %s\n", i.ip)
		} else {
			fmt.Printf("Not able to connect: %s\n", i.ip)
		}
	}
}

func (res *Result) check_ping(t chan bool) {
	fmt.Printf("Pinging :%s \n", res.ip)
	command := fmt.Sprintf("ping -c 3 %s > /dev/null && echo true || echo false", res.ip)
	fmt.Println(command)
	output, _ := exec.Command("/bin/bash", "-c", command).Output()
	val, _ := strconv.ParseBool(strings.TrimSpace(string(output)))
    fmt.Println(val)
	if val {
		res.ping = true
        t <- true
		fmt.Printf("Set ping to true\n")
	} else {
		fmt.Printf("Set ping to false\n")
        t <- false
		res.ping = false
	}
}
