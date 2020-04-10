SPLIC=bin/splic
SPLIS=bin/splis

LIB=lib
CORE=$(LIB)/core

build:
	go build -o $(SPLIC) cmd/splic/main.go
	go build -o $(SPLIS) cmd/splis/main.go

.PHONY: clean

clean:
	rm -f $(SPLIC) $(SPLIS)
