package models

import "strings"

type ComputerId int

type Computer struct {
	Id          ComputerId
	Serial      uint32
	Type        uint32
	Vendor      string
	Name        string
	Fingerprint uint32
}

func (computer *Computer) FullName() string {
	parts := make([]string, 0)
	if computer.Vendor != "" {
		parts = append(parts, computer.Vendor)
	}
	if computer.Name != "" {
		parts = append(parts, computer.Name)
	}
	return strings.Join(parts, " ")
}
