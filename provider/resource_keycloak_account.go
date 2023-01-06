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
				ForceNew: false,
			},
		},
	}

}

func mapFromDataToAccount(data *schema.ResourceData) *keycloak.Account {
	account := &keycloak.Account{
		Id:         data.Id(),
		RealmId:    data.Get("realm_id").(string),
		Name:       data.Get("name").(string),
	}
	return account
}

func mapFromAccountToData(data *schema.ResourceData, account *keycloak.Account) {
	
	data.SetId(account.Id)
	data.Set("realm_id", account.RealmId)
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


	return nil
}


func resourceKeycloakAccountRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	id := data.Id()

	account, err := keycloakClient.GetAccount(ctx, realmId, id)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	mapFromAccountToData(data, account)

	return nil
}


func resourceKeycloakAccountUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{{
		Summary:  "Updating an Account is not yet supported",
		Severity: diag.Warning,
	}}
}

func resourceKeycloakAccountDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	
	return diag.Diagnostics{{
		Summary:  "Updating an Account is not yet supported",
		Severity: diag.Warning,
	}}
}

func resourceKeycloakAccountImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d},nil
}

