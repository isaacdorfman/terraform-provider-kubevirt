package datavolume

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/common"
	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

func DataVolumeFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": common.NamespacedMetadataSchema("DataVolume", true),
		"spec":     dataVolumeSpecSchema(),
		"status":   dataVolumeStatusSchema(),
	}
}

func DataVolumeSchema() *schema.Schema {
	fields := DataVolumeFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: fmt.Sprintf("DataVolume provides a representation of our data volume."),
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func ExpandDataVolume(d *schema.ResourceData) (*cdiv1.DataVolume, error) {
	result := &cdiv1.DataVolume{}

	result.ObjectMeta = common.ExpandMetadata(d.Get("metadata").([]interface{}))
	spec, err := expandDataVolumeSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return result, err
	}
	result.Spec = *spec
	result.Status = *expandDataVolumeStatus(d.Get("status").([]interface{}))

	return result, nil
}

func FlattenDataVolume(dv cdiv1.DataVolume, d *schema.ResourceData) error {
	if err := d.Set("metadata", common.FlattenMetadata(dv.ObjectMeta, d)); err != nil {
		return err
	}
	if err := d.Set("spec", flattenDataVolumeSpec(dv.Spec)); err != nil {
		return err
	}
	if err := d.Set("status", flattenDataVolumeStatus(dv.Status)); err != nil {
		return err
	}

	return nil
}
