// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-mangopay/mangopay"
)

// Ensure MangopayProvider satisfies various provider interfaces.
var _ provider.Provider = &MangopayProvider{}
var _ provider.ProviderWithFunctions = &MangopayProvider{}

// MangopayProvider defines the provider implementation.
type MangopayProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MangopayProviderModel describes the provider data model.
type MangopayProviderModel struct {
	ClientID            types.String `tfsdk:"client_id"`
	ClientSecret        types.String `tfsdk:"client_secret"`
	MangopayEnvironment types.String `tfsdk:"environment"`
}

func (p *MangopayProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mangopay"
	resp.Version = p.version
}

func (p *MangopayProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {

	tflog.Info(ctx, "zobizoba")
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID for the Mangopay API",
				Required:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client Secret for the Mangopay API",
				Required:            true,
				Sensitive:           true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "Mangopay environment to use",
				Optional:            true,
			},
		},
	}
}

func (p *MangopayProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "wesh")
	var config MangopayProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "wesh2")

	if config.ClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown Mangopay API Client ID",
			"The provider cannot create the Mangopay API client as there is an unknown configuration value for the Mangopay API Client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MANGOPAY_CLIENT_ID environment variable.",
		)
	}

	tflog.Info(ctx, "wesh3")

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown Mangopay API Client Secret",
			"The provider cannot create the Mangopay API client as there is an unknown configuration value for the Mangopay API Client secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MANGOPAY_CLIENT_SECRET environment variable.",
		)
	}

	tflog.Info(ctx, "wesh4")

	if config.MangopayEnvironment.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment"),
			"Unknown Mangopay API Mangopay Environment",
			"The provider cannot create the Mangopay API client as there is an unknown configuration value for the Mangopay API Mangopay Environment. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the MANGOPAY_ENVIRONMENT environment variable.",
		)
	}

	tflog.Info(ctx, "wesh4")

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	clientID := os.Getenv("MANOPAY_CLIENT_ID")
	clientSecret := os.Getenv("MANOPAY_CLIENT_SECRET")
	environment := os.Getenv("MANOPAY_ENVIRONMENT")

	if !config.ClientID.IsNull() {
		clientID = config.ClientID.ValueString()
	}
	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}
	if !config.MangopayEnvironment.IsNull() {
		environment = config.MangopayEnvironment.ValueString()
	}

	if clientID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Mangopay API Client ID",
			"The provider cannot create the Mangopay API client as there is a missing or empty value for the Mangopay API Client ID. "+
				"Set the client_id value in the configuration, or use the MANGOPAY_CLIENT_ID environment variable."+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing Mangopay API Client Secret",
			"The provider cannot create the Mangopay API client as there is a missing value for the Mangopay API client secret."+
				"Set the client secret value in the configuration, or use the MANGOPAY_CLIENT_SECRET environment variable."+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if environment == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment"),
			"Missing Mangopay API Mangopay Environment",
			"The provider cannot create the Mangopay API client as there is a missing value for the Mangopay API Mangopay Environment. "+
				"Set the environment value in the configuration, or use the MANGOPAY_ENVIRONMENT environment variable."+
				"If either is already set, ensure the value is not empty.",
		)
	}

	// Example client configuration for data sources and resources
	client, err := mangopay.New(&clientID, &clientSecret, &environment)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Mangopay API Client",
			"An unexpected error occurred when creating the Mangopay API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Mangopay Client Error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Provider configured with mangopay client")
	if client.Host == "" {
		tflog.Info(ctx, "GROS NAZE")
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MangopayProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//NewHookResource,
	}
}

func (p *MangopayProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewClientsDataSource,
		NewHookDataSource,
	}
}

func (p *MangopayProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MangopayProvider{
			version: version,
		}
	}
}
