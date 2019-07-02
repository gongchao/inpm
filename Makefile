OS := $(shell uname -a|awk '{print $$1}'|tr '[A-Z]' '[a-z]')

build:
	GOOS=$(OS) go build

