module github.com/readwritepro/figtree

go 1.16

replace github.com/readwritepro/error-handler => ../error-handler

replace github.com/readwritepro/compare-test-results => ../compare-test-results

require (
	github.com/readwritepro/compare-test-results v0.0.0-00010101000000-000000000000
	github.com/readwritepro/error-handler v0.0.0-00010101000000-000000000000
)
