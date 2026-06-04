package storage

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"albumer/internal/domain"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var ErrNotFound = errors.New("not found")

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
		"CREATE INDEX person_legacy_id IF NOT EXISTS FOR (n:Person) ON (n.legacyId)",
		"CREATE INDEX band_legacy_id IF NOT EXISTS FOR (n:Band) ON (n.legacyId)",
		"CREATE INDEX music_legacy_id IF NOT EXISTS FOR (n:MusicItem) ON (n.legacyId)",
		"CREATE INDEX book_legacy_id IF NOT EXISTS FOR (n:Book) ON (n.legacyId)",
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
		"status": string(filter.Status),
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
			Status:       domain.NormalizeStatus(recordString(record, "status")),
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
			Status:       domain.NormalizeStatus(recordString(record, "status")),
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
		result, err := tx.Run(ctx, `
			MATCH (p:Person)
			OPTIONAL MATCH (music:MusicItem)-[:BY]->(p)
			WITH p, count(DISTINCT music) AS albumCount
			OPTIONAL MATCH (book:Book)-[:AUTHORED_BY]->(p)
			WITH p, albumCount, count(DISTINCT book) AS bookCount
			RETURN p.id AS id, p.name AS name, p.sortName AS sortName, p.notes AS notes,
			       CASE WHEN p.roles IS NULL
			            THEN CASE
			                 WHEN albumCount > 0 AND bookCount > 0 THEN ['artist', 'author']
			                 WHEN albumCount > 0 THEN ['artist']
			                 WHEN bookCount > 0 THEN ['author']
			                 ELSE []
			                 END
			            ELSE p.roles
			       END AS roles,
			       albumCount, bookCount
			ORDER BY toLower(p.name)
		`, nil)
		if err != nil {
			return nil, err
		}
		var people []domain.Person
		for result.Next(ctx) {
			record := result.Record()
			people = append(people, domain.Person{
				ID:         recordString(record, "id"),
				Name:       recordString(record, "name"),
				SortName:   recordString(record, "sortName"),
				Notes:      recordString(record, "notes"),
				Roles:      recordStrings(record, "roles"),
				AlbumCount: recordInt(record, "albumCount"),
				BookCount:  recordInt(record, "bookCount"),
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

func (s *Neo4jStore) ListMusicReferences(ctx context.Context, query, excludeID string) ([]domain.MusicReference, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	limit := 10000
	if strings.TrimSpace(query) != "" {
		limit = 20
	}

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (m:MusicItem)
WHERE ($q = '' OR m.searchTitle CONTAINS $q)
  AND ($excludeID = '' OR m.id <> $excludeID)
RETURN m.id AS id, m.title AS title, m.medium AS medium
ORDER BY toLower(m.title), m.medium
LIMIT $limit
`, map[string]any{"q": domain.SearchKey(query), "excludeID": excludeID, "limit": limit})
		if err != nil {
			return nil, err
		}
		var items []domain.MusicReference
		for result.Next(ctx) {
			record := result.Record()
			items = append(items, domain.MusicReference{
				ID:     recordString(record, "id"),
				Title:  recordString(record, "title"),
				Medium: domain.MusicMedium(recordString(record, "medium")),
			})
		}
		return items, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return value.([]domain.MusicReference), nil
}

func (s *Neo4jStore) GetPerson(ctx context.Context, id string) (domain.Person, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (p:Person {id: $id})
OPTIONAL MATCH (music:MusicItem)-[:BY]->(p)
WITH p, count(DISTINCT music) AS albumCount
OPTIONAL MATCH (book:Book)-[:AUTHORED_BY]->(p)
WITH p, albumCount, count(DISTINCT book) AS bookCount
RETURN p.id AS id, p.name AS name, p.sortName AS sortName, p.notes AS notes,
       CASE WHEN p.roles IS NULL
            THEN CASE
                 WHEN albumCount > 0 AND bookCount > 0 THEN ['artist', 'author']
                 WHEN albumCount > 0 THEN ['artist']
                 WHEN bookCount > 0 THEN ['author']
                 ELSE []
                 END
            ELSE p.roles
       END AS roles
`, map[string]any{"id": id})
		if err != nil {
			return domain.Person{}, err
		}
		if !result.Next(ctx) {
			return domain.Person{}, ErrNotFound
		}
		record := result.Record()
		return domain.Person{ID: recordString(record, "id"), Name: recordString(record, "name"), SortName: recordString(record, "sortName"), Notes: recordString(record, "notes"), Roles: recordStrings(record, "roles")}, result.Err()
	})
	if err != nil {
		return domain.Person{}, err
	}
	return value.(domain.Person), nil
}

