package tests

import (
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"testing"
)

func TestMatchTwoRelations(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)
	matcher := mentalese.NewRelationMatcher(log)
	tests := []struct {
		needle      string
		haystack    string
		binding     string
		wantBinding string
		wantMatch   bool
	}{
		{"parent(X, Y)", "parent('Luke', 'George')", "{}", "{X: 'Luke', Y: 'George'}", true},
		{"parent('Luke', 'George')", "parent(X, Y)", "{}", "{}", true},
		{"parent('Luke', Y)", "parent('Luke', 'George')", "{}", "{Y: 'George'}", true},
		{"parent('Luke', 'Richard')", "parent('Luke', 'George')", "{}", "{}", false},
		{"parent(X, Y)", "parent(A, B)", "{}", "{X:A, Y:B}", true},
		{"parent(X, Y)", "parent('Luke', 'George')", "{X: 'Luke'}", "{X: 'Luke', Y: 'George'}", true},
		{"parent(X, Y)", "parent('Luke', 'George')", "{X: 'Vincent'}", "{X: 'Vincent'}", false},
		{"quantification(X, [], Y, [ isa(Y, every) ])", "quantification(A, [], B, [ isa(B, every) ])", "{}", "{X: A, Y: B}", true},
		{"quantification(X, _, Y, [ isa(Y, every) ])", "quantification(A, [], B, [ isa(B, P) specification(B, S) isa(S, very)])", "{}", "{X: A, Y: B}", true},
		{"quantification(X, [], Y, Y1)", "quantification(A, E, B, [ isa(Y, every) ])", "{}", "{X: A, Y: B, Y1: [isa(Y, every)]}", true},
		{"quantification(X, _, Y, [ isa(Y, Q) ])", "quantification(A, [], B, [ isa(B, every) specification(B, S) isa(S, very)])", "{}", "{X: A, Q: every, Y: B}", true},
	}

	for _, test := range tests {

		needle := parser.CreateRelation(test.needle)
		haystack := parser.CreateRelation(test.haystack)
		binding := parser.CreateBinding(test.binding)
		wantBinding := parser.CreateBinding(test.wantBinding)
		wantMatch := test.wantMatch

		resultBinding, resultMatch := matcher.MatchTwoRelations(needle, haystack, binding)

		if !resultBinding.Equals(wantBinding) || resultMatch != wantMatch {
			t.Errorf("MatchTwoRelations(%v %v %v): got %v %v, want %v %v", needle, haystack, binding, resultBinding, resultMatch, wantBinding, wantMatch)
		}
	}
}

func TestMatchRelationToSet(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)
	matcher := mentalese.NewRelationMatcher(log)
	haystack := parser.CreateRelationSet("[gender('Luke', male) gender('George', male) parent('Luke', 'George') parent('Carry', 'Steven') gender('Carry', female)]")

	var tests = []struct {
		needle       string
		haystack     mentalese.RelationSet
		binding      string
		wantBindings string
		wantIndexes  []int
	}{
		{"parent(X, Y)", haystack, "{}", "[{X:'Luke', Y:'George'} {X:'Carry', Y:'Steven'}]", []int{2, 3}},
		{"parent(X, 'Henry')", haystack, "{}", "[]", []int{}},
		{"parent(X, 'Steven')", haystack, "{}", "[{X:'Carry'}]", []int{3}},
		{"parent(X, 'Steven')", haystack, "{X:'Carry'}", "[{X:'Carry'}]", []int{3}},
		{"parent('Carry', 'Marvin')", haystack, "{}", "[]", []int{}},
		{"parent(X, 'George')", haystack, "{X: A, Y: B}", "[]", []int{}},
	}

	for _, test := range tests {

		needle := parser.CreateRelation(test.needle)
		binding := parser.CreateBinding(test.binding)
		wantBindings := parser.CreateBindings(test.wantBindings)
		wantIndexes := test.wantIndexes

		resultBindings, resultIndexes := matcher.MatchRelationToSet(needle, haystack, binding)

		bindingsOk := len(wantBindings) == len(resultBindings)
		for i, resultBinding := range resultBindings {
			bindingsOk = bindingsOk && resultBinding.Equals(wantBindings[i])
		}

		indexesOk := len(wantIndexes) == len(resultIndexes)
		for i, resultIndex := range resultIndexes {
			indexesOk = indexesOk && resultIndex == wantIndexes[i]
		}

		if !bindingsOk || !indexesOk {
			t.Errorf("MatchRelationToSet(%v %v %v): got %v %v, want %v %v", needle, haystack, binding, resultBindings, resultIndexes, wantBindings, wantIndexes)
		}
	}
}

func TestMatchSequenceToSet(t *testing.T) {

	parser := importer.NewInternalGrammarParser()
	log := common.NewSystemLog(false)
	matcher := mentalese.NewRelationMatcher(log)
	haystack := parser.CreateRelationSet(`[
		gender('Luke', male)
		gender('George', male)
		gender('Jeanne', female)
		parent('Luke', 'George')
		parent('Carry', 'Steven')
		parent('Carry', 'Jeanne')
		gender('Carry', female)]
	`)

	var tests = []struct {
		needle       string
		haystack     mentalese.RelationSet
		binding      string
		wantBindings string
		wantIndexes  []int
		wantMatch    bool
	}{
		{"[parent(X, Y) gender(X, male)]", haystack, "{}", "[{Y:'George', X:'Luke'}]", []int{3, 0}, true},
		{"[parent(X, Y) gender(Y, female)]", haystack, "{X: 'Carry'}", "[{X: 'Carry', Y:'Jeanne'}]", []int{5, 2}, true},
		{"[parent(X, Y) gender(Y, female)]", haystack, "{X: 'Quincy'}", "[]", []int{}, false},
		{"[parent(X, Y) gender(X, female)]", haystack, "{}", "[{X:'Carry', Y:'Steven'} {X:'Carry', Y:'Jeanne'}]", []int{4, 6, 5}, true},
		{"[parent('Carry', Y) gender(Y, M)]", haystack, "{Q: 3}", "[{Q: 3, Y:'Jeanne', M: female}]", []int{5, 2}, true},
		{"[gender(Y, M) parent(X, Y) gender(X, M)]", haystack, "{}", "[{X:'Luke', Y:'George', M:male} {X:'Carry', Y:'Jeanne', M:female}]", []int{1, 3, 0, 2, 5, 6}, true},
	}

	for _, test := range tests {

		needle := parser.CreateRelationSet(test.needle)
		binding := parser.CreateBinding(test.binding)
		wantBindings := parser.CreateBindings(test.wantBindings)
		wantIndexes := test.wantIndexes
		wantMatch := test.wantMatch
		resultBindings, resultIndexes, _, resultMatch := matcher.MatchSequenceToSetWithIndexes(needle, haystack, binding)

		bindingsOk := len(wantBindings) == len(resultBindings)
		for i, resultBinding := range resultBindings {
			bindingsOk = bindingsOk && resultBinding.Equals(wantBindings[i])
		}

		indexesOk := len(wantIndexes) == len(resultIndexes)
		for i, resultIndex := range resultIndexes {
			indexesOk = indexesOk && resultIndex == wantIndexes[i]
		}

		if !bindingsOk || !indexesOk || wantMatch != resultMatch {
			t.Errorf("MatchSequenceToSet(%v %v %v): got %v %v %v, want %v %v %v", needle, haystack, binding, resultBindings, resultIndexes, resultMatch, wantBindings, wantIndexes, wantMatch)
		}
	}
}
