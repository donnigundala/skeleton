package providers

import (
	firebase "github.com/donnigundala/dg-framework/dg-firebase"
)

// FirebaseProvider returns a new Firebase service provider.
func FirebaseProvider() *firebase.FirebaseServiceProvider {
	return &firebase.FirebaseServiceProvider{}
}
