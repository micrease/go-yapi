package main

import (
	"bufio"
	"fmt"
	yapi "github.com/micrease/go-yapi"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

func getAPI() (string, string) {
	r := bufio.NewReader(os.Stdin)

	fmt.Print("API URL: ")
	apiURL, _ := r.ReadString('\n')
	apiURL = strings.TrimSpace(apiURL)

	fmt.Print("API Token: ")
	byteAPIToken, _ := terminal.ReadPassword(int(syscall.Stdin))
	apiToken := strings.TrimSpace(string(byteAPIToken))
	fmt.Println()
	return apiURL, apiToken
}

func main() {
	apiURL, apiToken := getAPI()
	yapiClient, err := yapi.NewClient(apiURL, apiToken)
	if err != nil {
		panic(err)
	}

	project, err := yapiClient.Project.Get()

	fmt.Println("project", project.ToString())
	if err != nil {
		panic(err)
	}

	catMenus, err := yapiClient.CatMenu.Get(project.Data.ID)
	fmt.Println("catMenus", catMenus.ToString())
	for _, catmenu := range catMenus.Data {
		interfaceListParam := new(yapi.InterfaceListParam)
		interfaceListParam.CatID = catmenu.ID
		interfaceListParam.Page = 1
		interfaceListParam.Limit = 1000
		interfaces, err := yapiClient.Interface.GetList(interfaceListParam)

		fmt.Println("interfaces", interfaces.ToString())
		fmt.Println(err)
		for _, i := range interfaces.Data.List {
			result, err := yapiClient.Interface.Get(i.ID)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("result", result.ToString())
		}
	}
}
