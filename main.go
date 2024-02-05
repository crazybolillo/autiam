package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const version = "v1.0.0"

var client = &http.Client{Timeout: 10 * time.Second}

type Credentials struct {
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	Token           string `json:"Token"`
}

func main() { os.Exit(cli()) }

func cli() int {
	showVersion := flag.Bool("version", false, "Print version and exit")
	expiration := flag.Int("expire", 300, "TTL for metadata token")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return 0
	}

	args := flag.Args()
	if len(args) != 1 {
		slog.Error("a single positional argument is required")
		return 1
	}

	metaToken, err := metadataToken(*expiration)
	if err != nil {
		slog.Error("failed to obtain metadata token", "error", err.Error())
		return 1
	}

	credentials, err := securityCredentials(args[0], metaToken)
	if err != nil {
		slog.Error("failed to obtain credentials", "error", err.Error())
	}

	fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", credentials.AccessKeyID)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", credentials.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=%s\n", credentials.Token)

	return 0
}

func securityCredentials(role string, token string) (Credentials, error) {
	targetUrl, _ := url.Parse("http://169.254.169.254/latest/meta-data/iam/security-credentials/" + role)
	request := &http.Request{
		Method: http.MethodGet,
		URL:    targetUrl,
		Header: map[string][]string{
			"X-aws-ec2-metadata-token": {token},
		},
	}

	resp, err := client.Do(request)
	if err != nil {
		return Credentials{}, err
	}

	defer resp.Body.Close()
	var credentials Credentials

	err = json.NewDecoder(resp.Body).Decode(&credentials)
	if err != nil {
		return Credentials{}, err
	}

	return credentials, nil
}

func metadataToken(ttl int) (string, error) {
	targetUrl, _ := url.Parse("http://169.254.169.254/latest/api/token")
	request := &http.Request{
		Method: http.MethodPut,
		URL:    targetUrl,
		Header: map[string][]string{
			"X-aws-ec2-metadata-token-ttl-seconds": {strconv.Itoa(ttl)},
		},
	}

	resp, err := client.Do(request)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}
