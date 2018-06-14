GINKGOBUILD=ginkgo build
INSTALL=go install


tests:
	$(GINKGOBUILD) ./pkg/benchmarkscale
	$(GINKGOBUILD) ./pkg/benchmarktiming
install:
	$(GINKGOBUILD) ./pkg/benchmarkscale
	$(GINKGOBUILD) ./pkg/benchmarktiming
	$(INSTALL) ./pkg/benchmarkrunner
