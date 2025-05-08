package scrapper

import (
	"fmt"
	"io"
	"net/http"
)

func GetHtml(url string) (string, error) {
	fmt.Printf("Visiting %s\n", url)

	res, err := http.Get(url)
	if err != nil {
		err := fmt.Errorf("Error while doing Get request: %v", err)
		return "", err
	}
	defer res.Body.Close()

	fmt.Printf("Response status: %s\n", res.Status)

	if res.StatusCode > 299 {
		err := fmt.Errorf("Error: %s", res.Status)
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		err := fmt.Errorf("Error while reading body: %v", err)
		return "", err
	}

	return string(body), nil
}
