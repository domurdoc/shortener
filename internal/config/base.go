package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

func (n NetAddress) String() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

func (n *NetAddress) Set(FlagValue string) error {
	parts := strings.Split(FlagValue, ":")
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

func (u *URL) Set(FlagValue string) error {
	parsedURL, err := url.Parse(FlagValue)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("schema and host must be provided")
	}
	u.parsedURL = parsedURL
	return nil
}

type Options struct {
	Addr    NetAddress
	BaseURL URL
}

func (o Options) String() string {
	return fmt.Sprintf("addr = %q; baseURL = %q", o.Addr, o.BaseURL)
}

func NewOptions(defaultAddr, defaultBaseURL string) *Options {
	baseURL := URL{}
	if err := baseURL.Set(defaultBaseURL); err != nil {
		panic(err)
	}
	addr := NetAddress{}
	if err := addr.Set(defaultAddr); err != nil {
		panic(err)
	}
	return &Options{Addr: addr, BaseURL: baseURL}
}
