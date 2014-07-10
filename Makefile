GO=go
SRCS=main.go

goagainst: $(SRCS)
	go build -o goagainst $(SRCS)
