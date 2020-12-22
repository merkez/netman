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
	IFUP                 = "ifup"
	INTERFACES_FILE      = "/etc/network/interfaces"
	DOWN_INTERFACES_FILE = "/etc/network/downinterfaces"
	TEMPLATE_PATH        = "/home/vagrant/netman/interfaces.tmpl"
)

type Interfaces struct {
	NetInt []NetInterface
}
type NetInterface struct {
	InetName string
}

func main() {

	var netInterfaces []NetInterface
	ineti, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, inter := range ineti {
		if !strings.Contains(inter.Flags.String(), "up") {
			netInterfaces = append(netInterfaces, NetInterface{InetName: inter.Name})
		}
	}
	downInterfaces := Interfaces{NetInt: netInterfaces}

	info := createInterfacesSetting(downInterfaces)
	fmt.Printf("%s\n", info)
	if err := createFile(INTERFACES_FILE, info); err != nil {
		panic(err)
	}
	var tempInterfaces []string
	for _, i := range downInterfaces.NetInt {
		tempInterfaces = append(tempInterfaces, i.InetName)
		if err := execUp(i.InetName); err != nil {
			panic(err)
		}
	}
	c := strings.Join(tempInterfaces, ",")
	if err := createFile(DOWN_INTERFACES_FILE, c); err != nil {
		panic(err)
	}

	if err := deleteDefaultRule(); err != nil {
		panic(err)
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
	tmpl := template.Must(template.ParseFiles(TEMPLATE_PATH))
	tmpl.Execute(&tpl, interfaces)
	return tpl.String()
}

func createFile(path string, content string) error {

	f, err := os.Create(path)
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

func deleteDefaultRule() error {
	//deleting the default rule set automatically in order to have Internet.
	cmd := exec.Command("sudo", "ip", "route", "del", "default")
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}
