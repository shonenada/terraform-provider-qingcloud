package qingcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingCloudKeyPair() *schema.Resource {
	return &schema.Resource{
		Read:   resourceQingCloudKeyPairRead,
		Create: resourceQingCloudKeyPairCreate,
		Update: resourceQingCloudKeyPairUpdate,
		Delete: resourceQingCloudKeyPairDelete,

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceQingCloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)

	keyPairService, err := client.service.KeyPair(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize key pair service: %s", err)
	}

	opts := new(qc.DescribeKeyPairsInput)

	opts.KeyPairs = []*string{qc.String(d.Id())}

	rv, err := keyPairService.DescribeKeyPairs(opts)

	if err != nil {
		return fmt.Errorf("Failed to read key pair: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to read key pair: %s", *rv.Message)
	}

	keyPair := rv.KeyPairSet[0]

	d.Set("key_id", keyPair.KeyPairID)
	d.Set("name", keyPair.KeyPairName)
	d.Set("description", keyPair.Description)
	d.Set("public_key", keyPair.PubKey)

	return nil
}

func resourceQingCloudKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	publicKey := d.Get("public_key").(string)

	keyPairService, err := client.service.KeyPair(zone)

	if err != nil {
		return fmt.Errorf("Failed to initialize key pair service: %s", err)
	}

	mode := "user"

	opts := new(qc.CreateKeyPairInput)
	opts.Mode = &mode
	opts.KeyPairName = &name
	opts.PublicKey = &publicKey

	rv, err := keyPairService.CreateKeyPair(opts)

	if err != nil {
		return fmt.Errorf("Failed to create key pair: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to create key pair: %s", *rv.Message)
	}

	log.Printf("[INFO] Created KeyPair: %s", qc.StringValue(rv.KeyPairID))

	d.SetId(qc.StringValue(rv.KeyPairID))

	return nil
}

func resourceQingCloudKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)
	keyId := d.Get("keyId").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	keyPairService, err := client.service.KeyPair(zone)

	if err != nil {
		return fmt.Errorf("Failed to create key pair: %s", err)
	}

	opts := new(qc.ModifyKeyPairAttributesInput)
	opts.KeyPair = &keyId
	opts.KeyPairName = &name
	opts.Description = &description

	rv, err := keyPairService.ModifyKeyPairAttributes(opts)

	if err != nil {
		return fmt.Errorf("Failed to update key pair: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to update key pair: %s", *rv.Message)
	}

	log.Printf("[INFO] Updated KeyPair: %s", qc.StringValue(&keyId))

	return nil
}

func resourceQingCloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*QingCloudClient)

	zone := d.Get("zone").(string)

	keyPairService, err := client.service.KeyPair(zone)

	opts := new(qc.DeleteKeyPairsInput)
	opts.KeyPairs = []*string{qc.String(d.Id())}

	rv, err := keyPairService.DeleteKeyPairs(opts)

	if err != nil {
		return fmt.Errorf("Failed to delete key pair: %s", err)
	}

	if qc.IntValue(rv.RetCode) != 0 {
		return fmt.Errorf("Remote server refused to delete key pair: %s", err)
	}

	return nil
}
