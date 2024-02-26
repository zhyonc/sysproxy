package util

import (
	"bufio"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"
)

func CreateFile(path string, bytes []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}

func ReadFileAll(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func ReadFileLine(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func DownloadFile(url string) ([]byte, error) {
	httpClient := &http.Client{}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func DecodeBase64(srcBytes []byte, seq string) ([]string, error) {
	dstBytes, err := base64.StdEncoding.DecodeString(string(srcBytes))
	if err != nil {
		return nil, err
	}
	dstStr := string(dstBytes)
	return strings.Split(dstStr, seq), nil
}
