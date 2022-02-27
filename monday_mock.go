package monday

import "io"

// MockClient is used for testing purposes
type MockClient struct{}

func (c MockClient) GetBoards() ([]Board, error) {
	return []Board{
		Board{
			ID:   "1234567890",
			Name: "Board1",
		},
		Board{
			ID:   "1234567891",
			Name: "Board2",
		},
	}, nil
}

func (c MockClient) GetItemsByColumnValues(boardID int, columnID string, columnValue string) ([]Item, error) {
	return []Item{
		Item{
			ID:   "1234567890",
			Name: "Item1",
			ColumnValues: []ColumnValue{
				ColumnValue{
					ID:    "status",
					Title: "Estado",
					Type:  "color",
					Value: "{\"index\":1,\"post_id\":null,\"changed_at\":\"2018-07-30T06:27:05.982Z\"}",
				},
			},
		},
	}, nil
}

func (c MockClient) AddItem(boardID int, itemName string, columnValues map[string]interface{}) (string, error) {
	return "1234567890", nil
}

func (c MockClient) AddFileToColumn(itemID int, columnID string, fileName string, file io.Reader) (string, error) {
	return "1234567890", nil
}

func (c MockClient) ChangeMultipleColumnValues(boardID int, itemID int, columnValues map[string]interface{}) (string, error) {
	return "1234567890", nil
}

func (c MockClient) DeleteItem(itemID int) (string, error) {
	return "1234567890", nil
}
