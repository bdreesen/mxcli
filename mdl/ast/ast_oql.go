// SPDX-License-Identifier: Apache-2.0

package ast

// ============================================================================
// OQL AST Types
// ============================================================================

// OQLSelectItem represents a single SELECT column in an OQL query.
type OQLSelectItem struct {
	Expression   string     // raw expression text (for display/round-trip)
	Alias        string     // AS alias
	IsAggregate  bool       // whether this is an aggregate function (COUNT, SUM, etc.)
	Subquery     *OQLParsed // nested scalar subquery (for inline SELECT subqueries)
	SubqueryText string     // original text of the scalar subquery (for display)
}

// OQLTableRef represents a FROM or JOIN table reference in an OQL query.
type OQLTableRef struct {
	Entity       string     // qualified entity name (e.g. "CRM.Account")
	Alias        string     // table alias
	JoinType     string     // "from", "join", "left join", "right join", "full join", "cross join"
	AssocPath    string     // full association path if present (e.g. "act/CRM.Activity_Account/CRM.Account")
	OnExpr       string     // ON condition text (for JOINs)
	Subquery     *OQLParsed // nested subquery (for derived tables)
	SubqueryText string     // original text of the subquery (for display)
}

// OQLParsed represents a fully parsed OQL query with structured access to clauses.
type OQLParsed struct {
	RawQuery string          // preserved for round-tripping
	Select   []OQLSelectItem // SELECT columns
	Tables   []OQLTableRef   // FROM + JOINs
	Where    string          // WHERE clause text
	GroupBy  string          // GROUP BY clause text
	Having   string          // HAVING clause text
	OrderBy  string          // ORDER BY clause text
}
