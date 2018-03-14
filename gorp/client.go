package gorp

import (
	"fmt"
	"gopkg.in/resty.v1"
	"strconv"
)

//Client is ReportPortal REST API Client
type Client struct {
	project string
	http    *resty.Client
}

//NewClient creates new instance of Client
//host - server hostname
//project - name of the project
//uuid - User Token (see user profile page)
func NewClient(host, project, uuid string) *Client {
	http := resty.New().
		SetHostURL(host).
		SetAuthToken(uuid).
		OnAfterResponse(func(client *resty.Client, rs *resty.Response) error {
			if (rs.StatusCode() / 100) >= 4 {
				return fmt.Errorf("status code error: %d\n%s", rs.StatusCode(), rs.String())
			}
			return nil
		})
	return &Client{
		project: project,
		http:    http,
	}
}

//GetLaunches retrieves latest launches
func (c *Client) GetLaunches() (*LaunchPage, error) {
	var launches LaunchPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

//GetLaunchesByFilter retrieves launches by filter
func (c *Client) GetLaunchesByFilter(filter map[string]string) (*LaunchPage, error) {
	var launches LaunchPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		SetQueryParams(filter).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

//GetLaunchesByFilterString retrieves launches by filter as string
func (c *Client) GetLaunchesByFilterString(filter string) (*LaunchPage, error) {
	var launches LaunchPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		SetQueryString(filter).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

//GetLaunchesByFilterName retrieves launches by filter name
func (c *Client) GetLaunchesByFilterName(name string) (*LaunchPage, error) {
	filter, err := c.GetFiltersByName(name)
	if nil != err {
		return nil, err
	}

	if filter.Page.Size < 1 {
		return nil, fmt.Errorf("no filter %s found", name)
	}

	var launches LaunchPage
	params := ConvertToFilterParams(filter.Content[0])
	_, err = c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		SetQueryParams(params).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

//GetFiltersByName retrieves filter by its name
func (c *Client) GetFiltersByName(name string) (*FilterPage, error) {
	var filter FilterPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project, "name": name}).
		SetQueryParam("filter.eq.name", name).
		SetResult(&filter).
		Get("/api/v1/{project}/filter")
	return &filter, err
}

//ConvertToFilterParams converts RP internal filter representation to query string
func ConvertToFilterParams(filter *FilterResource) map[string]string {
	params := map[string]string{}
	for _, f := range filter.Entities {
		params[fmt.Sprintf("filter.%s.%s", f.Condition, f.Field)] = f.Value
	}

	if nil != filter.SelectionParams {
		if 0 != filter.SelectionParams.PageNumber {
			params["page.page"] = strconv.Itoa(filter.SelectionParams.PageNumber)
		}
		if nil != filter.SelectionParams.Orders {
			for _, order := range filter.SelectionParams.Orders {
				params["page.sort"] = fmt.Sprintf("%s,%s", order.SortingColumn, directionToStr(order.Asc))
			}
		}

	}

	return params
}

func directionToStr(asc bool) string {
	if asc {
		return "ASC"
	}
	return "DESC"

}