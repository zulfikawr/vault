package storage

import (
	"context"
	"io"
)

type Storage interface {
	Save(ctx context.Context, path string, data io.Reader) error
	        Retrieve(ctx context.Context, path string) (io.ReadCloser, error)
	                Delete(ctx context.Context, path string) error
	                Rename(ctx context.Context, oldPath, newPath string) error
	                CreateDir(ctx context.Context, path string) error
	                Exists(ctx context.Context, path string) (bool, error)
	        }
	        
