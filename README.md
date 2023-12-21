# Bitbucket Command Line Interface

[bb](https://bitbucket.org/gildas_cherruel/bb) is the missing command line interface for Bitbucket.

## Installation

You can download the latest version of `bb` from the [downloads](https://bitbucket.org/gildas_cherruel/bb/downloads/) page.

Once you get the `bb` executable, you can install it anywhere in your `$PATH`:

## Usage

`bb` is a modern command line interface. It uses subcommands to perform actions. You can get help on any subcommand by running `bb <subcommand> --help`.

General help is also available by running `bb --help` or `bb help`.

By default `bb` works in the current git repository. You can specify a Bitbucket repository with the `--repository` flag.

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

### Profiles

`bb` uses profiles to store your Bitbucket credentials. You can create a profile with the `bb profile create` command:

```bash
bb profile create \
  --name myprofile \
  --client-id <your-client-id> \
  --client-secret <your-client-secret>
```

You can also pass the `--default` flag to make this profile the default one, or pass a `--output` flag to change the profile output format.

Profiles support the following authentications:

- [OAuth 2.0](https://support.atlassian.com/bitbucket-cloud/docs/use-oauth-on-bitbucket-cloud/) with the `--client-id` and `--client-secret` flags
- [App passwords](https://support.atlassian.com/bitbucket-cloud/docs/app-passwords/) with the `--username` and `--password` flags.
- [Repository Access Tokens](https://support.atlassian.com/bitbucket-cloud/docs/repository-access-tokens/), [Project Access Tokens](https://support.atlassian.com/bitbucket-cloud/docs/project-access-tokens/), [Workspace Access Tokens](https://support.atlassian.com/bitbucket-cloud/docs/workspace-access-tokens/) with the `--access-token` flags.

You can get the list of your profiles with the `bb profile list` command:

```bash
bb profile list
```

You can get the details of a profile with the `bb profile get` or `bb profile show` command:

```bash
bb profile get myprofile
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

You can list projects with the `bb project list` command, the `--workspace` flag is required for all project commands:

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
  --key MYPROJECT \
  --workspace myworkspace
```

You can update a project with the `bb project update` command:

```bash
bb project update myproject \
  --name myproject \
  --workspace myworkspace
```

You can delete a project with the `bb project delete` command:

```bash
bb project delete myproject --workspace myworkspace
```

#### Project Default Reviewers

You can list the default reviewers of a project with the `bb project reviewer list` command:

```bash
bb project reviewer list --workspace myworkspace --project myproject
```

You can add a default reviewer to a project with the `bb project reviewer add` command:

```bash
bb project reviewer add \
  --workspace myworkspace \
  --project myproject \
  userUUID
```

You can remove a default reviewer from a project with the `bb project reviewer remove` command:

```bash
bb project reviewer remove \
  --workspace myworkspace \
  --project myproject \
  userUUID
```

 You can get the details of a default reviewer with the `bb project reviewer get` or `bb project reviewer show` command:

```bash
bb project reviewer get \
  --workspace myworkspace \
  --project myproject \
  userUUID
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

### Issues

You can list issues with the `bb issue list` command:

```bash
bb issue list
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

Finally, you can delete an artifact with the `bb artifact delete` command:

```bash
bb artifact delete myartifact.zip
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

### TODO

We will add more commands in the future. If you have any suggestions, please open an issue.
