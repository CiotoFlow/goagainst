GO=go

goagainst:
	$(GO) build -o goagainst ./src

clean:
	rm -f goagainst

.PHONY: clean goagainst