package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

func worker(ports, results chan int, hostname string) { // the worker function accepts two channels (int) and a string hostname
	for p := range ports {
		var address string
		address = hostname + ":" + strconv.Itoa(p) // creating valid address to scan -> example.com:80
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0 // if port is closed then send 0
			continue
		}
		conn.Close()
		results <- p // if port is open then send the port (number)
	}
}

func main() {
	//Simple port scanning in GoLang (with using Parallellism)

	ports := make(chan int, 150) // ports channel holds 150 int (the larger number means faster but less reliable)
	results := make(chan int)    // create a seperate channel for communicating to the results from the worker to the main thread
	var openPorts []int          // openPorts slice to store the results

	// Get hostname from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Hostname to scan (example.com): ")
	hostname, _ := reader.ReadString('\n')  // We're not interested in any error for now. We assume that user gives a correct hostname
	hostname = strings.Trim(hostname, "\n") // remove '\n' from user's input

	// Get ports to scan from user
	fmt.Printf("How many ports do you want to scan? (Default is 1024) : ")
	var portNumber int = 0
	fmt.Scanf("%d", &portNumber)
	if portNumber == 0 {
		portNumber = 1024 // default
	}

	// Print a message
	fmt.Printf("\n+++ Port scanning started for %v from port 1 to %v ...\n", hostname, portNumber)

	// Create all the workers (150)
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, hostname)
	}

	// Communication between ports(which using by worker func) and portNumber(which using by main func) by using channels in goLang
	go func() {
		for i := 0; i <= portNumber; i++ {
			ports <- i
		}
	}()

	// Store open port in openPorts slice
	for i := 0; i < portNumber; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openPorts)
	fmt.Println("+++ Scanning has been done !")
	fmt.Println("+++ Results -->")
	fmt.Println("")
	fmt.Println("---------")
	fmt.Println("")
	for _, port := range openPorts {
		fmt.Printf("Port number %v is open\n", port)
	}
}
