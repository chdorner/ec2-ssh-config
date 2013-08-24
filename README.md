# ec2-ssh-config

A little utility for syncing down EC2 instances to your local SSH config.

```shell
$ ec2-ssh-config -h
Usage of ec2-ssh-config:
  -aws_access_key="": aws access key, defaults to $AWS_ACCESS_KEY_ID
  -aws_secret_key="": aws secret access key, defaults to $AWS_SECRET_ACCESS_KEY
  -backup=true: create a backup file
  -f="~/.ssh/config": path to ssh config file
  -prefix="": prefix for host ($prefix$instance_name, ex. acl-mywebserver)
  -region="": ec2 region, defaults to $AWS_EC2_REGION
  -user="": ssh user to use
  -version=false: print version and exit
```
