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
var _ datasource.DataSource = &HookDataSource{}

func NewHookDataSource() datasource.DataSource {
	return &HookDataSource{}
}

// HookDataSource defines the data source implementation.
type HookDataSource struct {
	client *mangopay.Client
}

// hookModel describes the data source data model.
type hookModel struct {
	ID           types.String `tfsdk:"id"`
	URL          types.String `tfsdk:"url"`
	Status       types.String `tfsdk:"status"`
	Validity     types.String `tfsdk:"validity"`
	EventType    types.String `tfsdk:"event_type"`
	Tag          types.String `tfsdk:"tag"`
	CreationDate types.Int64  `tfsdk:"creation_date"`
}

// hookModel describes the data source data model.
type hookDataSourceModel struct {
	Hooks []hookModel `tfsdk:"hooks"`
}

func (d *HookDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Info(ctx, "hooks datasource metadata")
	resp.TypeName = req.ProviderTypeName + "_hooks"
}

func (d *HookDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Hook data source",

		Attributes: map[string]schema.Attribute{
			"hooks": schema.ListNestedAttribute{
				MarkdownDescription: "List of hooks",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier of the hook",
							Computed:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "The URL (http or https) to which the notification is sent.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Whether the hook is enabled or not.",
							Computed:            true,
						},
						"validity": schema.StringAttribute{
							MarkdownDescription: "Whether the hook is valid or not. Once the hook is set to INVALID it can no longer be modified.",
							Computed:            true,
						},
						"event_type": schema.StringAttribute{
							MarkdownDescription: "The type of the event",
							Computed:            true,
						},
						"tag": schema.StringAttribute{
							MarkdownDescription: "A custom tag for that hook",
							Computed:            true,
						},
						"creation_date": schema.Int64Attribute{
							MarkdownDescription: "The date when the hook was created",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *HookDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	tflog.Info(ctx, "hooks datasource configuration")
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

func (d *HookDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	hooks, err := d.client.GetAllHooks()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Mangopay Hook",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state := hookDataSourceModel{}

	for _, hook := range hooks {
		hookState := hookModel{
			ID:           types.StringValue(hook.Id),
			URL:          types.StringValue(hook.Url),
			Status:       types.StringValue(hook.Status),
			Validity:     types.StringValue(hook.Validity),
			EventType:    types.StringValue(hook.EventType),
			Tag:          types.StringValue(hook.Tag),
			CreationDate: types.Int64Value(int64(hook.CreationDate)),
		}

		state.Hooks = append(state.Hooks, hookState)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Info(ctx, "read the hooks data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
