SPLIC=bin/splic
SPLIL=bin/splil
SPLIS=bin/splis

LIB=lib
CORE=$(LIB)/core

build:
	go build -o $(SPLIC) cmd/splic/main.go
	go build -o $(SPLIL) cmd/splil/main.go
	go build -o $(SPLIS) cmd/splis/main.go

lib: $(CORE)/prelude.splis
	$(SPLIC) $(CORE)/prelude.splis

.PHONY: clean

clean:
	rm -f bin/splic bin/splil bin/splis
