package entities

import (
	"encoding/json"

	"github.com/google/uuid"
)

// BasicMemberData represents basic member details.
type BasicMemberDataResponse struct {
	Data BasicMemberData `json:"data"`
}

type BasicMemberData struct {
	MemberID    uuid.UUID `json:"member_id"`
	Name        string    `json:"member_name"`
	PartnerName string    `json:"partner_name"`
	PartnerID   uuid.UUID `json:"partner_id"`
	Email       string    `json:"member_email"`
	MemberType  string    `json:"member_type"`
	MemberRoles []string  `json:"member_roles"`
	ProviderID  uuid.UUID `json:"provider_id"`
}

func (m *BasicMemberData) UnmarshalJSON(data []byte) error {
	type Alias BasicMemberData
	aux := &struct {
		MemberID   *string `json:"member_id"`
		PartnerID  string  `json:"partner_id"`
		ProviderID string  `json:"provider_id"`
		// Other fields...
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert the strings to UUID for MemberID and PartnerID
	id, err := uuid.Parse(*aux.MemberID)
	if err != nil {
		return err
	}
	m.MemberID = id

	partnerID, err := uuid.Parse(aux.PartnerID)
	if err != nil {
		return err
	}
	m.PartnerID = partnerID

	providerID, err := uuid.Parse(aux.ProviderID)
	if err != nil {
		return err
	}
	m.ProviderID = providerID

	return nil
}
