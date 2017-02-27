package main

type InstanceResource struct {
	ID            int           `json:"id"`
	PlanGUID      string        `json:"plan_guid"`
	ServiceGUID   string        `json:"service_guid"`
	DashboardURL  interface{}   `json:"dashboard_url"`
	State         string        `json:"state"`
	GUIDAtTenant  string        `json:"guid_at_tenant"`
	TenantID      string        `json:"tenant_id"`
	ProvisionedAt string        `json:"provisioned_at"`
	DeletedAt     string        `json:"deleted_at"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
	Credentials   []interface{} `json:"credentials"`
}

type Instances struct {
	TotalResults int                `json:"total_results"`
	TotalPages   int                `json:"total_pages"`
	CurrentPage  int                `json:"current_page"`
	PrevURL      interface{}        `json:"prev_url"`
	NextURL      interface{}        `json:"next_url"`
	Resources    []InstanceResource `json:"resources"`
}

type ServicePlan struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Metadata    interface{} `json:"metadata"`
	Free        bool        `json:"free"`
}

type CatalogService struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Bindable    bool          `json:"bindable"`
	Tags        []string      `json:"tags"`
	Plans       []ServicePlan `json:"plans"`
	Metadata    struct {
	} `json:"metadata"`
	Requires       []interface{} `json:"requires"`
	PlanUpdateable bool          `json:"plan_updateable"`
}

type Catalog struct {
	Services        []CatalogService `json:"services"`
	DashboardClient interface{}      `json:"dashboard_client"`
}
