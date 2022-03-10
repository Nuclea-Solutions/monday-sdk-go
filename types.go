package monday

// Board represents a Monday Board
type Board struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

// Group represents a Monday Group
type Group struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Column represents a Column in Monday
type Column struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`         // text, boolean, color, ...
	Settings string `json:"settings_str"` // used to get label index values for color(status) and dropdown column types
}

// ColumnMap provides easy access to a board's column info using column id. Key of map is column_id
type ColumnMap map[string]Column

// ColumnValue represents
type ColumnValue struct {
	ID    string `json:"id"` // column id
	Title string `json:"title"`
	Value string `json:"value"` // see func DecodeValue below
	Type  string `json:"type"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Asset
type Asset struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	OriginalGeometry string `json:"original_geometry"`
	CreatedAt        string `json:"created_at"`
	PublicURL        string `json:"public_url"`
	URL              string `json:"url"`
	URLThumbnail     string `json:"url_thumbnail"`
	FileExtension    string `json:"file_extension"`
	UploadedBy       User   `json:"uploaded_by"`
	FileSize         string `json:"file_size"`
}

// Item represents an Item in Monday
type Item struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	ColumnValues []ColumnValue `json:"column_values"`
	Assets       []Asset       `json:"assets"`
}

type ItemsByColumnValues struct {
	Items []Item `json:"items_by_column_values"`
}

// following types used to convert value from/to specific Monday value type
type DateTime struct {
	Date string `json:"date"`
	Time string `json:"time"`
}
type StatusIndex struct {
	Index int `json:"index"`
}
type PersonTeam struct {
	ID   int    `json:"id"`
	Kind string `json:"kind"` // "person" or "team"
}
type Checkbox struct {
	Checked string `json:"checked"`
}
