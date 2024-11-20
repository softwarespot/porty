package ports

type SortBy int

const (
	SortByUsernameAppName SortBy = iota
	SortByUsername
	SortByAppName
	SortByPort
	SortByCreatedAt
	SortByUpdatedAt
	SortByAccessedAt
)

var SortByStrings = []string{
	"username-appname",
	"username",
	"appname",
	"port",
	"created-at",
	"updated-at",
	"accessed-at",
}

func (s SortBy) String() string {
	return SortByStrings[s]
}
