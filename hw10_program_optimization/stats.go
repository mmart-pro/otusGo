package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	var user User
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		if err = user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return
		}
		result[i] = user
		i++
	}
	err = scanner.Err()
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	re, _ := regexp.Compile("\\." + domain)
	for _, user := range u {
		matched := re.MatchString(user.Email)
		if matched {
			domain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[domain]++
		}
	}
	return result, nil
}
