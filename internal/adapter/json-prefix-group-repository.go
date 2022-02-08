package adapter

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"lbmaster-advanced-groups-api/internal/domain"
)

type jsonPrefixGroupRepository struct {
	filename string
}

type lbmasterConfig struct {
	PrefixGroups []prefixGroup          `json:"prefixGroups"`
	Raw          map[string]interface{} `json:"-"`
}

type prefixGroup struct {
	Prefix  string   `json:"prefix"`
	Members []string `json:"members"`
}

func NewJsonPrefixGroupRepository(filename string) *jsonPrefixGroupRepository {
	return &jsonPrefixGroupRepository{filename: filename}
}

func (r jsonPrefixGroupRepository) config() (lbmasterConfig, error) {
	var config lbmasterConfig
	var others map[string]interface{}
	c, err := ioutil.ReadFile(r.filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(c, &config)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(c, &others)
	if err != nil {
		return config, err
	}
	config.Raw = others
	return config, nil
}

func (r jsonPrefixGroupRepository) List() ([]domain.PrefixGroup, error) {
	config, err := r.config()
	if err != nil {
		return []domain.PrefixGroup{}, err
	}

	var result []domain.PrefixGroup
	for i, group := range config.PrefixGroups {
		result = append(result, domain.PrefixGroup{
			Index:  i,
			Prefix: group.Prefix,
		})
	}
	return result, nil
}

func (r jsonPrefixGroupRepository) Members(group domain.PrefixGroup) ([]domain.SteamUID, error) {
	config, err := r.config()
	if err != nil {
		return []domain.SteamUID{}, err
	}

	for i, p := range config.PrefixGroups {
		if i == group.Index {
			var members []domain.SteamUID
			for _, member := range p.Members {
				members = append(members, domain.SteamUID(member))
			}
			return members, nil
		}
	}
	return []domain.SteamUID{}, domain.ErrPrefixGroupNotFound
}

func (r jsonPrefixGroupRepository) AddMember(group domain.PrefixGroup, member domain.SteamUID) error {
	return r.modifyPrefixGroup(group, func(members []interface{}) []interface{} {
		for _, m := range members {
			if m == member.String() {
				return members
			}
		}
		return append(members, member.String())
	})
}

type pgModifier func(members []interface{}) []interface{}

func (r jsonPrefixGroupRepository) modifyPrefixGroup(group domain.PrefixGroup, modifier ...pgModifier) error {
	config, err := r.config()
	if err != nil {
		return err
	}

	toSave := config.Raw
	groups := toSave["prefixGroups"]
	switch v := groups.(type) {
	case []interface{}:
		for i, p := range v {
			if i == group.Index {
				switch pg := p.(type) {
				case map[string]interface{}:
					switch m := pg["members"].(type) {
					case []interface{}:
						for _, mod := range modifier {
							pg["members"] = mod(m)
						}
					default:
						return errors.New("members has an unknown type")
					}

					c, err := json.MarshalIndent(toSave, "", "  ")
					if err != nil {
						return err
					}
					return ioutil.WriteFile(r.filename, c, 0644)
				default:
					return errors.New("prefixGroup has an unknown type")
				}
			}
		}
		break
	default:
		return errors.New("prefixGroups has an unknown type")
	}
	return domain.ErrPrefixGroupNotFound
}

func (r jsonPrefixGroupRepository) RemoveMember(group domain.PrefixGroup, member domain.SteamUID) error {
	return r.modifyPrefixGroup(group, func(members []interface{}) []interface{} {
		var result []interface{}
		for _, v := range members {
			if v != member.String() {
				result = append(result, v)
			}
		}
		return result
	})
}
