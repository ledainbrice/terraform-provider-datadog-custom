package datadog

import (
	"context"
	"strconv"
	"time"

	"terraform-provider-datadog/datadog/sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRestrictions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRestrictionsRead,
		Schema: map[string]*schema.Schema{
			"restrictions": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"query": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"roles": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceRestrictionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cl := m.(*sdk.ClientDatadog)
	var diags diag.Diagnostics
	restrictions, err := cl.ReadRestrictionQueries()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("restrictions", restrictions); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
