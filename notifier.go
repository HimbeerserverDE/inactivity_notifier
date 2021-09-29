/*
inactivity_notifier sends an email if there is no activity for a customisable amount of time.
Usage:
	inactivity_notifier <server <timeout> <mailserver> <sender> <password> <recipient> <message> | renew <remote>>

where timeout is the time of inactivity after which the email is sent,
recipient is the email address to send the message to,
message is the content of the email
and remote is the address of the inactivity_notifier server.
*/
package main

import (
	"fmt"
	"net"
	"net/smtp"
	"os"
	"time"
)

func usage() {
	fmt.Println("Usage: inactivity_notifier <server <timeout> <mailserver> <sender> <password> <recipient> <message> | renew <remote>>")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "server":
		if len(os.Args) != 8 {
			usage()
		}

		duration, err := time.ParseDuration(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		t := time.NewTimer(duration)
		go func() {
			<-t.C
			fmt.Println("inactivity timeout reached")

			auth := smtp.PlainAuth("", os.Args[4], os.Args[5], os.Args[3])
			msg := []byte("To: " + os.Args[6] + "\r\n" +
				"Subject: Inactivity notification\r\n" +
				"\r\n" +
				os.Args[7] + "\r\n")

			err := smtp.SendMail(os.Args[3]+":25", auth, os.Args[4], []string{os.Args[6]}, msg)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			os.Exit(0)
		}()

		laddr, err := net.ResolveUDPAddr("udp", ":15115")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		lc, err := net.ListenUDP("udp", laddr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer lc.Close()

		fmt.Println("ready")
		for {
			buf := make([]byte, 1024)
			if _, _, err := lc.ReadFrom(buf); err != nil {
				fmt.Println(err)
				continue
			}

			t.Reset(duration)
			fmt.Println("timeout reset")
		}
	case "renew":
		if len(os.Args) != 3 {
			usage()
		}

		addr, err := net.ResolveUDPAddr("udp", os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer conn.Close()

		buf := []byte{0x00}
		if _, err := conn.Write(buf); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
