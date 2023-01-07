package keycloak

import (
	"context"
	"fmt"
	"encoding/json"
)

type Account struct {
	AccountId   string              `json:"id,omitempty"`
	RealmId     string              `json:"-"`
	Name        string              `json:"name"`
	Attributes  map[string]string   `json:"attrs"`
	Apps        []string 			`json:"apps,omitempty"`
	CreatedOn   int					`json:"createdOn,omitempty"`
}

func (keycloakClient *KeycloakClient) NewAccount(ctx context.Context, account *Account) error {
	var createAccountUrl string = fmt.Sprintf("/realms/%s/api/v1/accounts", account.RealmId)
	var err error
	var body []byte

	body, _, err = keycloakClient.postWithoutAdmin(ctx, createAccountUrl, account)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &account)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) GetAccount(ctx context.Context, realmId, id string) (*Account, error) {
	var account Account

	err := keycloakClient.getWithoutAdmin(ctx, fmt.Sprintf("/realms/%s/api/v1/accounts/%s", realmId, id), &account, nil)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
