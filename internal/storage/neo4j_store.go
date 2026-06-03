package storage

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"albumer/internal/domain"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jStore struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jStore(ctx context.Context, uri, user, password string) (*Neo4jStore, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, password, ""))
	if err != nil {
		return nil, err
	}
	if err := driver.VerifyConnectivity(ctx); err != nil {
		_ = driver.Close(ctx)
		return nil, err
	}
	return &Neo4jStore{driver: driver}, nil
}

func (s *Neo4jStore) Close(ctx context.Context) error {
	return s.driver.Close(ctx)
}

func (s *Neo4jStore) Driver() neo4j.DriverWithContext {
	return s.driver
}

func (s *Neo4jStore) EnsureSchema(ctx context.Context) error {
	statements := []string{
		"CREATE CONSTRAINT person_id IF NOT EXISTS FOR (n:Person) REQUIRE n.id IS UNIQUE",
		"CREATE CONSTRAINT group_id IF NOT EXISTS FOR (n:Band) REQUIRE n.id IS UNIQUE",
		"CREATE CONSTRAINT music_item_id IF NOT EXISTS FOR (n:MusicItem) REQUIRE n.id IS UNIQUE",
		"CREATE CONSTRAINT book_id IF NOT EXISTS FOR (n:Book) REQUIRE n.id IS UNIQUE",
		"CREATE FULLTEXT INDEX names IF NOT EXISTS FOR (n:Person|Band) ON EACH [n.name, n.sortName]",
		"CREATE FULLTEXT INDEX catalog_titles IF NOT EXISTS FOR (n:MusicItem|Book) ON EACH [n.title]",
		"CREATE INDEX person_search_name IF NOT EXISTS FOR (n:Person) ON (n.searchName)",
		"CREATE INDEX band_search_name IF NOT EXISTS FOR (n:Band) ON (n.searchName)",
		"CREATE INDEX music_search_title IF NOT EXISTS FOR (n:MusicItem) ON (n.searchTitle)",
		"CREATE INDEX book_search_title IF NOT EXISTS FOR (n:Book) ON (n.searchTitle)",
	}

	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for _, statement := range statements {
		if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, statement, nil)
			if err != nil {
				return nil, err
			}
			return result.Consume(ctx)
		}); err != nil {
			return fmt.Errorf("%s: %w", statement, err)
		}
	}
	return nil
}

func (s *Neo4jStore) ListCatalog(ctx context.Context, filter domain.CatalogFilter) ([]domain.CatalogItem, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"q":      domain.SearchKey(filter.Query),
		"kind":   filter.Kind,
		"medium": string(filter.Medium),
		"status": filter.Status,
	}

	var items []domain.CatalogItem
	music, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return readMusicCatalog(ctx, tx, params)
	})
	if err != nil {
		return nil, err
	}
	items = append(items, music.([]domain.CatalogItem)...)

	books, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return readBookCatalog(ctx, tx, params)
	})
	if err != nil {
		return nil, err
	}
	items = append(items, books.([]domain.CatalogItem)...)

	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
	})
	return items, nil
}

func readMusicCatalog(ctx context.Context, tx neo4j.ManagedTransaction, params map[string]any) ([]domain.CatalogItem, error) {
	result, err := tx.Run(ctx, `
MATCH (m:MusicItem)
OPTIONAL MATCH (m)-[:BY]->(c)
WITH m, collect(DISTINCT {
  id: c.id,
  name: c.name,
  searchName: c.searchName,
  kind: CASE WHEN c:Band THEN 'group' ELSE 'person' END
}) AS contributors
WHERE ($kind = '' OR $kind = 'music')
  AND ($medium = '' OR m.medium = $medium)
  AND ($status = '' OR m.status = $status)
  AND ($q = '' OR m.searchTitle CONTAINS $q OR any(contributor IN contributors WHERE contributor.searchName CONTAINS $q))
RETURN m.id AS id, m.title AS title, m.medium AS medium, m.status AS status,
       m.coverPath AS coverPath, m.releaseYear AS year, contributors
`, params)
	if err != nil {
		return nil, err
	}

	var items []domain.CatalogItem
	for result.Next(ctx) {
		record := result.Record()
		contributors := readContributors(recordValue[[]any](record, "contributors"))
		items = append(items, domain.CatalogItem{
			ID:           recordString(record, "id"),
			Kind:         "music",
			Title:        recordString(record, "title"),
			Subtitle:     contributorNames(contributors),
			Medium:       domain.MusicMedium(recordString(record, "medium")),
			Status:       recordString(record, "status"),
			CoverPath:    recordString(record, "coverPath"),
			Year:         recordInt(record, "year"),
			Contributors: contributors,
		})
	}
	return items, result.Err()
}

