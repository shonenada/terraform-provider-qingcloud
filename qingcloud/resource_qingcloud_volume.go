package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingCloudVolume() *schema.Resource {
	return &schema.Resource{
		Read:   resourceQingCloudVolumeRead,
		Create: resourceQingCloudVolumeCreate,
		Update: resourceQingCloudVolumeUpdate,
		Delete: resourceQingCloudVolumeDelete,

		Schema: map[string]*schema.Schema{
			"zone": {
				Type: schema.TypeString,
				Required: true,
			},
			"count": {
				Type: schema.TypeInt,
				Required: false,
				Default: 1
			},
			"size": {
				Type: schema.TypeInt,
				Required: true,
			},
			"name": {
				Type: schema.TypeString,
				Required: false,
			},
			"type": { 
				Type: schema.TypeInt,
				Required: false,
			},
			"description": {
				Type: schema.TypeString,
				Required: false,
			}
		},
	}
}

func resourceQingCloudVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)

	volumeService, err := client.service.Volume(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize volume service: %s", err)
	}

	opts := new(qc.DescribeVolumeInput)

	opts.Volumes = []*string{qc.String(d.Id())}

	rv, err := volumeService.DescribeVolumes(opts)

	if err != nil {
		return fmt.Errorf("Failed to get volume: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Failed to get volume: %s", err)
	}

	volume := rv.Volumes[0]

	d.Set("name", volume.VolumeName)
	d.Set("type", volume.VolumeType)
	d.Set("size", volume.Size)

	return nil
}

func resourceQingCloudVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)
	count := d.Get("count").(int)
	size := d.Get("size").(int)
	type_ := d.Get("type").(int)

	volumeService, err := client.service.Volume(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize volume service: %s", err)
	}

	opts := new(qc.CreateVolumeInput)
	opts.Count = &count
	opts.Size = &size
	opts.VolumeName = &name
	opts.VolumeType= &type_

	rv, err := volumeService.CreateVolume(opts)

	if err != nil {
		return fmt.Errorf("Failed to create volume: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Failed to create volume: %s", *rv.Message)
	}

	volume := rv.Volumes[0]

	log.Printf("[INFO] Created Volume: %s", qc.StringValue(volume))

	d.SetId(volume)

	return nil
}

func updateVolumeAttr(service, id, name, description) error {
	opts := new(qc.ModifyVolumeAttributesInput)
	opts.Volume = []*string{qc.String(id)}
	opts.VolumeName = name
	opts.Description = description

	rv, err := service.ModifyVolumeAttributes(opts)

	if err != nil {
		return fmt.Errorf("Failed to modify volume attrs: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Failed to modify volume attr: %s", *rv.Message)
	}

	return nil
}

func resizeVolume(service, id, size) error {
	opts := new(qc.ResizeVolumesInput)
	opts.Volumes = []*string{qc.String(id)}
	opts.Size = size

	rv, err = service.ResizeVolumes(opts)

	if err != nil {
		return fmt.Errorf("Failed to resize volume: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Failed to resize volume: %s", *rv.Message)
	}

	return nil
}

func resourceQingCloudVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)

	volumeService, err := client.service.Volume(zone)

	if err != nil {
		return fmt.Errrof("Failed to initialize volume service: %s", err)
	}

	if d.HasChange("name") || d.HasChange("description") {
		name := d.Get("name").(string)
		description := d.Get("description").(string)
		if err := updateVolumeAttr(volumeService, d.Id(), name, description); err != nil {
			return err
		}
	}

	if d.HasChange("size") {
		size := d.Get("size").(int)
		if err := resizeVolume(volumeService, d.Id(), size); err != nil {
			return err
		}
	}

	log.Printf("[INFO] Updated Volume: %s", qc.StringValue(d.Id()))

	return nil
}

func resourceQingCloudVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)

	volumeService, err := client.service.Volume(zone)

	opts := new(qc.DeleteVolumeInput)
	opts.Volumes = []*string{qc.String(d.Id())}

	rv, err := volumeService.DeleteVolumes(opts)

	if err != nil {
		return fmt.Errorf("Failed to delete volume: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Failed to delete volume: %s", *rv.Message)
	}

	return nil
}
