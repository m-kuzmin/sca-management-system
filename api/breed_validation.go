package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

/*
ValidateCatBreed returns:

- true, nil - if the breed exists

- false, nil - the breed doesn't pass validation

- *, != nil - there was an error during an API call
*/
func ValidateCatBreed(id string) (bool, error) {
	// GET http://api.thecatapi.com/v1/breeds/{id}
	// - {id} is valid: 200, {"id": "{id}"}
	// - {id} not valid: 400, (ignore the body content)

	url := fmt.Sprintf("http://api.thecatapi.com/v1/breeds/%s", url.PathEscape(id))

	log.Printf("Making API call: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var breed struct {
			ID string `json:"id"`
		}

		err := json.NewDecoder(resp.Body).Decode(&breed)
		if err != nil {
			return false, err
		}

		if breed.ID == id {
			return true, nil
		} else {
			return false, nil
		}
	} else if resp.StatusCode == http.StatusBadRequest {
		return false, nil
	} else {
		return false, CatAPIBreedError{statusCode: resp.StatusCode}
	}
}

type CatAPIBreedError struct {
	statusCode int
}

func (e CatAPIBreedError) Error() string {
	return fmt.Sprintf("cat API error while calling /breed: status %d ", e.statusCode)
}
