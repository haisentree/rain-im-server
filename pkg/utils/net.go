package utils

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// ip.me

func GetPublicIP() string {
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://ipv4.icanhazip.com")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	ip, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(ip))
}
