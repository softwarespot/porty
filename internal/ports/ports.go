package ports

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/softwarespot/porty/internal/helpers"
)

const (
	DefaultMinPort Port = 8000
	DefaultMaxPort Port = 9000
)

var defaultUserPort = UserPort{}

type Ports struct {
	MinPort   Port       `json:"minPort"`
	MaxPort   Port       `json:"maxPort"`
	UserPorts []UserPort `json:"userPorts"`
}

type UserPort struct {
	Username   string    `json:"usernmae"`
	AppName    string    `json:"appName"`
	Port       Port      `json:"port"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	AccessedAt time.Time `json:"accessedAt"`
}

func New() *Ports {
	return &Ports{
		MinPort:   DefaultMinPort,
		MaxPort:   DefaultMaxPort,
		UserPorts: []UserPort{},
	}
}

func (p *Ports) All(sortBy SortBy) []UserPort {
	ups := make([]UserPort, len(p.UserPorts))
	copy(ups, p.UserPorts)
	sort.SliceStable(ups, getSorterFunc(ups, sortBy))

	return ups
}

func (p *Ports) AllByUsername(username string, sortBy SortBy) ([]UserPort, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}

	ups := make([]UserPort, 0, len(p.UserPorts))
	for _, up := range p.UserPorts {
		if up.Username == username {
			ups = append(ups, up)
		}
	}
	if len(ups) == 0 {
		return nil, fmt.Errorf("no ports have been registered for the username %q", username)
	}

	sort.SliceStable(ups, getSorterFunc(ups, sortBy))

	return ups, nil
}

func (p *Ports) Register(username, appName string) (UserPort, error) {
	if err := validateUsernameAppName(username, appName); err != nil {
		return defaultUserPort, err
	}

	if idx := p.getIndex(username, appName); idx >= 0 {
		return defaultUserPort, fmt.Errorf("port exists for the username %q and app name %q", username, appName)
	}

	port, err := p.Next()
	if err != nil {
		return defaultUserPort, fmt.Errorf("unable to get the next port: %w", err)
	}

	now := time.Now()
	up := UserPort{
		Username:   username,
		AppName:    appName,
		Port:       port,
		CreatedAt:  now,
		UpdatedAt:  now,
		AccessedAt: now,
	}
	p.UserPorts = append(p.UserPorts, up)
	sort.SliceStable(p.UserPorts, getSorterFunc(p.UserPorts, SortByPort))

	return up, nil
}

func (p *Ports) Unregister(username, appName string) (UserPort, error) {
	if err := validateUsernameAppName(username, appName); err != nil {
		return defaultUserPort, err
	}

	idx := p.getIndex(username, appName)
	if idx == -1 {
		return defaultUserPort, fmt.Errorf("no port exists for the username %q and app name %q", username, appName)
	}

	up := p.UserPorts[idx]
	p.UserPorts = append(p.UserPorts[:idx], p.UserPorts[idx+1:]...)

	return up, nil
}

func (p *Ports) Get(username, appName string) (UserPort, error) {
	if err := validateUsernameAppName(username, appName); err != nil {
		return defaultUserPort, err
	}

	idx := p.getIndex(username, appName)
	if idx == -1 {
		return defaultUserPort, fmt.Errorf("no port exists for the username %q and app name %q", username, appName)
	}

	up := p.UserPorts[idx]
	up.AccessedAt = time.Now()
	p.UserPorts[idx] = up

	return up, nil
}

func (p *Ports) GetByPort(port Port) (UserPort, error) {
	for _, up := range p.UserPorts {
		if up.Port == port {
			return up, nil
		}
	}
	return defaultUserPort, fmt.Errorf("no port exists of %d", port)
}

func (p *Ports) Next() (Port, error) {
	idx := 0
	for port := p.MinPort; port <= p.MaxPort; port++ {
		if idx >= len(p.UserPorts) {
			return port, nil
		}
		if p.UserPorts[idx].Port != port {
			return port, nil
		}
		idx++
	}
	return 0, fmt.Errorf("all ports between %d-%d, have been registered", p.MinPort, p.MaxPort)
}

func (p *Ports) ToPort(s string) (Port, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid port number of %s was provided, expected a number: %w", s, err)
	}
	port := Port(i)
	if err := p.validatePort(port); err != nil {
		return 0, err
	}
	return port, nil
}

func (p *Ports) ToSortBy(s string) (SortBy, error) {
	if s == "" {
		return SortByUsernameAppName, nil
	}
	for i, v := range SortByStrings {
		if s == v {
			return SortBy(i), nil
		}
	}
	return SortByUsernameAppName, fmt.Errorf("invalid sort by string of %q was provided, expected one of the following: %s", s, helpers.ShellQuoteJoin(SortByStrings))
}

func (p *Ports) getIndex(username, appName string) int {
	for i, up := range p.UserPorts {
		if up.Username == username &&
			up.AppName == appName {
			return i
		}
	}
	return -1
}

func (p *Ports) validatePort(port Port) error {
	if port < p.MinPort || port > p.MaxPort {
		return fmt.Errorf(`invalid port number of %d was provided, expected value between %d-%d`, port, p.MinPort, p.MaxPort)
	}
	return nil
}

func validateUsernameAppName(username, appName string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if appName == "" {
		return errors.New("appName cannot be empty")
	}
	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	return nil
}

func getSorterFunc(ups []UserPort, sortBy SortBy) func(i, j int) bool {
	return func(i, j int) bool {
		switch sortBy {
		case SortByUsernameAppName:
			if ups[i].Username == ups[j].Username {
				return ups[i].AppName < ups[j].AppName
			}
			return ups[i].Username < ups[j].Username
		case SortByUsername:
			return ups[i].Username < ups[j].Username
		case SortByAppName:
			return ups[i].AppName < ups[j].AppName
		case SortByPort:
			return ups[i].Port < ups[j].Port
		case SortByCreatedAt:
			return ups[i].CreatedAt.After(ups[j].CreatedAt)
		case SortByUpdatedAt:
			return ups[i].UpdatedAt.After(ups[j].UpdatedAt)
		case SortByAccessedAt:
			return ups[i].AccessedAt.After(ups[j].AccessedAt)
		default:
			return true
		}
	}
}
