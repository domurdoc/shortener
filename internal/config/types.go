package config

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

type LogLevel struct {
	text string
}

func (l LogLevel) String() string {
	return l.text
}

func (l *LogLevel) Set(value string) error {
	levels := []string{"debug", "info", "warn", "error", "fatal"}
	value = strings.ToLower(value)
	if slices.Contains(levels, value) {
		l.text = value
		return nil
	}
	return fmt.Errorf("must be one of (case-insensitive): %v", levels)
}

type NetAddress struct {
	Host string
	Port int
}

func (n NetAddress) String() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

func (n *NetAddress) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return fmt.Errorf("need address in a form ip:port")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	if port < 0 || port > (1<<16-1) {
		return fmt.Errorf("port must be greater than 0 and less than 65k")
	}
	host := strings.TrimSpace(parts[0])
	n.Port = port
	n.Host = host
	return nil
}

type URL struct {
	parsedURL *url.URL
}

func (u URL) String() string {
	return u.parsedURL.String()
}

func (u *URL) Set(value string) error {
	parsedURL, err := url.Parse(value)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("schema and host must be provided")
	}
	u.parsedURL = parsedURL
	return nil
}

type FilePath struct {
	path string
}

func (f FilePath) String() string {
	return f.path
}

func (f *FilePath) Set(value string) error {
	// TODO: add validation
	f.path = strings.TrimSpace(value)
	return nil
}

type DataSourceName struct {
	dsn string
}

func (d *DataSourceName) Set(value string) error {
	d.dsn = value
	return nil
}

func (d DataSourceName) String() string {
	return d.dsn
}
