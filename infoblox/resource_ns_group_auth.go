package infoblox

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/skyinfoblox"
	"github.com/sky-uk/skyinfoblox/api/nsgroupauth"
	"github.com/sky-uk/terraform-provider-infoblox/infoblox/util"
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
	if v, _ := d.GetOk("grid_default_group"); v != nil {
		gridDefaultGroup := v.(bool)
		nsGroupAuth.GridDefault = &gridDefaultGroup
	}
	if v, ok := d.GetOk("use_external_primary"); ok && v != nil {
		useExternalPrimary := v.(bool)
		nsGroupAuth.UseExternalPrimary = &useExternalPrimary
	}
	if v, ok := d.GetOk("external_primaries"); ok && v != nil {
		servers := []map[string]interface{}{}
		for _, server := range v.([]interface{}) {
			servers = append(servers, server.(map[string]interface{}))
		}
		nsGroupAuth.ExternalPrimaries = util.BuildExternalServerListFromT(servers)
	}
	if v, ok := d.GetOk("external_secondaries"); ok && v != nil {
		servers := []map[string]interface{}{}
		for _, server := range v.([]interface{}) {
			servers = append(servers, server.(map[string]interface{}))
		}
		nsGroupAuth.ExternalSecondaries = util.BuildExternalServerListFromT(servers)
	}
	if v, ok := d.GetOk("grid_primary"); ok && v != nil {
		servers := []map[string]interface{}{}
		for _, server := range v.([]interface{}) {
			servers = append(servers, server.(map[string]interface{}))
		}
		nsGroupAuth.GridPrimary = util.BuildMemberServerListFromT(servers)
	}
	if v, ok := d.GetOk("grid_secondaries"); ok && v != nil {
		servers := []map[string]interface{}{}
		for _, server := range v.([]interface{}) {
			servers = append(servers, server.(map[string]interface{}))
		}
		nsGroupAuth.GridSecondaries = util.BuildMemberServerListFromT(servers)
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
	d.Set("external_primaries", util.BuildExternalServersListFromIBX(response.ExternalPrimaries))
	d.Set("external_secondaries", util.BuildExternalServersListFromIBX(response.ExternalSecondaries))
	d.Set("grid_primary", util.BuildMemberServerListFromIBX(response.GridPrimary))
	d.Set("grid_secondaries", util.BuildMemberServerListFromIBX(response.GridSecondaries))

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
	if d.HasChange("grid_default_group") {
		if v, ok := d.GetOk("grid_default_group"); ok && v != nil {
			gridDefaultGroup := v.(bool)
			nsGroupAuth.GridDefault = &gridDefaultGroup
		}
		hasChanges = true
	}
	if d.HasChange("use_external_primary") {
		if v, ok := d.GetOk("use_external_primary"); ok && v != nil {
			useExternalPrimary := v.(bool)
			nsGroupAuth.UseExternalPrimary = &useExternalPrimary
		}
		hasChanges = true
	}
	if d.HasChange("external_primaries") {
		if v, ok := d.GetOk("external_primaries"); ok && v != nil {
			servers := []map[string]interface{}{}
			for _, server := range v.([]interface{}) {
				servers = append(servers, server.(map[string]interface{}))
			}
			nsGroupAuth.ExternalPrimaries = util.BuildExternalServerListFromT(servers)
		}
		hasChanges = true
	}
	if d.HasChange("external_secondaries") {
		if v, ok := d.GetOk("external_secondaries"); ok && v != nil {
			servers := []map[string]interface{}{}
			for _, server := range v.([]interface{}) {
				servers = append(servers, server.(map[string]interface{}))
			}
			nsGroupAuth.ExternalSecondaries = util.BuildExternalServerListFromT(servers)
		}
		hasChanges = true
	}
	if d.HasChange("grid_primary") {
		if v, ok := d.GetOk("grid_primary"); ok && v != nil {
			servers := []map[string]interface{}{}
			for _, server := range v.([]interface{}) {
				servers = append(servers, server.(map[string]interface{}))
			}
			nsGroupAuth.GridPrimary = util.BuildMemberServerListFromT(servers)
		}
		hasChanges = true
	}
	if d.HasChange("grid_secondaries") {
		if v, ok := d.GetOk("grid_secondaries"); ok && v != nil {
			servers := []map[string]interface{}{}
			for _, server := range v.([]interface{}) {
				servers = append(servers, server.(map[string]interface{}))
			}
			nsGroupAuth.GridSecondaries = util.BuildMemberServerListFromT(servers)
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
		d.Set("grid_default_group", *response.GridDefault)
		d.Set("use_external_primary", *response.UseExternalPrimary)
		d.Set("external_primaries", util.BuildExternalServersListFromIBX(response.ExternalPrimaries))
		d.Set("external_secondaries", util.BuildExternalServersListFromIBX(response.ExternalSecondaries))
		d.Set("grid_primary", util.BuildMemberServerListFromIBX(response.GridPrimary))
		d.Set("grid_secondaries", util.BuildMemberServerListFromIBX(response.GridSecondaries))
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
