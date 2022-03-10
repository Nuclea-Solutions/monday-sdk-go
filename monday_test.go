package monday

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeColumnValues(t *testing.T) {
	client := MockClient{}
	board, _ := client.GetBoard(324)

	boardID, err := strconv.Atoi(board.ID)
	if err != nil {
		t.Fatalf("processing boardID: %s\n", err)
	}
	items, _ := client.GetItemsByColumnValues(boardID, "someid", "somevalue")
	newItems, err := decodeColumnValues(board.Columns, items)

	want := []Item{
		{
			ID:   "1234567890",
			Name: "Item1",
			ColumnValues: []ColumnValue{
				{
					ID:    "estado2",
					Title: "Estado",
					Type:  "color",
					Value: "Disponible",
				},
			},
			Assets: []Asset{
				{
					PublicURL: "https://source.unsplash.com/random/150x150",
				},
			},
		},
	}

	if !cmp.Equal(want, newItems[0]) {
		t.Fatalf("\nexpected:\n%s\ngot:\n%s\n", want, newItems[0])
	}
}
