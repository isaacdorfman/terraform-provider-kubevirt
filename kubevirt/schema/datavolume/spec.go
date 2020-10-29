package datavolume

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/common"
	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

func dataVolumeSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"source": dataVolumeSourceSchema(),
		"PVC":    common.PersistentVolumeClaimSpecSchema(),
		"contentType": {
			Type:        schema.TypeString,
			Description: "ContentType options: \"kubevirt\", \"archive\".",
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"kubevirt",
				"archive",
			}, false),
		},
	}
}

func dataVolumeSpecSchema() *schema.Schema {
	fields := dataVolumeSpecFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: fmt.Sprintf("DataVolumeSpec defines our specification for a DataVolume type"),
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandDataVolumeSpec(dataVolumeSpec []interface{}) (*cdiv1.DataVolumeSpec, error) {
	result := &cdiv1.DataVolumeSpec{}

	if len(dataVolumeSpec) == 0 || dataVolumeSpec[0] == nil {
		return result, nil
	}

	in := dataVolumeSpec[0].(map[string]interface{})

	result.Source = *expandDataVolumeSource(in["source"].([]interface{}))
	p, err := common.ExpandPersistentVolumeClaimSpec(in["PVC"].([]interface{}))
	if err != nil {
		return result, err
	}
	result.PVC = p

	if v, ok := in["contentType"].(string); ok {
		result.ContentType = cdiv1.DataVolumeContentType(v)
	}

	return result, nil
}

func flattenDataVolumeSpec(spec cdiv1.DataVolumeSpec) []interface{} {
	att := make(map[string]interface{})
	att["source"] = flattenDataVolumeSource(spec.Source)
	att["PVC"] = common.FlattenPersistentVolumeClaimSpec(*spec.PVC)
	att["contentType"] = string(spec.ContentType)
	return []interface{}{att}
}
