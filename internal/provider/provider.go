// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
				Optional:            true,
				Sensitive:           true,
			},
			"auth0_client_secret": schema.StringAttribute{
				Description:         "Auth0 client secret",
				MarkdownDescription: "Auth0 client secret",
				Optional:            true,
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

	auth0_endpoint := os.Getenv("AUTH0_ENDPOINT")
	auth0_client_id := os.Getenv("AUTH0_CLIENT_ID")
	auth0_client_secret := os.Getenv("AUTH0_CLIENT_SECRET")
	person_endpoint := os.Getenv("PERSON_ENDPOINT")

	if data.Auth0Endpoint.ValueString() != "" {
		auth0_endpoint = data.Auth0Endpoint.ValueString()
	}
	if data.Auth0ClientID.ValueString() != "" {
		auth0_client_id = data.Auth0ClientID.ValueString()
	}
	if data.Auth0ClientSecret.ValueString() != "" {
		auth0_client_secret = data.Auth0ClientSecret.ValueString()
	}
	if data.PersonEndpoint.ValueString() != "" {
		person_endpoint = data.PersonEndpoint.ValueString()
	}

	if auth0_endpoint == "" {
		auth0_endpoint = "auth.mozilla.auth0.com"
	}
	if person_endpoint == "" {
		person_endpoint = "person.api.sso.mozilla.com"
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	if auth0_client_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth0_client_id"),
			"Missing Auth0 client ID",
			"Client ID not found in AUTH0_CLIENT_ID environment variable or provider configuration block auth0_client_id attribute.",
		)
	}
	if auth0_client_secret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth0_client_secret"),
			"Missing Auth0 client secret",
			"Client ID not found in AUTH0_CLIENT_SECRET environment variable or provider configuration block auth0_client_secret attribute.",
		)
	}

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
