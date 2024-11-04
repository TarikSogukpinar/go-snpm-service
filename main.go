package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gosnmp/gosnmp"
)

func main() {
	g := &gosnmp.GoSNMP{
		Target:    "127.0.0.1",
		Port:      161,
		Community: "public",
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Retries:   3,
	}

	err := g.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Conn.Close()

	oids := []string{
		// Sistem Bilgileri
		"1.3.6.1.2.1.1.1.0",       // Sistem Açıklaması
		"1.3.6.1.2.1.1.3.0",       // Çalışma Süresi
		"1.3.6.1.2.1.1.5.0",       // Cihaz Adı
		"1.3.6.1.2.1.1.4.0",       // Sistem İletişim Bilgisi
		"1.3.6.1.2.1.1.6.0",       // Sistem Konumu
		"1.3.6.1.4.1.2021.4.5.0",  // Toplam Bellek
		"1.3.6.1.4.1.2021.4.6.0",  // Kullanılabilir Bellek
		"1.3.6.1.4.1.2021.4.11.0", // Kullanılan Bellek
		"1.3.6.1.4.1.2021.11.9.0", // CPU Kullanımı
		// Ağ Arayüz Bilgileri
		"1.3.6.1.2.1.2.2.1.8",    // Interface Operational Status
		"1.3.6.1.2.1.2.2.1.10",   // Incoming Octets
		"1.3.6.1.2.1.2.2.1.16",   // Outgoing Octets
		"1.3.6.1.2.1.2.2.1.14",   // In Errors
		"1.3.6.1.2.1.2.2.1.20",   // Out Errors
		"1.3.6.1.2.1.31.1.1.1.1", // Interface Name
		"1.3.6.1.2.1.2.2.1.4",    // Interface MTU
		"1.3.6.1.2.1.2.2.1.5",    // Interface Speed
		"1.3.6.1.2.1.2.2.1.7",    // Interface Status
		// CPU Bilgileri
		"1.3.6.1.4.1.2021.10.1.3.1", // CPU Load Average
		"1.3.6.1.4.1.2021.10.1.3.2", // User CPU Time
		"1.3.6.1.4.1.2021.10.1.3.3", // System CPU Time
		"1.3.6.1.4.1.2021.10.1.3.4", // Idle CPU Time
		// Disk Kullanım Bilgileri
		"1.3.6.1.2.1.25.2.3.1.6", // Disk Toplamı
		"1.3.6.1.2.1.25.2.3.1.5", // Disk Kullanımı
		// Disk I/O
		"1.3.6.1.4.1.2021.1.5.1.0", // Disk I/O
		// Ek Ağ Bilgileri
		"1.3.6.1.2.1.31.1.1.1.6", // Interface Speed
		"1.3.6.1.2.1.31.1.1.1.8", // Interface Admin Status
	}

	result, err := g.Get(oids)
	if err != nil {
		log.Fatalf("Get() err: %v", err)
	}

	for _, variable := range result.Variables {
		fmt.Printf("OID: %s ", variable.Name)

		switch variable.Type {
		case gosnmp.OctetString:
			fmt.Printf("Value: %s\n", string(variable.Value.([]byte)))
		case gosnmp.Integer:

			value := variable.Value.(int)
			switch variable.Name {
			case ".1.3.6.1.2.1.1.3.0":
				fmt.Printf("Uptime: %d seconds (%d minutes)\n", value/100, (value/100)/60)
			case ".1.3.6.1.4.1.2021.4.5.0":
				fmt.Printf("Total RAM: %d KB (%d MB)\n", value, value/1024)
			case ".1.3.6.1.4.1.2021.4.6.0":
				fmt.Printf("Available RAM: %d KB (%d MB)\n", value, value/1024)
			case ".1.3.6.1.4.1.2021.4.11.0":
				fmt.Printf("Used RAM: %d KB (%d MB)\n", value, value/1024)
			case ".1.3.6.1.4.1.2021.11.9.0":
				fmt.Printf("CPU Usage: %d%%\n", value)
			case ".1.3.6.1.2.1.2.2.1.8":
				fmt.Printf("Interface Operational Status: %d\n", value)
			case ".1.3.6.1.2.1.2.2.1.14":
				fmt.Printf("Incoming Errors: %d\n", value)
			case ".1.3.6.1.2.1.2.2.1.20":
				fmt.Printf("Outgoing Errors: %d\n", value)
			case ".1.3.6.1.4.1.2021.10.1.3.1":
				fmt.Printf("CPU Load Average: %d%%\n", value)
			case ".1.3.6.1.4.1.2021.10.1.3.2":
				fmt.Printf("User CPU Time: %d%%\n", value)
			case ".1.3.6.1.4.1.2021.10.1.3.3":
				fmt.Printf("System CPU Time: %d%%\n", value)
			case ".1.3.6.1.4.1.2021.10.1.3.4":
				fmt.Printf("Idle CPU Time: %d%%\n", value)
			case ".1.3.6.1.4.1.2021.1.5.1.0":
				fmt.Printf("Disk I/O: %d\n", value)
			default:
				fmt.Printf("Value: %d\n", value)
			}
		default:
			fmt.Printf("Value: %v\n", gosnmp.ToBigInt(variable.Value))
		}
	}
}
