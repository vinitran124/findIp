package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

type AwsIp struct {
	SyncToken  string `json:"syncToken"`
	CreateDate string `json:"createDate"`
	Prefixes   []struct {
		IPPrefix string `json:"ip_prefix"`
		Region   string `json:"region"`
		Service  string `json:"service"`
	} `json:"prefixes"`
	Ipv6Prefixes []struct {
		Ipv6Prefix string `json:"ipv6_prefix"`
		Region     string `json:"region"`
		Service    string `json:"service"`
	} `json:"ipv6_prefixes"`
}

const allFile = "./all.txt"
const awsFile = "./ip-ranges.json"

func main() {
	ipList, err := getIpList()
	if err != nil {
		log.Println(err)
		return
	}

	awsIp, err := getAwsIpNetwork()
	if err != nil {
		log.Println(err)
		return
	}

	for _, ip := range ipList {
		if isInAWS(ip, awsIp) == false {
			log.Println(ip)
		}
	}

	log.Println("end")
	return
}

func getAwsIpNetwork() ([]string, error) {
	jsonFile, err := os.Open(awsFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	awsData := AwsIp{}
	if err := json.Unmarshal(jsonData, &awsData); err != nil {
		return nil, err
	}

	var ipNetworkList []string
	for _, ip := range awsData.Prefixes {
		ipNetworkList = append(ipNetworkList, getNetworkPartIp(ip.IPPrefix))
	}

	///add private ip
	ipNetworkList = append(ipNetworkList, "10")
	ipNetworkList = append(ipNetworkList, "172")
	ipNetworkList = append(ipNetworkList, "192")

	ipNetworkList = removeDuplicateStr(ipNetworkList)
	return ipNetworkList, err
}

func getIpList() ([]string, error) {
	file, err := os.Open(allFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func getNetworkPartIp(ip string) string {
	ips := strings.Split(ip, ".")
	return ips[0]
}

func isInAWS(ip string, awsIp []string) bool {
	ip = getNetworkPartIp(ip)
	for _, aws := range awsIp {
		if ip == aws {
			return true
		}
	}
	return false
}
