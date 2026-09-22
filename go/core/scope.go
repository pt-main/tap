package core

import "errors"

type ScopeType map[string]interface{}

func ScopeGet[T any](st ScopeType, what string) (T, error) {
	var nul T
	val, ok := st[what]
	if !ok {
		return nul, errors.New("Scope.Get: Invalid key: " + what)
	}
	res, ok := val.(T)
	if !ok {
		return nul, errors.New("Scope.Get: Invalid type for key: " + what)
	}
	return res, nil
}
