package parse

import "testing"

func TestEnum(t *testing.T) {
	tt := []Pair{
		{"enum Op { EOF, JMP, POP };", "(decl (enum Op EOF JMP POP))"},
		{"enum Op { EOF };", "(decl (enum Op EOF))"},
	}

	check(t, tt)
}
func TestTypeSpecifier(t *testing.T) {
	tt := []Pair{
		// will fail if parser isn't initialized with "bool" in p.types
		{"bool;", "(decl (type_specifier bool))"},
	}

	check(t, tt)
}
func TestDefaultTypeSpecifier(t *testing.T) {
	tt := []Pair{
		{"void;", "(decl (default_type_specifier void))"},
		{"char;", "(decl (default_type_specifier char))"},
		{"short;", "(decl (default_type_specifier short))"},
		{"int;", "(decl (default_type_specifier int))"},
		{"long;", "(decl (default_type_specifier long))"},
		{"float;", "(decl (default_type_specifier float))"},
		{"double;", "(decl (default_type_specifier double))"},
		{"signed;", "(decl (default_type_specifier signed))"},
		{"unsigned;", "(decl (default_type_specifier unsigned))"},
		{"unsigned int;", "(decl (default_type_specifier unsigned) (default_type_specifier int))"},
	}

	check(t, tt)
}
func TestTypeQualifer(t *testing.T) {
	tt := []Pair{
		{"const;", "(decl (type_qualifier const))"},
		{"volatile;", "(decl (type_qualifier volatile))"},
		{"const volatile;", "(decl (type_qualifier const) (type_qualifier volatile))"},
	}

	check(t, tt)
}

func TestStorageClass(t *testing.T) {
	tt := []Pair{
		{"typedef;", "(decl (storage_class typedef))"},
		{"extern;", "(decl (storage_class extern))"},
		{"static;", "(decl (storage_class static))"},
		{"auto;", "(decl (storage_class auto))"},
		{"register;", "(decl (storage_class register))"},
		{"extern register;", "(decl (storage_class extern) (storage_class register))"},
	}

	check(t, tt)
}
