SRC_FILES= $(shell find . -path './.*' -prune -o \( -name '*.go' -a ! -name '*_test.go' \) -print)

all: piot

piot: $(SRC_FILES)
	GOOS=linux GOARCH=amd64 go build  -o $@

piot.exe: $(SRC_FILES)
	GOOS=windows GOARCH=amd64 go build -o $@

.PHONY: clean
clean:
	rm -rfv piot piot.exe