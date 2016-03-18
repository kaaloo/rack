package client

import "fmt"

type Service struct {
	Name         string            `json:"name"`
	Status       string            `json:"status"`
	StatusReason string            `json:"status-reason"`
	Type         string            `json:"type"`
	Exports      map[string]string `json:"exports"`
	// DEPRECATED: should inject any data in Exports
	// we only set this on the outgoing response for old clients
	URL string `json:"url"`

	Outputs    map[string]string `json:"-"`
	Parameters map[string]string `json:"-"`
	Tags       map[string]string `json:"-"`
}

type Services []Service

func (c *Client) GetServices() (Services, error) {
	var services Services

	err := c.Get("/services", &services)

	if err != nil {
		return nil, err
	}

	return services, nil
}

func (c *Client) CreateService(kind string, options map[string]string) (*Service, error) {
	params := Params(options)
	params["type"] = kind
	var service Service

	err := c.Post("/services", params, &service)

	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (c *Client) GetService(name string) (*Service, error) {
	var service Service

	err := c.Get(fmt.Sprintf("/services/%s", name), &service)

	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (c *Client) DeleteService(name string) (*Service, error) {
	var service Service

	err := c.Delete(fmt.Sprintf("/services/%s", name), &service)

	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (c *Client) UpdateService(name string, options map[string]string) (*Service, error) {
	params := Params(options)
	var service Service

	err := c.Put(fmt.Sprintf("/services/%s", name), params, &service)

	if err != nil {
		return nil, err
	}

	return &service, nil
}
