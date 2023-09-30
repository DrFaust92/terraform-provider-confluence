---
layout: "confluence"
page_title: "Confluence: confluence_content"
sidebar_current: "docs-confluence-resource-content"
description: |-
  Provides content in Confluence
---

# confluence_content

Provides a piece of content on your Confluence site.

## Example Usage

```hcl
resource confluence_content "default" {
  space  = "my-space"
  title  = "Example Page"
  body   = "<p>This page was built with Terraform<p>"
  parent = "123456"
  labels = ["label1", "label2"]
}
```

## Argument Reference

The following arguments are supported:

* `body` - (Required) The actual content of the page in [Confluence Storage
  Format](https://confluence.atlassian.com/doc/confluence-storage-format-790796544.html)

* `parent` - (Optional) The content id of the parent page. If unset, terraform
  will leave it as-is, which could be the default location (under the space's
  home page) or any other location the page has been manually moved.

* `space` - (Required) The space key to create the content under. This can also
  be set via the `CONFLUENCE_SPACE` environment variable.

* `title` - (Required) The title of the page.

* `type` - (Optional) The content type (either "page" or "blogpost"). Default is page.

* `labels` - (Optional) A list of labels to apply to the content.

## Attributes Reference

This resource exports the following attributes:

* `body` - The actual content of the page in Confluence Storage Format.

* `parent` - The content id of the parent page.

* `space` - The space key the content is under.

* `title` - The title of the page.

* `type` - The content type (either "page" or "blogpost").

* `url` - The web link to the content.

* `version` - The version number of the content.

## Import

Content can be imported using the content id.

```
$ terraform import confluence_content.default {{id}}
```