func (s *Neo4jStore) GetGroup(ctx context.Context, id string) (domain.Group, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (g:Band {id: $id})
OPTIONAL MATCH (g)-[:HAS_MEMBER]->(p:Person)
RETURN g.id AS id, g.name AS name, g.sortName AS sortName, g.notes AS notes,
       collect(DISTINCT {id: p.id, name: p.name}) AS members
`, map[string]any{"id": id})
		if err != nil {
			return domain.Group{}, err
		}
		if !result.Next(ctx) {
			return domain.Group{}, ErrNotFound
		}
		record := result.Record()
		group := domain.Group{ID: recordString(record, "id"), Name: recordString(record, "name"), SortName: recordString(record, "sortName"), Notes: recordString(record, "notes")}
		for _, member := range readContributors(recordValue[[]any](record, "members")) {
			group.Members = append(group.Members, domain.Person{ID: member.ID, Name: member.Name})
		}
		return group, result.Err()
	})
	if err != nil {
		return domain.Group{}, err
	}
	return value.(domain.Group), nil
}

func (s *Neo4jStore) GetMusicItem(ctx context.Context, id string) (domain.MusicItem, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (m:MusicItem {id: $id})
OPTIONAL MATCH (m)-[:BY]->(c)
WITH m, collect(DISTINCT {id: c.id, name: c.name, kind: CASE WHEN c:Band THEN 'group' ELSE 'person' END}) AS contributors
OPTIONAL MATCH (m)-[:SAME_RELEASE]-(related:MusicItem)
RETURN m.id AS id, m.title AS title, m.medium AS medium, m.releaseYear AS year,
       m.status AS status, m.coverPath AS coverPath, m.notes AS notes,
       contributors,
       collect(DISTINCT {id: related.id, title: related.title, medium: related.medium}) AS related
`, map[string]any{"id": id})
		if err != nil {
			return domain.MusicItem{}, err
		}
		if !result.Next(ctx) {
			return domain.MusicItem{}, ErrNotFound
		}
		record := result.Record()
		item := domain.MusicItem{
			ID:        recordString(record, "id"),
			Title:     recordString(record, "title"),
			Medium:    domain.MusicMedium(recordString(record, "medium")),
			Year:      recordInt(record, "year"),
			Status:    domain.NormalizeStatus(recordString(record, "status")),
			CoverPath: recordString(record, "coverPath"),
			Notes:     recordString(record, "notes"),
			Artists:   readContributors(recordValue[[]any](record, "contributors")),
			Related:   readMusicReferences(recordValue[[]any](record, "related")),
		}
		return item, result.Err()
	})
	if err != nil {
		return domain.MusicItem{}, err
	}
	return value.(domain.MusicItem), nil
}

