package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainStat := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var user User

		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}

		if user.Email == "" {
			continue
		}

		if !strings.HasSuffix(user.Email, domain) {
			continue
		}

		domain := strings.ToLower(user.Email[strings.LastIndex(user.Email, "@")+1:])
		domainStat[domain]++
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("scanner failed: %w", scanner.Err())
	}

	return domainStat, nil
}
