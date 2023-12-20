package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/Aruminium/goecat/pkg/ethercat/command"
	"github.com/Aruminium/goecat/pkg/ethercat/datagram"
	"github.com/Aruminium/goecat/tools/packet"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	device       string                    = "en7"
	snapshot_len int32                     = 1024
	promiscuous  bool                      = false
	timeout      time.Duration             = 30 * time.Second
	options      gopacket.SerializeOptions = gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	handle *pcap.Handle
	err    error
)

func main() {
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packet, err := packet.NewEtherCATPacket(device)
	if err != nil {
		log.Fatal(err)
	}

	lrcm1 := datagram.NewLrcm(1, 0, 1)
	ecatDatagram1 := datagram.Datagram{
		Command: command.LRD,
		Index:   uint8(0x00),
		Address: uint32(0x00000000),
		LRCM:    lrcm1,
		IRQ:     uint16(0x0000),
		Data:    []byte{0x00},
		WKC:     uint16(0x0000),
	}

	lrcm2 := datagram.NewLrcm(1, 0, 1)
	ecatDatagram2 := datagram.Datagram{
		Command: command.LRD,
		Index:   uint8(0x00),
		Address: uint32(0x00000000),
		LRCM:    lrcm2,
		IRQ:     uint16(0x0000),
		Data:    []byte{0x00},
		WKC:     uint16(0x0000),
	}

	lrcm3 := datagram.NewLrcm(0, 0, 1)
	ecatDatagram3 := datagram.Datagram{
		Command: command.LRD,
		Index:   uint8(0x00),
		Address: uint32(0x00000000),
		LRCM:    lrcm3,
		IRQ:     uint16(0x0000),
		Data:    []byte{0x00},
		WKC:     uint16(0x0000),
	}

	packet.Ecat.AppendDatagram(ecatDatagram1)
	packet.Ecat.AppendDatagram(ecatDatagram2)
	packet.Ecat.AppendDatagram(ecatDatagram3)

	data, err := packet.Send(handle, options)
	if err != nil {
		fmt.Printf("[-] Error while sending: %s\n", err.Error())
	}
	fmt.Println(hex.Dump(data))
}