func readBookCatalog(ctx context.Context, tx neo4j.ManagedTransaction, params map[string]any) ([]domain.CatalogItem, error) {
	result, err := tx.Run(ctx, `
MATCH (b:Book)
OPTIONAL MATCH (b)-[:AUTHORED_BY]->(p:Person)
WITH b, collect(DISTINCT {id: p.id, name: p.name, searchName: p.searchName, kind: 'person'}) AS contributors
WHERE ($kind = '' OR $kind = 'book')
  AND ($status = '' OR b.status = $status)
  AND ($q = '' OR b.searchTitle CONTAINS $q OR any(contributor IN contributors WHERE contributor.searchName CONTAINS $q))
RETURN b.id AS id, b.title AS title, b.status AS status,
       b.coverPath AS coverPath, b.publishedYear AS year, contributors
`, params)
	if err != nil {
		return nil, err
	}

	var items []domain.CatalogItem
	for result.Next(ctx) {
		record := result.Record()
		contributors := readContributors(recordValue[[]any](record, "contributors"))
		items = append(items, domain.CatalogItem{
			ID:           recordString(record, "id"),
			Kind:         "book",
			Title:        recordString(record, "title"),
			Subtitle:     contributorNames(contributors),
			Status:       recordString(record, "status"),
			CoverPath:    recordString(record, "coverPath"),
			Year:         recordInt(record, "year"),
			Contributors: contributors,
		})
	}
	return items, result.Err()
}

func (s *Neo4jStore) ListPeople(ctx context.Context) ([]domain.Person, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, "MATCH (p:Person) RETURN p.id AS id, p.name AS name, p.sortName AS sortName, p.notes AS notes ORDER BY toLower(p.name)", nil)
		if err != nil {
			return nil, err
		}
		var people []domain.Person
		for result.Next(ctx) {
			record := result.Record()
			people = append(people, domain.Person{
				ID:       recordString(record, "id"),
				Name:     recordString(record, "name"),
				SortName: recordString(record, "sortName"),
				Notes:    recordString(record, "notes"),
			})
		}
		return people, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return value.([]domain.Person), nil
}

func (s *Neo4jStore) ListGroups(ctx context.Context) ([]domain.Group, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, "MATCH (g:Band) RETURN g.id AS id, g.name AS name, g.sortName AS sortName, g.notes AS notes ORDER BY toLower(g.name)", nil)
		if err != nil {
			return nil, err
		}
		var groups []domain.Group
		for result.Next(ctx) {
			record := result.Record()
			groups = append(groups, domain.Group{
				ID:       recordString(record, "id"),
				Name:     recordString(record, "name"),
				SortName: recordString(record, "sortName"),
				Notes:    recordString(record, "notes"),
			})
		}
		return groups, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return value.([]domain.Group), nil
}

func (s *Neo4jStore) CreatePerson(ctx context.Context, person domain.Person) (domain.Person, error) {
	if err := domain.ValidateNamed(person.Name); err != nil {
		return person, err
	}
	person.ID = newID()
	now := time.Now().UTC()
	person.CreatedAt = now
	person.UpdatedAt = now

	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
CREATE (:Person {id: $id, name: $name, sortName: $sortName, searchName: $searchName, notes: $notes, createdAt: $createdAt, updatedAt: $updatedAt})
`, nodeParams(person.ID, person.Name, person.SortName, person.Notes, now))
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return person, err
}

func (s *Neo4jStore) CreateGroup(ctx context.Context, group domain.Group, memberIDs []string) (domain.Group, error) {
	if err := domain.ValidateNamed(group.Name); err != nil {
		return group, err
	}
	group.ID = newID()
	now := time.Now().UTC()
	group.CreatedAt = now
	group.UpdatedAt = now

	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
CREATE (g:Band {id: $id, name: $name, sortName: $sortName, searchName: $searchName, notes: $notes, createdAt: $createdAt, updatedAt: $updatedAt})
WITH g
UNWIND $memberIDs AS memberID
MATCH (p:Person {id: memberID})
MERGE (g)-[:HAS_MEMBER]->(p)
`, map[string]any{
			"id": group.ID, "name": group.Name, "sortName": group.SortName, "notes": group.Notes,
			"createdAt": now, "updatedAt": now, "memberIDs": memberIDs,
		})
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return group, err
}

