package kubevirt

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/client"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/common"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/datavolume"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

func resourceKubevirtDataVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubevirtDataVolumeCreate,
		Read:   resourceKubevirtDataVolumeRead,
		Update: resourceKubevirtDataVolumeUpdate,
		Delete: resourceKubevirtDataVolumeDelete,
		Exists: resourceKubevirtDataVolumeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: datavolume.DataVolumeFields(),
	}
}

func resourceKubevirtDataVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	dv, err := datavolume.ExpandDataVolume(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating new data volume: %#v", dv)
	if err := cli.CreateDataVolume(*dv); err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new data volume: %#v", dv)

	d.SetId(utils.BuildId(dv.ObjectMeta))

	// Wait for data volume instance's status phase to be succeeded:
	name := dv.ObjectMeta.Name
	namespace := dv.ObjectMeta.Namespace

	stateConf := &resource.StateChangeConf{
		Pending: []string{"Creating"},
		Timeout: d.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			var err error
			dv, err = cli.ReadDataVolume(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Printf("[DEBUG] data volume %s is not created yet", name)
					return dv, "Creating", nil
				}
				return dv, "", err
			}

			switch dv.Status.Phase {
			case cdiv1.Succeeded:
				return dv, "", nil
			case cdiv1.Failed:
				return dv, "", fmt.Errorf("data volume failed to be created, finished with phase=\"failed\"")
			}

			log.Printf("[DEBUG] data volume %s is being created", name)
			return dv, "Creating", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}
	return datavolume.FlattenDataVolume(*dv, d)
}

func resourceKubevirtDataVolumeRead(d *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading data volume %s", name)

	dv, err := cli.ReadDataVolume(namespace, name)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received data volume: %#v", dv)

	return datavolume.FlattenDataVolume(*dv, d)
}

func resourceKubevirtDataVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(d.Id())
	if err != nil {
		return err
	}

	ops := common.PatchMetadata("metadata.0.", "/metadata/", d)
	data, err := ops.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal update operations: %s", err)
	}

	log.Printf("[INFO] Updating data volume: %s", ops)
	out := &cdiv1.DataVolume{}
	if err := cli.UpdateDataVolume(namespace, name, out, data); err != nil {
		return err
	}

	log.Printf("[INFO] Submitted updated persistent volume claim: %#v", out)

	return resourceKubevirtDataVolumeRead(d, meta)
}

func resourceKubevirtDataVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting data volume: %#v", name)
	if err := cli.DeleteDataVolume(namespace, name); err != nil {
		return err
	}

	// Wait for data volume instance to be removed:
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Timeout: d.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			dv, err := cli.ReadDataVolume(namespace, name)
			if err != nil {
				if errors.IsNotFound(err) {
					return nil, "", nil
				}
				return dv, "", err
			}

			log.Printf("[DEBUG] data volume %s is being deleted", dv.GetName())
			return dv, "Deleting", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("%s", err)
	}

	log.Printf("[INFO] data volume %s deleted", name)

	d.SetId("")
	return nil
}

func resourceKubevirtDataVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cli := (meta).(client.Client)

	namespace, name, err := utils.IdParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking data volume %s", name)
	if _, err := cli.ReadDataVolume(namespace, name); err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
		return true, err
	}
	return true, nil
}
