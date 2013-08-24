package sshconfig

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type SSHConfig struct {
	Path  string
	Hosts []*SSHHost
}

type SSHHost struct {
	Host  string
	Attrs map[string]string
}

func Parse(path string) (*SSHConfig, error) {
	c := &SSHConfig{
		Path:  path,
		Hosts: nil,
	}

	file, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		return c, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var current *SSHHost
	for scanner.Scan() {
		line := scanner.Text()
		r := regexp.MustCompile(`^[\ ]*([a-zA-z]+) (.*)$`)
		kv := r.FindStringSubmatch(line)
		if len(kv) != 3 {
			continue
		}
		k, v := kv[1], kv[2]

		if k == "Host" {
			current = NewHost(v)
			c.Hosts = append(c.Hosts, current)
		} else if current != nil {
			current.Attrs[k] = v
		}
	}

	return c, nil
}

func NewHost(host string) *SSHHost {
	return &SSHHost{
		Host:  host,
		Attrs: make(map[string]string),
	}
}

func (c *SSHConfig) Store() {
	f, err := os.Create(c.Path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	for _, h := range c.Hosts {
		s := fmt.Sprintf("Host %s\n", h.Host)
		for k, v := range h.Attrs {
			s += fmt.Sprintf("  %s %s\n", k, v)
		}
		s += "\n"
		_, err = f.WriteString(s)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func (c *SSHConfig) FindHost(host string) *SSHHost {
	for _, h := range c.Hosts {
		if h.Host == host {
			return h
		}
	}
	return nil
}
