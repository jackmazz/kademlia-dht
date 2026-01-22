GO := go
PROTOC := protoc

# Build the following commands.  This assumes that each command
# CMDNAME is in the directory cmd/CMDNAME, and can be built by
# changing to that directory and running go build.
#
# No commands are provided for you in this project source, but you may
# implement any commands that you like for testing purposes.
COMMANDS :=

# This rule turns COMMANDS into executable filenames, do not change.
# You don't need to understand this.
CMDFILES := $(shell for word in $(COMMANDS); do echo cmd/$$word/$$word; done)

# Running the command `make` with no arguments should build the
# commands specified in $(COMMANDS).  This is for your convenience,
# you can also build them with `go build`.
#
# The body of this rule is a shell script that loops through every
# command defined in COMMANDS and builds it with go build.  Shell
# scripts embedded in Makefiles have somewhat strange parsing rules
# due to the way that Make works; see `info make` for more
# information.
all: go.sum
	@for cmd in $(COMMANDS); do                             \
	    echo "Building $$cmd";                              \
	    (cd cmd/$$cmd; go build);                           \
        done

submission:
	tar cf kdht.tar \
	    $(shell for word in $(CMDFILES); do echo "--exclude $$word"; done) \
	    --exclude '.*' --exclude '.DS_Store' \
	    Makefile impl tests cmd

giventest: all
	go test cse586.kdht/given/keys

test: all
	go test -test.v cse586.kdht/tests
	go test -test.v cse586.kdht/impl

go.sum: api/kdht/messages.pb.go
	go get cse586.kdht/api/kdht

clean:
	rm -f $(CMDFILES) kdht.tar
	find . -name '*~' -delete

# Build a protobuf implementation from a protocol description
build-protobuf:
	$(PROTOC) --go_out=. --go_opt=paths=source_relative api/kdht/messages.proto

.PHONY: all clean submission giventest test build-protobuf
