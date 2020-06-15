SPLIC=../../bin/splic
SPLIS=../../bin/splis

LIB=lib
CORE=$(LIB)/core

all: build run

build:
	$(SPLIC) prelude.test.splis
	$(SPLIC) math.test.splis

run:
	$(SPLIS) prelude.test.nob
	$(SPLIS) math.test.nob

.PHONY: clean

clean:
	rm -f
