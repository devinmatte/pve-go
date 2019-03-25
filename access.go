package main


// TODO:
//  - Add the rest of `/access`
//  - Add custom structs for domains, roles, and users
//  - Implement passing said struct to create/edit

import (
    "net/url"
    "github.com/mitchellh/mapstructure"
)

type Group struct {
    Name    string
    Comment string
    Members []string
}

type Domain struct {
    Comment string
    Digest  string
    Plugin  string
    Type    string
}

type Role struct {
    RoleID  string
    Privs   string
}

// Returns the valid subdirectories for `/access`. For example, because `roles`
// is returned, `/access/roles` is a valid API path
func (proxmox Proxmox) Access() ([]string, error) {
    data, err := proxmox.Get("/access")
    if err != nil {
        return nil, err
    }

    var subdirs []string
    dataArr := data.([]interface{})
    for _, element := range dataArr {
        elementMap := element.(map[string] interface{})
        subdirs = append(subdirs, elementMap["subdir"].(string))
    }

    return subdirs, nil
}

// Returns the authentication domain index. Literally just returns a list of
// Domain structs for each domain.
// NOTE: although no permissions are required to get the list of domains,
// either Realm.Allocate or Sys.Audit perms are required to get the details of
// a specific domain so the authenticated user must have one of those perms.
// TODO:
//  - Fix that permissions issue?
//    - Maybe just return a list of realm names?
func (proxmox Proxmox) GetAccessDomains() ([]*Domain, error) {
    var domains []*Domain
    data, err := proxmox.Get("/access/domains")
    if err != nil {
        return nil, err
    }
    dataArr := data.([]interface{})

    for _, element := range dataArr {
        elementMap := element.(map[string] interface{})
        domain, err := proxmox.GetAccessDomain(elementMap["realm"].(string))
        if err != nil {
            return nil, err
        }
        domains = append(domains, domain)
    }
    return domains, nil
}

// TODO:
//  - Pass in options as Domain struct instead of form
func (proxmox Proxmox) AddAccessDomain(domain Domain) error {
    return nil
}

// Gets the auth server configuration for the relevant domain
func (proxmox Proxmox) GetAccessDomain(name string) (*Domain, error) {
    var domain Domain
    data, err := proxmox.Get("/access/domains/" + name)
    if err != nil {
        return nil, err
    }
    data = data.(map[string]interface{})

    err = mapstructure.Decode(data, &domain)
    if err != nil {
        return nil, err
    }
    return &domain, nil
}

// TODO:
//  - Pass in Domain struct instead of form
func (proxmox Proxmox) EditAccessDomain (domain string,
    form url.Values) (map[string] interface{}, error) {
    data, err := proxmox.PostForm("/access/domains/" + domain, form)
    if err != nil {
        return nil, err
    }
    dataMap := data.(map[string]interface{})
    return dataMap, nil
}

// Untested
func (proxmox Proxmox) DeleteAccessDomain(domain string) error {
    _, err := proxmox.Delete("/access/domains/" + domain)
    if err != nil {
        return err
    }
    return nil
}

// Returns the group index (effectively a list of Group structs)
// NOTE: The available groups are restricted to groups where the authenticated
// user has User.Modify, Sys.Audit, or Group.Allocate permissions.
func (proxmox Proxmox) GetAccessGroups() ([]*Group, error) {
    data, err := proxmox.Get("/access/groups")
    if err != nil {
        return nil, err
    }
    dataMap := data.([]interface{})

    var groups []*Group
    for _, element := range dataMap {
        elementMap := element.(map[string] interface{})
        group, err := proxmox.GetAccessGroup(elementMap["groupid"].(string))
        if err != nil {
            return nil, err
        }
        groups = append(groups, group)
    }
    return groups, nil
}

// TODO
//  - Implement using Group struct instead of form
func (proxmox Proxmox) AddAccessGroup(group Group) error {
    return nil
}

// Returns an individual group configuration
func (proxmox Proxmox) GetAccessGroup(name string) (*Group, error) {
    var group Group

    data, err := proxmox.Get("/access/groups/" + name)
    if err != nil {
        return nil, err
    }
    data = data.(map[string] interface{})

    err = mapstructure.Decode(data, &group)
    if err != nil {
        return nil, err
    }
    group.Name = name

    return &group, nil
}

// TODO:
//  - Implementing passing Group struct
func (proxmox Proxmox) EditAccessGroup (group Group) error {
    return nil
}

// Untested
func (proxmox Proxmox) DeleteAccessGroup(group string) error {
    _, err := proxmox.Delete("/access/domains/" + group)
    if err != nil {
        return err
    }
    return nil
}

func (proxmox Proxmox) GetRoles() ([]*Role, error) {
    var roles []*Role
    data, err := proxmox.Get("/access/roles")
    if err != nil {
        return nil, err
    }
    dataArr := data.([]interface{})

    for _, element := range dataArr {
        elementMap := element.(map[string] interface{})
        role, err := proxmox.GetRole(elementMap["roleid"].(string))
        if err != nil {
            return nil, err
        }
        roles = append(roles, role)
    }
    return nil, nil
}

func (proxmox Proxmox) AddRole(role Role) error {
    return nil
}

func (proxmox Proxmox) GetRole(roleid string) (*Role, error) {
    var role Role
    data, err := proxmox.Get("/access/roles/" + roleid)
    if err != nil {
        return nil, err
    }
    data = data.(map[string] interface{})

    err = mapstructure.Decode(data, &role)
    if err != nil {
        return nil, err
    }
    return &role, nil
}

func (proxmox Proxmox) EditRole(role Role) error {
    return nil
}

func (proxmox Proxmox) DeleteRole(roleid string) error {
    _, err := proxmox.Delete("/access/roles/" + roleid)
    if err != nil {
        return err
    }
    return nil
}
