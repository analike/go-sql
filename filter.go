/**
 * @package elr-crawler (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Sept 07, 2026; 1:33 PM
 */

package sql

import (
	"fmt"
	"strings"
)

type Filter struct {
	paramCount *uint
	params     []any
	filters    []string
}

// AddParam add to the param stack and returns the placeholder text for
// the given value.
//
// Useful in advance query cases where reusing a particular param is necessary.
func (f *Filter) AddParam(value any) string {
	f.params = append(f.params, value)
	return f.incrementCounter().getPlaceholder()
}

func (f *Filter) WithFilter(field string, value any) *Filter {
	f.filters = append(
		f.filters,
		fmt.Sprintf(
			"%s = %s",
			field,
			f.incrementCounter().getPlaceholder(),
		),
	)
	f.params = append(f.params, value)
	return f
}

func (f *Filter) WithFilterReversed(field string, value any) *Filter {
	f.filters = append(
		f.filters,
		fmt.Sprintf(
			"%s = %s",
			f.incrementCounter().getPlaceholder(),
			field,
		),
	)
	f.params = append(f.params, value)
	return f
}

func (f *Filter) WithFilterOperator(field string, value any, operator string) *Filter {
	f.filters = append(
		f.filters,
		fmt.Sprintf(
			"%s %s %s",
			field,
			operator,
			f.incrementCounter().getPlaceholder(),
		),
	)
	f.params = append(f.params, value)
	return f
}

func (f *Filter) WithFilterOperatorReversed(field string, value any, operator string) *Filter {
	f.filters = append(
		f.filters,
		fmt.Sprintf(
			"%s %s %s",
			f.incrementCounter().getPlaceholder(),
			operator,
			field,
		),
	)
	f.params = append(f.params, value)
	return f
}

func (f *Filter) WithRawFilter(filter string) *Filter {
	f.filters = append(f.filters, filter)
	return f
}

func (f *Filter) WithValueFormat(field string, value any, format string) *Filter {
	f.filters = append(
		f.filters,
		fmt.Sprintf(
			"%s %s %s",
			field,
			"=",
			fmt.Sprintf(format, f.incrementCounter().getPlaceholder()),
		),
	)
	f.params = append(f.params, value)
	return f
}

func (f *Filter) WithLikeNoCase(field string, value any) *Filter {
	f.WithFilterOperatorReversed(field, value, "ILIKE")
	return f
}

func (f *Filter) WithLike(field string, value any) *Filter {
	f.WithFilterOperatorReversed(field, value, "LIKE")
	return f
}

func (f *Filter) WithFilterOperatorEmbed(field string, value string, operator string) *Filter {
	f.filters = append(
		f.filters,
		fmt.Sprintf(
			"%s %s %s",
			field,
			operator,
			value,
		),
	)
	return f
}

// WithGroup adds a bracketed sub-clause, e.g. (field5 = $5 OR field6 = $6).
//
// Param: `build` - func that receives a *Filter, generates the where clause
// while sharing the parent's paramCount pointer so placeholders stay correctly
// numbered across the entire tree. Its resulting params are merged back
// into the parent in call order.
//
// Param: `joiner` - typically "AND" or "OR".
func (f *Filter) WithGroup(joiner string, build func(*Filter)) *Filter {
	f.ensureCounter()

	group := &Filter{paramCount: f.paramCount}
	build(group)

	if len(group.filters) == 0 {
		return f
	}

	clause := strings.Join(group.filters, " "+joiner+" ")
	f.filters = append(f.filters, "("+clause+")")
	f.params = append(f.params, group.params...)
	return f
}

// WithOrGroup is a convenient wrapper for WithGroup using OR as joiner
func (f *Filter) WithOrGroup(build func(*Filter)) *Filter {
	return f.WithGroup("OR", build)
}

// WithAndGroup is a convenient wrapper for WithGroup using AND as joiner.
func (f *Filter) WithAndGroup(build func(*Filter)) *Filter {
	return f.WithGroup("AND", build)
}

// GetWhere returns the where clause of the SQL query with correct placeholders
// in the same order that WithFilter* was called
//
// e.g: field1 = $1 AND field2 = $2
func (f *Filter) GetWhere() string {
	if len(f.filters) < 1 {
		f.WithFilter("1", "1")
	}
	return strings.Join(f.filters, " AND ")
}

func (f *Filter) GetParams() []any {
	return f.params
}

func (f *Filter) incrementCounter() *Filter {
	f.ensureCounter()
	*f.paramCount++
	return f
}

func (f *Filter) getPlaceholder() string {
	return fmt.Sprintf("$%d", *f.paramCount)
}

// ensureCounter initializes counter is nil
func (f *Filter) ensureCounter() {
	if f.paramCount == nil {
		var count uint
		f.paramCount = &count
	}
}
