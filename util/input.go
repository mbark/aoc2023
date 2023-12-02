package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetInput(day int) string {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/2023/day/%d/input", day), nil)
	must(err)
	req.Header.Set("cookie", "session="+os.Getenv("SESSION"))

	res, err := http.DefaultClient.Do(req)
	must(err)

	if res.StatusCode != 200 {
		panic(fmt.Sprintf("non-200 status code: %d", res.StatusCode))
	}

	bytes, err := io.ReadAll(res.Body)
	must(err)

	return string(bytes)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
