/*
 * @Descripttion:
 * @version:
 * @Author: sueRimn
 * @Date: 2021-11-04 23:31:13
 * @LastEditors: sueRimn
 * @LastEditTime: 2021-11-05 13:43:03
 */
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Default(""),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
