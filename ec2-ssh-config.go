package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"time"

	"github.com/chdorner/ec2-ssh-config/sshconfig"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/ec2"
)

var (
	Version        string
	aws_access_key string
	aws_secret_key string
	ec2_region     string
	ssh_user       string
	prefix         string
	config_file    string
	backup         bool
	version        bool

	auth aws.Auth
	c    *ec2.EC2
)

type Instance struct {
	Host string
	Name string
}

func main() {
	if version {
		fmt.Println("ec2-ssh-config", Version)
		os.Exit(0)
	}

	c, err := sshconfig.Parse(config_file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, i := range getInstances() {
		host := fmt.Sprintf("%s%s", prefix, i.Name)
		h := c.FindHost(host)
		if h == nil {
			h = sshconfig.NewHost(host)
			c.Hosts = append(c.Hosts, h)
		}
		h.Attrs["HostName"] = i.Host
		if ssh_user != "" {
			h.Attrs["User"] = ssh_user
		}
	}

	if backup {
		err = doBackup()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	c.Store()
}

func getInstances() []*Instance {
	auth := aws.Auth{aws_access_key, aws_secret_key}
	c = ec2.New(auth, aws.Regions[ec2_region])

	filter := ec2.NewFilter()
	filter.Add("instance-state-name", "running")
	resp, err := c.Instances(nil, filter)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	instances := []*Instance{}
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			name := ""
			for _, tag := range instance.Tags {
				if tag.Key == "Name" {
					name = tag.Value
					break
				}
			}
			if name != "" {
				i := &Instance{instance.DNSName, name}
				instances = append(instances, i)
			}
		}
	}

	return instances
}

func doBackup() error {
	src, err := os.Open(config_file)
	if err != nil && os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(fmt.Sprintf("%s.%s", config_file, time.Now().Format("20060102150405")))
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return errors.New("Could not backup config file, aborting..")
	}

	return nil
}

func init() {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	default_config := fmt.Sprintf("%s/.ssh/config", user.HomeDir)

	flag.StringVar(&aws_access_key, "aws_access_key", os.Getenv("AWS_ACCESS_KEY_ID"), "aws access key, defaults to $AWS_ACCESS_KEY_ID")
	flag.StringVar(&aws_secret_key, "aws_secret_key", os.Getenv("AWS_SECRET_ACCESS_KEY"), "aws secret access key, defaults to $AWS_SECRET_ACCESS_KEY")
	flag.StringVar(&ec2_region, "region", os.Getenv("AWS_EC2_REGION"), "ec2 region, defaults to $AWS_EC2_REGION")
	flag.StringVar(&prefix, "prefix", "", "prefix for host ($prefix$instance_name, ex. acl-mywebserver)")
	flag.StringVar(&ssh_user, "user", "", "ssh user to use")
	flag.StringVar(&config_file, "f", default_config, "path to ssh config file")
	flag.BoolVar(&backup, "backup", true, "create a backup file")
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.Parse()
}
