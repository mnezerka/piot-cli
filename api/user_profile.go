package api

import "fmt"

type UserProfile struct {
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	OrgId   string `json:"org_id"`
	Orgs    []Org  `json:"orgs"`
}

func (p *UserProfile) GetActiveOrg() (*Org, error) {

	for i := 0; i < len(p.Orgs); i++ {
		if p.OrgId == p.Orgs[i].Id {
			return &p.Orgs[i], nil
		}
	}

	return nil, fmt.Errorf("No active organization in current profile.")
}
