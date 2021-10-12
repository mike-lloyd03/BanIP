package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var command string
	if len(os.Args) > 1 {
		if os.Args[1] == "ban" || os.Args[1] == "unban" {
			command = os.Args[1]
		} else {
			log.Fatalf("command \"%s\" not recognized", os.Args[1])
		}
	} else {
		command = "ban"
	}

	ipList, _ := getIPList("./ban_list")
	flushBanIPRules()

	for _, ip := range ipList {
		if command == "ban" {
			fmt.Println("banning ip address: ", ip)
			banIP(ip)
		} else if command == "unban" {
			fmt.Println("unbanning ip address: ", ip)
		}
	}
}

func getIPList(ipListPath string) ([]string, error) {
	file, err := os.Open(ipListPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ipList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ipList = append(ipList, scanner.Text())
	}
	return ipList, scanner.Err()
}

func banIP(sourceIP string) {
	cmd := exec.Command("iptables", "-I", "INPUT", "-s", sourceIP, "-j", "DROP", "-m", "comment", "--comment", "BanIP Rule")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func flushBanIPRules() {
	cmdString := "iptables --line-number -nL INPUT | grep 'BanIP Rule' | awk '{print $1}' | tac"
	out, err := exec.Command("bash", "-c", cmdString).Output()
	if err != nil {
		log.Fatal(err)
	}

	ruleNums := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	for _, ruleNum := range ruleNums {
		if ruleNum != "" {
			deleteRule(ruleNum)
		}
	}
}

func deleteRule(ruleNum string) {
	cmd := exec.Command("iptables", "-D", "INPUT", ruleNum)
	// fmt.Println(cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
