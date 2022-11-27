/*
 * Otto user service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

import (
	"encoding/json"
	"net/http"
	"time"
)

type User struct {
	Uuid string `json:"uuid"`

	Name string `json:"name,omitempty"`

	Username string `json:"username"`

	ProfilePic string `json:"profile_pic,omitempty"`

	BioMessage string `json:"bio_message,omitempty"`

	Email string `json:"email,omitempty"`

	Password string `json:"password,omitempty"`

	Role Role `json:"role,omitempty"`

	IsBanned bool `json:"is_banned,omitempty"`

	Birthday string `json:"birthday,omitempty"`

	AddressStreet string `json:"address_street,omitempty"`

	AddressCity string `json:"address_city,omitempty"`

	AddressZip string `json:"address_zip,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`

	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func DecodeRequestToUser(r *http.Request) (*User, error) {
	decoder := json.NewDecoder(r.Body)
	var data *User
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
