package knowledge

import (
	"nli-go/lib/mentalese"
	"nli-go/lib/common"
)

type InMemoryFactBase struct {
	facts mentalese.RelationSet
	ds2db []mentalese.DbMapping
	matcher *mentalese.RelationMatcher
}

func NewInMemoryFactBase(facts mentalese.RelationSet, ds2db []mentalese.DbMapping) mentalese.FactBase {
	return InMemoryFactBase{facts: facts, ds2db: ds2db, matcher: mentalese.NewRelationMatcher()}
}

// Note! An internal fact base would use the same predicates as the domain language;
// This is an simulation of an external database
func (factBase InMemoryFactBase) Bind(goal mentalese.Relation) []mentalese.Binding {

	common.LogTree("Factbase Bind", goal);

	subgoalBindings := []mentalese.Binding{}

	for _, ds2db := range factBase.ds2db {

		// gender(14, G), gender(A, male) => externalBinding: G = male
		externalBinding, match := factBase.matcher.MatchTwoRelations(goal, ds2db.DsSource, mentalese.Binding{})
		if match {

			// gender(14, G), gender(A, male) => internalBinding: A = 14
			internalBinding, _ := factBase.matcher.MatchTwoRelations(ds2db.DsSource, goal, mentalese.Binding{})

			// create a version of the conditions with bound variables
			boundConditions := factBase.matcher.BindRelationSetSingleBinding(ds2db.DbTarget, internalBinding)
			// match this bound version to the database
			internalBindings, _, match := factBase.matcher.MatchSequenceToSet(boundConditions, factBase.facts, mentalese.Binding{})

			if match {
				for _, binding := range internalBindings {
					subgoalBindings = append(subgoalBindings, externalBinding.Intersection(binding))
				}
			}
		}
	}

	common.LogTree("Factbase Bind", subgoalBindings);

	return subgoalBindings
}
