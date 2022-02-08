package domain

import "errors"

var (
	ErrPrefixGroupNotFound = errors.New("the prefix group does not exist")
)

type PrefixGroupRepository interface {
	List() ([]PrefixGroup, error)
	Members(group PrefixGroup) ([]SteamUID, error)
	AddMember(group PrefixGroup, member SteamUID) error
	RemoveMember(group PrefixGroup, member SteamUID) error
}

type SteamUID string

func (u SteamUID) String() string {
	return string(u)
}

type PrefixGroup struct {
	Index  int
	Prefix string
}
