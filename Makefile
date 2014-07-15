GO=go

goagainst:
	$(GO) build -o goagainst ./src

clean:
	rm -f goagainst

check:
	go test ./src/trollan

.PHONY: clean goagainst check