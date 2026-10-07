// Command quarel-probe is a minimal loopback HTTP probe for distroless images.
package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: quarel-probe identity|community|http(s)://localhost[:port]/health")
		os.Exit(2)
	}
	target := os.Args[1]
	if target == "identity" || target == "community" {
		mode := os.Getenv("QUAREL_TLS")
		if mode == "" {
			dataDir := os.Getenv("QUAREL_DATA_DIR")
			if dataDir == "" {
				dataDir = "./data"
			}
			data, err := os.ReadFile(filepath.Join(dataDir, "settings.json"))
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				fail("load settings: " + err.Error())
			}
			values := map[string]string{}
			if err == nil {
				if err := json.Unmarshal(data, &values); err != nil {
					fail("load settings: " + err.Error())
				}
			}
			mode = values["QUAREL_TLS"]
		}
		if mode == "" {
			if target == "identity" {
				mode = "off"
			} else {
				mode = "self-signed"
			}
		}
		scheme := "https"
		if mode == "off" {
			scheme = "http"
		}
		port := "8080"
		if target == "community" {
			port = "8090"
		}
		target = scheme + "://127.0.0.1:" + port + "/v1/health"
	}
	u, err := url.Parse(target)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" {
		fail("invalid health URL")
	}
	ip := net.ParseIP(u.Hostname())
	if u.Hostname() != "localhost" && (ip == nil || !ip.IsLoopback()) {
		fail("health URL must target loopback")
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if u.Scheme == "https" {
		// This binary is only used for in-container loopback probes against
		// locally generated self-signed certificates.
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	client := &http.Client{
		Timeout:   2 * time.Second,
		Transport: transport,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(u.String())
	if err != nil {
		fail(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		fail(fmt.Sprintf("health endpoint returned HTTP %d", resp.StatusCode))
	}
}

func fail(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
