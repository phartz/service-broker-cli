package sbcli

type InstanceResource struct {
	ID             int          `json:"id"`
	PlanGUID       string       `json:"plan_guid"`
	ServiceGUID    string       `json:"service_guid"`
	Metadata       Metadata     `json:"metadata"`
	DashboardURL   interface{}  `json:"dashboard_url"`
	DeploymentName interface{}  `json:"deployment_name"`
	State          string       `json:"state"`
	GUIDAtTenant   string       `json:"guid_at_tenant"`
	TenantID       string       `json:"tenant_id"`
	ProvisionedAt  string       `json:"provisioned_at"`
	DeletedAt      string       `json:"deleted_at"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
	Credentials    []Credential `json:"credentials"`
	VMDetails      []VMDetails  `json:"vm_details"`
}

type VMDetails struct {
	VMIdentifier   string      `json:"vm_identifier"`
	CPU            int         `json:"cpu"`
	EphemeralDisk  int         `json:"ephemeral_disk"`
	PersistentDisk int         `json:"persistent_disk"`
	Memory         int         `json:"memory"`
	InstanceType   string      `json:"instance_type"`
	Hostname       string      `json:"hostname"`
	Role           interface{} `json:"role"`
}

type Metadata struct {
	InstanceGUIDAtTenant string      `json:"instance_guid_at_tenant"`
	UserParams           interface{} `json:"user_params"`
	PlanGUID             string      `json:"plan_guid"`
	TenantID             string      `json:"tenant_id"`
	OrganizationGUID     string      `json:"organization_guid"`
	SpaceGUID            string      `json:"space_guid"`
}

type Credential struct {
	ID           int    `json:"id"`
	InstanceID   int    `json:"instance_id"`
	GUIDAtTenant string `json:"guid_at_tenant"`
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

type LastState struct {
	State       string `json:"state"`
	Description string `json:"description"`
}

type ProvisonPayload struct {
	OrganizationGUID string        `json:"organization_guid"`
	PlanID           string        `json:"plan_id"`
	ServiceID        string        `json:"service_id"`
	SpaceGUID        string        `json:"space_guid"`
	Parameters       interface{}   `json:"parameters"`
	Context          ContextPaylod `json:"context"`
}

type UpdatePayload struct {
	ServiceID      string               `json:"service_id"`
	PlanID         string               `json:"plan_id"`
	Parameters     interface{}          `json:"parameters"`
	PreviousValues PreviousUpdateValues `json:"previous_values"`
	Context        ContextPaylod        `json:"context"`
}

type ContextPaylod struct {
	OrganizationID string `json:"organization_id"`
	SpaceID        string `json:"space_id"`
}

type PreviousUpdateValues struct {
	PlanID         string `json:"plan_id"`
	ServiceID      string `json:"service_id"`
	OrganizationID string `json:"organization_id"`
	SpaceID        string `json:"space_id"`
}

type BindPayload struct {
	ServiceID  string      `json:"service_id"`
	PlanID     string      `json:"plan_id"`
	Parameters interface{} `json:"parameters"`
}

type CustomPaylod struct {
	Parameters interface{} `json:"parameters"`
}

type SBError struct {
	Description string `json:"description"`
	Error       string `json:"error"`
}
