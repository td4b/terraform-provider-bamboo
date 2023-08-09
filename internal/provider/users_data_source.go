package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type UserModel struct {
	ID         types.Int64  `tfsdk:"id"`
	EmployeeID types.Int64  `tfsdk:"employeeid"`
	FirstName  types.String `tfsdk:"firstname"`
	LastName   types.String `tfsdk:"lastname"`
	Email      types.String `tfsdk:"email"`
	Status     types.String `tfsdk:"status"`
	LastLogin  types.String `tfsdk:"lastlogin"`
}

// coffeesDataSourceModel maps the data source schema data.
type UserDataSourceModel struct {
	Users []UserModel `tfsdk:"users"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &UsersDataSource{}
	_ datasource.DataSourceWithConfigure = &UsersDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewUsersDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

// coffeesDataSource is the data source implementation.
type UsersDataSource struct {
	client *BambooClient
}

// Metadata returns the data source type name.
func (d *UsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *UsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"employeeid": schema.Int64Attribute{
							Computed: true,
						},
						"firstname": schema.StringAttribute{
							Computed: true,
						},
						"lastname": schema.StringAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"lastlogin": schema.StringAttribute{
							Computed: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state UserDataSourceModel
	tflog.Info(ctx, "<diag> Trying to get users via UsersModel!")
	users, err := d.client.Getusers()
	if err != nil {
		tflog.Info(ctx, "<diag> got an error!"+err.Error())
	}
	// Map response body to model
	for _, user := range users {
		userState := UserModel{
			ID:         types.Int64Value(int64(user.ID)),
			EmployeeID: types.Int64Value(int64(user.EmployeeID)),
			FirstName:  types.StringValue(user.FirstName),
			LastName:   types.StringValue(user.LastName),
			Email:      types.StringValue(user.Email),
			Status:     types.StringValue(user.Status),
			LastLogin:  types.StringValue(user.LastLogin),
		}
		state.Users = append(state.Users, userState)
	}
	tflog.Debug(ctx, "<diag> Got State: "+fmt.Sprintf("%s", state.Users))
	//Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *UsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*BambooClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *bamboogo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
