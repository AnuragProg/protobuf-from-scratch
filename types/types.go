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
	Name        string
	Description string
	Timestamp   uint64
	Tags        []string
}
