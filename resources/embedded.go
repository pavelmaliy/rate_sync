package resources

import "embed"

// EmbeddedFS embeds static resources in the running go application.
// nolint
//
//go:embed firestore/*
var EmbeddedFS embed.FS