func (s *Neo4jStore) GetBook(ctx context.Context, id string) (domain.Book, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	value, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (b:Book {id: $id})
OPTIONAL MATCH (b)-[:AUTHORED_BY]->(p:Person)
RETURN b.id AS id, b.title AS title, b.isbn AS isbn, b.publishedYear AS year,
       b.status AS status, b.coverPath AS coverPath, b.notes AS notes,
       collect(DISTINCT {id: p.id, name: p.name}) AS authors
`, map[string]any{"id": id})
		if err != nil {
			return domain.Book{}, err
		}
		if !result.Next(ctx) {
			return domain.Book{}, ErrNotFound
		}
		record := result.Record()
		book := domain.Book{
			ID:        recordString(record, "id"),
			Title:     recordString(record, "title"),
			ISBN:      recordString(record, "isbn"),
			Year:      recordInt(record, "year"),
			Status:    domain.NormalizeStatus(recordString(record, "status")),
			CoverPath: recordString(record, "coverPath"),
			Notes:     recordString(record, "notes"),
		}
		for _, author := range readContributors(recordValue[[]any](record, "authors")) {
			book.Authors = append(book.Authors, domain.Person{ID: author.ID, Name: author.Name})
		}
		return book, result.Err()
	})
	if err != nil {
		return domain.Book{}, err
	}
	return value.(domain.Book), nil
}

func (s *Neo4jStore) CreatePerson(ctx context.Context, person domain.Person) (domain.Person, error) {
	if err := domain.ValidatePerson(person); err != nil {
		return person, err
	}
	person.SortName = domain.DefaultSortName(person.Name, person.SortName)
	person.ID = newID()
	now := time.Now().UTC()
	person.CreatedAt = now
	person.UpdatedAt = now

	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
CREATE (:Person {id: $id, name: $name, sortName: $sortName, searchName: $searchName, notes: $notes, roles: $roles, createdAt: $createdAt, updatedAt: $updatedAt})
`, personParams(person, now))
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
	group.SortName = domain.DefaultSortName(group.Name, group.SortName)
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
WITH DISTINCT m
UNWIND $relatedIDs AS relatedID
MATCH (related:MusicItem {id: relatedID})
WHERE related.id <> m.id
MERGE (m)-[:SAME_RELEASE]->(related)
`, map[string]any{
			"id": item.ID, "title": item.Title, "medium": string(item.Medium), "year": item.Year,
			"status": string(item.Status), "searchTitle": domain.SearchKey(item.Title), "coverPath": item.CoverPath, "notes": item.Notes,
			"createdAt": now, "updatedAt": now, "contributors": contributorParams(item.Artists),
			"relatedIDs": relatedIDs(item),
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
			"status": string(book.Status), "searchTitle": domain.SearchKey(book.Title), "coverPath": book.CoverPath, "notes": book.Notes,
			"createdAt": now, "updatedAt": now, "authorIDs": authorIDs,
		})
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return book, err
}

func (s *Neo4jStore) UpdatePerson(ctx context.Context, person domain.Person) (domain.Person, error) {
	if err := domain.ValidatePerson(person); err != nil {
		return person, err
	}
	person.SortName = domain.DefaultSortName(person.Name, person.SortName)
	now := time.Now().UTC()

	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (p:Person {id: $id})
SET p.name = $name,
    p.sortName = $sortName,
    p.searchName = $searchName,
    p.notes = $notes,
    p.roles = $roles,
    p.updatedAt = $updatedAt
`, personParams(person, now))
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return person, err
}

