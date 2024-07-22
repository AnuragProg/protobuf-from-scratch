package types

type Type string

const (
	PROJECT_TYPE Type = "project"

	PROJECT_TYPE_NAME_FIELD_NO        uint64 = 0
	PROJECT_TYPE_DESCRIPTION_FIELD_NO uint64 = 1
	PROJECT_TYPE_TIMESTAMP_FIELD_NO   uint64 = 2
	PROJECT_TYPE_TAGS_FIELD_NO        uint64 = 3
)

type ProjectType struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Timestamp   uint64   `json:"timestamp"`
	Tags        []string `json:"tags"`
}
