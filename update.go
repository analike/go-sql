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

type UpdateField map[string]string

type Update struct {
	places []UpdateField
	values []any
	count  int
}

func (s *Update) WithValue(key string, val any) *Update {
	s.count++
	s.places = append(s.places, UpdateField{key: fmt.Sprintf("$%d", s.count)})
	s.values = append(s.values, val)
	return s
}

func (s *Update) GetUpdates() string {
	var res []string
	for _, place := range s.places {
		for k, v := range place {
			res = append(res, fmt.Sprintf("%s = %s", k, v))
		}
	}
	return strings.Join(res, ", ")
}

func (s *Update) GetValues() []any {
	return s.values
}

// WithNoSetValue adds an entry to the count stack without touching the key:val pair.
//
// Useful for cases when adding a parameter that is excluded from SET ...
//
// Returns the positional argument for the value
func (s *Update) WithNoSetValue(val any) string {
	s.count++
	s.values = append(s.values, val)
	return fmt.Sprintf("$%d", s.count)
}

func (s *Update) GetNextCount() int {
	s.count++
	return s.count
}
