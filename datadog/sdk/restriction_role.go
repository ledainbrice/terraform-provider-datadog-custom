package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Role struct {
	Id   string
	Name string
}

func flatRoles(roles []Role) []string {
	res := make([]string, 0)
	for _, role := range roles {
		if role.Id != "" {
			res = append(res, role.Id)
		}
	}
	return res
}

// func CreateRestrictionRoleId(restriction_id string, role_id string) string {
// 	return fmt.Sprintf("%s#%s", restriction_id, role_id)
// }

// func GetRoleIdFromRestrictionRoleId(restriction_role_id string) string {
// 	return strings.Split(restriction_role_id, "#")[1]
// }

// func GetRestrictionIdFromRestrictionRoleId(restriction_role_id string) string {
// 	return strings.Split(restriction_role_id, "#")[0]
// }
// func (cl *ClientDatadog) GetRestrictionRole(restriction_role_id string) (Role, error) {

// 	restriction_id := GetRestrictionIdFromRestrictionRoleId(restriction_role_id)
// 	restrictions, err := cl.ReadRestrictionRoles(restriction_id)
// 	if err != nil {
// 		return Role{}, err
// 	}
// 	for _, restriction := range restrictions {
// 		if restriction.Id == "id" {
// 			return restriction, nil
// 		}
// 	}

// 	return Role{}, errors.New("restriction not found")
// }

func check2xx(status int) bool {
	return status < 300 && status > 199
}

func (cl *ClientDatadog) ReadRestrictionRoles(restriction_id string) ([]Role, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/logs/config/restriction_queries/%s/roles", cl.url, restriction_id), nil)
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

	res := rqs["data"].([]interface{})
	aux := make([]Role, 0)
	for _, ra := range res {
		r := ra.(map[string]interface{})
		fmt.Println(r["attributes"])
		attributes := r["attributes"].(map[string]interface{})
		role := Role{
			Id:   r["id"].(string),
			Name: attributes["name"].(string),
		}
		aux = append(aux, role)
	}
	time.Sleep(8 * time.Second)

	return aux, nil
}

func (cl *ClientDatadog) CreateRestrictionRole(restriction_id string, role_id string) error {
	body := []byte(`{
		"data": {
			"id": "` + role_id + `",
		 	"type": "roles"
		}
	}`)
	fmt.Println(restriction_id)
	fmt.Println(string(body))
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/logs/config/restriction_queries/%s/roles", cl.url, restriction_id), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	cl.auth(req)

	r, err := cl.client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	status := r.StatusCode
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	responseString := string(responseData)
	log.Printf("[INFO] '%s'", responseString)
	log.Printf("[INFO] " + strconv.Itoa(status))
	if !check2xx(r.StatusCode) {
		return errors.New("[CreateRestrictionRole] Request status: " + strconv.Itoa(status) + " , response: " + responseString)
	}
	time.Sleep(8 * time.Second)

	return nil
}

func (cl *ClientDatadog) DeleteRestrictionRole(restriction_id string, role_id string) error {
	body := []byte(`{
		"data": {
			"id": "` + role_id + `",
		 	"type": "roles"
		}
	}`)
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v2/logs/config/restriction_queries/%s/roles", cl.url, restriction_id), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	cl.auth(req)

	r, err := cl.client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	status := r.StatusCode
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	responseString := string(responseData)
	log.Printf("[INFO] '%s'", responseString)
	log.Printf("[INFO] " + strconv.Itoa(status))
	if !check2xx(r.StatusCode) {
		return errors.New("[DeleteRestrictionRole] Request status: " + strconv.Itoa(status) + " , response: " + responseString)
	}
	time.Sleep(8 * time.Second)

	return nil
}
