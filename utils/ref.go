package utils

import "strings"

func Ref[T any](v T) *T {
	return &v
}

func ValueFromRef[T any](v *T) T {
	if v == nil {
		d := new(T)
		return *d
	} else {
		return *v
	}
}

func IsDigits(s string) bool {
	for _, b := range s {
		if b < '0' || b > '9' {
			return false
		}
	}
	return len(s) > 0
}

func IsBlankStrPtr(s *string) bool {
	if s == nil {
		return true
	}
	return IsBlankStr(*s)
}

func GetNonEmptyStr(s ...string) string {
	for _, p := range s {
		if !IsBlankStr(p) {
			return p
		}
	}
	return ""
}

func IsBlankStr(s string) bool {
	return strings.TrimSpace(s) == ""
}
