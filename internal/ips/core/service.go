package core

import (
	"context"
	"encoding/json"
	"fmt"
	"ips-lacpass-backend/internal/ips/client"
	customErrors "ips-lacpass-backend/pkg/errors"
	authMiddleware "ips-lacpass-backend/pkg/middleware"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

type IpsService struct {
	Repository *client.IpsClient
}

func NewService(r *client.IpsClient) IpsService {
	return IpsService{
		Repository: r,
	}
}

func (is *IpsService) GetIps(ctx context.Context) (map[string]interface{}, error) {
	userId, err := authMiddleware.GetUserDocIDFromContext(ctx)
	if err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 401,
			Body:       []map[string]interface{}{{"error": "user_identifier_not_found", "message": "User identifier not found in request context"}},
			Err:        err,
		}
	}

	bundle, err := is.Repository.GetDocumentReference(userId)
	if err != nil {
		fmt.Printf("Error fetching document reference: %v\n", err)
		return nil, err
	}
	entries := bundle.Entry
	if len(entries) == 0 {
		return nil, &customErrors.HttpError{
			StatusCode: 404,
			Body:       []map[string]interface{}{{"error": "not_found", "message": "No IPS found for the user"}},
			Err:        fmt.Errorf("no IPS found for the user"),
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Resource.Meta.LastUpdated > entries[j].Resource.Meta.LastUpdated
	})

	ipsBundle, err := is.Repository.GetIpsBundle(entries[0].Resource.Content[0].Attachment.URL)
	if err != nil {
		return nil, err
	}

	return ipsBundle, nil

}

// Will return the IPS composition sections
func getIPSComposition(entries []Entry) (*Composition, error) {
	i := slices.IndexFunc(entries, func(e Entry) bool {
		if e.Resource["resourceType"] != "Composition" {
			return false
		}
		var composition Composition
		if err := mapstructure.Decode(e.Resource, &composition); err != nil {
			return false
		}
		return composition.Type.Coding[0]["code"] == "60591-5"
	})
	if i == -1 {
		return nil, fmt.Errorf(`no composition found`)
	}
	comp := entries[i].Resource
	var composition Composition
	if err := mapstructure.Decode(comp, &composition); err != nil {
		return nil, fmt.Errorf(`error decoding composition: %v`, err)
	}
	composition.URL = entries[i].FullURL
	// Remove empty sections
	var result []Section
	for _, s := range composition.Section {
		if s.Code.Coding != nil {
			result = append(result, s)
		}
	}
	composition.Section = result

	return &composition, nil
}

func getEntry(reference string, current []Entry, newIpsEntries []Entry) *Entry {
	indInCurrent := slices.IndexFunc(current, func(e Entry) bool {
		return e.FullURL == reference
	})
	if indInCurrent != -1 {
		return &current[indInCurrent]
	}

	indInNew := slices.IndexFunc(newIpsEntries, func(e Entry) bool {
		return e.FullURL == reference
	})
	if indInNew != -1 {
		return &newIpsEntries[indInNew]
	}
	return nil
}

func findAllKeysContainingString(m map[string]interface{}, substring string) []string {
	var matchingKeys []string
	slower := strings.ToLower(substring)
	for key := range m {
		if strings.Contains(strings.ToLower(key), slower) {
			matchingKeys = append(matchingKeys, key)
		}
	}
	return matchingKeys
}

func removeDuplicates(entries []Entry) []Entry {
	encountered := map[string]bool{}
	var result []Entry
	for _, e := range entries {
		if !encountered[e.FullURL] {
			encountered[e.FullURL] = true
			result = append(result, e)
		}
	}
	return result
}

