# benchmark
Tools to test scalability &amp; performance of Dispatch

# Running the tests
Makefile *should* allow for easy setup
```
make install
```
should install the test binary
```
tester
```
is the command to run tests. By default all tests will be run. However, the tester accepts regular expression arguments that allow you to select specific tests.
Command line flags are also available.
```
-sample
```
allows you to configure the number of samples of each test. Probably the only really important flag. Set to 1 by default
## Configuring the tests
No more configuration. Tests are controlled by the command line.

Please let me know if there are any essential config options you think we need.
