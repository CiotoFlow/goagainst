GO=go

goagainst: $(SRCS)
	$(GO) build -o goagainst ./src

clean:
	rm -f goagainst

.PHONY: clean