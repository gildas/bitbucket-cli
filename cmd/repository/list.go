package repository

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all public repositories",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Role       *flags.EnumFlag
	Workspace  *flags.EnumFlag
	MainBranch string
	Project    string
	ProjectKey string
	Language   string
	HasIssues  bool
	HasWiki    bool
	IsPrivate  bool
	Columns    *flags.EnumSliceFlag
	SortBy     *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Role = flags.NewEnumFlag("all", "+owner", "admin", "contributor", "member")
	listOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns...)
	listOptions.SortBy = flags.NewEnumFlag(sortBy...)
	listCmd.Flags().Var(listOptions.Role, "role", "Role of the user in the repository (all, owner, admin, contributor, member), Default: owner")
	listCmd.Flags().Var(listOptions.Workspace, "workspace", "Workspace to list repositories from")
	listCmd.Flags().StringVar(&listOptions.Project, "project", "", "Project to list repositories from (optional)")
	listCmd.Flags().StringVar(&listOptions.ProjectKey, "project-key", "", "Project key to list repositories from (optional)")
	listCmd.Flags().BoolVar(&listOptions.HasIssues, "has-issues", false, "Filter repositories that have issues enabled (optional)")
	listCmd.Flags().BoolVar(&listOptions.HasWiki, "has-wiki", false, "Filter repositories that have wiki enabled (optional)")
	listCmd.Flags().BoolVar(&listOptions.IsPrivate, "is-private", false, "Filter repositories that are private (optional)")
	listCmd.Flags().StringVar(&listOptions.Language, "language", "", "Filter repositories by language (optional)")
	listCmd.Flags().StringVar(&listOptions.MainBranch, "main-branch", "", "Filter repositories by main branch name (optional)")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Workspace.CompletionFunc("workspace"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Role.CompletionFunc("role"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	query := url.Values{}
	wantFilter := cmd.Flags().Changed("role") || cmd.Flags().Changed("workspace") ||
		cmd.Flags().Changed("has-issues") || cmd.Flags().Changed("has-wiki") ||
		cmd.Flags().Changed("is-private") || cmd.Flags().Changed("language") ||
		cmd.Flags().Changed("main-branch") || cmd.Flags().Changed("project") ||
		cmd.Flags().Changed("project-key")

	if wantFilter {
		if listOptions.Role.Value == "all" {
			return errors.Errorf("You must specify one role when using filter flags (--project, --project-key, --main-branch, --language, --has-issues, --has-wiki, --is-private, --workspace)")
		}
		query.Add("role", listOptions.Role.Value)
		var filters []string

		if cmd.Flags().Changed("has-issues") {
			filters = append(filters, fmt.Sprintf("has_issues=%t", listOptions.HasIssues))
		}
		if cmd.Flags().Changed("has-wiki") {
			filters = append(filters, fmt.Sprintf("has_wiki=%t", listOptions.HasWiki))
		}
		if cmd.Flags().Changed("is-private") {
			filters = append(filters, fmt.Sprintf("is_private=%t", listOptions.IsPrivate))
		}
		if cmd.Flags().Changed("language") && len(listOptions.Language) > 0 {
			filters = append(filters, fmt.Sprintf("language=\"%s\"", listOptions.Language))
		}
		if cmd.Flags().Changed("main-branch") && len(listOptions.MainBranch) > 0 {
			filters = append(filters, fmt.Sprintf("mainbranch.name=\"%s\"", listOptions.MainBranch))
		}
		if cmd.Flags().Changed("project-key") && len(listOptions.ProjectKey) > 0 {
			filters = append(filters, fmt.Sprintf("project.key=\"%s\"", listOptions.ProjectKey))
		}
		if cmd.Flags().Changed("project") && len(listOptions.Project) > 0 {
			filters = append(filters, fmt.Sprintf("project.name=\"%s\"", listOptions.Project))
		}
		if len(filters) > 0 {
			query.Add("q", strings.Join(filters, " AND "))
		}
	} else if listOptions.Role.Value != "all" {
		query.Add("role", listOptions.Role.Value)
	}

	uripath := path.Join("/repositories", listOptions.Workspace.Value)
	if len(query) > 0 {
		uripath += "?" + query.Encode()
	}

	log.Infof("Listing all repositories, workspace %s, role %s", listOptions.Workspace, listOptions.Role)
	repositories, err := profile.GetAll[Repository](
		cmd.Context(),
		cmd,
		uripath,
	)
	if err != nil {
		return err
	}
	if len(repositories) == 0 {
		log.Infof("No repository found")
		return nil
	}
	core.Sort(repositories, func(a, b Repository) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	})
	return profile.Current.Print(cmd.Context(), cmd, Repositories(repositories))
}
