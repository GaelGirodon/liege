package model

// QueryParam is a query parameter required for a route to match.
type QueryParam struct {
	// Name is the query parameter name.
	Name string `json:"name"`
	// Value is the query parameter value
	// (an empty string means that any value will match).
	Value string `json:"value"`
}
