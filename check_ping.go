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
}

func main() {
    T := make(chan Result)
	ips := []*Result{{"192.168.0.2", false}, {"127.0.0.1", false}}
	for _, i := range ips {
		go i.check_ping(T)
	}
	for _, i := range ips {
        x := <- T
        fmt.Println(i)
		if x.ping {
			fmt.Printf("Able to connect: %s\n", x.ip)
		} else {
			fmt.Printf("Not able to connect: %s\n", x.ip)
		}
	}
}

func (res *Result) check_ping(t chan Result) {
	fmt.Printf("Pinging :%s \n", res.ip)
	command := fmt.Sprintf("ping -c 3 %s > /dev/null && echo true || echo false", res.ip)
	output, _ := exec.Command("/bin/bash", "-c", command).Output()
	val, _ := strconv.ParseBool(strings.TrimSpace(string(output)))
	if val {
        t <- Result{res.ip, true}
	} else {
        t <- Result{res.ip, false}
	}
}
