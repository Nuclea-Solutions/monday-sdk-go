// Package monday implements all the communication to its v2 API. Check https://go.dev/play/p/aZ7tgaqFxWP
package monday

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Nuclea-Solutions/graphql"
)

// Service represents the Monday Service interface
type Service interface {
	GetBoards() ([]Board, error)
	GetItemsByColumnValues(boardID int, columnID string, columnValue string) ([]Item, error)
	AddItem(boardID int, itemName string, columnValues map[string]interface{}) (string, error)
	AddSubItem(parentItemID int, itemName string, columnValues map[string]interface{}) (string, error)
	AddFileToColumn(itemID int, columnID string, fileName string, file io.Reader) (string, error)
	ChangeMultipleColumnValues(boardID int, itemID int, columnValues map[string]interface{}) (string, error)
	DeleteItem(itemID int) (string, error)
}

// Client represents the Monday client
type Client struct {
	client      *graphql.Client
	filesClient *graphql.Client
}

// NewClient creates a graphql client (safe to share across requests)
func NewClient() Client {
	return Client{
		client:      graphql.NewClient("https://api.monday.com/v2/"),
		filesClient: graphql.NewClient("https://api.monday.com/v2/file"),
	}
}

// GetBoards returns []Board for all boards.
func (c Client) GetBoards() ([]Board, error) {
	req := graphql.NewRequest(`
	    query {
            boards {
                id name
            }
        }
    `)
	var response struct {
		Boards []Board `json:"boards"`
	}
	err := c.runRequest(req, &response)
	return response.Boards, err
}

// ItemsByColumnValues returns []Item filtered by a the value of a column
func (c Client) GetItemsByColumnValues(boardID int, columnID string, columnValue string) ([]Item, error) {
	req := graphql.NewRequest(`
    query ($boardID: Int!, $columnID: String!, $columnValue: String!) {
      items_by_column_values (board_id: $boardID, column_id: $columnID, column_value: $columnValue) {
        id
        name
        column_values {
          id
          value
          type
        }
      }
    }
    `)
	req.Var("boardID", boardID)
	req.Var("columnID", columnID)
	req.Var("columnValue", columnValue)
	var response ItemsByColumnValues
	err := c.runRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}

// Example of creating columnValues for AddItem
// map entry key is column id; run GetColumns to get column id's
/*
	columnValues := map[string]interface{}{
		"text":   "have a nice day",
		"date":   monday.BuildDate("2019-05-22"),
		"status": monday.BuildStatusIndex(2),
		"people": monday.BuildPeople(123456, 987654),   // parameters are user ids
	}
*/

// AddItem adds 1 item to specified board/group. The id of the added item is returned.
func (c Client) AddItem(boardID int, itemName string, columnValues map[string]interface{}) (string, error) {
	var req *graphql.Request
	if columnValues == nil {
		req = graphql.NewRequest(`
          mutation ($boardID: Int!, $itemName: String!) {
              create_item (board_id: $boardID, item_name: $itemName) {
                  id
              }
          }
      `)
	} else {
		req = graphql.NewRequest(`
        mutation ($boardID: Int!, $itemName: String!, $colValues: JSON!) {
            create_item (board_id: $boardID, item_name: $itemName, column_values: $colValues ) {
                id
            }
        }
    `)
		jsonValues, err := json.Marshal(&columnValues)
		if err != nil {
			log.Println("Processing AddItem Marshaling json")
			return "", err
		}
		// log.Printf("%s\n", string(jsonValues))
		req.Var("colValues", string(jsonValues))
	}
	// log.Printf("boardID: %d", boardID)
	// log.Println(string(jsonValues))

	req.Var("boardID", boardID)
	req.Var("itemName", itemName)

	var response struct {
		CreateItem Item `json:"create_item"`
	}
	err := c.runRequest(req, &response)
	log.Println(response)
	return response.CreateItem.ID, err
}

// AddItemUpdate adds an update entry to specified item.
func (c Client) AddItemUpdate(itemID string, msg string) error {
	intItemID, err := strconv.Atoi(itemID)
	if err != nil {
		log.Println("AddItemUpdate - bad itemID", err)
		return err
	}
	req := graphql.NewRequest(`
		mutation ($itemID: Int!, $body: String!) {
			create_update (item_id: $itemID, body: $body ) {
				id
			}
		}
	`)
	req.Var("itemID", intItemID)
	req.Var("body", msg)

	type UpdateID struct {
		ID string `json:"id"`
	}
	var response struct {
		CreateUpdate UpdateID `json:"create_update"`
	}
	err = c.runRequest(req, &response)
	return err
}

