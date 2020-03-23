run:
	./splis examples/sicp/ch01.splis

build:
	go build -o bin/splic cmd/splic/main.go
	go build -o bin/splil cmd/splil/main.go
	go build -o bin/splis cmd/splis/main.go

lib: lib/core/prelude.splis
	bin/splic lib/core/prelude.splis

test-0: examples/tests/test-0.splis
	bin/splic examples/tests/test-0.splis
	bin/splil -o examples/tests/test-0.splis.bin examples/tests/test-0.splis.lib lib/core/prelude.splis.lib
