package server

type Permissions struct {
}

// MissingPermissions is used when a specific request has no permissions associated with it
var MissingPermissions = &Permissions{}
