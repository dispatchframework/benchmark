# benchmark
Tools to test scalability &amp; performance of Dispatch

# Running the tests
Makefile *should* allow for easy setup
```
make install
```
should install the script that runs all the tests based on a configuration file, after running this you should be able to run
```
benchmarkrunner ./config.yaml
```
adjusted based on where your configuration file is
## Configuring the tests
Currently configuration options are limited. There are two tests, timing and scalability. The configuration file allows you to specify which ones to run, where they output to (this doesn't work yet lol), and the location of the binaries of the tests. To run everything assuming standard go configuration, the following should work.
```yaml
Scale:
  enabled: true
  Location: "./pkg/benchmarkscale"
Timing:
  enabled: true
  Location: "./pkg/benchmarktiming"

```
