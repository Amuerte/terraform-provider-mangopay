// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-mangopay/mangopay"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ClientsDataSource{}

func NewClientsDataSource() datasource.DataSource {
	return &ClientsDataSource{}
}

// ClientsDataSource defines the data source implementation.
type ClientsDataSource struct {
	client *mangopay.Client
}

// clientsDataSourceModel describes the data source data model.
type clientsDataSourceModel struct {
	PlatformType            types.String                `tfsdk:"platform_type"`
	ClientId                types.String                `tfsdk:"client_id"`
	Name                    types.String                `tfsdk:"name"`
	RegisteredName          types.String                `tfsdk:"registered_name"`
	TechEmails              types.List                  `tfsdk:"tech_emails"`
	AdminEmails             types.List                  `tfsdk:"admin_emails"`
	BillingEmails           types.List                  `tfsdk:"billing_emails"`
	FraudEmails             types.List                  `tfsdk:"fraud_emails"`
	HeadquartersAddress     addressModel                `tfsdk:"headquarters_address"`
	HeadquartersPhoneNumber types.String                `tfsdk:"headquarters_phone_number"`
	TaxNumber               types.String                `tfsdk:"tax_number"`
	PlatformCategorization  platformCategorizationModel `tfsdk:"platform_categorization"`
	PlatformURL             types.String                `tfsdk:"platform_url"`
	PlatformDescription     types.String                `tfsdk:"platform_description"`
	CompanyReference        types.String                `tfsdk:"company_reference"`
	PrimaryThemeColour      types.String                `tfsdk:"primary_theme_colour"`
	PrimaryButtonColour     types.String                `tfsdk:"primary_button_colour"`
	Logo                    types.String                `tfsdk:"logo"`
	CompanyNumber           types.String                `tfsdk:"company_number"`
	MCC                     types.String                `tfsdk:"mcc"`
}

type addressModel struct {
	AddressLine1 types.String `tfsdk:"address_line1"`
	AddressLine2 types.String `tfsdk:"address_line2"`
	City         types.String `tfsdk:"city"`
	Region       types.String `tfsdk:"region"`
	PostalCode   types.String `tfsdk:"postal_code"`
	Country      types.String `tfsdk:"country"`
}

type platformCategorizationModel struct {
	BusinessType types.String `tfsdk:"business_type"`
	Sector       types.String `tfsdk:"sector"`
}

func (d *ClientsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Info(ctx, "clients datasource metadata")
	resp.TypeName = req.ProviderTypeName + "_clients"
}

func (d *ClientsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Clients data source",

		Attributes: map[string]schema.Attribute{
			"platform_type": schema.StringAttribute{
				MarkdownDescription: "The type of the platform",
				Computed:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier associated with the API key, giving access to either the Sandbox or Production environment.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The trading name of the company operating the platform",
				Computed:            true,
			},
			"registered_name": schema.StringAttribute{
				MarkdownDescription: "The registered legal name of the company operating the platform",
				Computed:            true,
			},

			"tech_emails": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of email addresses to contact the platform for technical matters",
				Computed:            true,
			},
			"admin_emails": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of email addresses to contact the platform for administrative or commercial matters",
				Computed:            true,
			},
			"billing_emails": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of email addresses to contact the platform for billing matters",
				Computed:            true,
			},
			"fraud_emails": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of email addresses to contact the platform for fraud and compliance matters",
				Computed:            true,
			},
			"headquarters_address": schema.SingleNestedAttribute{
				MarkdownDescription: "The address of the platform operator’s headquarters. This parameter must be provided for the platform’s payouts to be processed",
				Attributes: map[string]schema.Attribute{
					"address_line1": schema.StringAttribute{
						MarkdownDescription: "The first line of the address",
						Computed:            true,
					},
					"address_line2": schema.StringAttribute{
						MarkdownDescription: "The second line of the address",
						Computed:            true,
					},
					"city": schema.StringAttribute{
						MarkdownDescription: "The city of the address.",
						Computed:            true,
					},
					"region": schema.StringAttribute{
						MarkdownDescription: "The region of the address.",
						Computed:            true,
					},
					"postal_code": schema.StringAttribute{
						MarkdownDescription: "The postal code of the address.",
						Computed:            true,
					},
					"country": schema.StringAttribute{
						MarkdownDescription: "The country of the address.",
						Computed:            true,
					},
				},
				Computed: true,
			},
			"headquarters_phone_number": schema.StringAttribute{
				MarkdownDescription: "The phone number of the platform operator’s headquarters.",
				Computed:            true,
			},
			"tax_number": schema.StringAttribute{
				MarkdownDescription: "The tax (or VAT) number for the company operating the platform.",
				Computed:            true,
			},

			"platform_categorization": schema.SingleNestedAttribute{
				MarkdownDescription: "The categorization of the platform in terms of business and sector of activity",
				Attributes: map[string]schema.Attribute{
					"business_type": schema.StringAttribute{
						MarkdownDescription: "The business type of the platform",
						Computed:            true,
					},
					"sector": schema.StringAttribute{
						MarkdownDescription: "The sector of activity of the platform",
						Computed:            true,
					},
				},
				Computed: true,
			},

			"platform_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the platform’s website",
				Computed:            true,
			},
			"platform_description": schema.StringAttribute{
				MarkdownDescription: "The description of what the platform does",
				Computed:            true,
			},
			"company_reference": schema.StringAttribute{
				MarkdownDescription: "The unique reference for the platform, which should be used when contacting Mangopay",
				Computed:            true,
			},
			"primary_theme_colour": schema.StringAttribute{
				MarkdownDescription: "The primary color of your branding, which is displayed on some payment pages (e.g., mandate confirmation)",
				Computed:            true,
			},
			"primary_button_colour": schema.StringAttribute{
				MarkdownDescription: "The primary color of your branding, which is displayed in call-to-action buttons on some payment pages (e.g., mandate confirmation).",
				Computed:            true,
			},
			"logo": schema.StringAttribute{
				MarkdownDescription: "The URL of the platform’s logo. Logos may be added by using the Upload a Client Logo endpoint.",
				Computed:            true,
			},
			"company_number": schema.StringAttribute{
				MarkdownDescription: "The registration number of the company operating the platform, assigned by the relevant national authority.",
				Computed:            true,
			},
			"mcc": schema.StringAttribute{
				MarkdownDescription: "4-digit merchant category code. The MCC is used to classify a business by the types of goods or services it provides.",
				Computed:            true,
			},
		},
	}
}

