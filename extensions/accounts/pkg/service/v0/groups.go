package service

import (
	"context"
	"path"
	"strconv"

	accountsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/accounts/v0"
	accountssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/accounts/v0"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/owncloud/ocis/extensions/accounts/pkg/storage"
	merrors "go-micro.dev/v4/errors"
	p "google.golang.org/protobuf/proto"
)

func (s Service) expandMembers(g *accountsmsg.Group) {
	if g == nil {
		return
	}
	expanded := []*accountsmsg.Account{}
	for i := range g.Members {
		// TODO resolve by name, when a create or update is issued they may not have an id? fall back to searching the group id in the index?
		a := &accountsmsg.Account{}
		if err := s.repo.LoadAccount(context.Background(), g.Members[i].Id, a); err == nil {
			expanded = append(expanded, a)
		} else {
			// log errors but con/var/tmp/ocis-accounts-store-408341811tinue execution for now
			s.log.Error().Err(err).Str("id", g.Members[i].Id).Msg("could not load account")
		}
	}
	g.Members = expanded
}

// deflateMembers replaces the users of a group with an instance that only contains the id
func (s Service) deflateMembers(g *accountsmsg.Group) {
	if g == nil {
		return
	}
	deflated := []*accountsmsg.Account{}
	for i := range g.Members {
		if g.Members[i].Id != "" {
			deflated = append(deflated, &accountsmsg.Account{Id: g.Members[i].Id})
		} else {
			// TODO fetch and use an id when group only has a name but no id
			s.log.Error().Str("id", g.Id).Interface("account", g.Members[i]).Msg("resolving members by name is not implemented yet")
		}
	}
	g.Members = deflated
}

// ListGroups implements the GroupsServiceHandler interface
func (s Service) ListGroups(ctx context.Context, in *accountssvc.ListGroupsRequest, out *accountssvc.ListGroupsResponse) (err error) {
	if in.Query == "" {
		err = s.repo.LoadGroups(ctx, &out.Groups)
		if err != nil {
			s.log.Err(err).Msg("failed to load all groups from storage")
			return merrors.InternalServerError(s.id, "failed to load all groups")
		}
		for i := range out.Groups {
			a := out.Groups[i]

			// TODO add accounts only if requested
			// if in.FieldMask ...
			s.expandMembers(a)

		}
		return nil
	}

	searchResults, err := s.findGroupsByQuery(ctx, in.Query)
	out.Groups = make([]*accountsmsg.Group, 0, len(searchResults))

	for _, hit := range searchResults {
		g := &accountsmsg.Group{}
		if err = s.repo.LoadGroup(ctx, hit, g); err != nil {
			s.log.Error().Err(err).Str("group", hit).Msg("could not load group, skipping")
			continue
		}
		s.log.Debug().Interface("group", g).Msg("found group")

		// TODO add accounts if requested
		// if in.FieldMask ...
		s.expandMembers(g)

		out.Groups = append(out.Groups, g)
	}

	return
}
func (s Service) findGroupsByQuery(ctx context.Context, query string) ([]string, error) {
	return s.index.Query(ctx, &accountsmsg.Group{}, query)
}

// GetGroup implements the GroupsServiceHandler interface
func (s Service) GetGroup(c context.Context, in *accountssvc.GetGroupRequest, out *accountsmsg.Group) (err error) {
	var id string
	if id, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up group id: %v", err.Error())
	}

	if err = s.repo.LoadGroup(c, id, out); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "group not found: %v", err.Error())
		}

		s.log.Error().Err(err).Str("id", id).Msg("could not load group")
		return merrors.InternalServerError(s.id, "could not load group: %v", err.Error())
	}
	s.log.Debug().Interface("group", out).Msg("found group")

	// TODO only add accounts if requested
	// if in.FieldMask ...
	s.expandMembers(out)

	return
}

