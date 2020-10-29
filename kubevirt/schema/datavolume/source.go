package datavolume

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

func dataVolumeSourceFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"HTTP": dataVolumeSourceHTTPSchema(),
		"PVC":  dataVolumeSourcePVCSchema(),
	}
}

func dataVolumeSourceSchema() *schema.Schema {
	fields := dataVolumeSourceFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: fmt.Sprintf("Source is the src of the data for the requested DataVolume."),
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func dataVolumeSourceHTTPFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"URL": {
			Type:        schema.TypeString,
			Description: "URL is the URL of the http source.",
			Required:    true,
		},
		"secretRef": {
			Type:        schema.TypeString,
			Description: "SecretRef provides the secret reference needed to access the HTTP source.",
			Optional:    true,
		},
		"certConfigMap": {
			Type:        schema.TypeString,
			Description: "CertConfigMap provides a reference to the Registry certs.",
			Optional:    true,
		},
	}
}

func dataVolumeSourceHTTPSchema() *schema.Schema {
	fields := dataVolumeSourceHTTPFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "DataVolumeSourceHTTP provides the parameters to create a Data Volume from an HTTP source.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func dataVolumeSourcePVCFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"namespace": {
			Type:        schema.TypeString,
			Description: "The namespace which the PVC located in.",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the PVC.",
			Optional:    true,
		},
	}
}

func dataVolumeSourcePVCSchema() *schema.Schema {
	fields := dataVolumeSourcePVCFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "DataVolumeSourcePVC provides the parameters to create a Data Volume from an existing PVC.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

// Expanders

func expandDataVolumeSource(dataVolumeSource []interface{}) *cdiv1.DataVolumeSource {
	result := &cdiv1.DataVolumeSource{}

	if len(dataVolumeSource) == 0 || dataVolumeSource[0] == nil {
		return result
	}

	in := dataVolumeSource[0].(map[string]interface{})

	result.HTTP = expandDataVolumeSourceHTTP(in["HTTP"].([]interface{}))
	result.PVC = expandDataVolumeSourcePVC(in["PVC"].([]interface{}))

	return result
}

func expandDataVolumeSourceHTTP(dataVolumeSourceHTTP []interface{}) *cdiv1.DataVolumeSourceHTTP {
	result := &cdiv1.DataVolumeSourceHTTP{}

	if len(dataVolumeSourceHTTP) == 0 || dataVolumeSourceHTTP[0] == nil {
		return result
	}

	in := dataVolumeSourceHTTP[0].(map[string]interface{})

	if v, ok := in["URL"].(string); ok {
		result.URL = v
	}
	if v, ok := in["secretRef"].(string); ok {
		result.SecretRef = v
	}
	if v, ok := in["certConfigMap"].(string); ok {
		result.CertConfigMap = v
	}

	return result
}

func expandDataVolumeSourcePVC(dataVolumeSourcePVC []interface{}) *cdiv1.DataVolumeSourcePVC {
	result := &cdiv1.DataVolumeSourcePVC{}

	if len(dataVolumeSourcePVC) == 0 || dataVolumeSourcePVC[0] == nil {
		return result
	}

	in := dataVolumeSourcePVC[0].(map[string]interface{})

	if v, ok := in["namespace"].(string); ok {
		result.Namespace = v
	}
	if v, ok := in["name"].(string); ok {
		result.Name = v
	}

	return result
}

// Flatteners

func flattenDataVolumeSource(in cdiv1.DataVolumeSource) []interface{} {
	att := make(map[string]interface{})

	if in.HTTP != nil {
		att["HTTP"] = flattenDataVolumeSourceHTTP(*in.HTTP)
	}
	if in.PVC != nil {
		att["PVC"] = flattenDataVolumeSourcePVC(*in.PVC)
	}

	return []interface{}{att}
}

func flattenDataVolumeSourceHTTP(in cdiv1.DataVolumeSourceHTTP) []interface{} {
	att := make(map[string]interface{})

	att["URL"] = in.URL
	att["secretRef"] = in.SecretRef
	att["certConfigMap"] = in.CertConfigMap

	return []interface{}{att}
}

func flattenDataVolumeSourcePVC(in cdiv1.DataVolumeSourcePVC) []interface{} {
	att := make(map[string]interface{})

	att["namespace"] = in.Namespace
	att["Name"] = in.Name

	return []interface{}{att}
}
