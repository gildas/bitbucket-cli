package tag

import (
	"context"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Tags []Tag

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (tags Tags) GetHeaders(cmd *cobra.Command) []string {
	return Tag{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (tags Tags) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(tags) {
		return []string{}
	}
	return tags[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (tags Tags) Size() int {
	return len(tags)
}

// GetTags gets the tags of a repository
func GetTags(context context.Context, cmd *cobra.Command) (tags []Tag, err error) {
	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return nil, err
	}

	return profile.GetAll[Tag](context, cmd, repository.GetPath("refs", "tags"))
}

// GetTagNames gets the tag names of a repository
func GetTagNames(context context.Context, cmd *cobra.Command, args []string, toComplete string) (names []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "gettags")
	log.Infof("Getting tags for profile %v", profile.Current)
	tags, err := GetTags(context, cmd)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, err
	}
	names = core.Map(tags, func(tag Tag) string { return tag.Name })
	core.Sort(names, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return common.FilterValidArgs(names, args, toComplete), nil
}