// CreateGroup implements the GroupsServiceHandler interface
func (s Service) CreateGroup(c context.Context, in *accountssvc.CreateGroupRequest, out *accountsmsg.Group) (err error) {
	if in.Group == nil {
		return merrors.InternalServerError(s.id, "invalid group: empty")
	}
	p.Merge(out, in.Group)

	if out.Id == "" {
		out.Id = uuid.Must(uuid.NewV4()).String()
	}

	if _, err = cleanupID(out.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	s.deflateMembers(out)

	if err = s.repo.WriteGroup(c, out); err != nil {
		s.log.Error().Err(err).Interface("group", out).Msg("could not persist new group")
		return merrors.InternalServerError(s.id, "could not persist new group: %v", err.Error())
	}

	indexResults, err := s.index.Add(out)
	if err != nil {
		s.rollbackCreateGroup(c, out)
		return merrors.InternalServerError(s.id, "could not index new group: %v", err.Error())
	}

	for _, r := range indexResults {
		if r.Field == "GidNumber" {
			gid, err := strconv.Atoi(path.Base(r.Value))
			if err != nil {
				s.rollbackCreateGroup(c, out)
				return err
			}
			out.GidNumber = int64(gid)
			return s.repo.WriteGroup(context.Background(), out)
		}
	}

	return
}

// rollbackCreateGroup tries to rollback changes made by `CreateGroup` if parts of it failed.
func (s Service) rollbackCreateGroup(ctx context.Context, group *accountsmsg.Group) {
	err := s.index.Delete(group)
	if err != nil {
		s.log.Err(err).Msg("failed to rollback group from indices")
	}
	err = s.repo.DeleteGroup(ctx, group.Id)
	if err != nil {
		s.log.Err(err).Msg("failed to rollback group from repo")
	}
}

// UpdateGroup implements the GroupsServiceHandler interface
func (s Service) UpdateGroup(c context.Context, in *accountssvc.UpdateGroupRequest, out *accountsmsg.Group) (err error) {
	return merrors.InternalServerError(s.id, "not implemented")
}

// DeleteGroup implements the GroupsServiceHandler interface
func (s Service) DeleteGroup(c context.Context, in *accountssvc.DeleteGroupRequest, out *empty.Empty) (err error) {
	var id string
	if id, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up group id: %v", err.Error())
	}

	g := &accountsmsg.Group{}
	if err = s.repo.LoadGroup(c, id, g); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "group not found: %v", err.Error())
		}
		return merrors.InternalServerError(s.id, "could not load group: %v", err.Error())
	}

	// delete memberof relationship in users
	for i := range g.Members {
		err = s.RemoveMember(c, &accountssvc.RemoveMemberRequest{
			AccountId: g.Members[i].Id,
			GroupId:   id,
		}, g)
		if err != nil {
			s.log.Error().Err(err).Str("groupid", id).Str("accountid", g.Members[i].Id).Msg("could not remove account memberof, skipping")
		}
	}

	if err = s.repo.DeleteGroup(c, id); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "group not found: %v", err.Error())
		}

		return merrors.InternalServerError(s.id, "could not load group: %v", err.Error())
	}

	if err = s.index.Delete(g); err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("could not remove group from index")
		return merrors.InternalServerError(s.id, "could not remove group from index: %v", err.Error())
	}

	s.log.Info().Str("id", id).Msg("deleted group")
	return
}

// AddMember implements the GroupsServiceHandler interface
func (s Service) AddMember(c context.Context, in *accountssvc.AddMemberRequest, out *accountsmsg.Group) (err error) {
	// cleanup ids
	var groupID string
	if groupID, err = cleanupID(in.GroupId); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up group id: %v", err.Error())
	}

	var accountID string
	if accountID, err = cleanupID(in.AccountId); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	// load structs
	a := &accountsmsg.Account{}
	if err = s.repo.LoadAccount(c, accountID, a); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "group not found: %v", err.Error())
		}
		return merrors.InternalServerError(s.id, "could not load group: %v", err.Error())
	}

	g := &accountsmsg.Group{}
	if err = s.repo.LoadGroup(c, groupID, g); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "could not load group: %v", err.Error())
		}
		return merrors.InternalServerError(s.id, "could not load group: %v", err.Error())
	}

	// check if we need to add the account to the group
	alreadyRelated := false
	for i := range g.Members {
		if g.Members[i].Id == a.Id {
			alreadyRelated = true
		}
	}
	aref := &accountsmsg.Account{
		Id: a.Id,
	}
	if !alreadyRelated {
		g.Members = append(g.Members, aref)
	}

	// check if we need to add the group to the account
	alreadyRelated = false
	for i := range a.MemberOf {
		if a.MemberOf[i].Id == g.Id {
			alreadyRelated = true
			break
		}
	}
	// only store the reference to prevent recursion when marshaling json
	gref := &accountsmsg.Group{
		Id: g.Id,
	}
	if !alreadyRelated {
		a.MemberOf = append(a.MemberOf, gref)
	}

	if err = s.repo.WriteAccount(c, a); err != nil {
		return merrors.InternalServerError(s.id, "could not persist account: %v", err.Error())
	}
	if err = s.repo.WriteGroup(c, g); err != nil {
		return merrors.InternalServerError(s.id, "could not persist group: %v", err.Error())
	}
	// FIXME update index!
	// TODO rollback changes when only one of them failed?
	// TODO store relation in another file?
	// TODO return error if they are already related?
	return nil
}

