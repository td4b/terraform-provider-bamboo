package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &bambooProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &bambooProvider{
			version: version,
		}
	}
}

// hashicupsProvider is the provider implementation.
type bambooProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// hashicupsProviderModel maps provider schema data to a Go type.
type bambooProviderModel struct {
	Host    types.String `tfsdk:"host"`
	Company types.String `tfsdk:"company"`
	Apikey  types.String `tfsdk:"apikey"`
}

// Metadata returns the provider type name.
func (p *bambooProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bamboo"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *bambooProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"company": schema.StringAttribute{
				Optional: true,
			},
			"apikey": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *bambooProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	tflog.Info(ctx, "<diag> Creating Bamboo client")
	var config bambooProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Bamboo API Host",
			"The provider cannot create the Bamboo API client as there is an unknown configuration value for the Bamboo API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAMBOO_HOST environment variable.",
		)
	}

	if config.Company.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("company"),
			"Unknown Bamboo API Company",
			"The provider cannot create the Bamboo API client as there is an unknown configuration value for the Bamboo API company. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAMBOO_COMPANY environment variable.",
		)
	}

	if config.Apikey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown Bamboo API Key",
			"The provider cannot create the Bamboo API client as there is an unknown configuration value for the Bamboo API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAMBOO_APIKEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "<diag> Checked if params are set")

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("BAMBOO_HOST")
	company := os.Getenv("BAMBOO_COMPANY")
	apikey := os.Getenv("BAMBOO_APIKEY")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Company.IsNull() {
		company = config.Company.ValueString()
	}

	if !config.Apikey.IsNull() {
		apikey = config.Apikey.ValueString()
	}

	tflog.Info(ctx, "<diag> checked if valuestrings are set")

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Bamboo API Host",
			"The provider cannot create the Bamboo API client as there is a missing or empty value for the Bamboo API host. "+
				"Set the host value in the configuration or use the BAMBOO_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if company == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("company"),
			"Missing Bamboo API Company",
			"The provider cannot create the Bamboo API client as there is a missing or empty value for the Bamboo API company. "+
				"Set the company value in the configuration or use the BAMBOO_COMPANY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing Bamboo API Key",
			"The provider cannot create the Bamboo API client as there is a missing or empty value for the Bamboo API key. "+
				"Set the apikey value in the configuration or use the BAMBOO_APIKEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	tflog.Info(ctx, "<diag> checked if stringempty")

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new HashiCups client using the configuration values
	tflog.Info(ctx, "<diag> trying to create bamboo client")
	client, err := NewClient(&host, &company, &apikey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Bamboo API Client",
			"An unexpected error occurred when creating the Bamboo API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Bamboo Client Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "<diag> bamboo client was created!")
}

// DataSources defines the data sources implemented in the provider.
func (p *bambooProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUsersDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *bambooProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
