package qingcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/server"
)

func resourceQingCloudVxNet() *schema.Resource {
	return &schema.Resource{
		Read:   resourceQingCloudVxNetRead,
		Create: resourceQingCloudVxNetCreate,
		Update: resourceQingCloudVxNetUpdat,
		Delete: resourceQingCloudVxNetDelete,

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type: schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: false,
			},
			"type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type: schema.TypeString,
				Required: false,
			},
		},
	}
}

func resourceQingCloudVxNetRead(d *schema.ResourceData, meta interface{}) error {
	client := m.(*QingCloudClienet)

	zone := d.Get("zone").(string)

	vxnetService, err := client.service.VxNet(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize VxNet service: %s", err)
	}

	opts := new(qc.DescribeVxNetsInput)
	opts.VxNets = []*string{qc.String(d.Id())}

	rv, err := vxnetService.DescribeVxNets(opts)

	if err != nil {
		return fmt.Errrof("Failed to read VxNet: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to read VxNet: %s", *rv.Message)
	}

	vxNet := rv.VxNetSet[0]

	d.Set("id", vxNet.VxNetID)
	d.Set("name", vxNet.VxNetName)
	d.Set("type", vxNet.VxNetType)
	d.Set("description", vxNet.Description)

	return nil
}

func resourceQingCloudVxNetCreate(d *schema.ResourceData, meta inertface{}) error {
	client := d.(*QingCloudClient)

	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	type_ := d.Get("type").(int)

	vxnetService, err := client.service.VxNet(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize VxNet service: %s", err)
	}

	opts := new(qc.CreateVxNetsInput)
	opts.Count = 1
	opts.VxNetName = name
	opts.VxNetType =  type_

	rv, err := vxnetService.CreateVxNets(opts)

	if err != nil {
		return fmt.Errorf("Failed to create VxNet: %s", err)
	}

	if qc.IntValue(rc.RetCode) !+ 0 {
		return fmt.Errorf("Remote server refused to create VxNet: %s", *rv.Message)
	}

	vxnet := rv.VxNets[0]

	d.SetId(qc.StringValue(vxnet.VxNetID))

	return nil
}

func resourceQingCloudVxNetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := m.(*QingCloudClient)

	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	vxnetService, err := client.service.VxNet(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize VxNet service: %s", err)
	}

	opts := new(qc.ModifyVxNetAttributesInput)
	opts.VxNet = d.GetId()
	opts.VxNetName = name
	opts.Description = description

	rv, err := vxnetService.ModifyVxNetAttributes(opts)

	if err != nil {
		return fmt.Errorf("Failed to update VxNet: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to update VxNet: %s", *rv.Message)
	}

	log.Printf("[INFO] Updated VxNet: %s", qc.StringValue(d.GetId()))

	return nil
}

func resourceQingCloudVxNetDelete(d *schema.ResourceData, meta interface{}) error {
	client := m.(*QingCloudClient)

	zone := d.Get("zone").(string)

	vxnetService, err := client.service.VxNet(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize VxNet service: %s", err)
	}

	opts := new(qc.DeleteVxNetsInput)
	opts.VxNets = []*string{d.GetId()}

	rv, err := vxnetService.DeleteVxNet(opts)

	if err != nil {
		return fmt.Errorf("Failed to delete VxNet: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to delete VxNet: %s", &rv.Message)
	}

	return nil
}
