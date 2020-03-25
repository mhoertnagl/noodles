SPLIC=bin/splic
SPLIL=bin/splil
SPLIS=bin/splis

LIB=lib
CORE=$(LIB)/core

TEST=examples/tests

run:
	./splis examples/sicp/ch01.splis

build:
	go build -o $(SPLIC) cmd/splic/main.go
	go build -o $(SPLIL) cmd/splil/main.go
	go build -o $(SPLIS) cmd/splis/main.go

lib: $(CORE)/prelude.splis
	$(SPLIC) $(CORE)/prelude.splis

test-0: $(TEST)/test-0/test-0.splis
	$(SPLIC) $(TEST)/test-0/test-0.splis
	$(SPLIL) -o $(TEST)/test-0/test-0.splis.bin \
	         $(TEST)/test-0/test-0.splis.lib
	$(SPLIS) $(TEST)/test-0/test-0.splis.bin

test-1: $(TEST)/test-1/test-1.splis
	$(SPLIC) $(TEST)/test-1/prelude.splis
	$(SPLIC) $(TEST)/test-1/test-1.splis
	$(SPLIL) -o $(TEST)/test-1/test-1.splis.bin \
					 $(TEST)/test-1/prelude.splis.lib \
					 $(TEST)/test-1/test-1.splis.lib
	$(SPLIS) $(TEST)/test-1/test-1.splis.bin

test-2: $(TEST)/test-2/test-2.splis
	$(SPLIC) $(TEST)/test-2/test-2.splis
	$(SPLIL) -o $(TEST)/test-2/test-2.splis.bin \
		       $(TEST)/test-2/test-2.splis.lib
	$(SPLIS) $(TEST)/test-2/test-2.splis.bin
