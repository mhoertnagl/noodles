NOODLEC=bin/noodlec
NOODLES=bin/noodles

LIB=lib
CORE=$(LIB)/core

build:
	go build -o $(NOODLEC) cmd/noodlec/main.go
	go build -o $(NOODLES) cmd/noodles/main.go

.PHONY: clean

clean:
	rm -f $(NOODLEC) $(NOODLES)
