package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"albumer/internal/domain"
	"albumer/internal/ids"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type LegacyData struct {
	Artists []LegacyArtist `json:"artists"`
	Albums  []LegacyAlbum  `json:"albums"`
}

type LegacyArtist struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	ArtistsIDs []string `json:"artistsIds"`
}

type LegacyAlbum struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Status     int      `json:"status"`
	ArtistsIDs []string `json:"artistsIds"`
}

type LegacyImportSummary struct {
	People     int
	Groups     int
	GroupLinks int
	CDs        int
	CDLinks    int
}

type SearchBackfillSummary struct {
	People int
	Groups int
	Music  int
	Books  int
}

func LoadLegacyData(path string) (LegacyData, error) {
	file, err := os.Open(path)
	if err != nil {
		return LegacyData{}, err
	}
	defer file.Close()

	var data LegacyData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return LegacyData{}, err
	}
	return data, nil
}

func ImportLegacyCDs(ctx context.Context, driver neo4j.DriverWithContext, data LegacyData) (LegacyImportSummary, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	value, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		now := time.Now().UTC()
		summary := LegacyImportSummary{}
		idMap := map[string]string{}

		for _, artist := range data.Artists {
			if artist.ID == "" || artist.Name == "" {
				continue
			}
			newID := ids.NewUUID()
			idMap[artist.ID] = newID
			if len(artist.ArtistsIDs) == 0 {
				if err := createLegacyPerson(ctx, tx, newID, artist, now); err != nil {
					return summary, err
				}
				summary.People++
				continue
			}

			if err := createLegacyBand(ctx, tx, newID, artist, now); err != nil {
				return summary, err
			}
			summary.Groups++
		}

		for _, artist := range data.Artists {
			bandID := idMap[artist.ID]
			if bandID == "" || len(artist.ArtistsIDs) == 0 {
				continue
			}
			for _, memberID := range artist.ArtistsIDs {
				if memberID == "" {
					continue
				}
				mappedMemberID := idMap[memberID]
				if mappedMemberID == "" {
					continue
				}
				if err := createLegacyBandMember(ctx, tx, bandID, mappedMemberID); err != nil {
					return summary, err
				}
				summary.GroupLinks++
			}
		}

		for _, album := range data.Albums {
			if album.ID == "" || album.Name == "" {
				continue
			}
			newID := ids.NewUUID()
			if err := createLegacyCD(ctx, tx, newID, album, now); err != nil {
				return summary, err
			}
			summary.CDs++

			for _, artistID := range album.ArtistsIDs {
				if artistID == "" {
					continue
				}
				mappedArtistID := idMap[artistID]
				if mappedArtistID == "" {
					continue
				}
				if err := createLegacyCDContributor(ctx, tx, newID, mappedArtistID); err != nil {
					return summary, err
				}
				summary.CDLinks++
			}
		}

		return summary, nil
	})
	if err != nil {
		return LegacyImportSummary{}, err
	}
	return value.(LegacyImportSummary), nil
}

func BackfillSearchKeys(ctx context.Context, driver neo4j.DriverWithContext) (SearchBackfillSummary, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	value, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		summary := SearchBackfillSummary{}

		people, err := fetchNamedNodes(ctx, tx, "Person")
		if err != nil {
			return summary, err
		}
		for _, node := range people {
			if err := setSearchValue(ctx, tx, "Person", node.ID, "searchName", domain.SearchKey(node.Value)); err != nil {
				return summary, err
			}
			summary.People++
		}

		groups, err := fetchNamedNodes(ctx, tx, "Band")
		if err != nil {
			return summary, err
		}
		for _, node := range groups {
			if err := setSearchValue(ctx, tx, "Band", node.ID, "searchName", domain.SearchKey(node.Value)); err != nil {
				return summary, err
			}
			summary.Groups++
		}

		music, err := fetchTitledNodes(ctx, tx, "MusicItem")
		if err != nil {
			return summary, err
		}
		for _, node := range music {
			if err := setSearchValue(ctx, tx, "MusicItem", node.ID, "searchTitle", domain.SearchKey(node.Value)); err != nil {
				return summary, err
			}
			summary.Music++
		}

		books, err := fetchTitledNodes(ctx, tx, "Book")
		if err != nil {
			return summary, err
		}
		for _, node := range books {
			if err := setSearchValue(ctx, tx, "Book", node.ID, "searchTitle", domain.SearchKey(node.Value)); err != nil {
				return summary, err
			}
			summary.Books++
		}

		return summary, nil
	})
	if err != nil {
		return SearchBackfillSummary{}, err
	}
	return value.(SearchBackfillSummary), nil
}

func LegacyStatus(value int) string {
	if value == 1 {
		return "owned"
	}
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("legacy:%d", value)
}

