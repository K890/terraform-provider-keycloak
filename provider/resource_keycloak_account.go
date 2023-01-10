package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakAccountCreate,
		ReadContext:   resourceKeycloakAccountRead,
		DeleteContext: resourceKeycloakAccountDelete,
		UpdateContext: resourceKeycloakAccountUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakAccountImport,
		},

		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": {
				Type: schema.TypeString,
				Required: false,
				Computed: true,
				ForceNew: false,
			},
			"attrs": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"apps": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"created_on": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}

}

func mapFromDataToAccount(data *schema.ResourceData) *keycloak.Account {
	attributes := map[string]string{}
	if v, ok := data.GetOk("attrs"); ok {
		for key, value := range v.(map[string]interface{}) {
			attributes[key] = value.(string)
		}
	}
	account := &keycloak.Account{
		RealmId:    data.Get("realm_id").(string),
		Name:       data.Get("name").(string),
		Attributes: attributes,
	}
	return account
}

func mapFromAccountToData(data *schema.ResourceData, account *keycloak.Account) {
	attributes := map[string]string{}
	for k, v := range account.Attributes {
		attributes[k] = v
	}
	data.SetId(account.AccountId)
	data.Set("account_id", account.AccountId)
	data.Set("name", account.Name)
}


func resourceKeycloakAccountCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	account := mapFromDataToAccount(data)
	err := keycloakClient.NewAccount(ctx, account)
	if err != nil {
		return diag.FromErr(err)
	}
	mapFromAccountToData(data, account)

	return resourceKeycloakAccountRead(ctx, data, meta)
}


func resourceKeycloakAccountRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	id := data.Get("account_id").(string)

	account, err := keycloakClient.GetAccount(ctx, realmId, id)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}
	mapFromAccountToData(data, account)

	return nil
}


func resourceKeycloakAccountUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	account := mapFromDataToAccount(data)

	err := keycloakClient.UpdateAccount(ctx, account)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromAccountToData(data, account)

	return nil
}

func resourceKeycloakAccountDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realm_id := data.Get("realm_id").(string)
	id := data.Get("account_id").(string)
	return diag.FromErr(keycloakClient.DeleteAccount(ctx,realm_id,id))
}

func resourceKeycloakAccountImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d},nil
}

