package utils

import (
	"regexp"
)

func IsValidVietnamesePhoneNumber(phone string) bool {
	re := regexp.MustCompile(`^((\+84|84|0)[3|5|7|8|9])+([0-9]{8})$`)
	return re.MatchString(phone)
}