package storage

import (
	"context"

	"albumer/internal/domain"
)

type Store interface {
	EnsureSchema(ctx context.Context) error
	Close(ctx context.Context) error
	ListCatalog(ctx context.Context, filter domain.CatalogFilter) ([]domain.CatalogItem, error)
	ListPeople(ctx context.Context) ([]domain.Person, error)
	ListGroups(ctx context.Context) ([]domain.Group, error)
	ListMusicReferences(ctx context.Context, query, excludeID string) ([]domain.MusicReference, error)
	GetPerson(ctx context.Context, id string) (domain.Person, error)
	GetGroup(ctx context.Context, id string) (domain.Group, error)
	GetMusicItem(ctx context.Context, id string) (domain.MusicItem, error)
	GetBook(ctx context.Context, id string) (domain.Book, error)
	CreatePerson(ctx context.Context, person domain.Person) (domain.Person, error)
	CreateGroup(ctx context.Context, group domain.Group, memberIDs []string) (domain.Group, error)
	CreateMusicItem(ctx context.Context, item domain.MusicItem) (domain.MusicItem, error)
	CreateBook(ctx context.Context, book domain.Book) (domain.Book, error)
	UpdatePerson(ctx context.Context, person domain.Person) (domain.Person, error)
	UpdateGroup(ctx context.Context, group domain.Group, memberIDs []string) (domain.Group, error)
	UpdateMusicItem(ctx context.Context, item domain.MusicItem) (domain.MusicItem, error)
	UpdateBook(ctx context.Context, book domain.Book) (domain.Book, error)
}
