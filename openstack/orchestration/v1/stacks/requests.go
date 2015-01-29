package stacks

import (
	"errors"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToStackCreateMap() (map[string]interface{}, error)
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	DisableRollback *bool
	Environment     string
	Files           map[string]interface{}
	Name            string
	Parameters      map[string]string
	Template        string
	TemplateURL     string
	Timeout         int
}

// ToStackCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToStackCreateMap() (map[string]interface{}, error) {
	s := make(map[string]interface{})

	if opts.Name == "" {
		return s, errors.New("Required field 'Name' not provided.")
	}
	s["stack_name"] = opts.Name

	if opts.Template != "" {
		s["template"] = opts.Template
	} else if opts.TemplateURL != "" {
		s["template_url"] = opts.TemplateURL
	} else {
		return s, errors.New("Either Template or TemplateURL must be provided.")
	}

	if opts.DisableRollback != nil {
		s["disable_rollback"] = &opts.DisableRollback
	}

	if opts.Environment != "" {
		s["environment"] = opts.Environment
	}
	if opts.Files != nil {
		s["files"] = opts.Files
	}
	if opts.Parameters != nil {
		s["parameters"] = opts.Parameters
	}

	if opts.Timeout != 0 {
		s["timeout_mins"] = opts.Timeout
	}

	return s, nil
}

// Create accepts a CreateOpts struct and creates a new stack using the values
// provided.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToStackCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	// Send request to API
	_, res.Err = perigee.Request("POST", createURL(c), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{201},
	})
	return res
}

// AdoptOptsBuilder is the interface options structs have to satisfy in order
// to be used in the Adopt function in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type AdoptOptsBuilder interface {
	ToStackAdoptMap() (map[string]interface{}, error)
}

// AdoptOpts is the common options struct used in this package's Adopt
// operation.
type AdoptOpts struct {
	AdoptStackData  string
	DisableRollback *bool
	Environment     string
	Files           map[string]interface{}
	Name            string
	Parameters      map[string]string
	Template        string
	TemplateURL     string
	Timeout         int
}

// ToStackAdoptMap casts a CreateOpts struct to a map.
func (opts AdoptOpts) ToStackAdoptMap() (map[string]interface{}, error) {
	s := make(map[string]interface{})

	if opts.Name == "" {
		return s, errors.New("Required field 'Name' not provided.")
	}
	s["stack_name"] = opts.Name

	if opts.Template != "" {
		s["template"] = opts.Template
	} else if opts.TemplateURL != "" {
		s["template_url"] = opts.TemplateURL
	} else {
		return s, errors.New("Either Template or TemplateURL must be provided.")
	}

	if opts.AdoptStackData == "" {
		return s, errors.New("Required field 'AdoptStackData' not provided.")
	}
	s["adopt_stack_data"] = opts.AdoptStackData

	if opts.DisableRollback != nil {
		s["disable_rollback"] = &opts.DisableRollback
	}

	if opts.Environment != "" {
		s["environment"] = opts.Environment
	}
	if opts.Files != nil {
		s["files"] = opts.Files
	}
	if opts.Parameters != nil {
		s["parameters"] = opts.Parameters
	}

	if opts.Timeout != 0 {
		s["timeout_mins"] = opts.Timeout
	}

	return map[string]interface{}{"stack": s}, nil
}

// Adopt accepts an AdoptOpts struct and creates a new stack using the resources
// from another stack.
func Adopt(c *gophercloud.ServiceClient, opts AdoptOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToStackAdoptMap()
	if err != nil {
		res.Err = err
		return res
	}

	// Send request to API
	_, res.Err = perigee.Request("POST", adoptURL(c), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{201},
	})
	return res
}

// SortDir is a type for specifying in which direction to sort a list of stacks.
type SortDir string

// SortKey is a type for specifying by which key to sort a list of stacks.
type SortKey string

var (
	// SortAsc is used to sort a list of stacks in ascending order.
	SortAsc SortDir = "asc"
	// SortDesc is used to sort a list of stacks in descending order.
	SortDesc SortDir = "desc"
	// SortName is used to sort a list of stacks by name.
	SortName SortKey = "name"
	// SortStatus is used to sort a list of stacks by status.
	SortStatus SortKey = "status"
	// SortCreatedAt is used to sort a list of stacks by date created.
	SortCreatedAt SortKey = "created_at"
	// SortUpdatedAt is used to sort a list of stacks by date updated.
	SortUpdatedAt SortKey = "updated_at"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToStackListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the network attributes you want to see returned. SortKey allows you to sort
// by a particular network attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Status  string  `q:"status"`
	Name    string  `q:"name"`
	Marker  string  `q:"marker"`
	Limit   int     `q:"limit"`
	SortKey SortKey `q:"sort_keys"`
	SortDir SortDir `q:"sort_dir"`
}

// ToStackListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToStackListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List returns a Pager which allows you to iterate over a collection of
// stacks. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToStackListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	createPage := func(r pagination.PageResult) pagination.Page {
		return StackPage{pagination.SinglePageBase(r)}
	}
	return pagination.NewPager(c, url, createPage)
}

// Get retreives a stack based on the stack name and stack ID.
func Get(c *gophercloud.ServiceClient, stackName, stackID string) GetResult {
	var res GetResult

	// Send request to API
	_, res.Err = perigee.Request("GET", getURL(c, stackName, stackID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the Update operation in this package.
type UpdateOptsBuilder interface {
	ToStackUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	Environment string
	Files       map[string]interface{}
	Parameters  map[string]string
	Template    string
	TemplateURL string
	Timeout     int
}

// ToStackUpdateMap casts a CreateOpts struct to a map.
func (opts UpdateOpts) ToStackUpdateMap() (map[string]interface{}, error) {
	s := make(map[string]interface{})

	if opts.Template != "" {
		s["template"] = opts.Template
	} else if opts.TemplateURL != "" {
		s["template_url"] = opts.TemplateURL
	} else {
		return s, errors.New("Either Template or TemplateURL must be provided.")
	}

	if opts.Environment != "" {
		s["environment"] = opts.Environment
	}

	if opts.Files != nil {
		s["files"] = opts.Files
	}

	if opts.Parameters != nil {
		s["parameters"] = opts.Parameters
	}

	if opts.Timeout != 0 {
		s["timeout_mins"] = opts.Timeout
	}

	return s, nil
}

// Update accepts an UpdateOpts struct and updates an existing stack using the values
// provided.
func Update(c *gophercloud.ServiceClient, stackName, stackID string, opts UpdateOptsBuilder) UpdateResult {
	var res UpdateResult

	reqBody, err := opts.ToStackUpdateMap()
	if err != nil {
		res.Err = err
		return res
	}

	// Send request to API
	_, res.Err = perigee.Request("PUT", updateURL(c, stackName, stackID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		OkCodes:     []int{202},
	})
	return res
}

// Delete deletes a stack based on the stack name and stack ID.
func Delete(c *gophercloud.ServiceClient, stackName, stackID string) DeleteResult {
	var res DeleteResult

	// Send request to API
	_, res.Err = perigee.Request("DELETE", deleteURL(c, stackName, stackID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		OkCodes:     []int{204},
	})
	return res
}

// PreviewOptsBuilder is the interface options structs have to satisfy in order
// to be used in the Preview operation in this package.
type PreviewOptsBuilder interface {
	ToStackPreviewMap() (map[string]interface{}, error)
}

// PreviewOpts contains the common options struct used in this package's Preview
// operation.
type PreviewOpts struct {
	DisableRollback *bool
	Environment     string
	Files           map[string]interface{}
	Name            string
	Parameters      map[string]string
	Template        string
	TemplateURL     string
	Timeout         int
}

// ToStackPreviewMap casts a PreviewOpts struct to a map.
func (opts PreviewOpts) ToStackPreviewMap() (map[string]interface{}, error) {
	s := make(map[string]interface{})

	if opts.Name == "" {
		return s, errors.New("Required field 'Name' not provided.")
	}
	s["stack_name"] = opts.Name

	if opts.Template != "" {
		s["template"] = opts.Template
	} else if opts.TemplateURL != "" {
		s["template_url"] = opts.TemplateURL
	} else {
		return s, errors.New("Either Template or TemplateURL must be provided.")
	}

	if opts.DisableRollback != nil {
		s["disable_rollback"] = &opts.DisableRollback
	}

	if opts.Environment != "" {
		s["environment"] = opts.Environment
	}
	if opts.Files != nil {
		s["files"] = opts.Files
	}
	if opts.Parameters != nil {
		s["parameters"] = opts.Parameters
	}

	if opts.Timeout != 0 {
		s["timeout_mins"] = opts.Timeout
	}

	return s, nil
}

// Preview accepts a PreviewOptsBuilder interface and creates a preview of a stack using the values
// provided.
func Preview(c *gophercloud.ServiceClient, opts PreviewOptsBuilder) PreviewResult {
	var res PreviewResult

	reqBody, err := opts.ToStackPreviewMap()
	if err != nil {
		res.Err = err
		return res
	}

	// Send request to API
	_, res.Err = perigee.Request("POST", previewURL(c), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		ReqBody:     &reqBody,
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}

// Abandon deletes the stack with the provided stackName and stackID, but leaves its
// resources intact, and returns data describing the stack and its resources.
func Abandon(c *gophercloud.ServiceClient, stackName, stackID string) AbandonResult {
	var res AbandonResult

	// Send request to API
	_, res.Err = perigee.Request("POST", abandonURL(c, stackName, stackID), perigee.Options{
		MoreHeaders: c.AuthenticatedHeaders(),
		Results:     &res.Body,
		OkCodes:     []int{200},
	})
	return res
}
