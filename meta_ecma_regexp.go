package navigator

import (
	"time"

	"github.com/dlclark/regexp2"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

// metaPatternMatchTimeout caps regexp2 evaluation so pathological patterns
// cannot stall validation on untrusted documents.
const metaPatternMatchTimeout = 250 * time.Millisecond

// ecmascriptRegexp adapts github.com/dlclark/regexp2 (ECMAScript mode) to
// jsonschema/v6's Regexp interface. JSON Schema Draft 2020-12 uses ECMA-262
// regular expressions; Go's standard regexp package (RE2) cannot compile
// patterns emitted by Zod for .email() / .url(), etc.
type ecmascriptRegexp regexp2.Regexp

func compileMetaJSONSchemaPattern(pattern string) (jsonschema.Regexp, error) {
	re, err := regexp2.Compile(pattern, regexp2.ECMAScript)
	if err != nil {
		return nil, err
	}
	re.MatchTimeout = metaPatternMatchTimeout
	return (*ecmascriptRegexp)(re), nil
}

func (re *ecmascriptRegexp) MatchString(s string) bool {
	r := (*regexp2.Regexp)(re)
	matched, err := r.MatchString(s)
	return err == nil && matched
}

func (re *ecmascriptRegexp) String() string {
	return (*regexp2.Regexp)(re).String()
}
