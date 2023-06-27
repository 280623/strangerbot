package service

import (
	"context"
	"strings"

	"strangerbot/repository"
	"strangerbot/vars"
)

func ServiceValidWhiteEmail(ctx context.Context, email string) (bool, error) {

	domain := false
	if vars.WhiteDomainEnabled {
		if len(vars.WhiteDomain) > 0 {
			domainList := strings.Split(vars.WhiteDomain, ",")
			for _, domainEntry := range domainList {
				if strings.Contains(email, domainEntry) {
					domain = true
					break
				}
			}
		}
	}

	emailW := false
	if vars.WhiteEmailEnabled {

		repo := repository.GetRepository()
		emailList, err := repo.GetWhiteEmailAll(ctx)
		if err != nil {
			return false, err
		}

		for _, emailEntry := range emailList {
			if email == emailEntry.Email {
				emailW = true
				break
			}
		}
	}

	if vars.WhiteDomainEnabled || vars.WhiteEmailEnabled {
		if domain || emailW {
			return true, nil
		}
		return false, nil
	}

	return true, nil
}
