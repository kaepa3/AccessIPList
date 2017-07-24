package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

//Config 設定ファイル
type Config struct {
	WhiteList []string
}

var ConfigPath = "./config.toml"
var config Config

func main() {
	readConfig()
	//ディレクトリループ
	root := "./"
	list := listFiles(root, root)
	ipList := make([]string, 0, 100)
	for _, val := range list {
		if strings.Contains(val, ".go") || strings.Contains(val, ".toml") {
			continue
		}
		ips := analyzeFile(val)
		for _, ip := range ips {
			ipList = append(ipList, ip)
		}
	}
	fmt.Println(strings.Join(ipList, "\n"))
}

func readConfig() {
	if Exists(ConfigPath) {
		_, err := toml.DecodeFile(ConfigPath, &config)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
func listFiles(rootPath, searchPath string) []string {
	fis, err := ioutil.ReadDir(searchPath)
	if err != nil {
		panic(err)
	}
	var list []string
	for _, fi := range fis {
		fullPath := filepath.Join(searchPath, fi.Name())
		if fi.IsDir() {
			result := listFiles(rootPath, fullPath)
			for _, val := range result {
				list = append(list, val)
			}
		} else {
			list = append(list, fullPath)
		}
	}
	return list
}

func analyzeFile(path string) []string {
	var list []string
	fp, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 4096)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		ip := pulloutIp(string(line))
		if ip == "" || isIncludeWhiteList(ip) {
			continue
		}
		if isIncludeIp(list, ip) {
			list = append(list, ip)
		}
	}
	return list
}
func isIncludeWhiteList(ip string) bool {
	for _, val := range config.WhiteList {
		if val == ip {
			return true
		}
	}
	return false
}

// ipの正規表現
var rep = regexp.MustCompile(`^([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})`)

func pulloutIp(line string) string {
	result := rep.FindAllString(line, -1)
	if 0 == len(result) {
		return ""
	}
	return result[0]
}

func isIncludeIp(list []string, ip string) bool {
	for _, val := range list {
		if ip == val {
			return false
		}
	}
	return true
}
