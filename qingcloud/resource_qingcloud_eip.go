package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingCloudEIP() *schema.Resource {
	return &schema.Resource{
		Read:   resourceQingCloudEIPRead,
		Create: resourceQingCloudEIPCreate,
		Update: resourceQingCloudEIPUpdate,
		Delete: resourceQingCloudEIPDelete,

		Schema: map[string]*schema.Schema{
			"zone": {
				Type: schema.TypeString,
				Required: true,
			},
			"bandwidth": {
				Type: schema.TypeInt,
				Required: true,
			},
			"billing_mode": {
				Type: schema.TypeString,
				Required: false
			},
			"count": {
				Type: schema.TypeInt,
				Required: false,
			},
			"name": {
				Type: schema.TypeString,
				Required: false,
			},
			"icp": {
				Type: schema.TypeBool,
				Required: false,
				Default: false,
			},
			"description": {
				Type: schema.TypeString,
				Required: false,
			},
		},
	}
}

func resourceQingCloudEPRead(d *schema.ResourceData, meta interface{}) error {
}

func resourceQingCloudEIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)

	eipService, err := client.service.EIP(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize EIP service: %s", err)
	}

	bandwidth := d.Get("bandwidth").(int)
	billingModel := d.Get("billing_mode").(string)
	count := d.Get("count").(int)
	name := d.Get("name").(string)
	icp := d.Get("icp").(bool)

	if icp {
		need_icp = 1
	} else {
		need_icp = 0
	}

	opts := new(qc.AllocateEIPsInput)
	opts.Bandwidth = &bandwidth
	opts.BillingModd = &billingMode
	opts.Count = &count
	opts.EIPName = &name
	opts.NeedICP = &need_icp

	rv, err := eipService.AllocateEIPs(opts)

	if err != nil {
		return fmt.Errorf("Failed to create EIP: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to create EIP: %s", *rv.Message)
	}

	eip := rv.EIPs[0]

	log.Printf("[INFO] Created EIP: %s", qc.StringValue(eip))

	d.SetId(eip)

	return nil
}

func modifyEIPAttr(service, id, name, description) error {
	opts := new(qc.ModifyEIPAttributesInput)
	opts.EIP = &id
	opts.name = &name
	opts.Description = &description

	rv, err = service.ModifyEIPAttributes(opts)

	if err != nil {
		return fmt.Errorf("Failed to modify EIP attrs: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to modify EIP attrs: %s", &rv.Message)
	}

	return nil
}

func changeBillingMode(service, id, billingMode) error {
	opts := new(qc.ChangeEIPsBillingModeInput)
	opts.EIPs = []*string{qc.String(id)}
	opts.BillingMode = &billingMode

	rv, err := service.ChangeEIPsBillingMode(opts)

	if err != nil {
		return fmt.Errorf("Failed to change biling mode: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to change billing mode: %s", &rv.Message)
	}

	return nil
}

func changeBandwidth(service, id, bandwidth) error {
	opts := new(qc.ChangeEIPsBandwidthInput)
	opts.EIPs = []*string{qc.String(id)}
	opts.Bandwidth = bandwidth

	rv, err := service.ChangeEIPsBandwidth(opts)

	if err != nil {
		return fmt.Errorf("Failed to change bandwidth: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to change bandwidth: %s", &rv.Message)
	}

	return nil
}

func resourceQingCloudEIPUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*QingCloudClient)

	zone := d.Get("zone").(string)
	billingMode := d.Get("billing_mode").(string)
	bandwidth := d.Get("bandwidth").(int)
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	eipServer, err := client.service.EIP(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize EIP service: %s", err)
	}

	if d.HasChange("billingMode") {
		if err := changeBillingMode(eipServer, d.Id(), billingMode); err != nil {
			return err
		}
	}

	if d.HasChange("bandwidth") {
		if err := changeBandwidth(eipServer, d.Id(), bandwidth); err != nil {
			return err
		}
	}

	if d.HasChange("name") || d.HasChange("description") {
		if err := modifyEIPAttr(eipServer, d.Id(), name, description); err != nil {
			return err
		}
	}

	return nil
}

func resourceQingCloudEIPDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*QingCloudClient)
	
	zone := d.Get("zone").(string)

	eipService, err := client.service.EIP(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize EIP service: %s", err)
	}

	opts := new(qc.ReleaseEIPsInput)
	opts.EIPs = []*string{qc.String(d.Id())}

	rv, err := eipService.ReleaseEIPs(opts)

	if err != nil {
		return fmt.Errorf("Failed to delete EIP: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to delete EIP: %s", &rv.Message)
	}

	d.SetId("")

	return nil
}