func (d *ClientsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	tflog.Info(ctx, "clients datasource configuration")
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*mangopay.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *mangopay.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ClientsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	platformClient, err := d.client.GetPlatformClient()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Mangopay Clients",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state := clientsDataSourceModel{
		PlatformType:            types.StringValue(platformClient.PlatformType),
		ClientId:                types.StringValue(platformClient.ClientId),
		Name:                    types.StringValue(platformClient.Name),
		RegisteredName:          types.StringValue(platformClient.RegisteredName),
		HeadquartersPhoneNumber: types.StringValue(platformClient.HeadquartersPhoneNumber),
		TaxNumber:               types.StringValue(platformClient.TaxNumber),
		PlatformURL:             types.StringValue(platformClient.PlatformURL),
		PlatformDescription:     types.StringValue(platformClient.PlatformDescription),
		CompanyReference:        types.StringValue(platformClient.CompanyReference),
		PrimaryThemeColour:      types.StringValue(platformClient.PrimaryThemeColour),
		PrimaryButtonColour:     types.StringValue(platformClient.PrimaryButtonColour),
		Logo:                    types.StringValue(platformClient.Logo),
		CompanyNumber:           types.StringValue(platformClient.CompanyNumber),
		MCC:                     types.StringValue(platformClient.MCC),
	}

	techEmails, _ := types.ListValueFrom(ctx, types.StringType, platformClient.TechEmails)
	adminEmails, _ := types.ListValueFrom(ctx, types.StringType, platformClient.AdminEmails)
	billingEmails, _ := types.ListValueFrom(ctx, types.StringType, platformClient.BillingEmails)
	fraudEmails, _ := types.ListValueFrom(ctx, types.StringType, platformClient.FraudEmails)

	state.TechEmails = techEmails
	state.AdminEmails = adminEmails
	state.BillingEmails = billingEmails
	state.FraudEmails = fraudEmails

	headquartersAddressState := addressModel{
		AddressLine1: types.StringValue(platformClient.HeadquartersAddress.AddressLine1),
		AddressLine2: types.StringValue(platformClient.HeadquartersAddress.AddressLine2),
		City:         types.StringValue(platformClient.HeadquartersAddress.City),
		Region:       types.StringValue(platformClient.HeadquartersAddress.Region),
		PostalCode:   types.StringValue(platformClient.HeadquartersAddress.PostalCode),
		Country:      types.StringValue(platformClient.HeadquartersAddress.Country),
	}

	state.HeadquartersAddress = headquartersAddressState
	platformCategorizationState := platformCategorizationModel{
		BusinessType: types.StringValue(platformClient.PlatformCategorization.BusinessType),
		Sector:       types.StringValue(platformClient.PlatformCategorization.Sector),
	}

	state.PlatformCategorization = platformCategorizationState

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Info(ctx, "read the clients data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
