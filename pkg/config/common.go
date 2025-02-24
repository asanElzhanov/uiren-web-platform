package config

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type value struct {
	value  any
	exists bool
}

type Value interface {
	// LookupString retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupString() (string, bool)
	// String retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupString.
	String() string
	// LookupInt retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupInt() (int, bool)
	// Int retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupInt.
	Int() int
	// LookupInt32 retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupInt32() (int32, bool)
	// Int32 retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupInt32.
	Int32() int32
	// LookupInt64 retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupInt64() (int64, bool)
	// Int64 retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupInt64.
	Int64() int64
	// LookupBoolean retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupBoolean() (bool, bool)
	// Boolean retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupBoolean.
	Boolean() bool
	// LookupDuration retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupDuration() (time.Duration, bool)
	// Duration retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupDuration.
	Duration() time.Duration
	// LookupFloat32 retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupFloat32() (float32, bool)
	// Float32 retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupFloat32.
	Float32() float32
	// LookupFloat64 retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupFloat64() (float64, bool)
	// Float64 retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupFloat64.
	Float64() float64
	// LookupStrings retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupStrings() ([]string, bool)
	// Strings retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupStrings.
	Strings() []string
	// LookupInts retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupInts() ([]int, bool)
	// Ints retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupInts.
	Ints() []int
	// LookupInts32 retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupInts32() ([]int32, bool)
	// Ints32 retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupInts32.
	Ints32() []int32
	// LookupInts64 retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupInts64() ([]int64, bool)
	// Ints64 retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupInts64.
	Ints64() []int64
	// LookupFloats32 retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupFloats32() ([]float32, bool)
	// Floats32 retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupFloats32.
	Floats32() []float32
	// LookupFloats64 retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupFloats64() ([]float64, bool)
	// Floats64 retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupFloats64.
	Floats64() []float64
	// LookupBooleans retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupBooleans() ([]bool, bool)
	// Booleans retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupBooleans.
	Booleans() []bool
	// LookupDurations retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupDurations() ([]time.Duration, bool)
	// Durations retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupDurations.
	Durations() []time.Duration
	// LookupMap retrieves the value of the config variable.
	// If the variable is present in the config
	// the value (which may be empty) is returned and the boolean
	// is true. Otherwise the returned value will be empty and the
	// boolean will be false.
	LookupMap() (map[string]interface{}, bool)
	// Map retrieves the value of the config variable.
	// It returns the value, which will be empty if the variable is not present.
	// To distinguish between an empty value and an unset value, use LookupMap.
	Map() map[string]interface{}
}

func (v value) LookupString() (string, bool) {
	if !v.exists {
		return "", false
	}
	s, ok := v.value.(string)
	if !ok {
		return "", false
	}
	return s, true
}

func (v value) String() string {
	val, _ := v.LookupString()
	return val
}

func (v value) LookupInt() (int, bool) {
	if !v.exists {
		return 0, false
	}
	i, ok := v.value.(int)
	if !ok {
		return 0, false
	}
	return i, true
}

func (v value) Int() int {
	val, _ := v.LookupInt()
	return val
}

func (v value) LookupInt32() (int32, bool) {
	if !v.exists {
		return 0, false
	}
	i, ok := v.value.(int32)
	if !ok {
		return 0, false
	}
	return i, true
}

func (v value) Int32() int32 {
	val, _ := v.LookupInt32()
	return val
}

func (v value) LookupInt64() (int64, bool) {
	if !v.exists {
		return 0, false
	}
	i, ok := v.value.(int64)
	if !ok {
		return 0, false
	}
	return i, true
}

func (v value) Int64() int64 {
	val, _ := v.LookupInt64()
	return val
}

func (v value) LookupBoolean() (bool, bool) {
	if !v.exists {
		return false, false
	}
	b, ok := v.value.(bool)
	if !ok {
		return false, false
	}
	return b, true
}

func (v value) Boolean() bool {
	val, _ := v.LookupBoolean()
	return val
}

func (v value) LookupFloat32() (float32, bool) {
	if !v.exists {
		return 0, false
	}
	i, ok := v.value.(float32)
	if !ok {
		return 0, false
	}
	return i, true
}

func (v value) Float32() float32 {
	val, _ := v.LookupFloat32()
	return val
}

func (v value) LookupFloat64() (float64, bool) {
	if !v.exists {
		return 0, false
	}
	i, ok := v.value.(float64)
	if !ok {
		return 0, false
	}
	return i, true
}

func (v value) Float64() float64 {
	val, _ := v.LookupFloat64()
	return val
}

