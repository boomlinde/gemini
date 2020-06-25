package main

import (
	"github.com/boomlinde/gemini/client"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	c := client.NewClient(filepath.Join(usr.HomeDir, ".gget", "pinned"))

	uri := os.Args[1]
	action := "get"
	if os.Args[1] == "pin" {
		uri = os.Args[2]
		action = "pin"
	}

	switch action {
	case "get":
		conn, err := c.Request(uri)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		header, err := client.GetHeader(conn)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(header)
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			log.Fatal(err)
		}
	case "pin":
		if err := c.Pin(uri); err != nil {
			log.Fatal(err)
		}
	}
}
