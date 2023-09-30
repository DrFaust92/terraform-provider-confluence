package confluence

import (
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceContent() *schema.Resource {
	return &schema.Resource{
		Create: resourceContentCreate,
		Read:   resourceContentRead,
		Update: resourceContentUpdate,
		Delete: resourceContentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "page",
			},
			"space": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("CONFLUENCE_SPACE", nil),
			},
			"body": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: resourceContentDiffBody,
			},
			"title": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				DiffSuppressFunc: resourceContentDiffParent,
			},
			"labels": {
				Type:        schema.TypeSet,
                Optional:    true,
                Description: "An optional list of labels",
                Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceContentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	contentRequest := contentFromResourceData(d)
	contentResponse, err := client.CreateContent(contentRequest)
	if err != nil {
		return err
	}
	d.SetId(contentResponse.Id)
	return resourceContentRead(d, m)
}

func resourceContentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	contentResponse, err := client.GetContent(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}
	return updateResourceDataFromContent(d, contentResponse, client)
}

func resourceContentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	contentRequest := contentFromResourceData(d)
	_, err := client.UpdateContent(contentRequest)
	if err != nil {
		d.SetId("")
		return err
	}
	return resourceContentRead(d, m)
}

func resourceContentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	err := client.DeleteContent(d.Id())
	if err != nil {
		return err
	}
	// d.SetId("") is automatically called assuming delete returns no errors
	return nil
}

func contentFromResourceData(d *schema.ResourceData) *Content {
	result := &Content{
		Id:   d.Id(),
		Type: d.Get("type").(string),
		Space: &Space{
			Key: d.Get("space").(string),
		},
		Body: &Body{
			Storage: &Storage{
				Value:          d.Get("body").(string),
				Representation: "storage",
			},
		},
		Title: d.Get("title").(string),

	}
	version := d.Get("version").(int) // Get returns 0 if unset
	if version > 0 {
		result.Version = &Version{Number: version}
	}
	parent := d.Get("parent").(string)
	if parent != "" {
		result.Ancestors = []*Content{
			{
				Id:   parent,
				Type: "page",
			},
		}
	}
	result.MetaData = &MetaData {
	    Labels: &Labels{
            LabelResults: getLabelsFromResourceData(d),
        },
    }
	return result
}

func getLabelsFromResourceData(d *schema.ResourceData) []*Label {
    tfLabels := d.Get("labels").(*schema.Set).List()
    Labels := []*Label{}
    for _, label := range tfLabels {
        l := new(Label)
        l.Name = label.(string)
        l.Prefix = "global"
        Labels = append(Labels, l)
   }
   return Labels
}

func updateResourceDataFromContent(d *schema.ResourceData, content *Content, client *Client) error {
	d.SetId(content.Id)
	m := map[string]interface{}{
		"type":    content.Type,
		"space":   content.Space.Key,
		"body":    content.Body.Storage.Value,
		"title":   content.Title,
		"version": content.Version.Number,
		"url":     client.URL(content.Links.Context + content.Links.WebUI),
	}
	if len(content.Ancestors) > 1 {
		m["parent"] = content.Ancestors[len(content.Ancestors)-1].Id
	}

    var labels []string
    for _, x := range content.MetaData.Labels.LabelResults {
        labels = append(labels, x.Name)
    }
    m["labels"] = labels

	for k, v := range m {
		err := d.Set(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Body was showing as requiring changes when there weren't any. It appears there
// are some whitespace differences between the old and new. This supresses the
// false differences by comparing the trimmed strings
func resourceContentDiffBody(k, old, new string, d *schema.ResourceData) bool {
	old = strings.TrimSpace(old)
	new = strings.TrimSpace(new)

	// Ignore macro-id which changes on page update
	reg := regexp.MustCompile(` *ac:macro-id="[\w\d\-]*" *`)
	old = reg.ReplaceAllString(old, "${1}")
	new = reg.ReplaceAllString(new, "${1}")
	return old == new

}

// If the parent was not set, running diff will show the actual value for old and
// the empty value for new. This case is supressed.
func resourceContentDiffParent(k, old, new string, d *schema.ResourceData) bool {
	return (new == "") || (old == new)
}
