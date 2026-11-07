/**
 * @package elr-crawler (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Sept 07, 2026; 1:33 PM
 */

package sql

func NewFilter() *Filter {
	return &Filter{
		paramCount: new(uint),
	}
}

func NewInsert() *Insert {
	return &Insert{}
}

func NewUpdate() *Update {
	return &Update{}
}
