package monday

import "io"

// MockClient is used for testing purposes
type MockClient struct{}

func (c MockClient) GetBoard(boardID int) (Board, error) {
  return Board{
    ID:   "1234567890",
    Name: "Board1",
    Columns: {
      {
        ID: "estado2",
        Title: "Estado",
        Type:"color",
        Settings: "{\"done_colors\":[1],\"color_mapping\":{\"0\":1,\"1\":106,\"2\":0,\"6\":15,\"9\":2,\"15\":160,\"106\":6,\"108\":9,\"160\":108},\"labels\":{\"0\":\"Disponible\",\"1\":\"Escriturado\",\"2\":\"Bloqueado\",\"3\":\"Proceso de escrituración\",\"160\":\"Apartado\"},\"labels_positions_v2\":{\"0\":0,\"1\":4,\"2\":1,\"3\":3,\"5\":5,\"160\":2},\"labels_colors\":{\"0\":{\"color\":\"#00c875\",\"border\":\"#00B461\",\"var_name\":\"green-shadow\"},\"1\":{\"color\":\"#68a1bd\",\"border\":\"#68a1bd\",\"var_name\":\"river\"},\"2\":{\"color\":\"#fdab3d\",\"border\":\"#E99729\",\"var_name\":\"orange\"},\"3\":{\"color\":\"#0086c0\",\"border\":\"#3DB0DF\",\"var_name\":\"blue-links\"},\"160\":{\"color\":\"#4eccc6\",\"border\":\"#4eccc6\",\"var_name\":\"australia\"}}}"
      }
    }
  }
}

func (c MockClient) GetBoards() ([]Board, error) {
	return []Board{
		{
			ID:   "1234567890",
			Name: "Board1",
		},
		{
			ID:   "1234567891",
			Name: "Board2",
		},
	}, nil
}

func (c MockClient) GetItemsByColumnValues(boardID int, columnID string, columnValue string) ([]Item, error) {
	return {
		{
			ID:   "1234567890",
			Name: "Item1",
			ColumnValues: []ColumnValue{
				{
					ID:    "estado2",
					Title: "Estado",
					Type:  "color",
					Value: "{\"index\":0,\"post_id\":null,\"changed_at\":\"2018-07-30T06:27:05.982Z\"}",
				},
			},
      Assets: []Asset{
        {
          PublicURL: "https://source.unsplash.com/random/150x150",
        },
      },
		},
	}, nil
}

func (c MockClient) AddItem(boardID int, itemName string, columnValues map[string]interface{}) (string, error) {
	return "1234567890", nil
}

func (c MockClient) AddSubItem(parentItemID int, itemName string, columnValues map[string]interface{}) (string, error) {
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
