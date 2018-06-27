GINKGOBUILD=ginkgo build
INSTALL=go install
GET=go get -u
PKGLOC=$(GOPATH)/src/dispatchframework/benchmark

install:
	$(INSTALL) ./pkg/tester
install-full:
	dep ensure
	$(INSTALL) ./pkg/tester
