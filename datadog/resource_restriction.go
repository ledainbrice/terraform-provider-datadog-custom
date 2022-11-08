package datadog

import (
	"context"
	"terraform-provider-datadog/datadog/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatadog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatadogCreate,
		ReadContext:   resourceDatadogRead,
		UpdateContext: resourceDatadogUpdate,
		DeleteContext: resourceDatadogDelete,
		Schema: map[string]*schema.Schema{
			"query": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceDatadogCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cl := m.(*sdk.ClientDatadog)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	query := d.Get("query").(string)
	a, exist := d.GetOk("roles")
	aux := make([]interface{}, 0)
	if exist {
		aux = a.([]interface{})
	}

	roles := make([]sdk.Role, 0)
	for _, ra := range aux {
		r := ra.(map[string]interface{})
		role := sdk.Role{
			Id: r["role_id"].(string),
		}
		roles = append(roles, role)
	}

	restriction_id, err := cl.CreateRestrictionQuery(query, roles)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(restriction_id)
	resourceDatadogRead(ctx, d, m)

	return diags
}

func resourceDatadogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cl := m.(*sdk.ClientDatadog)
	var diags diag.Diagnostics
	restriction_id := d.Id()

	restriction, err := cl.GetRestriction(restriction_id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("query", restriction.Query); err != nil {
		return diag.FromErr(err)
	}
	roles := make([]map[string]string, 0)
	for _, r := range restriction.Roles {
		a := map[string]string{
			"role_id": r.Id,
		}
		roles = append(roles, a)
	}
	if err := d.Set("roles", roles); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDatadogUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cl := m.(*sdk.ClientDatadog)
	restriction_id := d.Id()
	if d.HasChange("query") || d.HasChange("roles") {
		query := d.Get("query").(string)
		roles := d.Get("roles").([]interface{})

		role_ids := make([]sdk.Role, 0)
		for _, r := range roles {
			role := r.(map[string]interface{})
			role_ids = append(role_ids, sdk.Role{
				Id: role["role_id"].(string),
			})
		}
		err := cl.UpdateRestrictionQuery(restriction_id, query, role_ids)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDatadogRead(ctx, d, m)
}

func resourceDatadogDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cl := m.(*sdk.ClientDatadog)
	var diags diag.Diagnostics
	restriction_id := d.Id()

	err := cl.DeleteRestrictionQuery(restriction_id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}

// b, err := json.Marshal(roles)
// 		if err != nil {
// 			return diag.FromErr(err)
// 		}
// 		err = errors.New(string(b))
// 		if err != nil {
// 			return diag.FromErr(err)
// 		}
