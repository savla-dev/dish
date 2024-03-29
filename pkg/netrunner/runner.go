package netrunner

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"dish/pkg/config"
	"dish/pkg/socket"
)

// RawConnect function for direct host:port socket check
func RawConnect(socket socket.Socket) error {
	endpoint := net.JoinHostPort(socket.Host, strconv.Itoa(socket.Port))
	timeout := time.Duration(time.Second * time.Duration(config.Timeout))

	if config.Verbose {
		log.Println("runner: rawconnect: " + endpoint)
	}

	// open the socket
	conn, err := net.DialTimeout("tcp", endpoint, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

// checkHTTPCode function for response and expected HTTP codes comparison
// panics if it fails to convert expected code to int
func checkHTTPCode(responseCode int, expectedCodes []int) bool {
	for _, code := range expectedCodes {
		if responseCode == code {
			return true
		}
	}
	return false
}

// CheckSite executes test over HTTP/S endpoints exclusively
func CheckSite(socket socket.Socket) (bool, int, error) {
	// config http client
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}
	url := socket.Host + ":" + strconv.Itoa(socket.Port) + socket.PathHTTP

	if config.Verbose {
		log.Println("runner: checksite:", url)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, 0, err
	}
	req.Header.Set("User-Agent", "savla-dish/1.5")

	// open socket --- Head to url
	//resp, err := client.Head(url)
	resp, err := client.Do(req)
	if err != nil {
		return false, 0, err
	}

	// fetch StatusCode for HTTP expected code comparison
	if resp != nil {
		defer resp.Body.Close()
		return checkHTTPCode(resp.StatusCode, socket.ExpectedHTTPCodes), resp.StatusCode, nil
	}

	return true, resp.StatusCode, nil
}
