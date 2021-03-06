package heroku

import (
	"context"
	"fmt"
	"log"

	"github.com/cyberdelia/heroku-go/v3"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHerokuDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceHerokuDomainCreate,
		Read:   resourceHerokuDomainRead,
		Delete: resourceHerokuDomainDelete,

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"app": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceHerokuDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	app := d.Get("app").(string)
	hostname := d.Get("hostname").(string)

	log.Printf("[DEBUG] Domain create configuration: %#v, %#v", app, hostname)

	do, err := client.DomainCreate(context.TODO(), app, heroku.DomainCreateOpts{Hostname: hostname})
	if err != nil {
		return err
	}

	d.SetId(do.ID)
	d.Set("hostname", do.Hostname)
	d.Set("cname", fmt.Sprintf("%s.herokuapp.com", app))

	log.Printf("[INFO] Domain ID: %s", d.Id())
	return nil
}

func resourceHerokuDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	log.Printf("[INFO] Deleting Domain: %s", d.Id())

	// Destroy the domain
	_, err := client.DomainDelete(context.TODO(), d.Get("app").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting domain: %s", err)
	}

	return nil
}

func resourceHerokuDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*heroku.Service)

	app := d.Get("app").(string)
	do, err := client.DomainInfo(context.TODO(), app, d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving domain: %s", err)
	}

	d.Set("hostname", do.Hostname)
	d.Set("cname", fmt.Sprintf("%s.herokuapp.com", app))

	return nil
}
