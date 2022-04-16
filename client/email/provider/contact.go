package provider

import "strings"

type Contact struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type ContactList []*Contact

func (cl ContactList) String() string {
	ret := make([]string, 0, len(cl))
	for _, contact := range cl {
		ret = append(ret, contact.String())
	}
	return strings.Join(ret, ", ")
}

func (c Contact) String() string {
	if c.Name == "" {
		return c.Address
	} else {
		return c.Name + " <" + c.Address + ">"
	}
}