// RemoveMember implements the GroupsServiceHandler interface
func (s Service) RemoveMember(c context.Context, in *accountssvc.RemoveMemberRequest, out *accountsmsg.Group) (err error) {

	// cleanup ids
	var groupID string
	if groupID, err = cleanupID(in.GroupId); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up group id: %v", err.Error())
	}

	var accountID string
	if accountID, err = cleanupID(in.AccountId); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	// load structs
	a := &accountsmsg.Account{}
	if err = s.repo.LoadAccount(c, accountID, a); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "could not load account: %v", err.Error())
		}
		s.log.Error().Err(err).Str("id", accountID).Msg("could not load account")
		return merrors.InternalServerError(s.id, "could not load account: %v", err.Error())
	}

	g := &accountsmsg.Group{}
	if err = s.repo.LoadGroup(c, groupID, g); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "could not load group: %v", err.Error())
		}
		s.log.Error().Err(err).Str("id", groupID).Msg("could not load group")
		return merrors.InternalServerError(s.id, "could not load group: %v", err.Error())
	}

	//remove the account from the group if it exists
	newMembers := []*accountsmsg.Account{}
	for i := range g.Members {
		if g.Members[i].Id != a.Id {
			newMembers = append(newMembers, g.Members[i])
		}
	}
	g.Members = newMembers

	// remove the group from the account if it exists
	newGroups := []*accountsmsg.Group{}
	for i := range a.MemberOf {
		if a.MemberOf[i].Id != g.Id {
			newGroups = append(newGroups, a.MemberOf[i])
		}
	}
	a.MemberOf = newGroups

	if err = s.repo.WriteAccount(c, a); err != nil {
		s.log.Error().Err(err).Interface("account", a).Msg("could not persist account")
		return merrors.InternalServerError(s.id, "could not persist account: %v", err.Error())
	}
	if err = s.repo.WriteGroup(c, g); err != nil {
		s.log.Error().Err(err).Interface("group", g).Msg("could not persist group")
		return merrors.InternalServerError(s.id, "could not persist group: %v", err.Error())
	}
	// FIXME update index!
	// TODO rollback changes when only one of them failed?
	// TODO store relation in another file?
	// TODO return error if they are not related?
	return nil
}

// ListMembers implements the GroupsServiceHandler interface
func (s Service) ListMembers(c context.Context, in *accountssvc.ListMembersRequest, out *accountssvc.ListMembersResponse) (err error) {
	// cleanup ids
	var groupID string
	if groupID, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up group id: %v", err.Error())
	}

	g := &accountsmsg.Group{}
	if err = s.repo.LoadGroup(c, groupID, g); err != nil {
		if storage.IsNotFoundErr(err) {
			return merrors.NotFound(s.id, "group not found: %v", err.Error())
		}
		s.log.Error().Err(err).Str("id", groupID).Msg("could not load group")
		return merrors.InternalServerError(s.id, "could not load group: %v", err.Error())
	}

	// TODO only expand accounts if requested
	// if in.FieldMask ...
	s.expandMembers(g)
	out.Members = g.Members
	return
}
