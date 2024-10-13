// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure CISProvider satisfies various provider interfaces.
var _ provider.Provider = &CISProvider{}
var _ provider.ProviderWithFunctions = &CISProvider{}

// CISProvider defines the provider implementation.
type CISProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CISProviderModel describes the provider data model.
type CISProviderModel struct {
	Auth0Endpoint     types.String `tfsdk:"auth0_endpoint"`
	Auth0ClientID     types.String `tfsdk:"auth0_client_id"`
	Auth0ClientSecret types.String `tfsdk:"auth0_client_secret"`
	PersonEndpoint    types.String `tfsdk:"person_endpoint"`
}

func (p *CISProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cis"
	resp.Version = p.version
}

func (p *CISProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"auth0_endpoint": schema.StringAttribute{
				Description:         "Auth0 endpoint",
				MarkdownDescription: "Auth0 endpoint",
				Optional:            true,
			},
			"auth0_client_id": schema.StringAttribute{
				Description:         "Auth0 client ID",
				MarkdownDescription: "Auth0 client ID",
				Required:            true,
				Sensitive:           true,
			},
			"auth0_client_secret": schema.StringAttribute{
				Description:         "Auth0 client secret",
				MarkdownDescription: "Auth0 client secret",
				Required:            true,
				Sensitive:           true,
			},
			"person_endpoint": schema.StringAttribute{
				Description:         "CIS person endpoint",
				MarkdownDescription: "CIS person endpoint",
				Optional:            true,
			},
		},
	}
}

func (p *CISProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CISProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CISProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *CISProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func (p *CISProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CISProvider{
			version: version,
		}
	}
}
