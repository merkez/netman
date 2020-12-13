package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	IFUP = "ifup"
)

type Interfaces struct {
	NetInt []NetInterface
}
type NetInterface struct {
	InetName string
}

func main()  {

	var netInterfaces []NetInterface
	 ineti, err :=  net.Interfaces()
	 if err != nil {
	 	panic(err)
	 }
	 for _, inter := range ineti {
		 if !strings.Contains(inter.Flags.String(), "up"){
			netInterfaces = append(netInterfaces, NetInterface{InetName:inter.Name})
		 }
	 }
	 downInterfaces := Interfaces{NetInt: netInterfaces}
	 info := createInterfacesSetting(downInterfaces)
	 fmt.Printf("%s\n", info)
	 createInterfacesFile(info)

	 for _, i := range downInterfaces.NetInt {
	 	if err := execUp(i.InetName); err != nil {
	 		panic(err)
		}
	 }

}

func execUp(iname string) error {
	cmd := exec.Command(IFUP, iname)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()

	return err
}

func createInterfacesSetting(interfaces Interfaces) string {
	var tpl bytes.Buffer
	tmpl := template.Must(template.ParseFiles("/home/ubuntu/netmanage/interfaces.tmpl"))
	tmpl.Execute(&tpl, interfaces)
	return tpl.String()
}

func createInterfacesFile(content string) error {

	f, err := os.Create("/etc/network/interfaces")
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}