// GetItems returns []Item for all items in specified board.
func (c Client) GetItems(boardID int) ([]Item, error) {
	req := graphql.NewRequest(`	
		query ($boardID: [Int]) {
			boards (ids: $boardID){
				# items (limit: 10) {
				items () {
					id
					group {	id }
					name
					# column_values (ids: ["text", "status", "check"]) {  -- to get specific columns  
					column_values { 
						id value
					}
				}	
			}
		}	
	`)
	req.Var("boardID", []int{boardID})

	type group struct {
		ID string `json:"id"`
	}
	type itemData struct {
		ID           string        `json:"id"`
		Group        group         `json:"group"`
		Name         string        `json:"name"`
		ColumnValues []ColumnValue `json:"column_values"`
	}
	type boardItems struct {
		Items []itemData `json:"items"`
	}
	var response struct {
		Boards []boardItems `json:"boards"`
	}
	items := make([]Item, 0, 1000)
	err := c.runRequest(req, &response)
	if err != nil {
		return items, err
	}
	var responseItems []itemData = response.Boards[0].Items
	for _, responseItem := range responseItems {
		items = append(items, Item{
			ID:           responseItem.ID,
			Name:         responseItem.Name,
			ColumnValues: responseItem.ColumnValues,
		})
	}
	return items, nil
}

// ChangeMultipleColumnValues update the values of the columns specified in columnValues
// It can be used to update all types of columns, e.g.:
//   Connect: columnValues := map[string]interface{"connect_boards2": {"item_ids" : []int{12345, 23456, 34567}}}
//   Country: columnValues := map[string]interface{"country_2": {"countryCode": "MX", "countryName": "Mexico"}}
func (c Client) ChangeMultipleColumnValues(boardID int, itemID int, columnValues map[string]interface{}) (string, error) {
	req := graphql.NewRequest(`
    mutation ($boardID:Int!, $itemID:Int!, $columnValues:JSON!)
      {
        change_multiple_column_values(item_id:$itemID, board_id:$boardID, column_values: $columnValues)
          { id }
      }`)
	req.Var("boardID", boardID)
	req.Var("itemID", itemID)
	req.Var("columnValues", columnValues)

	var response struct {
		ChangeMultipleColumnValues Item `json:"change_multiple_column_values"`
	}

	err := c.runRequest(req, &response)
	return response.ChangeMultipleColumnValues.ID, err
}

// DeleteItem deletes an Item from Monday
func (c Client) DeleteItem(itemID int) (string, error) {
	req := graphql.NewRequest(`
    mutation($itemID: Int!) {
      delete_item(item_id: $itemID) {
        id
      }
    }
  `)

	req.Var("itemID", itemID)

	var response struct {
		DeleteItem Item `json:"delete_item"`
	}

	err := c.runRequest(req, &response)
	return response.DeleteItem.ID, err

}

// AddSubitem adds a subitem to a board on Monday
func (c Client) AddSubItem(parentItemID int, itemName string, columnValues map[string]interface{}) (string, error) {
	req := graphql.NewRequest(`
    mutation ($parentItemID: Int!, $itemName: String!, $columnValues: JSON!) {
      create_subitem(parent_item_id: $parentItemID, item_name: $itemName, column_values: $columnValues) {
        id
      }
    }
  `)

	req.Var("parentItemID", parentItemID)
	req.Var("itemName", itemName)

	jsonValues, err := json.Marshal(&columnValues)
	if err != nil {
		return "", err
	}

	req.Var("columnValues", string(jsonValues))

	var response struct {
		CreateSubItem Item `json:"create_subitem"`
	}

	err = c.runRequest(req, &response)

	return response.CreateSubItem.ID, err
}

// AddFileToColumn adds a file to a File column on Monday
func (c Client) AddFileToColumn(itemID int, columnID string, fileName string, file io.Reader) (string, error) {
	req := graphql.NewRequest(`
    mutation ($itemID: Int!, $columnID: String!, $file: File!) {
      add_file_to_column(item_id: $itemID, column_id: $columnID, file: $file) {
          id
      }
    }
  `)

	req.Var("itemID", itemID)
	req.Var("columnID", columnID)
	req.File("image", fileName, file)

	type AssetID struct {
		ID string `json:"id"`
	}
	var response struct {
		AddFileToColumn AssetID `json:"add_file_to_column"`
	}

	if err := c.runRequestWithFile(req, &response); err != nil {
		return "", err
	}

	return response.AddFileToColumn.ID, nil
}

func (c Client) runRequest(req *graphql.Request, response interface{}) error {
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", os.Getenv("MONDAY_API_TOKEN"))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.Background()
	err := c.client.Run(ctx, req, response)
	return err
}

func (c Client) runRequestWithFile(req *graphql.Request, response interface{}) error {
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", os.Getenv("MONDAY_API_TOKEN"))
	req.Header.Set("Content-Type", "multipart/form-data")
	ctx := context.Background()
	// Specify graphql package we are using MultiPartForm
	graphql.UseMultipartForm()(c.filesClient)
	err := c.filesClient.Run(ctx, req, response)
	return err
}
