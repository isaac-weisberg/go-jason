//go:build !prod

package gojason

import "strings"

type parseJsonStateChain struct {
	states []string
}

func newParseJsonStateChain() parseJsonStateChain {
	return parseJsonStateChain{
		states: []string{},
	}
}

func (parseJsonStateChain parseJsonStateChain) push(parseJsonObjectState parseJsonObjectState) {
	parseJsonStateChain.states = append(parseJsonStateChain.states, parseJsonObjectState.String())
}

func (parseJsonStateChain parseJsonStateChain) string() string {
	return strings.Join(parseJsonStateChain.states, " ")
}
