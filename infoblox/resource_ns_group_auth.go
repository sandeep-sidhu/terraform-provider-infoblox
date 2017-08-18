package infoblox

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/terraform-provider-infoblox/infoblox/util"
	"github.com/sky-uk/skyinfoblox/api/nsgroupauth"
	"github.com/sky-uk/skyinfoblox"
	"fmt"
	"net/http"
)

func resourceNSGroupAuth() *schema.Resource {
	return &schema.Resource{
		Create: resourceNSGroupAuthCreate,
		Read:   resourceNSGroupAuthRead,
		Update: resourceNSGroupAuthUpdate,
		Delete: resourceNSGroupAuthDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of the NS Group",
				Required:     true,
				ValidateFunc: util.ValidateZoneAuthCheckLeadingTrailingSpaces,
			},
			"comment": {
				Type:         schema.TypeString,
				Description:  "Comment field",
				Optional:     true,
				ValidateFunc: util.ValidateZoneAuthCheckLeadingTrailingSpaces,
			},
			"grid_default_group": {
				Type:        schema.TypeBool,
				Description: "Determines if this name server group is the Grid default",
				Optional:    true,
			},
			"use_external_primary": {
				Type:        schema.TypeBool,
				Description: "This flag controls whether the group is using an external primary",
				Optional:    true,
			},
			"external_primaries":   util.ExternalServerListSchema(true, false),
			"external_secondaries": util.ExternalServerListSchema(true, false),
			"grid_primary":         util.MemberServerListSchema(),
			"grid_secondaries":     util.MemberServerListSchema(),
		},
	}
}

func resourceNSGroupAuthCreate(d *schema.ResourceData, m interface{}) error {

	var nsGroupAuth nsgroupauth.NSGroupAuth
	client := m.(*skyinfoblox.InfobloxClient)

	if v, ok := d.GetOk("name"); ok && v != "" {
		nsGroupAuth.Name = v.(string)
	}
	if v, ok := d.GetOk("comment"); ok && v != "" {
		nsGroupAuth.Comment = v.(string)
	}

	createAPI := nsgroupauth.NewCreate(nsGroupAuth)
	err := client.Do(createAPI)
	httpStatus := createAPI.StatusCode()
	if err != nil || httpStatus < http.StatusOK || httpStatus >= http.StatusBadRequest {
		return fmt.Errorf("Infoblox NS Group Auth Create for %s failed with status code %d and error: %+v", nsGroupAuth.Name, httpStatus, err)
	}
	nsGroupAuth.Reference = *createAPI.ResponseObject().(*string)

	d.SetId(nsGroupAuth.Reference)
	return resourceNSGroupAuthRead(d, m)
}

func resourceNSGroupAuthRead(d *schema.ResourceData, m interface{}) error {

	returnFields := []string{"comment", "external_primaries", "external_secondaries", "grid_primary", "grid_secondaries", "is_grid_default", "name", "use_external_primary"}
	reference := d.Id()
	client := m.(*skyinfoblox.InfobloxClient)

	getNSGroupAuthAPI := nsgroupauth.NewGet(reference, returnFields)
	err := client.Do(getNSGroupAuthAPI)
	httpStatus := getNSGroupAuthAPI.StatusCode()
	if httpStatus == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil || httpStatus < http.StatusOK || httpStatus >= http.StatusBadRequest {
		return fmt.Errorf("Infoblox NS Group Auth Read for %s failed with status code %d and error: %+v", reference, httpStatus, err)
	}

	response := *getNSGroupAuthAPI.ResponseObject().(*nsgroupauth.NSGroupAuth)

	d.SetId(response.Reference)
	d.Set("name", response.Name)
	d.Set("comment", response.Comment)
	d.Set("grid_default_group", *response.GridDefault)
	d.Set("use_external_primary", *response.UseExternalPrimary)
	d.Set("external_primaries", response.ExternalPrimaries)
	d.Set("external_secondaries", response.ExternalSecondaries)
	d.Set("grid_primary", response.GridPrimary)
	d.Set("grid_secondaries", response.GridSecondaries)

	return nil
}

func resourceNSGroupAuthUpdate(d *schema.ResourceData, m interface{}) error {

	var nsGroupAuth nsgroupauth.NSGroupAuth
	hasChanges := false

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok && v != "" {
			nsGroupAuth.Name = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok && v != "" {
			nsGroupAuth.Comment = v.(string)
		}
		hasChanges = true
	}

	if hasChanges {

		client := m.(*skyinfoblox.InfobloxClient)
		returnFields := []string{"comment", "external_primaries", "external_secondaries", "grid_primary", "grid_secondaries", "is_grid_default", "name", "use_external_primary"}
		nsGroupAuth.Reference = d.Id()

		nsGroupAuthUpdateAPI := nsgroupauth.NewUpdate(nsGroupAuth, returnFields)
		err := client.Do(nsGroupAuthUpdateAPI)
		httpStatus := nsGroupAuthUpdateAPI.StatusCode()

		if err != nil || httpStatus < http.StatusOK || httpStatus >= http.StatusBadRequest {
			return fmt.Errorf("Infoblox NS Group Auth Update for %s failed with status code %d and error: %+v", nsGroupAuth.Name, httpStatus, err)
		}
		response := *nsGroupAuthUpdateAPI.ResponseObject().(*nsgroupauth.NSGroupAuth)

		d.SetId(response.Reference)
		d.Set("name", response.Name)
		d.Set("comment", response.Comment)
	}

	return resourceNSGroupAuthRead(d, m)
}

func resourceNSGroupAuthDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*skyinfoblox.InfobloxClient)
	reference := d.Id()

	nsGroupAuthDeleteAPI := nsgroupauth.NewDelete(reference)
	err := client.Do(nsGroupAuthDeleteAPI)
	httpStatus := nsGroupAuthDeleteAPI.StatusCode()

	if httpStatus == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if err != nil || httpStatus < http.StatusOK || httpStatus >= http.StatusBadRequest {
		return fmt.Errorf("Infoblox NS Group Auth Delete for %s failed with status code %d and error: %+v", reference, httpStatus, err)
	}

	d.SetId("")
	return nil
}
