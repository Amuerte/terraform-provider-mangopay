// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-mangopay/mangopay"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &HookResource{}
var _ resource.ResourceWithImportState = &HookResource{}

func NewHookResource() resource.Resource {
	return &HookResource{}
}

// HookResource defines the resource implementation.
type HookResource struct {
	client *mangopay.Client
}

// HookResourceModel describes the resource data model.
type hookResourceModel struct {
	ID           types.String `tfsdk:"id"`
	URL          types.String `tfsdk:"url"`
	Status       types.String `tfsdk:"status"`
	Validity     types.String `tfsdk:"validity"`
	EventType    types.String `tfsdk:"event_type"`
	Tag          types.String `tfsdk:"tag"`
	CreationDate types.Int64  `tfsdk:"creation_date"`
}

func (r *HookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hook"
}

func (r *HookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This is the Hook resource that enables to configure Mangopay webhooks.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the hook",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL (http or https) to which the notification is sent.",
				Required:            true,
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
				Required:            true,
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "A custom tag for that hook",
				Optional:            true,
			},
			"creation_date": schema.Int64Attribute{
				MarkdownDescription: "The date when the hook was created",
				Computed:            true,
			},
		},
	}
}

func (r *HookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*mangopay.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *HookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan hookResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Creating a hook resource")

	hook, err := r.client.CreateHook(mangopay.Hook{
		Url:       plan.URL.ValueString(),
		EventType: plan.EventType.ValueString(),
		Tag:       plan.Tag.ValueStringPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating hook",
			"Could not create hook, unexpected error: "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(string(hook.Id))
	plan.Status = types.StringValue(string(hook.Status))
	plan.Validity = types.StringValue(string(hook.Validity))
	plan.CreationDate = types.Int64Value(hook.CreationDate)
	plan.Tag = types.StringPointerValue(hook.Tag)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *HookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data hookResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data hookResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data hookResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *HookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}