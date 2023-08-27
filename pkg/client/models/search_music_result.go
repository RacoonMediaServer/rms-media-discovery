// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SearchMusicResult search music result
//
// swagger:model SearchMusicResult
type SearchMusicResult struct {

	// album
	Album string `json:"album,omitempty"`

	// albums count
	AlbumsCount int64 `json:"albumsCount,omitempty"`

	// artist
	Artist string `json:"artist,omitempty"`

	// genres
	Genres []string `json:"genres"`

	// picture
	Picture string `json:"picture,omitempty"`

	// release year
	ReleaseYear int64 `json:"releaseYear,omitempty"`

	// title
	// Required: true
	Title *string `json:"title"`

	// tracks count
	TracksCount int64 `json:"tracksCount,omitempty"`

	// type
	// Required: true
	// Enum: [artist album track]
	Type *string `json:"type"`
}

// Validate validates this search music result
func (m *SearchMusicResult) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTitle(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SearchMusicResult) validateTitle(formats strfmt.Registry) error {

	if err := validate.Required("title", "body", m.Title); err != nil {
		return err
	}

	return nil
}

var searchMusicResultTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["artist","album","track"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		searchMusicResultTypeTypePropEnum = append(searchMusicResultTypeTypePropEnum, v)
	}
}

const (

	// SearchMusicResultTypeArtist captures enum value "artist"
	SearchMusicResultTypeArtist string = "artist"

	// SearchMusicResultTypeAlbum captures enum value "album"
	SearchMusicResultTypeAlbum string = "album"

	// SearchMusicResultTypeTrack captures enum value "track"
	SearchMusicResultTypeTrack string = "track"
)

// prop value enum
func (m *SearchMusicResult) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, searchMusicResultTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *SearchMusicResult) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", *m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this search music result based on context it is used
func (m *SearchMusicResult) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SearchMusicResult) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SearchMusicResult) UnmarshalBinary(b []byte) error {
	var res SearchMusicResult
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
