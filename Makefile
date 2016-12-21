#
# Shorthand a text label expander.
#
PROJECT = shorthand

VERSION = $(shell grep 'Version = ' $(PROJECT).go | cut -d \" -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

build:
	go build -o bin/$(PROJECT) cmds/$(PROJECT)/$(PROJECT).go
	$(PROJECT) build.shorthand

lint:
	gofmt -w $(PROJECT).go && golint $(PROJECT).go
	gofmt -w $(PROJECT)_test.go && golint $(PROJECT)_test.go
	gofmt -w cmds/$(PROJECT)/$(PROJECT).go && golint cmds/$(PROJECT)/$(PROJECT).go

test:
	go test

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROJECT)-$(VERSION)-release.zip ]; then /bin/rm $(PROJECT)-$(VERSION)-release.zip; fi

install:
	GOBIN=$(HOME)/bin go install cmds/$(PROJECT)/$(PROJECT).go

uninstall:
	if [ -f $(GOBIN)/$(PROJECT) ]; then /bin/rm $(GOBIN)/$(PROJECT); fi

doc:
	$(PROJECT) build.shorthand

save:
	git commit -am "Quick Save"
	git push origin $(BRANCH)

release:
	./mk-release.sh

publish:
	./publish.sh
