package grep_test

type TestCase struct {
	fileName     string
	pattern      string
	grepCount    int
	grepLastLine int32
	onlyFiles    bool
}
