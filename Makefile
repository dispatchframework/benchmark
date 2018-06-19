GINKGOBUILD=ginkgo build
INSTALL=go install
GET=go get -u
PKGLOC=$(GOPATH)/src/dispatchframework/benchmark


tests:
	$(GINKGOBUILD) $(PKGLOC)/pkg/benchmarkscale
	$(GINKGOBUILD) $(PKGLOC)/pkg/benchmarktiming
	$(GINKGOBUILD) $(PKGLOC)/pkg/benchmarkapi
install:
	$(GINKGOBUILD) $(PKGLOC)/pkg/benchmarkscale
	$(GINKGOBUILD) $(PKGLOC)/pkg/benchmarktiming
	$(GINKGOBUILD) $(PKGLOC)/pkg/benchmarkapi
	$(INSTALL) $(PKGLOC)/pkg/benchmarkrunner

install-ginkgo:
	$(GET) github.com/onsi/ginkgo/ginkgo
	$(GET) github.com/onsi/gomega/...

install-full:
	$(GET) github.com/onsi/ginkgo/ginkgo
	$(GET) github.com/onsi/gomega/...
	$(GINKGOBUILD) $(PKGLOC)/pkg/benchmarkscale
	$(GINKGOBUILD) $(PKGLOC)/pkg/benchmarktiming
	$(GINKGOBUILD) $(PKGLOC)/pkg/benchmarkapi
	$(INSTALL) $(PKGLOC)/pkg/benchmarkrunner
