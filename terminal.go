package main

import (
	"bufio"
	//"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"github.com/olekukonko/tablewriter"
	"github.com/tj/go-spin"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ChannelMap struct {
	Start, End int
	CH500      int
}

var sub_band_channel_map = map[int]ChannelMap{
	1: ChannelMap{0, 7, 64},
	2: ChannelMap{8, 15, 65},
	3: ChannelMap{16, 23, 66},
	4: ChannelMap{24, 31, 67},
	5: ChannelMap{32, 39, 68},
	6: ChannelMap{40, 47, 69},
	7: ChannelMap{48, 55, 70},
	8: ChannelMap{56, 63, 71},
}

func term_read(r io.Reader) {
	for {
		if _, err := io.Copy(os.Stdout, r); err != nil {
			log.Fatal(err)
		}
	}
}

func term_write(w io.Writer) {
	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		cmd := fmt.Sprintf("%s\r\n", strings.TrimSpace(line))
		if _, err := w.Write([]byte(cmd)); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func send_cmd(dev io.ReadWriteCloser, s string) (string, error) {
	cmd := fmt.Sprintf("%s\r\n", strings.TrimSpace(s))
	dev.Write([]byte(cmd))

	time.Sleep(time.Millisecond * 500)

	buf := make([]byte, 32)
	n, err := dev.Read(buf)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading from serial port: ", err)
			return "", err
		}
	}
	return string(buf[:n]), nil
}

func printCMD(dev io.ReadWriteCloser, d string, cmd string) {

	c, _ := send_cmd(dev, cmd)

	fmt.Printf("%-20s: %s\n", d, strings.TrimSpace(c))
}

func send_cmd_ok(dev io.ReadWriteCloser, cmd string) error {
	s, _ := send_cmd(dev, cmd)
	if strings.TrimSpace(s) != "ok" {
		return fmt.Errorf("Error sending command: [%s] %s", strings.TrimSpace(cmd), s)
	}
	return nil
}

func set_subband(dev io.ReadWriteCloser, band int) error {
	if band < 0 || band > 8 {
		return errors.New("Invalid range for sub-band's")
	}

	for k, c := range sub_band_channel_map {
		on := "off"
		for i := c.Start; i <= c.End; i++ {

			if k == band {
				on = "on"
			}
			cmd := fmt.Sprintf("mac set ch status %d %s\r\n", i, on)
			err := send_cmd_ok(dev, cmd)
			if err != nil {
				return err
			}
		}
		cmd := fmt.Sprintf("mac set ch status %d %s\r\n", c.CH500, on)
		err := send_cmd_ok(dev, cmd)
		if err != nil {
			return err
		}
	}
	cmd := fmt.Sprintf("mac save\r\n")
	err := send_cmd_ok(dev, cmd)
	if err != nil {
		return err
	}
	return nil
}

func usage() {

	fmt.Println("LoraWAN Microchip RN2903 Terminal")
	fmt.Println("by Pawel Pastuszak @ gmail . com")
	fmt.Println("lora-term usage:")
	flag.PrintDefaults()
	os.Exit(-1)
}

func main() {

	port := flag.String("port", "/dev/tty.usbmodem1421", "serial port to test (/dev/tty.usbmodem1421, /dev/ttyUSB0)")
	baud := flag.Uint("baud", 57600, "Baud rate")
	stopbits := flag.Uint("stopbits", 1, "Stop bits")
	databits := flag.Uint("databits", 8, "Data bits")
	minread := flag.Uint("minread", 1, "Minimum read count")
	term := flag.Bool("term", false, "Terminal Emu mode")
	subband := flag.Int("sub-band", -1, "Set specific sub band channels on")

	flag.Parse()

	if *port == "" {
		fmt.Println("Must specify port")
		usage()
	}

	parity := serial.PARITY_NONE

	options := serial.OpenOptions{
		PortName:        *port,
		BaudRate:        *baud,
		DataBits:        *databits,
		StopBits:        *stopbits,
		MinimumReadSize: *minread,
		ParityMode:      parity,
	}

	dev, err := serial.Open(options)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	if *term == true {
		go term_read(dev)
		term_write(dev)
	} else {

		if *subband != -1 {
			err := set_subband(dev, *subband)
			if err != nil {
				fmt.Println(err)
			}
			os.Exit(0)
		}

		printCMD(dev, "Version", "sys get ver")
		printCMD(dev, "Dev Addr", "mac get devaddr")
		printCMD(dev, "Dev EUI", "mac get deveui")
		printCMD(dev, "APP EUI", "mac get appeui")
		//	printCMD(dev, "Network Session Key", "mac get nwkskey")
		//	printCMD(dev, "App Sesion Key", "mac get appskey")
		//	printCMD(dev, "App Key", "mac get appkey")
		fmt.Println("---------------------------------------------------------")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Sub Chan", "0", "1", "2", "3", "4", "5", "6", "7", "500MHz"})
		table.SetBorder(false) // Set Border to false
		spinner := spin.New()
		spinner.Set(spin.Spin1)
		// To store the keys in slice in sorted order
		var keys []int
		for k, _ := range sub_band_channel_map {
			keys = append(keys, k)

		}
		sort.Ints(keys)

		for _, k := range keys {
			row := make([]string, 10)
			row[0] = strconv.Itoa(k)
			v := 1
			for i := sub_band_channel_map[k].Start; i <= sub_band_channel_map[k].End; i++ {
				fmt.Printf("\r  \033[36mcomputing sub channels\033[m %s ", spinner.Next())

				cmd := fmt.Sprintf("mac get ch status %d\r\n", i)
				s, _ := send_cmd(dev, cmd)
				row[v] = s
				v++
			}
			cmd := fmt.Sprintf("mac get ch status %d\r\n", sub_band_channel_map[k].CH500)
			s, _ := send_cmd(dev, cmd)
			row[9] = s
			table.Append(row)
		}
		fmt.Println("")

		table.Render()
	}
}