func (s *Neo4jStore) CreateMusicItem(ctx context.Context, item domain.MusicItem) (domain.MusicItem, error) {
	if err := domain.ValidateMusicItem(item); err != nil {
		return item, err
	}
	item.ID = newID()
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now

	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
CREATE (m:MusicItem {
  id: $id, title: $title, medium: $medium, releaseYear: $year, status: $status,
  searchTitle: $searchTitle, coverPath: $coverPath, notes: $notes, createdAt: $createdAt, updatedAt: $updatedAt
})
WITH m
UNWIND $contributors AS contributor
OPTIONAL MATCH (p:Person {id: contributor.id})
OPTIONAL MATCH (g:Band {id: contributor.id})
WITH m, CASE WHEN contributor.kind = 'group' THEN g ELSE p END AS c
WHERE c IS NOT NULL
MERGE (m)-[:BY]->(c)
`, map[string]any{
			"id": item.ID, "title": item.Title, "medium": string(item.Medium), "year": item.Year,
			"status": item.Status, "searchTitle": domain.SearchKey(item.Title), "coverPath": item.CoverPath, "notes": item.Notes,
			"createdAt": now, "updatedAt": now, "contributors": contributorParams(item.Artists),
		})
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return item, err
}

func (s *Neo4jStore) CreateBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	if err := domain.ValidateBook(book); err != nil {
		return book, err
	}
	book.ID = newID()
	now := time.Now().UTC()
	book.CreatedAt = now
	book.UpdatedAt = now

	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	authorIDs := make([]string, 0, len(book.Authors))
	for _, author := range book.Authors {
		authorIDs = append(authorIDs, author.ID)
	}
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
CREATE (b:Book {
  id: $id, title: $title, isbn: $isbn, publishedYear: $year, status: $status,
  searchTitle: $searchTitle, coverPath: $coverPath, notes: $notes, createdAt: $createdAt, updatedAt: $updatedAt
})
WITH b
UNWIND $authorIDs AS authorID
MATCH (p:Person {id: authorID})
MERGE (b)-[:AUTHORED_BY]->(p)
`, map[string]any{
			"id": book.ID, "title": book.Title, "isbn": book.ISBN, "year": book.Year,
			"status": book.Status, "searchTitle": domain.SearchKey(book.Title), "coverPath": book.CoverPath, "notes": book.Notes,
			"createdAt": now, "updatedAt": now, "authorIDs": authorIDs,
		})
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return book, err
}

func nodeParams(id, name, sortName, notes string, now time.Time) map[string]any {
	return map[string]any{
		"id": id, "name": name, "sortName": sortName, "searchName": domain.SearchKey(name), "notes": notes,
		"createdAt": now, "updatedAt": now,
	}
}

func contributorParams(contributors []domain.Contributor) []map[string]any {
	params := make([]map[string]any, 0, len(contributors))
	for _, contributor := range contributors {
		params = append(params, map[string]any{"id": contributor.ID, "kind": contributor.Kind})
	}
	return params
}

func readContributors(values []any) []domain.Contributor {
	contributors := make([]domain.Contributor, 0, len(values))
	for _, value := range values {
		row, ok := value.(map[string]any)
		if !ok || row["id"] == nil {
			continue
		}
		contributors = append(contributors, domain.Contributor{
			ID:   asString(row["id"]),
			Name: asString(row["name"]),
			Kind: asString(row["kind"]),
		})
	}
	return contributors
}

func contributorNames(contributors []domain.Contributor) string {
	names := make([]string, 0, len(contributors))
	for _, contributor := range contributors {
		if contributor.Name != "" {
			names = append(names, contributor.Name)
		}
	}
	return strings.Join(names, ", ")
}

func recordString(record *neo4j.Record, key string) string {
	return asString(recordValue[any](record, key))
}

func recordInt(record *neo4j.Record, key string) int {
	switch value := recordValue[any](record, key).(type) {
	case int:
		return value
	case int64:
		return int(value)
	default:
		return 0
	}
}

func recordValue[T any](record *neo4j.Record, key string) T {
	value, _ := record.Get(key)
	typed, _ := value.(T)
	return typed
}

func asString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return ""
	}
}
