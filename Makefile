GINKGOBUILD=ginkgo build
INSTALL=go install
GET=go get -u

tests:
	$(GINKGOBUILD) ./pkg/benchmarkscale
	$(GINKGOBUILD) ./pkg/benchmarktiming
	$(GINKGOBUILD) ./pkg/benchmarkapi
install:
	$(GINKGOBUILD) ./pkg/benchmarkscale
	$(GINKGOBUILD) ./pkg/benchmarktiming
	$(GINKGOBUILD) ./pkg/benchmarkapi
	$(INSTALL) ./pkg/benchmarkrunner

install-full:
	$(GET) github.com/onsi/ginkgo/ginkgo
	$(GET) github.com/onsi/gomega/...
	$(GINKGOBUILD) ./pkg/benchmarkscale
	$(GINKGOBUILD) ./pkg/benchmarktiming
	$(GINKGOBUILD) ./pkg/benchmarkapi
	$(INSTALL) ./pkg/benchmarkrunner
