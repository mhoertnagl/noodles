ROOT=../../..

SPLIC=$(ROOT)/bin/splic
SPLIL=$(ROOT)/bin/splil
SPLIS=$(ROOT)/bin/splis

LIB=$(ROOT)/lib
CORE=$(LIB)/core

run: test-2.splis
	$(SPLIC) test-2.splis
	$(SPLIL) -o test-2.splis.bin test-2.splis.lib
	$(SPLIS) test-2.splis.bin

.PHONY: clean

clean:
	rm -f *.splis.lib *.splis.bin
