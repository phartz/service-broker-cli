package main

import "time"

type ResourceMetadata struct {
	InstanceGUIDAtTenant string      `json:"instance_guid_at_tenant"`
	UserParams           interface{} `json:"user_params"`
	PlanGUID             string      `json:"plan_guid"`
	TenantID             string      `json:"tenant_id"`
}

type InstanceResource struct {
	ID            int                `json:"id"`
	PlanGUID      string             `json:"plan_guid"`
	ServiceGUID   string             `json:"service_guid"`
	Metadata      []ResourceMetadata `json:"metadata"`
	DashboardURL  interface{}        `json:"dashboard_url"`
	State         string             `json:"state"`
	GUIDAtTenant  string             `json:"guid_at_tenant"`
	TenantID      string             `json:"tenant_id"`
	ProvisionedAt time.Time          `json:"provisioned_at"`
	DeletedAt     time.Time          `json:"deleted_at"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	Credentials   []interface{}      `json:"credentials"`
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
