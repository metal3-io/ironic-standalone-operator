# Fuzzing

Fuzz tests for the ironic-standalone-operator using Go's built-in fuzzing support.

## Running Fuzz Tests

### Quick Start with Makefile

```bash
# Run fuzz tests as regression tests (seed corpus only, fast)
make fuzz

# Run fuzz tests with the fuzzer enabled (default: 30 seconds)
make fuzz-run

# Run fuzz tests for a custom duration (e.g. 5 minutes)
make fuzz-run FUZZ_TIME=5m
```

## Crash Corpus and Regression Testing

When fuzzing discovers a crash, Go automatically saves the failing input to
`testdata/fuzz/<FuzzTestName>/` in the test directory. These crash files should
be committed to the repository:

```bash
git add test/fuzz/testdata/
git commit -m "Add fuzz crash corpus"
```

Once committed, the crashes are automatically replayed as regression tests when
running `make fuzz` (or `go test` without `-fuzz`).

## Resources

- [Go Fuzzing Documentation](https://go.dev/doc/fuzz/)
- [Go Fuzzing Tutorial](https://go.dev/security/fuzz/)
- [Reference: metal3-io/baremetal-operator fuzz tests](https://github.com/metal3-io/baremetal-operator/tree/main/test/fuzz)
