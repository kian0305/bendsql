package api

import "time"

type AccountInfoDTO struct {
	ID              uint64    `json:"id"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	State           string    `json:"state"`
	AvatarURL       string    `json:"avatarURL"`
	DefaultOrgSlug  string    `json:"defaultOrgSlug"`
	PasswordEnabled bool      `json:"passwordEnabled"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type OrgInfoDTO struct {
	Name         string    `json:"name"`
	MembersCount int64     `json:"memberCount"`
	TenantID     string    `json:"tenantID"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type WarehouseStatusDTO struct {
	Name           string `json:"id,omitempty"`
	ReadyInstances int64  `json:"readyInstances,omitempty"`
	Size           string `json:"size,omitempty"`
	State          string `json:"state,omitempty"`
	TotalInstances int64  `json:"totalInstances,omitempty"`
}
