package utils

import "regexp"

func compileOmittedHeaders(omitHeaders []string) []*regexp.Regexp {
	compiled := make([]*regexp.Regexp, len(omitHeaders))
	for i, value := range omitHeaders {
		compiled[i] = regexp.MustCompile("(?i)" + value)
	}
	return compiled
}

var (
	inputUrlRegex         = regexp.MustCompile(`(?i)^http://1\.1\.\d+\.\d+/bmi/(https?://)?`)
	omittedHeadersRegexes = compileOmittedHeaders(BHP_EXTERNAL_REQUEST_OMIT_HEADERS)
)