func (s *Neo4jStore) UpdateGroup(ctx context.Context, group domain.Group, memberIDs []string) (domain.Group, error) {
	if err := domain.ValidateNamed(group.Name); err != nil {
		return group, err
	}
	group.SortName = domain.DefaultSortName(group.Name, group.SortName)
	now := time.Now().UTC()

	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (g:Band {id: $id})
SET g.name = $name,
    g.sortName = $sortName,
    g.searchName = $searchName,
    g.notes = $notes,
    g.updatedAt = $updatedAt
WITH g
OPTIONAL MATCH (g)-[old:HAS_MEMBER]->(:Person)
DELETE old
WITH g
UNWIND $memberIDs AS memberID
MATCH (p:Person {id: memberID})
MERGE (g)-[:HAS_MEMBER]->(p)
`, map[string]any{
			"id": group.ID, "name": group.Name, "sortName": group.SortName, "searchName": domain.SearchKey(group.Name), "notes": group.Notes,
			"updatedAt": now, "memberIDs": memberIDs,
		})
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return group, err
}

func (s *Neo4jStore) UpdateMusicItem(ctx context.Context, item domain.MusicItem) (domain.MusicItem, error) {
	if err := domain.ValidateMusicItem(item); err != nil {
		return item, err
	}
	now := time.Now().UTC()

	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (m:MusicItem {id: $id})
SET m.title = $title,
    m.medium = $medium,
    m.releaseYear = $year,
    m.status = $status,
    m.searchTitle = $searchTitle,
    m.coverPath = $coverPath,
    m.notes = $notes,
    m.updatedAt = $updatedAt
WITH m
OPTIONAL MATCH (m)-[old:BY]->()
DELETE old
WITH m
UNWIND $contributors AS contributor
OPTIONAL MATCH (p:Person {id: contributor.id})
OPTIONAL MATCH (g:Band {id: contributor.id})
WITH m, CASE WHEN contributor.kind = 'group' THEN g ELSE p END AS c
WHERE c IS NOT NULL
MERGE (m)-[:BY]->(c)
WITH DISTINCT m
OPTIONAL MATCH (m)-[oldRelated:SAME_RELEASE]-()
DELETE oldRelated
WITH m
UNWIND $relatedIDs AS relatedID
MATCH (related:MusicItem {id: relatedID})
WHERE related.id <> m.id
MERGE (m)-[:SAME_RELEASE]->(related)
`, map[string]any{
			"id": item.ID, "title": item.Title, "medium": string(item.Medium), "year": item.Year,
			"status": string(item.Status), "searchTitle": domain.SearchKey(item.Title), "coverPath": item.CoverPath, "notes": item.Notes,
			"updatedAt": now, "contributors": contributorParams(item.Artists),
			"relatedIDs": relatedIDs(item),
		})
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	return item, err
}

func (s *Neo4jStore) UpdateBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	if err := domain.ValidateBook(book); err != nil {
		return book, err
	}
	now := time.Now().UTC()
	authorIDs := make([]string, 0, len(book.Authors))
	for _, author := range book.Authors {
		authorIDs = append(authorIDs, author.ID)
	}

	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
MATCH (b:Book {id: $id})
SET b.title = $title,
    b.isbn = $isbn,
    b.publishedYear = $year,
    b.status = $status,
    b.searchTitle = $searchTitle,
    b.coverPath = $coverPath,
    b.notes = $notes,
    b.updatedAt = $updatedAt
WITH b
OPTIONAL MATCH (b)-[old:AUTHORED_BY]->(:Person)
DELETE old
WITH b
UNWIND $authorIDs AS authorID
MATCH (p:Person {id: authorID})
MERGE (b)-[:AUTHORED_BY]->(p)
`, map[string]any{
			"id": book.ID, "title": book.Title, "isbn": book.ISBN, "year": book.Year,
			"status": string(book.Status), "searchTitle": domain.SearchKey(book.Title), "coverPath": book.CoverPath, "notes": book.Notes,
			"updatedAt": now, "authorIDs": authorIDs,
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

func personParams(person domain.Person, now time.Time) map[string]any {
	params := nodeParams(person.ID, person.Name, person.SortName, person.Notes, now)
	params["roles"] = person.Roles
	return params
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

func readMusicReferences(values []any) []domain.MusicReference {
	items := make([]domain.MusicReference, 0, len(values))
	for _, value := range values {
		row, ok := value.(map[string]any)
		if !ok || row["id"] == nil {
			continue
		}
		items = append(items, domain.MusicReference{
			ID:     asString(row["id"]),
			Title:  asString(row["title"]),
			Medium: domain.MusicMedium(asString(row["medium"])),
		})
	}
	return items
}

func relatedIDs(item domain.MusicItem) []string {
	ids := make([]string, 0, len(item.Related))
	for _, related := range item.Related {
		if related.ID != "" && related.ID != item.ID {
			ids = append(ids, related.ID)
		}
	}
	return ids
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

func recordStrings(record *neo4j.Record, key string) []string {
	switch values := recordValue[any](record, key).(type) {
	case []string:
		return values
	case []any:
		result := make([]string, 0, len(values))
		for _, value := range values {
			if text := asString(value); text != "" {
				result = append(result, text)
			}
		}
		return result
	default:
		return nil
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
