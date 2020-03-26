ROOT=../../..

SPLIC=$(ROOT)/bin/splic
SPLIL=$(ROOT)/bin/splil
SPLIS=$(ROOT)/bin/splis

LIB=$(ROOT)/lib
CORE=$(LIB)/core

run: test-1.splis
	$(SPLIC) prelude.splis
	$(SPLIC) test-1.splis
	$(SPLIL) -o test-1.splis.bin prelude.splis.lib test-1.splis.lib
	$(SPLIS) test-1.splis.bin

.PHONY: clean

clean:
	rm -f *.splis.lib *.splis.bin
