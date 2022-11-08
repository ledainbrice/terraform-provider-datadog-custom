package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Restriction struct {
	Id    string
	Query string
	Roles []Role
}

func (cl *ClientDatadog) GetRestriction(id string) (Restriction, error) {
	restrictions, err := cl.ReadRestrictionQueries()
	if err != nil {
		return Restriction{}, err
	}
	for _, restriction := range restrictions {
		if restriction.Id == id {
			return restriction, nil
		}
	}

	return Restriction{}, errors.New("restriction not found")
}

func (cl *ClientDatadog) ReadRestrictionQueries() ([]Restriction, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/logs/config/restriction_queries", cl.url), nil)
	if err != nil {
		return nil, err
	}
	cl.auth(req)

	r, err := cl.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	rqs := make(map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&rqs)
	if err != nil {
		return nil, err
	}
	res := rqs["data"]
	aux := make([]Restriction, 0)
	for _, restriction := range res.([]interface{}) {
		restriction2 := restriction.(map[string]interface{})
		id := restriction2["id"].(string)
		attributes := restriction2["attributes"].(map[string]interface{})
		query := attributes["restriction_query"].(string)

		roles, err := cl.ReadRestrictionRoles(id)
		if err != nil {
			return nil, err
		}
		r := Restriction{
			Id:    id,
			Query: query,
			Roles: roles,
		}

		aux = append(aux, r)
	}

	return aux, nil
}

func (cl *ClientDatadog) CreateRestrictionQuery(query string, roles []Role) (string, error) {
	body := []byte(`{
		"data": {
			"attributes": {
				"restriction_query": "` + query + `" 
		 	},
		 	"type": "logs_restriction_queries"
		}
	}`)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/logs/config/restriction_queries", cl.url), bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	cl.auth(req)

	r, err := cl.client.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	rq := make(map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		return "", err
	}

	aux := rq["data"].(map[string]interface{})
	id := aux["id"].(string)

	for _, role := range roles {
		err := cl.CreateRestrictionRole(id, role.Id)
		if err != nil {
			return "", err
		}
	}
	return id, nil
}

func (cl *ClientDatadog) UpdateRestrictionQuery(id string, query string, roles []Role) error {
	body := []byte(`{
		"data": {
			"attributes": {
				"restriction_query": "` + query + `" 
		 	},
		 	"type": "logs_restriction_queries"
		}
	}`)
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v2/logs/config/restriction_queries/%s", cl.url, id), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	cl.auth(req)

	r, err := cl.client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	current_roles, err := cl.ReadRestrictionRoles(id)
	if err != nil {
		return err
	}
	current_role_ids := flatRoles(current_roles)
	new_role_ids := flatRoles(roles)

	remove_roles, add_roles := diff(new_role_ids, current_role_ids)

	for _, role_id := range add_roles {
		err = cl.CreateRestrictionRole(id, role_id)
		if err != nil {
			return err
		}
	}

	for _, role_id := range remove_roles {
		err = cl.DeleteRestrictionRole(id, role_id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cl *ClientDatadog) DeleteRestrictionQuery(id string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v2/logs/config/restriction_queries/%s", cl.url, id), nil)
	if err != nil {
		return err
	}
	cl.auth(req)

	r, err := cl.client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}
