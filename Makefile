#
# Shorthand a text label expander.
#
PROJECT = shorthand

VERSION = $(shell grep 'Version = ' $(PROJECT).go | cut -d \" -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

build:
	go build -o bin/$(PROJECT) cmds/$(PROJECT)/$(PROJECT).go

website:
	$(PROJECT) build.shorthand
	git add index.html install.html license.html
	git add shorthand.md shorthand.html

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

status:
	git status

save:
	git commit -am "Quick Save"
	git push origin $(BRANCH)

publish:
	./publish.sh


dist/linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/shorthand cmds/shorthand/shorthand.go

dist/windows-amd64:
	env GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/shorthand.exe cmds/shorthand/shorthand.go

dist/macosx-amd64:
	env GOOS=darwin	GOARCH=amd64 go build -o dist/macosx-amd64/shorthand cmds/shorthand/shorthand.go

dist/raspbian-arm7:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/shorthand cmds/shorthand/shorthand.go

release: dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7
	mkdir -p dist
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v shorthand.md dist/
	zip -r $(PROJECT)-$(VERSION)-release.zip .md dist/*

