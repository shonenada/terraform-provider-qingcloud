package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingCloudRouter() *schema.Resource {
	return &schema.Resource{
		Read:   resourceQingCloudRouterRead,
		Create: resourceQingCloudRouterCreate,
		Update: resourceQingCloudRouterUpdate,
		Delete: resourceQingCloudRouterDelete,

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: false,
			},
			"type": {
				Type:     schema.TypeInt,
				Required: false,
			},
			"security_group": {
				Type:     schema.TypeString,
				Required: false,
			},
			"vpc_network": {
				Type:     schema.TypeString,
				Required: false,
			},
			"description": {
				Type:     schema.TypeString,
				Required: false,
			},
			"eip": {
				Type:     schema.TypeString,
				Required: false,
			},
		},
	}
}

func resourceQingCloudRouterRead(d *schema.ResourceData, m interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)

	routerService, err := cilent.service.Router(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize router service: %s", err)
	}

	opts := new(qc.DescribeRoutersInput)
	opts.Routers = []*string{qc.String(d.Id())}

	rv, err := routerService.DescribeRouters(opts)

	if err != nil {
		return fmt.Errorf("Failed to read router: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to read router: %s", *rv.Message)
	}

	router := rv.RouterSet[0]

	d.Set("name", router.RouterName)
	d.Set("type", router.RouterType)
	d.Set("description", router.Description)
	d.Set("security_group", router.SecurityGroupID)
	d.Set("eip", router.EIP)

	return nil
}

func resourceQingCloudRouterCreate(d *schema.ResourceData, m interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	type_ := d.Get("type").(int)
	secuityGroup := d.Get("secuity_group").(string)
	vpc_network := d.Get("vpc_network").(string)

	routerService, err := client.service.Router(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize router service: %s", err)
	}

	opts := new(qc.CreateRoutersInput)
	opts.RouterName = name
	opts.RouterType = type_
	opts.SecurityGroup = secuityGroup
	opts.VpcNetwork = vpc_network

	rv, err := routerService.CreateRouters(opts)

	if err != nil {
		return fmt.Errorf("Failed to create router service: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to create router: %s", *rv.Message)
	}

	router = qc.StringValue(rv.Routers[0])

	log.Printf("[INFO]: Created router: %s", router)

	d.SetId(router)

	return nil
}

func resourceQingCloudRouterUpdate(d *schema.ResourceData, m interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	secuityGroup := d.Get("secuity_group").(string)
	vx_net := d.Get("description").(string)
	description := d.Get("description").(string)
	eip := d.Get("eip").(string)

	routerService, err := client.service.Router(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize router service: %s", err)
	}

	opts := new(qc.ModifyRouterAttributesInput)
	opts.Router = d.GetId()
	opts.Description = description
	opts.EIP = eip
	opts.RouterName = name
	opts.SecurityGroup = secuityGroup

	rv, err := routerService.ModifyRouterAttributes(opts)

	if err != nil {
		return fmt.Errof("Failed to update router: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to update router: %s", *rv.Message)
	}

	log.Printf("[INFO] Updated router: %s", qc.StringValue(d.GetId()))

	return nil
}

func resourceQingCloudRouterDelete(d *schema.ResourceData, m interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)

	routerService, err := client.service.Router(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize router srevice: %s", err)
	}

	opts := new(qc.DeleteRoutersInput)
	opts.Routers = []*string{qc.String(d.Id())}

	rv, err := routerService.DeleteRouters(opts)

	if err != nil {
		return fmt.Errorf("Failed to delete router: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to delete router: %s", *rv.Message)
	}

	log.Printf("[INFO] Delete router: %s", qc.StringValue(d.Id()))

	d.SetId("")

	return nil
}
