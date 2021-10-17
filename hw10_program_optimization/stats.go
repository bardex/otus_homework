package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domain = "." + domain
	scanner := bufio.NewScanner(r)
	var user struct{ Email string }
	for scanner.Scan() {
		user.Email = ""
		if err := jsoniter.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}
		if strings.HasSuffix(user.Email, domain) {
			d := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[d]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
