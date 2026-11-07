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

type Insert struct {
	values []any
	places []int
}

func (s *Insert) WithValue(val any) *Insert {
	s.values = append(s.values, val)
	return s
}

func (s *Insert) WithValues(vals ...any) *Insert {
	if len(vals) > 0 {
		s.values = append(s.values, vals...)
	}
	return s
}

func (s *Insert) GetValues() []any {
	return s.values
}

func (s *Insert) GetPlaceholders() string {
	var (
		count = 1
		bag   []string
	)
	for range s.values {
		bag = append(bag, fmt.Sprintf("$%d", count))
		count++
	}
	return strings.Join(bag, ", ")
}
