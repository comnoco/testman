TMPDIR ?= /tmp/testman-integration

integration:
	# test with testman
	@mkdir -p $(TMPDIR)
	testman test -timeout=10s -run ^TestStable ./...  # should always work
	testman test -timeout=60s -run ^TestUnstable -retry=50 ./...  # should work at least 1/50
	(testman test -timeout=10s -continue-on-error -retry=10 ./...; echo "EXIT CODE: $?") > $(TMPDIR)/got.continue-on-error.log	  # should fail, but also run all tests
	testman test -v -test.timeout=10s -test.race -test.v -test.count=42 -run ^TestStable ./... 2>$(TMPDIR)/got.opts.log  # should get all the passed options in the verbose output
	# cat $(TMPDIR)/got.opts.log | grep -q -- 'moul.io~testman~examples~testpkg -test.count=42 -test.timeout=10s -test.v -test.run ^TestStableAlwaysSucceed'
	# cat $(TMPDIR)/got.opts.log | grep -q -- 'test -c -v -race -o'
	diff ./testdata/expected.continue-on-error.log $(TMPDIR)/got.continue-on-error.log

	# test with default tools
	go install moul.io/retry
	go get moul.io/retry
	go test -run ^TestStable -count=20 ./... >/dev/null                                                     # should always work
	retry -m=5  --interval=0 -- "(go test -run ^TestBroken -count=1 ./... >/dev/null)" && exit 1 || exit 0  # should always fail
	retry -m=50 --interval=0 -- "(go test -run ^TestUnstable -count 1 ./... >/dev/null)"                    # should work at least 1/50
	go mod tidy

	@echo "SUCCESS."
	@rm -rf $(TMPDIR)
