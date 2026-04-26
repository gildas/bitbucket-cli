package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/project"
	"bitbucket.org/gildas_cherruel/bb/cmd/project/reviewer"
	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type Repository struct {
	ID                   common.UUID         `json:"uuid"                  mapstructure:"uuid"`
	Name                 string              `json:"name,omitempty"                  mapstructure:"name"`
	FullName             string              `json:"full_name,omitempty"             mapstructure:"full_name"`
	Slug                 string              `json:"slug,omitempty"                  mapstructure:"slug"`
	Owner                user.User           `json:"owner,omitempty"                 mapstructure:"owner"`
	Workspace            workspace.Workspace `json:"workspace,omitempty"             mapstructure:"workspace"`
	Project              project.Project     `json:"project,omitempty"               mapstructure:"project"`
	HasIssues            bool                `json:"has_issues"            mapstructure:"has_issues"`
	HasWiki              bool                `json:"has_wiki"              mapstructure:"has_wiki"`
	IsPrivate            bool                `json:"is_private"            mapstructure:"is_private"`
	ForkPolicy           string              `json:"fork_policy,omitempty" mapstructure:"fork_policy"`
	Size                 int64               `json:"size,omitempty"                  mapstructure:"size"`
	Language             string              `json:"language,omitempty"    mapstructure:"language"`
	MainBranch           string              `json:"-"                     mapstructure:"-"`
	DefaultMergeStrategy string              `json:"-"                     mapstructure:"-"`
	BranchingModel       string              `json:"-"                     mapstructure:"-"`
	Parent               *Repository         `json:"parent,omitempty"      mapstructure:"parent"`
	Links                common.Links        `json:"links"                 mapstructure:"links"`
	CreatedOn            time.Time           `json:"created_on"            mapstructure:"created_on"`
	UpdatedOn            time.Time           `json:"updated_on"            mapstructure:"updated_on"`
}

/*
type repositorySettings struct {
	DefaultMergeStrategy bool `json:"default_merge_strategy" mapstructure:"default_merge_strategy"`
	BranchingModel       bool `json:"branching_model"        mapstructure:"branching_model"`
}
*/