func createLegacyPerson(ctx context.Context, tx neo4j.ManagedTransaction, id string, artist LegacyArtist, now time.Time) error {
	result, err := tx.Run(ctx, `
CREATE (p:Person {id: $id})
SET p.name = $name,
    p.sortName = $name,
    p.searchName = $searchName,
    p.createdAt = $now,
    p.updatedAt = $now
`, map[string]any{"id": id, "name": artist.Name, "searchName": domain.SearchKey(artist.Name), "now": now})
	if err != nil {
		return err
	}
	_, err = result.Consume(ctx)
	return err
}

func createLegacyBand(ctx context.Context, tx neo4j.ManagedTransaction, id string, artist LegacyArtist, now time.Time) error {
	result, err := tx.Run(ctx, `
CREATE (g:Band {id: $id})
SET g.name = $name,
    g.sortName = $name,
    g.searchName = $searchName,
    g.createdAt = $now,
    g.updatedAt = $now
`, map[string]any{"id": id, "name": artist.Name, "searchName": domain.SearchKey(artist.Name), "now": now})
	if err != nil {
		return err
	}
	_, err = result.Consume(ctx)
	return err
}

func createLegacyBandMember(ctx context.Context, tx neo4j.ManagedTransaction, bandID, memberID string) error {
	result, err := tx.Run(ctx, `
MATCH (g:Band {id: $bandID})
MATCH (p:Person {id: $memberID})
MERGE (g)-[:HAS_MEMBER]->(p)
`, map[string]any{"bandID": bandID, "memberID": memberID})
	if err != nil {
		return err
	}
	_, err = result.Consume(ctx)
	return err
}

func createLegacyCD(ctx context.Context, tx neo4j.ManagedTransaction, id string, album LegacyAlbum, now time.Time) error {
	result, err := tx.Run(ctx, `
CREATE (m:MusicItem {id: $id})
SET m.title = $title,
    m.searchTitle = $searchTitle,
    m.medium = $medium,
    m.status = $status,
    m.createdAt = $now,
    m.updatedAt = $now
`, map[string]any{
		"id": id, "title": album.Name, "searchTitle": domain.SearchKey(album.Name), "medium": string(domain.MediumCD),
		"status": LegacyStatus(album.Status), "now": now,
	})
	if err != nil {
		return err
	}
	_, err = result.Consume(ctx)
	return err
}

func createLegacyCDContributor(ctx context.Context, tx neo4j.ManagedTransaction, albumID, artistID string) error {
	result, err := tx.Run(ctx, `
MATCH (m:MusicItem {id: $albumID})
OPTIONAL MATCH (p:Person {id: $artistID})
OPTIONAL MATCH (g:Band {id: $artistID})
WITH m, coalesce(g, p) AS contributor
WHERE contributor IS NOT NULL
MERGE (m)-[:BY]->(contributor)
`, map[string]any{"albumID": albumID, "artistID": artistID})
	if err != nil {
		return err
	}
	_, err = result.Consume(ctx)
	return err
}

type legacySearchNode struct {
	ID    string
	Value string
}

func fetchNamedNodes(ctx context.Context, tx neo4j.ManagedTransaction, label string) ([]legacySearchNode, error) {
	result, err := tx.Run(ctx, fmt.Sprintf("MATCH (n:%s) RETURN n.id AS id, n.name AS value", label), nil)
	if err != nil {
		return nil, err
	}
	return readSearchNodes(ctx, result)
}

func fetchTitledNodes(ctx context.Context, tx neo4j.ManagedTransaction, label string) ([]legacySearchNode, error) {
	result, err := tx.Run(ctx, fmt.Sprintf("MATCH (n:%s) RETURN n.id AS id, n.title AS value", label), nil)
	if err != nil {
		return nil, err
	}
	return readSearchNodes(ctx, result)
}

func readSearchNodes(ctx context.Context, result neo4j.ResultWithContext) ([]legacySearchNode, error) {
	var nodes []legacySearchNode
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("id")
		value, _ := record.Get("value")
		nodes = append(nodes, legacySearchNode{ID: stringValue(id), Value: stringValue(value)})
	}
	return nodes, result.Err()
}

func setSearchValue(ctx context.Context, tx neo4j.ManagedTransaction, label, id, property, value string) error {
	result, err := tx.Run(ctx, fmt.Sprintf("MATCH (n:%s {id: $id}) SET n.%s = $value", label, property), map[string]any{"id": id, "value": value})
	if err != nil {
		return err
	}
	_, err = result.Consume(ctx)
	return err
}

func stringValue(value any) string {
	if typed, ok := value.(string); ok {
		return typed
	}
	return ""
}
