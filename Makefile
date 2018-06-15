GINKGOBUILD=ginkgo build
INSTALL=go install


tests:
	$(GINKGOBUILD) ./pkg/benchmarkscale
	$(GINKGOBUILD) ./pkg/benchmarktiming
	$(GINKGOBUILD) ./pkg/benchmarkapi
install:
	$(GINKGOBUILD) ./pkg/benchmarkscale
	$(GINKGOBUILD) ./pkg/benchmarktiming
	$(GINKGOBUILD) ./pkg/benchmarkapi
	$(INSTALL) ./pkg/benchmarkrunner