type branch struct {
	Type string `json:"type" mapstructure:"type"`
	Name string `json:"name" mapstructure:"name"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:     "repo",
	Aliases: []string{"repository"},
	Short:   "Manage repositories",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Workspace requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

var columns = common.Columns[Repository]{
	{Name: "name", DefaultSorter: true, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	}},
	{Name: "full_name", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.FullName), strings.ToLower(b.FullName)) == -1
	}},
	{Name: "slug", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.Slug), strings.ToLower(b.Slug)) == -1
	}},
	{Name: "owner", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.Owner.Name), strings.ToLower(b.Owner.Name)) == -1
	}},
	{Name: "workspace", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.Workspace.Name), strings.ToLower(b.Workspace.Name)) == -1
	}},
	{Name: "project", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.Project.Name), strings.ToLower(b.Project.Name)) == -1
	}},
	{Name: "main_branch", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.MainBranch), strings.ToLower(b.MainBranch)) == -1
	}},
	{Name: "has_issues", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return a.HasIssues == b.HasIssues
	}},
	{Name: "has_wiki", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return a.HasWiki == b.HasWiki
	}},
	{Name: "is_private", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return a.IsPrivate == b.IsPrivate
	}},
	{Name: "fork_policy", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.ForkPolicy), strings.ToLower(b.ForkPolicy)) == -1
	}},
	{Name: "size", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return a.Size < b.Size
	}},
	{Name: "language", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.Language), strings.ToLower(b.Language)) == -1
	}},
	{Name: "default_merge_strategy", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.DefaultMergeStrategy), strings.ToLower(b.DefaultMergeStrategy)) == -1
	}},
	{Name: "branching_model", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.BranchingModel), strings.ToLower(b.BranchingModel)) == -1
	}},
	{Name: "parent", DefaultSorter: false, Compare: func(a, b Repository) bool {
		if a.Parent == nil && b.Parent == nil {
			return false
		} else if a.Parent == nil {
			return true
		} else if b.Parent == nil {
			return false
		}
		return strings.Compare(strings.ToLower(a.Parent.FullName), strings.ToLower(b.Parent.FullName)) == -1
	}},
	{Name: "created_on", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return a.CreatedOn.Before(b.CreatedOn)
	}},
	{Name: "updated_on", DefaultSorter: false, Compare: func(a, b Repository) bool {
		return a.UpdatedOn.Before(b.UpdatedOn)
	}},
}

var RepositoryCache = common.NewCache[Repository]()

// GetType gets the type of this repository
//
// implements core.TypeCarrier
func (repository Repository) GetType() string {
	return "repository"
}

// GetID gets the ID of the repository
//
// implements core.Identifiable
func (repository Repository) GetID() uuid.UUID {
	return uuid.UUID(repository.ID)
}

// GetName gets the name of the repository
//
// implements core.Named
func (repository Repository) GetName() string {
	return repository.Name
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (repository Repository) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"ID", "Name", "Full Name"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (repository Repository) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, repository.ID.String())
		case "name":
			row = append(row, repository.Name)
		case "full name":
			row = append(row, repository.FullName)
		case "slug":
			row = append(row, repository.Slug)
		case "owner":
			row = append(row, repository.Owner.Name)
		case "workspace":
			row = append(row, repository.Workspace.Name)
		case "project":
			row = append(row, repository.Project.Name)
		case "main branch":
			row = append(row, repository.MainBranch)
		case "issues", "has issues":
			row = append(row, strconv.FormatBool(repository.HasIssues))
		case "wiki", "has wiki":
			row = append(row, strconv.FormatBool(repository.HasWiki))
		case "is private":
			row = append(row, strconv.FormatBool(repository.IsPrivate))
		case "fork policy":
			row = append(row, repository.ForkPolicy)
		case "size":
			row = append(row, strconv.FormatInt(repository.Size, 10))
		case "language":
			row = append(row, repository.Language)
		case "default merge strategy":
			row = append(row, repository.DefaultMergeStrategy)
		case "branching model":
			row = append(row, repository.BranchingModel)
		case "parent":
			if repository.Parent != nil {
				row = append(row, repository.Parent.FullName)
			} else {
				row = append(row, " ")
			}
		case "created on", "created-on", "created_on", "created":
			row = append(row, repository.CreatedOn.Format("2006-01-02 15:04:05"))
		case "updated on", "updated-on", "updated_on", "updated":
			if !repository.UpdatedOn.IsZero() {
				row = append(row, repository.UpdatedOn.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, " ")
			}
		}
	}
	return row
}

// FetchWorkspace fetches the workspace of the repository
//
// Deprecated
func (repository *Repository) FetchWorkspaceX(context context.Context, cmd *cobra.Command, profile *profile.Profile) (*workspace.Workspace, error) {
	if repository == nil {
		return nil, errors.ArgumentMissing.With("repository")
	}
	if !repository.Workspace.ID.IsNil() {
		return &repository.Workspace, nil
	}
	workspacename := strings.Split(repository.FullName, "/")[0]
	workspace, err := workspace.GetWorkspaceBySlugOrID(context, cmd, workspacename)
	if err == nil {
		repository.Workspace = *workspace
	}
	return workspace, err
}

// GetRepositoryName gets the name of the repository from the command line or from the git config
func GetRepositoryName(context context.Context, cmd *cobra.Command) (repositoryName string, err error) {
	if cmd.Flag("repository") != nil {
		if repositoryName = cmd.Flag("repository").Value.String(); len(repositoryName) > 0 {
			return
		}
	}
	if remote, err := remote.GetRemote(context, cmd); err == nil {
		return remote.RepositoryName(), nil
	}
	return "", errors.ArgumentMissing.With("repository")
}

// GetRepository gets a repository by its slug
func GetRepository(ctx context.Context, cmd *cobra.Command) (repository *Repository, err error) {
	name, err := GetRepositoryName(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return GetRepositoryBySlugOrID(ctx, cmd, name)
}

// GetRepositoryBySlugOrID gets a repository by its slug name
//
// If the slug is in the format "workspace/repository", the workspace is used to get the repository.
//
// Otherwise, the workspace is determined by the git config or the default workspace in the profile.
func GetRepositoryBySlugOrID(ctx context.Context, cmd *cobra.Command, slugOrID string) (repository *Repository, err error) {
	log := logger.Must(logger.FromContext(ctx)).Child("repository", "get_by_slug_or_id", "repository", slugOrID)
	var ws *workspace.Workspace

	if components := strings.Split(slugOrID, "/"); len(components) == 2 {
		log.Debugf("Repository slug %s contains a workspace, extracting workspace and repository name", slugOrID)
		slugOrID = components[1]
		ws, err = workspace.GetWorkspaceBySlugOrID(ctx, cmd, components[0])
	} else {
		log.Debugf("Repository slug %s does not contain a workspace, using git config or default workspace", slugOrID)
		ws, err = workspace.GetWorkspace(ctx, cmd)
	}
	if err != nil {
		return nil, err
	}

	// In case we got a real UUID, get the Bitbucket UUID
	if id, err := common.ParseUUID(slugOrID); err == nil {
		slugOrID = id.String()
	}

	if repository, err = RepositoryCache.Get(fmt.Sprintf("%s/%s", ws.Slug, slugOrID)); err == nil {
		log.Debugf("Repository %s/%s found in cache", ws.Slug, slugOrID)
		return
	}

	log.Infof("Getting repository %s in workspace %s", slugOrID, ws.Slug)
	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return nil, err
	}

	err = profile.Get(
		ctx,
		cmd,
		fmt.Sprintf("/repositories/%s/%s", ws.Slug, slugOrID),
		&repository,
	)
	if err == nil {
		_ = RepositoryCache.Set(*repository, fmt.Sprintf("%s/%s", ws.Slug, slugOrID))
	}
	return
}

// GetForks gets the forks of the repository
func (repository Repository) GetForks(ctx context.Context, cmd *cobra.Command) (forks []Repository, err error) {
	log := logger.Must(logger.FromContext(ctx)).Child("repository", "forks")

	log.Infof("Getting forks of repository %s/%s", repository.Workspace.Slug, repository.Slug)
	return profile.GetAll[Repository](
		ctx,
		cmd,
		fmt.Sprintf("/repositories/%s/%s/forks", repository.Workspace.Slug, repository.Slug),
	)
}

// GetEffectiveDefaultReviewers gets the effective default reviewers for a repository
func (repository Repository) GetEffectiveDefaultReviewers(ctx context.Context, cmd *cobra.Command) (reviewers []reviewer.Reviewer, err error) {
	log := logger.Must(logger.FromContext(ctx)).Child("repository", "effective-default-reviewers")

	log.Infof("Getting effective default reviewers of repository %s/%s", repository.Workspace.Slug, repository.Slug)
	return profile.GetAll[reviewer.Reviewer](
		ctx,
		cmd,
		fmt.Sprintf("/repositories/%s/%s/effective-default-reviewers", repository.Workspace.Slug, repository.Slug),
	)
}

// GetRepositoryFromGit gets a repository from a git origin
func GetRepositoryFromGit(context context.Context, cmd *cobra.Command, profile *profile.Profile) (repository *Repository, err error) {
	log := logger.Must(logger.FromContext(context)).Child("repository", "fromgit")

	remote, err := remote.GetRemoteFromGitConfig(context, "origin")
	if err != nil {
		return nil, err
	}
	if repository, err = RepositoryCache.Get(remote.RepositoryName()); err == nil {
		log.Debugf("Repository %s found in cache", remote.RepositoryName())
		return
	}
	err = profile.Get(
		context,
		cmd,
		fmt.Sprintf("/repositories/%s", remote.RepositoryName()),
		&repository,
	)
	if err == nil {
		_ = RepositoryCache.Set(*repository, remote.RepositoryName())
	}
	return
}

// String returns the string representation of the repository
//
// implements fmt.Stringer
func (repository Repository) String() string {
	return repository.FullName
}

// Validate validates a Repository
func (repository *Repository) Validate() error {
	var merr errors.MultiError

	if repository.ID.IsNil() {
		merr.Append(errors.ArgumentMissing.With("uuid"))
	}
	if len(repository.Name) == 0 {
		merr.Append(errors.ArgumentMissing.With("name"))
	}
	if len(repository.FullName) == 0 {
		merr.Append(errors.ArgumentMissing.With("full_name"))
	}
	if len(repository.Slug) == 0 {
		repository.Slug = repository.Name
	}

	return merr.AsError()
}

// MarshalJSON implements the json.Marshaler interface.
//
// Implements json.Marshaler
func (repository Repository) MarshalJSON() (data []byte, err error) {
	type surrogate Repository
	var owner *user.User
	var wspace *workspace.Workspace
	var proj *project.Project
	var br *branch
	var createdOn string
	var updatedOn string
	var hasIssues *bool
	var hasWiki *bool
	var isPrivate *bool

	if !repository.Owner.ID.IsNil() {
		owner = &repository.Owner
	}
	if !repository.Workspace.ID.IsNil() {
		wspace = &repository.Workspace
	}
	if !repository.Project.ID.IsNil() {
		proj = &repository.Project
	}
	if len(repository.MainBranch) > 0 {
		br = &branch{Type: "branch", Name: repository.MainBranch}
	}
	if !repository.CreatedOn.IsZero() {
		createdOn = repository.CreatedOn.Format("2006-01-02T15:04:05.999999999-07:00")
	}
	if !repository.UpdatedOn.IsZero() {
		updatedOn = repository.UpdatedOn.Format("2006-01-02T15:04:05.999999999-07:00")
	}
	if repository.HasIssues {
		hasIssues = &repository.HasIssues
	}
	if repository.HasWiki {
		hasWiki = &repository.HasWiki
	}
	if repository.IsPrivate {
		isPrivate = &repository.IsPrivate
	}
	if repository.Slug == repository.Name {
		repository.Slug = ""
	}

	data, err = json.Marshal(struct {
		Type string `json:"type"`
		surrogate
		Owner      *user.User           `json:"owner,omitempty"`
		Workspace  *workspace.Workspace `json:"workspace,omitempty"`
		Project    *project.Project     `json:"project,omitempty"`
		MainBranch *branch              `json:"mainbranch,omitempty"`
		CreatedOn  string               `json:"created_on,omitempty"`
		UpdatedOn  string               `json:"updated_on,omitempty"`
		HasIssues  *bool                `json:"has_issues,omitempty"`
		HasWiki    *bool                `json:"has_wiki,omitempty"`
		IsPrivate  *bool                `json:"is_private,omitempty"`
	}{
		Type:       repository.GetType(),
		surrogate:  surrogate(repository),
		Owner:      owner,
		Workspace:  wspace,
		Project:    proj,
		MainBranch: br,
		CreatedOn:  createdOn,
		UpdatedOn:  updatedOn,
		HasIssues:  hasIssues,
		HasWiki:    hasWiki,
		IsPrivate:  isPrivate,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
//
// Implements json.Unmarshaler
func (repository *Repository) UnmarshalJSON(data []byte) (err error) {
	type surrogate Repository
	var inner struct {
		Type string `json:"type"`
		surrogate
		MainBranch branch `json:"mainbranch"`
	}
	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	if inner.Type != repository.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, repository.GetType()))
	}
	*repository = Repository(inner.surrogate)
	repository.MainBranch = inner.MainBranch.Name
	return errors.JSONUnmarshalError.Wrap(repository.Validate())
}
