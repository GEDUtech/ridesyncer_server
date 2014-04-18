package utils

import (
	"fmt"
	"regexp"
)

var (
	alphanumericRegex *regexp.Regexp
	emailRegex        *regexp.Regexp
)

type Validation struct {
	ErrorCallback func(string, string)
}

func (this *Validation) Apply(field string, validation func() bool, message ...string) bool {
	if !validation() {
		this.ErrorCallback(field, message[0])
		return false
	}
	return true
}

func (this *Validation) NotEmpty(field, v string, message ...string) bool {
	if len(message) == 0 {
		message = append(message, "Required")
	}
	return this.Apply(field, func() bool {
		return len(v) > 0
	}, message...)
}

func (this *Validation) Alphanumeric(field, v string, message ...string) bool {
	if len(message) == 0 {
		message = append(message, "Alphanumeric")
	}
	return this.Apply(field, func() bool {
		return alphanumericRegex.MatchString(v)
	}, message...)
}

func (this *Validation) Between(field, v string, min, max int, message ...string) bool {
	if len(message) == 0 {
		message = append(message, fmt.Sprintf("Length should be between %d and %d", min, max))
	}
	return this.Apply(field, func() bool {
		return len(v) >= min && len(v) <= max
	}, message...)
}

func (this *Validation) Match(field, v1, v2 string, message ...string) bool {
	if len(message) == 0 {
		message = append(message, "Values do not match")
	}
	return this.Apply(field, func() bool {
		return v1 == v2
	}, message...)
}

func (this *Validation) Email(field, v string, message ...string) bool {
	if len(message) == 0 {
		message = append(message, "Invalid email")
	}
	return this.Apply(field, func() bool {
		return emailRegex.MatchString(v)
	}, message...)
}

func Range(v, min, max int) bool {
	return v >= min && v <= max
}

func (this *Validation) Bitwise(field string, v uint, min uint, max uint, messages ...string) bool {
	messages = defaultMessage("Invalid option", messages)
	return this.Apply(field, func() bool {
		if v < min || v > max {
			return false
		}
		return (v & ^(v - 1)) == v
	}, messages...)
}

func defaultMessage(message string, messages []string) []string {
	if len(messages) == 0 {
		messages = append(messages, message)
	}
	return messages
}

func init() {
	alphanumericRegex, _ = regexp.Compile("^[\\p{Ll}\\p{Lm}\\p{Lo}\\p{Lt}\\p{Lu}\\p{Nd}]+$/Du")
	emailRegex, _ = regexp.Compile("^[^ @]+@[^ @]+\\.edu$")
}
