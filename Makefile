.DEFAULT: build

VERSION := 0.1.0
TARGET  := ec2-ssh-config
LDFLAGS := -ldflags "-X main.Version $(VERSION)"

build: $(TARGET)

$(TARGET):
	go build -o $(TARGET) $(LDFLAGS)

.PHONY: $(TARGET)
