package helpers

// BoolPtr returns a pointer to the given bool value.
// This is a helper function for creating pointers to bool literals,
// which is commonly needed when working with API request structs.
func BoolPtr(b bool) *bool {
	return &b
}