func (v value) LookupStrings() ([]string, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []string
		str   string
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]string, 0, len(slice))

	for _, val = range slice {
		str, ok = val.(string)
		if !ok {
			return nil, false
		}
		res = append(res, str)
	}
	return res, true
}

func (v value) Strings() []string {
	val, _ := v.LookupStrings()
	return val
}

func (v value) LookupInts() ([]int, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []int
		i     int
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]int, 0, len(slice))

	for _, val = range slice {
		i, ok = val.(int)
		if !ok {
			return nil, false
		}
		res = append(res, i)
	}
	return res, true
}

func (v value) Ints() []int {
	val, _ := v.LookupInts()
	return val
}

func (v value) LookupInts32() ([]int32, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []int32
		i     int32
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]int32, 0, len(slice))
	for _, val = range slice {
		i, ok = val.(int32)
		if !ok {
			return nil, false
		}
		res = append(res, i)
	}
	return res, true
}

func (v value) Ints32() []int32 {
	val, _ := v.LookupInts32()
	return val
}

func (v value) LookupInts64() ([]int64, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []int64
		i     int64
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]int64, 0, len(slice))
	for _, val = range slice {
		i, ok = val.(int64)
		if !ok {
			return nil, false
		}
		res = append(res, i)
	}
	return res, true
}

func (v value) Ints64() []int64 {
	val, _ := v.LookupInts64()
	return val
}

func (v value) LookupFloats32() ([]float32, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []float32
		i     float32
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]float32, 0, len(slice))
	for _, val = range slice {
		i, ok = val.(float32)
		if !ok {
			return nil, false
		}
		res = append(res, i)
	}
	return res, true
}

func (v value) Floats32() []float32 {
	val, _ := v.LookupFloats32()
	return val
}

func (v value) LookupFloats64() ([]float64, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []float64
		i     float64
		ok    bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]float64, 0, len(slice))
	for _, val = range slice {
		i, ok = val.(float64)
		if !ok {
			return nil, false
		}
		res = append(res, i)
	}
	return res, true
}

func (v value) Floats64() []float64 {
	val, _ := v.LookupFloats64()
	return val
}

func (v value) LookupBooleans() ([]bool, bool) {
	var (
		slice []interface{}
		val   interface{}
		res   []bool
		b, ok bool
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]bool, 0, len(slice))
	for _, val = range slice {
		b, ok = val.(bool)
		if !ok {
			return nil, false
		}
		res = append(res, b)
	}
	return res, true
}

func (v value) Booleans() []bool {
	val, _ := v.LookupBooleans()
	return val
}

func (v value) LookupDuration() (time.Duration, bool) {
	var str string
	var ok bool
	if !v.exists {
		return 0, false
	}
	str, ok = v.value.(string)
	if !ok {
		return 0, false
	}

	var (
		err      error
		duration time.Duration
	)
	if str[len(str)-1:] == "y" {
		years, err := extractYears(str)
		if err != nil {
			return 0, false
		}
		duration = time.Duration(years) * (time.Hour * 24 * 365)

	} else {
		duration, err = time.ParseDuration(str)
		if err != nil {
			return 0, false
		}
	}

	return duration, true
}

func extractYears(input string) (int, error) {
	re := regexp.MustCompile(`(\d+)y`)
	match := re.FindStringSubmatch(input)

	if len(match) < 2 {
		return 0, fmt.Errorf("no match found")
	}

	years, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}

	return years, nil
}

func (v value) Duration() time.Duration {
	val, _ := v.LookupDuration()
	return val
}

func (v value) LookupDurations() ([]time.Duration, bool) {
	var (
		slice []interface{}
		val   interface{}
		str   string
		res   []time.Duration
		d     time.Duration
		ok    bool
		err   error
	)
	if !v.exists {
		return nil, false
	}
	slice, ok = v.value.([]interface{})
	if !ok {
		return nil, false
	}
	res = make([]time.Duration, 0, len(slice))
	for _, val = range slice {
		str, ok = val.(string)
		if !ok {
			return nil, false
		}
		d, err = time.ParseDuration(str)
		if err != nil {
			return nil, false
		}
		res = append(res, d)
	}
	return res, true
}

func (v value) Durations() []time.Duration {
	val, _ := v.LookupDurations()
	return val
}

func (v value) LookupMap() (map[string]interface{}, bool) {
	var m map[string]interface{}
	var ok bool
	if !v.exists {
		return nil, false
	}
	m, ok = v.value.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return m, true
}

func (v value) Map() map[string]interface{} {
	ret, _ := v.LookupMap()
	return ret
}
