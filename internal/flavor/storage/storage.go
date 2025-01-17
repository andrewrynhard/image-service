// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package storage implements a storage for flavor data.
package storage

import (
	"context"
)

// Storage is the flavor storage.
type Storage interface {
	Head(ctx context.Context, id string) error
	Get(ctx context.Context, id string) ([]byte, error)
	Put(ctx context.Context, id string, data []byte) error
}

// ErrNotFoundTag tags the errors when the flavor is not found.
type ErrNotFoundTag = struct{}