func (is *IpsService) MergeIPS(ctx context.Context, currentIpsBundle map[string]interface{}, newIpsBundle map[string]interface{}) (map[string]interface{}, error) {
	var currIPS, newIPS Bundle
	if err := mapstructure.Decode(currentIpsBundle, &currIPS); err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 400,
			Body:       []map[string]interface{}{{"error": "bad_request", "message": "Malformed current IPS"}},
			Err:        fmt.Errorf("malformed current IPS"),
		}
	}

	if err := mapstructure.Decode(newIpsBundle, &newIPS); err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 400,
			Body:       []map[string]interface{}{{"error": "bad_request", "message": "Malformed new IPS"}},
			Err:        fmt.Errorf("malformed new IPS"),
		}
	}

	curComp, err := getIPSComposition(currIPS.Entry)
	if err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 400,
			Body:       []map[string]interface{}{{"error": "bad_request", "message": "Current IPS does not have its composition"}},
			Err:        err,
		}
	}

	newComp, err := getIPSComposition(newIPS.Entry)
	if err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 400,
			Body:       []map[string]interface{}{{"error": "bad_request", "message": "Current IPS does not have its composition"}},
			Err:        err,
		}
	}

	// Merge composition for IPSs
	mergedComp := curComp
	for _, section := range newComp.Section {
		code := section.Code.Coding[0]["code"]
		if code == nil {
			continue
		}
		sectionIndex := slices.IndexFunc(mergedComp.Section, func(s Section) bool {
			return len(s.Code.Coding) > 0 && s.Code.Coding[0] != nil && s.Code.Coding[0]["code"] == code
		})

		if sectionIndex == -1 {
			// New IPS section is not present on current IPS
			mergedComp.Section = append(mergedComp.Section, section)
		} else {
			// Sections exists, add entries that do not exist in the current IPS
			for _, newEntry := range section.Entry {
				exists := false
				for _, oldEntry := range mergedComp.Section[sectionIndex].Entry {
					if newEntry["reference"] == oldEntry["reference"] {
						exists = true
						break
					}
				}
				if !exists {
					mergedComp.Section[sectionIndex].Entry = append(mergedComp.Section[sectionIndex].Entry, newEntry)
				}

			}
		}
	}

	fullURL := mergedComp.URL
	mergedComp.URL = ""

	jsonData, err := json.Marshal(mergedComp)
	if err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 500,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to convert composition to JSON"}},
			Err:        err,
		}
	}

	var mergedResource map[string]interface{}
	if err := json.Unmarshal(jsonData, &mergedResource); err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 500,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to convert composition JSON to a map"}},
			Err:        err,
		}
	}

	// Build the merge ips with the merge Composition
	mergedIPS := Bundle{
		ID:           uuid.NewString(),
		Identifier:   currIPS.Identifier,
		Meta:         currIPS.Meta,
		ResourceType: currIPS.ResourceType,
		Signature:    nil,
		Timestamp:    time.Now().UTC().String(),
		Type:         currIPS.Type,
		Entry:        []Entry{{FullURL: fullURL, Resource: mergedResource}},
	}
	for _, section := range mergedComp.Section {
		for _, secEntry := range section.Entry {
			newEntry := getEntry(secEntry["reference"].(string), currIPS.Entry, newIPS.Entry)

			if newEntry == nil {
				break
			}
			mergedIPS.Entry = append(mergedIPS.Entry, *newEntry)

			if newEntry.Resource == nil {
				break
			}
			// Check for any resource that contains more reference in its representation
			// If we find any reference we added it to the IPS
			rk := findAllKeysContainingString(newEntry.Resource, "reference")
			for _, k := range rk {
				v, ok := newEntry.Resource[k]
				if !ok {
					break
				}
				var ref Reference
				if err := mapstructure.Decode(v, &ref); err != nil {
					return nil, fmt.Errorf(`error decoding codeable reference: %v`, err)
				}
				if ref.Reference != "" {
					newEntry = getEntry(ref.Reference, currIPS.Entry, newIPS.Entry)
					mergedIPS.Entry = append(mergedIPS.Entry, *newEntry)
				}

			}
		}
	}
	mergedIPS.Entry = removeDuplicates(mergedIPS.Entry)

	jsonData, err = json.Marshal(mergedIPS)
	if err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 500,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to convert composition to JSON"}},
			Err:        err,
		}
	}

	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 500,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to convert composition JSON to a map"}},
			Err:        err,
		}
	}

	return data, nil
}
