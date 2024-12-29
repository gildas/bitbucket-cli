# Bitbucket Command Line Interface

[bb](https://bitbucket.org/gildas_cherruel/bb) is the missing command line interface for Bitbucket.

The [Bitbucket Command Line Interface](https://bitbucket.org/gildas_cherruel/bb) brings the power of the Bitbucket platform to your command line. Creating and merging Pull Requests, cloning repositories, and more are now just a few keystrokes away.

## Installation

### Linux

You can grab the latest Debian/Ubuntu package on the [Downloads](https://bitbucket.org/gildas_cherruel/bb/downloads) pages.

If you use [Homebrew](https://brew.sh), you can install `bb` with:

```bash
brew install gildas/tap/bitbucket-cli
```

You can also install `bb` with snap:

```bash
sudo snap install bitbucket-cli
sudo snap alias bitbucket-cli bb
```

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/bitbucket-cli)

### macOS

You can get `bb` from [Homebrew](https://brew.sh) with:

```bash
brew install gildas/tap/bitbucket-cli
```

### Windows

You can get `bb` from [Chocolatey](https://chocolatey.org) with:

```bash
choco install bitbucket-cli
```

### Binaries

You can download the latest version of `bb` from the [downloads](https://bitbucket.org/gildas_cherruel/bb/downloads/) page.

Once you get the `bb` executable, you can install/copy it anywhere in your `$PATH`.

## Usage

`bb` is a modern command line interface. It uses subcommands to perform actions. You can get help on any subcommand by running `bb <subcommand> --help`.

General help is also available by running `bb --help` or `bb help`.

By default `bb` works in the current git repository. You can specify a Bitbucket repository with the `--repository` flag.

See the [Completion](#completion) section for more information about completion. Many commands and flags are dynamically auto-completed.

Most `delete`, `upload`, and `download` commands support multiple arguments. You can pass a list of arguments or a file with one argument per line:

```bash
bb repo delete myrepository1 myrepository2 myrepository3
```

You can tell `bb` to stop on the first error, warn on errorsm or ignore errors when processing multiple arguments with the `--stop-on-error`, `--warn-on-error`, or `--ignore-errors` flags.

All commands that would modify something on Bitbucket now allow you to preview the changes before applying them. You can use the `--dry-run` flag to see what would happen.

```bash
bb repo delete myrepository3 --dry-run
```

### Output

`bb` outputs a table by default and get be set per profile. You can also use the `--output` flag to change the output format manually. The following formats are supported:

- `csv`: CSV
- `json`: JSON
- `yaml`: YAML
- `tsv`: TSV
- `table`: Table

For example:

```bash
bb --output json workspace list
```

Or

```bash
bb workspace list --output json
```

You can also set the output format with the environment variable `BB_OUTPUT_FORMAT`:

```bash
export BB_OUTPUT_FORMAT=json
```

### Profiles

`bb` uses profiles to store your Bitbucket credentials. You can create a profile with the `bb profile create` command:

```bash
bb profile create \
  --name myprofile \
  --client-id <your-client-id> \
  --client-secret <your-client-secret>
```

You can also pass the `--default` flag to make this profile the default one, or pass a `--output` flag to change the profile output format.

You can also pass the `--default-workspace` and/or `--default-project` flags to set the default workspace and/or project for this profile.

You can also pass the `--progress` flag to display a progress bar when upload/downloading artifacts and attachments.

Profiles support the following authentications:

- [OAuth 2.0](https://support.atlassian.com/bitbucket-cloud/docs/use-oauth-on-bitbucket-cloud/) with the `--client-id` and `--client-secret` flags
- [App passwords](https://support.atlassian.com/bitbucket-cloud/docs/app-passwords/) with the `--user` and `--password` flags.
- [Repository Access Tokens](https://support.atlassian.com/bitbucket-cloud/docs/repository-access-tokens/), [Project Access Tokens](https://support.atlassian.com/bitbucket-cloud/docs/project-access-tokens/), [Workspace Access Tokens](https://support.atlassian.com/bitbucket-cloud/docs/workspace-access-tokens/) with the `--access-token` flags.

You can get the list of your profiles with the `bb profile list` command:

```bash
bb profile list
```

You can get the details of a profile with the `bb profile get` or `bb profile show` command:

```bash
bb profile get myprofile
```

You can ge the details of the current profile:

```bash
bb profile get --current
```

Or:

```bash
bb profile which
```

You can update a profile with the `bb profile update` command:

```bash
bb profile update myprofile \
  --client-id <your-client-id> \
  --client-secret <your-client-secret>
```

You can delete a profile with the `bb profile delete` command:

```bash
bb profile delete myprofile
```

You can set the default profile with the `bb profile use` command:

```bash
bb profile use myprofile
```

You can also set the profile with the environment variable `BB_PROFILE`:

```bash
export BB_PROFILE=myprofile
```

The profile can also come from your current `.git/config` file. You can set the `bb.profile` variable in the `[bitbucket "cli"]` section of your `.git/config` file:

```ini
[bitbucket "cli"]
  profile = myprofile
```

```bash
git config --local bitbucket.cli.profile myprofile
```

The current profile comes in order from:

- the `--profile` flag
- the `BB_PROFILE` environment variable
- the `profile` variable in the `[bitbucket "cli"]` section of your `.git/config` file,  
  if the profile does not exist, the command will print a warning and use the default profile
- the profile marked `default` in the configuration file
- the first profile in the configuration file

Profiles are stored in the configuration file. By default, the configuration file is located:

- on Linux: `$XDG_CONFIG_HOME/bitbucket/config-cli.json`, or `~/.config/bitbucket/config-cli.json`, then `~/.bitbucket-cli`
- on macOS: `$HOME/Library/Application Support/bitbucket/config-cli.json`, then `~/.bitbucket-cli`
- on Windows: `%AppData%\bitbucket\config-cli.json`, then `$HOME/.bitbucket-cli`
- on Plan 9: `$home/lib/bitbucket/config-cli.json`, then `~/.bitbucket-cli`

You can also override the location of the configuration file with the environment variable `BB_CONFIG` or the `--config` flag:

```bash
export BB_CONFIG=~/.bb/config.json
```

```bash
bb --config ~/.bb/config.json workspace list
```

### Workspaces

You can list workspaces with the `bb workspace list` command:

```bash
bb workspace list
```

With the `--membership` flag, you can see the kind of membership you have in each workspace:

```bash
bb workspace list --membership
```

You can also get the details of a workspace with the `bb workspace get` or `bb workspace show` command:

```bash
bb workspace get myworkspace
```

Additionally, you can get the members of a workspace with the flag `--members`:

```bash
bb workspace get myworkspace --members
```

Or, even, a specific member with the flag `--member`:

```bash
bb workspace get myworkspace --member mymember
```

### Projects

You can list projects with the `bb project list` command. If the `--workspace` flag is not provided, the default workspace of the profile is used (if the profile does not have a default workspace, the command will fail):

```bash
bb project list --workspace myworkspace
```

The `--workspace` flag is also dynamically auto-completed with the workspaces you have access to.

You can also get the details of a project with the `bb project get` or `bb project show` command:

```bash
bb project get myproject --workspace myworkspace
```

You can create a project with the `bb project create` command:

```bash
bb project create \
  --name myproject \
  --key MYPROJECT
```

You can update a project with the `bb project update` command:

```bash
bb project update myproject \
  --name myproject
```

You can delete a project with the `bb project delete` command:

```bash
bb project delete myproject
```

#### Project Default Reviewers

You can list the default reviewers of a project with the `bb project reviewer list` command. In addition to the `--workspace`, if the `--project` flag is not provided, the default project of the workspace is used (if the workspace does not have a default project, the command will fail):

```bash
bb project reviewer list --workspace myworkspace --project myproject
```

You can add a default reviewer to a project with the `bb project reviewer add` command:

```bash
bb project reviewer add userUUID
```

The `{}` around the `userUUID` are optional.

You can remove a default reviewer from a project with the `bb project reviewer remove` command:

```bash
bb project reviewer remove userUUID
```

 You can get the details of a default reviewer with the `bb project reviewer get` or `bb project reviewer show` command:

```bash
bb project reviewer get \
  --workspace myworkspace \
  --project myproject \
  userUUID
```

### Repositories

You can list repositories with the `bb repo list` command:

```bash
bb repo list --workspace myworkspace
```

If you do not provide a workspace, the command will attempt to list all repositories you have access to, which can take a very long time.

You can also get the details of a repository with the `bb repo get` or `bb repo show` command. If the `--workspace` flag is not provided, the default workspace of the profile is used (if the profile does not have a default workspace, the command will fail):

```bash
bb repo get --workspace myworkspace myrepository
```

You can clone a repository with the `bb repo clone` command:

```bash
bb repo clone myworkspace/myrepository
```

or, with the `--workspace` flag:

```bash
bb repo clone --workspace myworkspace myrepository
```

Or, using the profile's default workspace:

```bash
bb repo clone myrepository
```

By default, the repository is cloned in a folder with the same name as the repository. You can specify a different folder with the `--destination` flag:

```bash
bb repo clone --workspace myworkspace --destination myfolder myrepository
```

You can create a repository with the `bb repo create` command:

```bash
bb repo create myrepository_slug \
  --name      myrepository \
  --project   myproject \
  --workspace myworkspace
```

If the `--project` flag is not provided, the repository will be created in the default project of the profile.

You can update a repository with the `bb repo update` command:

```bash
bb repo update --workspace myworkspace myrepository \
  --private \
  --fork-policy no_public_forks
```

You can delete a repository with the `bb repo delete` command:

```bash
bb repo delete --workspace myworkspace myrepository
```

You can fork a repository with the `bb repo fork` command:

```bash
bb repo fork myrepository \
  --workspace myworkspace \
  --project   myproject \
  --name      myfork
```

You can list the forks of a repository with the `bb repo get --forks` command:

```bash
bb repo get myrepository \
  --workspace myworkspace \
  --forks
```

### Pull Requests

You can list pull requests with the `bb pullrequest list` command:

```bash
bb pullrequest list
```

You can create a pull request with the `bb pullrequest create` command:

```bash
bb pullrequest create \
  --title "My pull request" \
  --source "my-branch" \
  --destination "master"
```

You can get the details of a pull request with the `bb pullrequest get` or `bb pullrequest show` command:

```bash
bb pullrequest get 1
```

You can `approve`or `unapprove` a pull request with the `bb pullrequest approve` or `bb pullrequest unapprove` command:

```bash
bb pullrequest approve 1
```

You can `decline` a pull request with the `bb pullrequest decline` command:

```bash
bb pullrequest decline 1
```

You can `merge` a pull request with the `bb pullrequest merge` command:

```bash
bb pullrequest merge 1
```

If no pull request is provided, the command will try to merge the pull request with the current branch.

You can list the comments of a pull request with the `bb pullrequest comment list` command:

```bash
bb pullrequest comment list --pullrequest 1
```

You can add a comment to a pull request with the `bb pullrequest comment create` or `bb pullrequest comment add` command:

```bash
bb pullrequest comment add --pullrequest 1 \
  --comment "My comment" \
  --file    README.md \
  --line    404
```

You can resolve a comment with the `bb pullrequest comment resolve` command:

```bash
bb pullrequest comment resolve --pullrequest 1 452466
```

You can re-open a comment with the `bb pullrequest comment reopen` command:

```bash
bb pullrequest comment reopen --pullrequest 1 452466
```

You can get the details of a comment with the `bb pullrequest comment get` or `bb pullrequest comment show` command:

```bash
bb pullrequest comment get --pullrequest 1 452466
```

You can update a comment with the `bb pullrequest comment update` command:

```bash
bb pullrequest comment update --pullrequest 1 452466 \
  --comment "My comment"
```

You can delete a comment with the `bb pullrequest comment delete` command:

```bash
bb pullrequest comment delete --pullrequest 1 452466
```

### Issues

You can list issues with the `bb issue list` command:

```bash
bb issue list
```

By default, all open and new issues are listed. You can use the `--state` flag to filter the issues by state:

```bash
bb issue list --state open
```

The flag `--state` can be used multiple times to filter by multiple states:

```bash
bb issue list --state open --state new --state resolved,wontfix
```

You can create an issue with the `bb issue create` command:

```bash
bb issue create \
  --title "My issue" \
  --content "My issue content"
```

You can get the details of an issue with the `bb issue get` or `bb issue show` command:

```bash
bb issue get 1
```

You can update an issue with the `bb issue update` command:

```bash
bb issue update 1 \
  --title "My issue" \
  --content "My issue content"
```

You can delete an issue with the `bb issue delete` command:

```bash
bb issue delete 1
```

You can vote for an issue with the `bb issue vote` command:

```bash
bb issue vote 1
```

You can unvote for an issue with the `bb issue unvote` command:

```bash
bb issue unvote 1
```

You can watch an issue with the `bb issue watch` command:

```bash
bb issue watch 1
```

You can unwatch an issue with the `bb issue unwatch` command:

```bash
bb issue unwatch 1
```

You can add a comment to an issue with the `bb issue comment create` or `bb issue comment add` command:

```bash
bb issue comment add --issue 1 \
  --content "My comment"
```

You can get the details of a comment with the `bb issue comment get` or `bb issue comment show` command:

```bash
bb issue comment get --issue 1 7643545
```

You can update a comment with the `bb issue comment update` command:

```bash
bb issue comment update --issue 1 7643545 \
  --content "My comment"
```

You can delete a comment with the `bb issue comment delete` command:

```bash
bb issue comment delete --issue 1 7643545
```

You can list the attachments of an issue with the `bb issue attachment list` command:

```bash
bb issue attachment list --issue 1
```

You can upload an attachment to an issue with the `bb issue attachment upload` command:

```bash
bb issue attachment upload --issue 1 myattachment.zip
```

You can download an attachment with the `bb issue attachment download` command:

```bash
bb issue attachment download --issue 1 myattachment.zip
```

You can delete an attachment with the `bb issue attachment delete` command:

```bash
bb issue attachment delete --issue 1 myattachment.zip
```

### Artifacts (Downloads)

You can list artifacts with the `bb artifact list` command:

```bash
bb artifact list
```

By default the current repository is used, you can specify a repository with the `--repository` flag.

You can also upload an artifact with the `bb artifact upload` command:

```bash
bb artifact upload myartifact.zip
```

At the moment, only one file at a time is supported (no folders or stdin). The artifact name is the file name.

You can download an artifact with the `bb artifact download` command:

```bash
bb artifact download myartifact.zip
```

You can provide a `--destination` flag to specify the destination folder. If the folder does not exist, it will be created.

You can also pass the `--progress` flag to display a progress bar when upload/downloading artifacts. This override the default value set at the Profile level.

Finally, you can delete an artifact with the `bb artifact delete` command:

```bash
bb artifact delete myartifact.zip
```

### GPG Keys

You can list GPG keys with the `bb key list` command:

```bash
bb key list
```

By default, the keys are listed for the current user. You can specify a user with the `--user` flag.

You can also get the details of a GPG key with the `bb key get` or `bb key show` command:

```bash
bb key get <fingerprint>
```

By default, the key is retrieved for the current user. You can specify a user with the `--user` flag.

You can create a GPG key with the `bb key create` command:

```bash
bb key create \
  --user <user> \
  --key <key>
```

You can instead provide the key in a file with the `--key-file` flag. If the filename is `-`, the key is read from stdin.

You can delete a GPG key with the `bb key delete` command:

```bash
bb key delete <fingerprint>
```

### Cache

The bitbucket-cli caches some data to speed up the commands. The following items are cached:

- workspaces
- projects
- users

The cache is stored in the [os.UserCacheDir](https://pkg.go.dev/os#UserCacheDir) directory, under `bitbucket`. The values are stored for a duration of 5 minutes, you can override this value with the environment variable `BITBUCKET_CLI_CACHE_DURATION` (for the format please follow [core.ParseDuration](https://pkg.go.dev/github.com/gildas/go-core#ParseDuration)). By default, the items are stored as JSON files unencrypted. To encrypt these files, you can set the environment variable `BITBUCKET_CLI_CACHE_ENCRYPTIONKEY` with an AES-256 key. The key must follow the [crypto/aes](https://pkg.go.dev/crypto/aes) requirements.

You can clear the cache with the `bb cache clear` command:

```bash
bb cache clear
```

### Completion

`bb` supports completion for Bash, fish, Powershell, and zsh.

#### Bash

To enable completion, run the following command:

```bash
source <(bb completion bash)
```

You can also add this line to your `~/.bashrc` file to enable completion for every new shell.

```bash
bb completion bash > ~/.bashrc
```

#### fish

To enable completion, run the following command:

```bash
bb completion fish | source
```

You can also add this line to your `~/.config/fish/config.fish` file to enable completion for every new shell.

```bash
bb completion fish > ~/.config/fish/completions/bb.fish
```

#### Powershell

To enable completion, run the following command:

```pwsh
bb completion powershell | Out-String | Invoke-Expression
```

You can also add the output of the above command to your `$PROFILE` file to enable completion for every new shell.

#### zsh

To enable completion, run the following command:

```bash
source <(bb completion zsh)
```

You can also add this line to your functions folder to enable completion for every new shell.

```bash
bb completion zsh > "~/${fpath[1]}/_bb"
```

On macOS, you can add the completion to the brew functions:

```bash
bb completion zsh > "$(brew --prefix)/share/zsh/site-functions/_bb"
```

## TODO

We will add more commands in the future. If you have any suggestions, please open an issue.
