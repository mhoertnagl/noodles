ROOT=../../..

SPLIC=$(ROOT)/bin/splic
SPLIL=$(ROOT)/bin/splil
SPLIS=$(ROOT)/bin/splis

LIB=$(ROOT)/lib
CORE=$(LIB)/core

run: test-0.splis
	$(SPLIC) test-0.splis
	$(SPLIL) -o test-0.splis.bin test-0.splis.lib
	$(SPLIS) test-0.splis.bin

.PHONY: clean

clean:
	rm -f *.splis.lib *.splis.bin
