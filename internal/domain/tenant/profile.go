package tenant

import (
	"time"

	"github.com/qkitzero/fitness-service/internal/domain/address"
	"github.com/qkitzero/fitness-service/internal/domain/contact"
)

type Profile interface {
	TenantID() TenantID
	Address() address.Address
	Phone() *contact.Phone
	Email() *contact.Email
	HomepageURL() *HomepageURL
	Note() *Note
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(addr address.Address, phone *contact.Phone, email *contact.Email, homepageURL *HomepageURL, note *Note)
}

type profile struct {
	tenantID    TenantID
	address     address.Address
	phone       *contact.Phone
	email       *contact.Email
	homepageURL *HomepageURL
	note        *Note
	createdAt   time.Time
	updatedAt   time.Time
}

func (p profile) TenantID() TenantID {
	return p.tenantID
}

func (p profile) Address() address.Address {
	return p.address
}

func (p profile) Phone() *contact.Phone {
	if p.phone == nil {
		return nil
	}
	ph := *p.phone
	return &ph
}

func (p profile) Email() *contact.Email {
	if p.email == nil {
		return nil
	}
	e := *p.email
	return &e
}

func (p profile) HomepageURL() *HomepageURL {
	if p.homepageURL == nil {
		return nil
	}
	h := *p.homepageURL
	return &h
}

func (p profile) Note() *Note {
	if p.note == nil {
		return nil
	}
	n := *p.note
	return &n
}

func (p profile) CreatedAt() time.Time {
	return p.createdAt
}

func (p profile) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *profile) Update(addr address.Address, phone *contact.Phone, email *contact.Email, homepageURL *HomepageURL, note *Note) {
	p.address = addr
	p.phone = nil
	if phone != nil {
		ph := *phone
		p.phone = &ph
	}
	p.email = nil
	if email != nil {
		e := *email
		p.email = &e
	}
	p.homepageURL = nil
	if homepageURL != nil {
		h := *homepageURL
		p.homepageURL = &h
	}
	p.note = nil
	if note != nil {
		n := *note
		p.note = &n
	}
	p.updatedAt = time.Now().UTC()
}

func NewProfile(
	tenantID TenantID,
	addr address.Address,
	phone *contact.Phone,
	email *contact.Email,
	homepageURL *HomepageURL,
	note *Note,
	createdAt time.Time,
	updatedAt time.Time,
) Profile {
	p := &profile{
		tenantID:  tenantID,
		address:   addr,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
	if phone != nil {
		ph := *phone
		p.phone = &ph
	}
	if email != nil {
		e := *email
		p.email = &e
	}
	if homepageURL != nil {
		h := *homepageURL
		p.homepageURL = &h
	}
	if note != nil {
		n := *note
		p.note = &n
	}
	return p
}
