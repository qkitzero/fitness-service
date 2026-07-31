package address

type Address struct {
	postalCode *PostalCode
	prefecture *Prefecture
	city       *City
	street     *Street
	building   *Building
}

func (a Address) PostalCode() *PostalCode {
	if a.postalCode == nil {
		return nil
	}
	pc := *a.postalCode
	return &pc
}

func (a Address) Prefecture() *Prefecture {
	if a.prefecture == nil {
		return nil
	}
	p := *a.prefecture
	return &p
}

func (a Address) City() *City {
	if a.city == nil {
		return nil
	}
	c := *a.city
	return &c
}

func (a Address) Street() *Street {
	if a.street == nil {
		return nil
	}
	st := *a.street
	return &st
}

func (a Address) Building() *Building {
	if a.building == nil {
		return nil
	}
	b := *a.building
	return &b
}

func NewAddress(
	postalCode *PostalCode,
	prefecture *Prefecture,
	city *City,
	street *Street,
	building *Building,
) Address {
	a := Address{}
	if postalCode != nil {
		pc := *postalCode
		a.postalCode = &pc
	}
	if prefecture != nil {
		p := *prefecture
		a.prefecture = &p
	}
	if city != nil {
		c := *city
		a.city = &c
	}
	if street != nil {
		st := *street
		a.street = &st
	}
	if building != nil {
		b := *building
		a.building = &b
	}
	return a
}
