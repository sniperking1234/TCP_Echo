package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var (
	port   int
	prefix string

	rootCmd = &cobra.Command{
		Use:   "echo",
		Short: "echo ist a tcp reply server",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			log.Printf("starting listening %s", port)
			echo()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&prefix, "prxfix", "", "return string prefix")
	rootCmd.Flags().IntVar(&port, "port", 0, "listening port")
	rootCmd.MarkFlagRequired("port")
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func echo() {
	listener, err := net.Listen("tcp", string(port))
	if err != nil {
		log.Fatal("Error listening", err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	reader := bufio.NewReader(conn)
	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				log.Println("failed to read data, err:", err.Error())
			}
			return
		}
		log.Println("requests is", bytes)

		line := fmt.Sprintf("%s%s", prefix, reader)
		log.Println("response is", line)
		conn.Write([]byte(line))
	}
}
