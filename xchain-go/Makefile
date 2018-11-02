BINARY := xchain
PROG := $(BINARY)
PROJ_HOME=${PWD}

VERSIONFILE := version.go
BUILD_DATE := $(shell date '+%Y-%m-%d [%H:%M:%S]')
GIT_VERSION := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git symbolic-ref -q --short HEAD)
GIT_SHOW:= $(shell git show --pretty='%h <%an> %ai | %s' -s HEAD)
# GIT_INFO:= $(GIT_BRANCH):$(GIT_VERSION)


build:
	make version
	@echo "====> Go build"
	@go build -o $(BINARY)


start:
	./xchain printusage


version:
	@echo "package main" > $(VERSIONFILE)
	@echo "const (" >> $(VERSIONFILE)
	@echo "  GitVersion = \"$(GIT_SHOW)\"" >> $(VERSIONFILE)
	@echo "  BuildDate = \"$(BUILD_DATE)\"" >> $(VERSIONFILE)
	@echo ")" >> $(VERSIONFILE)



	
