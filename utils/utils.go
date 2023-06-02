package utils

import (
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func SetPaging(queryMap *url.Values, tx *gorm.DB) {
	page := queryMap.Get("page")
	rowsPerPage := queryMap.Get("rows_per_page")
	if page == "" || rowsPerPage == "" {
		return
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return
	}

	rowsPerPageInt, err := strconv.Atoi(rowsPerPage)
	if err != nil {
		return
	}

	if pageInt < 0 || rowsPerPageInt < 0 {
		return
	}

	pageInt -= 1
	tx.Limit(rowsPerPageInt).Offset(pageInt * rowsPerPageInt)
}

func GetRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RemoveDuplicateValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func SliceContains[T comparable](s []T, e T) bool {
    for _, v := range s {
        if v == e {
            return true
        }
    }
    return false
}
