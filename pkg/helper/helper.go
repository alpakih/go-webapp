package helper

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"strconv"
	"strings"
)

// IndexOf get the index of the given value in the given slice,
// or -1 if not found.
func IndexOf(slice interface{}, value interface{}) int {
	list := reflect.ValueOf(slice)
	length := list.Len()
	for i := 0; i < length; i++ {
		if list.Index(i).Interface() == value {
			return i
		}
	}
	return -1
}

// Contains check if a slice contains a value.
func Contains(slice interface{}, value interface{}) bool {
	return IndexOf(slice, value) != -1
}

// IndexOfStr get the index of the given value in the given string slice,
// or -1 if not found.
// Prefer using this helper instead of IndexOf for better performance.
func IndexOfStr(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

// ContainsStr check if a string slice contains a value.
// Prefer using this helper instead of Contains for better performance.
func ContainsStr(slice []string, value string) bool {
	return IndexOfStr(slice, value) != -1
}

// SliceEqual check if two generic slices are the same.
func SliceEqual(first interface{}, second interface{}) bool {
	l1 := reflect.ValueOf(first)
	l2 := reflect.ValueOf(second)
	length := l1.Len()
	if length != l2.Len() {
		return false
	}

	for i := 0; i < length; i++ {
		if l1.Index(i).Interface() != l2.Index(i).Interface() {
			return false
		}
	}
	return true
}

// ToFloat64 convert a numeric value to float64.
func ToFloat64(value interface{}) (float64, error) {
	return strconv.ParseFloat(ToString(value), 64)
}

// ToString convert a value to string.
func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

// HeaderValue represent a value and its quality value (priority)
// in a multi-values HTTP header.
type HeaderValue struct {
	Value    string
	Priority float64
}

// EscapeLike escape "%" and "_" characters in the given string
// for use in SQL "LIKE" clauses.
func EscapeLike(str string) string {
	escapeChars := []string{"%", "_"}
	for _, v := range escapeChars {
		str = strings.ReplaceAll(str, v, "\\"+v)
	}
	return str
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}
