package validation

import (
	"errors"
	"regexp"
)

var (
	ErrBannedName  = errors.New("package name contains invalid characters")
	ErrInvalidName = errors.New("invalid version format")
)

type PackagesValidator struct {
	nameRegex    *regexp.Regexp
	versionRegex *regexp.Regexp
}

func NewPackagesValidator() *PackagesValidator {
	return &PackagesValidator{
		nameRegex:    regexp.MustCompile(`^[A-Za-z][A-Za-z@_-]*$`),
		versionRegex: regexp.MustCompile(`[0-9]`),
	}
}

func (v *PackagesValidator) ValidatePackageUpload(name, version string) error {
	if err := v.validadeName(name); err != nil {
		return err
	}

	if err := v.validateVersion(version); err != nil {
		return err
	}

	return nil
}

func (v *PackagesValidator) ValidatePackageDownload(name, version string) error {
	if err := v.validadeName(name); err != nil {
		return err
	}

	if err := v.validateVersion(version); err != nil {
		return err
	}

	return nil
}

func (v *PackagesValidator) validadeName(name string) error {
	if !v.nameRegex.MatchString(name) {
		return ErrBannedName
	}
	return nil
}

func (v *PackagesValidator) validateVersion(version string) error {
	if !v.versionRegex.MatchString(version) {
		return ErrInvalidName
	}
	return nil
}
