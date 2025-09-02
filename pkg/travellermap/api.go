package travellermap

import (
	"fmt"
	"io"
	"net/http"
)

func Get(url string) ([]byte, error) {
	fmt.Println("send get to", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get returned: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("responce received")
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http returned status: %s", resp.Status)
	}
	fmt.Println("read responce")
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, fmt.Errorf("response read error: %v", err)
	}

	return data, nil
}
