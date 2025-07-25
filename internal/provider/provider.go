package provider

import (
	"context"
	"os"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &hashicupsProvider{}
)

// New is a helper function to simplify provider server and hashicups implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hashicupsProvider{
			version: version,
		}
	}
}

// hashicupsProviderModel maps provider schema data to a Go type.
type hashicupsProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// hashicupsProvider is the provider implementation.
type hashicupsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// hashicups.
	version string
}

// Metadata returns the provider type name.
func (p *hashicupsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hashicups"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
// Schema defines the provider-level schema for configuration data.
func (p *hashicupsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a hashicups API client for data sources and resources.
func (p *hashicupsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config hashicupsProviderModel
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
			"Unknown hashicups API Host",
			"The provider cannot create the hashicups API client as there is an unknown configuration value for the hashicups API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the hashicups_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown hashicups API Username",
			"The provider cannot create the hashicups API client as there is an unknown configuration value for the hashicups API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the hashicups_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown hashicups API Password",
			"The provider cannot create the hashicups API client as there is an unknown configuration value for the hashicups API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the hashicups_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("hashicups_HOST")
	username := os.Getenv("hashicups_USERNAME")
	password := os.Getenv("hashicups_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing hashicups API Host",
			"The provider cannot create the hashicups API client as there is a missing or empty value for the hashicups API host. "+
				"Set the host value in the configuration or use the hashicups_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing hashicups API Username",
			"The provider cannot create the hashicups API client as there is a missing or empty value for the hashicups API username. "+
				"Set the username value in the configuration or use the hashicups_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing hashicups API Password",
			"The provider cannot create the hashicups API client as there is a missing or empty value for the hashicups API password. "+
				"Set the password value in the configuration or use the hashicups_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new hashicups client using the configuration values
	client, err := hashicups.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create hashicups API Client",
			"An unexpected error occurred when creating the hashicups API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"hashicups Client Error: "+err.Error(),
		)
		return
	}

	// Make the hashicups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *hashicupsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCoffeesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *hashicupsProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
