package main

type resolver struct {
	scopes     []map[string]bool
	distances  map[Expr]int
	inFunction bool
}

func createResolver() *resolver {
	return &resolver{
		scopes:     []map[string]bool{},
		distances:  map[Expr]int{},
		inFunction: false,
	}
}
