0.16.1
=============
2025-12-17

* Issue #25: Added Support for git worktrees (7a781d8)

0.16.0
=============
2025-08-15

* Issue #22: Profile cannot execute Bitbucket API endpoint with queries (b1e605b)
* Issue #21: Added --sort and --columns to get/list commands
* Issue #20: Moved error flags to persistent flags (111dd05)

0.15.0
=============
2025-08-04

* Issue #16: Documentation: Added output sample (6e54fa8e)
* Issue #15: Fixed Documentation about output (ad79f303)
* Issue #12: Adding API Root (a4ba3bf2)
* Issue #6: Removed the vault key on Windows (edc74848)
* Issue #6: Removed the vault key on Windows and cleaned up (a8687d65)
* More Windows Makefile compatibility (08260514)
* Issue #13: Use filepath to back-recurse folders (9381dac1)
* Issue #8: Proper Windows shell detection (c050145a)
* Issue #8: fork git for private repositories via https (8d5bf387)
* Issue #8: new options (b13a3388)
* Issue #6: Use the profile vault key when loading the profiles (55186967)
* logging (57b43e5d)
* Issue #8: The default vault key for cloning should be the vault key (e024b7c1)
* Issue #8: Use the default user/vault-key if the clone user/vault-key is not present when cloning repositories (fa3035cd)
* Removed the move warning in the documentation (82541e01)
* Update documentation (87041104)
* Profile should redact its sensitive data in logs (ff8e4511)
* Issue #6: Store and Load secrets to the keyring (7da75cd0)
* get/list profiles should display validated profiles (ee974dc0)
* Issue #8: Added clone protocol, clone vault key and username for private repositories (1e9dbfed)
* Issue #3 and Issue #10: Get the workspace from the repository and add default reviewers (b7d957a3)
* Issue #8: added GetDefaultReviewers func (14968429)
* Issue #3: use the default workspace if we fail to get the workspace from the flags or git (074f28bc)
* Issue #3: Added missing remote origin url format (f024a40a)
* Issue #8: Added --username flag for URL authentication (33301fdc)
* Issue #9: set default folder to repo slug (467ddb38)
* Make will now install in your $GOPATH/bin folder first (eaaae7a7)
* Issue #5: Added --project, --project-key, --is-private, --has-issues, --has-wiki, --langugage, and --main-branch flags (b8ee0565)
* Makefile: changed the bin structure when building (0ba21360)
* Makefile: use gh to upload artifacts (df592a4e)
* Moving code to github (9acb472e)
* Moved the project to github (974f153a)
* Added auto-generated changelog (8a1701c0)

0.14.0
=============
2025-05-02

* Added Activity Comment (a31cceaa)
* Issue #19: Sort list command outputs (95df7350)
* Issue #18: Sort completion responses and match them to the toComplete pattern (9178d936)
* Added activity subcommand to fetch PR status (to check if a PR has been approved) (8a9b0dda)

0.13.0
=============
2025-03-22

* Added SSH Key support (a0eb12f3)
* Renamed key command (3049d76c)
* Issue #17: Added Documentation (8cdff8a8)
* Issue #17: Added Authorization Code Grant (1270e6b9)
* Added a Verbose function (c1d85ba3)

0.12.1
=============
2025-03-20

* Issue #16: Added --dry-run and --*error to key commands (4f0427d0)
* Issue #15: delete key was missing owner id (ca6be401)
* removed superfluous String() (25cb133a)
* Issue #14: Added --name option to key create command (6cb210c3)
* Cache Workspace (de849cac)
* Better error management (39bdf4b5)

0.12.0
=============
2024-12-29

* Issue #12: bb pr merge tries to find the PR# if no arg is given (810087b2)
* Removed obsolete flag (e9e842a3)
* Issue #13: Added some cache (1456f442)
* Issue #11: Added GPG Key management (8f2bac61)
* Unified user.Account and user.User (a4310365)
* Updated to Go 1.23 (aa4b9534)

0.11.0
=============
2024-10-20

* Added package name (b823d167)
* go-flags does not have a Contains func anymore (825afdbc)
* Profile for a command should be computed from the current command arguments first (d8aa0421)
* Root command requires subcommands (c77717f7)
* Source/Destination flags for pullrequest creation shoul support completion (dfb47b0b)
* Issue #10: reviewer flag should accept completion (12db2eca)
* [FIX] Fix documentation for `profile create` (9cba8a52)
* cmd can be nil (c1f3e3c5)
* Issue #10: reviewer flag accepts user's uuid, accountId, Nickname, or name (2c8ed5e4)
* Remote tools (0e6c0066)
* Repository tools (f9a47ea1)
* Workspace tools (e00682f0)
* Issue #9: workspace get should complete with slugs (f4638503)
* New Error JSON format (63bba299)
* Issue #8: Create should implement reviewer argument (08d1a174)

0.10.2
=============
2024-06-14

* Issue #7: moved WhatIf to the common package (ee32fe2d)
* wrong variable (b8a64e91)
* Makefile: new publish rules for snapcraft (4f94a872)

0.10.1
=============
2024-02-12

* Makefile: more archive formats (d9cb914b)
* Issue #6: Display a progress bar when uploading/downloading (59d4a38f)
* Makefile: archive name updates (d6fbcc35)

0.10.0
=============
2024-02-03

* Use the new gildas/go-flags package (72b545cb)
* issue list command should support multiple --state (4ca85e3d)
* New short command to tell the current profile (2254d1bd)
* Issue #5: Added documentation (ef6c7dff)
* The Chocolatey package for approved (77636eee)
* Issue #5: fetch the profile from .git/config if any (5f2bcd41)
* Package: updated chocolatey checksum (f69427e8)
* Added a changelog (5d815278)
* Makefile: Snapcraft packaging (f8fe5f9f)
* Documentation: new commands (d5d2329f)
*  (ace10dbf)

0.9.0
=============
2024-01-15

* Added PullRequest Comments (9968db68)
* Moved the inline anchor to the common objects (097c9400)
* Documentation: removed obsolete information (85046467)
* Makefile: install rule to install on the current machine (0cbd7812)
* Makefile: Better OS and ARCH detection (06dfe97b)
* Print func should get the output format from the command then the profile (8a8fe692)
* Issue #3: Added multi-positional arguments with error processing management (c9edbd4d)
* Issue #4: All commands honor --dry-run (f706533c)
* Makefile: use nfpm to build Debian and RPM packages (8bcd2b17)
* Makefile: better package naming convention (ff263344)
* Package: set the version in snap (abdf039b)
* Package: add a note about non-affiliation wit Atlassian in Snap (02895d6d)
* Chocolatey fixes for review process (1f168f12)
* Documentation: chocolatey installation (2dcee864)
* Added chocolatey package for Windows (fe188212)
* Makefile: added Windows arm64 (09f3a5df)
* Moved distribution packages together (42d39dc3)
* Package: deleted unused chocolatey files (eb840a18)
* Package: added chocolatey skeleton (88d362e3)
* Moved distribution packages together (9da5627f)
* Documentation: more on installation (fd835eb3)
* Installation: removed Docker targets (3006a288)
* Installation: makefile rules for Debian packages (673f3c26)
* Documentation: Added snapcraft badge (cdc23c34)
* Installation: more snap information (9af32acc)

0.8.0
=============
2023-12-30

* Installation: Added snap (3cccdedc)
* Use the latest Go 1.21 (f25f5278)
* role flag should register its completion (c38d15cc)
* Issue list can be filtered by state (43aae4ff)
* Get repository's fullname (ade75395)

0.7.0
=============
2023-12-29

* Profiles can have a default workspace and/or project (b9bff9ca)
* Remote Flag's AllowedFunc should require the positional arguments (40f42042)
* Documentation: repositories (27b739bc)
* typo (d3f4329d)
* help should show information about their positional arguments (82840ee1)
* Added repo(sitory) command (42c09a6a)
* Fixed omit empty fields in JSON (b7772722)
* New links (61b7fff2)
* Added references (15d053a2)
* Added references (7bcdb868)
* More Project fetch funcs (76985ad4)
* Issue #1: Find git config in the parent folders (8d22f3f5)
* Issue #2: use environment variables to set the profile, etc (caae620f)
* Get the current profile (bc4a28f3)

0.6.0
=============
2023-12-21

* Makefile: use bb to upload artifacts (5e194fb9)
* Added new command aliases (7427b819)
* Added issue attachment command (de51fb27)
* use core.Map (c480aa8a)
* Documentation: issues and comments (393ddeb2)
* Added new issue comment commands (b1837e92)
* Added new command aliases (ca1a9ec3)
* List issue comments should not show "Updated On" column is none of the comments have been updated (d9b63d9e)
* GetRowAt should get headers (2e000851)
* Do not show "updated on" field if the comment has never been updated (40c549b7)
* list comment feature should complete its issue flag (727163e5)
* Added new command aliases (ecaafae6)

0.5.0
=============
2023-12-20

* Makefile: use bb to upload artifacts (8a5d2dca)
* Artifact download should use an io.Writer (7964b0cf)
* Use bigger timeouts for downloads/uploads (f2fc474d)
* Artifact download/delete name completion (d80f5b53)
* Download/Upload should use the private send func (215ebe84)
* repository should get resolved inside profile.send (d4421364)
* Added a new func to discover the current repository (8ae4ccd0)
* Current command should be provide to remote flag (f36b80d4)
* Added Issue comments (7f5ce8eb)
* Get the list of changes for a given issue (8710021f)
* Added issues (df362e84)
* Added User command (15d078b2)
* Use new output (a3cd921d)
* Renamed struct (befbf41f)
* Added Components (146c1565)
* Added entity (1d872a15)
* Moved Link/Links to common (80e17512)
* log 16k of response body (e8688c26)

0.4.0
=============
2023-12-14

* New artifact command (623b5894)
* Documentation: project reviewers (3ed00a81)
* Missing completion function registrations (a5451846)
* Added Project Reviewers (85aab32c)
* New UUID type for Bitbucket (3511a6f0)
* Moved struct (866b1195)
* Branch and commit should use the new output (ffee6dff)

0.3.1
=============
2023-12-12

* Documentation: output formats (2efd5c11)
* Added tsv output format (22c201ec)

0.3.0
=============
2023-12-12

* Output formats (4028f264)
* Duplication of profile names should be checked within the create command (e2679018)
* log errors with their stack (c4076db6)
* Makefile: amd64 and arm64 builds for Linux (77aaa6be)
* Makefile: amd64 and arm64 builds for Linux (830b1f53)

0.2.0
=============
2023-12-11

* Documentation: workspaces and projects (0b685144)
* Version comes from main folder (81f98bc2)
* Workspace argument completion (4551515d)
* Added new flag struct (be44f6d9)
* Added flag completion for EnumFlag (5d7c52bc)
* Better Project key collector (b23f4a10)
* Added Project update (ee3f20bd)
* Added Project Key completion (6007be31)
* Added Project creation (bdf020f3)
* Added Project deletion (f72ecfde)
* Projects should have a Workspace (89b71f31)
* Listing projects should belong to the project command (111b04d3)
* repository argument is not needed in threse commands (320e11e5)
* Updated Bibucket Error struct (c61926e8)
* Added workspaces (a89b1540)
* Added projects (4e77df2e)
* Updated Bibucket Error struct (0ccb414c)
* Allow URIs to start with a / (247cf389)
* fixed help (4e99033e)
* Better Makefile (1b69c9a2)

0.1.0
=============
2023-12-04

* Added documentation (12940100)
* Pull Request dynamic list funcs vary slightly (57ebaac2)
* Unified Pull Request dynamic list (f57da77f)
* Profile subcommand completion should use a dynamic list of profile names (c5e54409)
* New pullrequest subcommand (97561c13)
* Failure should fail the app (2599f437)
* get the open pull requests with profile.GetAll (f73b2a80)
* new alias (d1cfc546)
* New profile subcommands (9c3366fc)
* Created EnuFlag and added EnumFlag to PullRequest list (43958279)
* List should fetch all resources (98d41401)
* Pass the command context to the profile client (6541b10f)
* New Pullrequest commands: merge, decline (1d8f54d4)
* Added some testdata (9edb8e48)
* the Logger should be transmitted via contexts (29fb50a9)
* New Pullrequest commands: approve. unapprove (cb6555e9)
* Display simpler data after pullrequest creation (d2e6767c)
* Added OAUTH 2 Client Credential authentication (3247fdf2)
* Do not show usage on errors (f3be2e5c)
* Create Profile should create config file as needed (85102c98)
* Refactored current profile assignment (f243cbd1)
* Refactored initialization (7662ed17)
* Removed Verbose/Error/Fatal (3c3ce362)
* Moved Logger initialization (378493e4)
* Commands should print their error (1f39c641)
* Profiles should handle the http requests (35ac374c)
* missing structs (6a64f1f0)
* Added Branches (c5eee178)
* Endpoint should use the proper packages (9f5ac811)
* Added Commits (d37a7415)
* Moved user (9ebe4229)
* Moved links (95a8c25c)
* Renamed version variables (0e85c95e)
* use the current git config by default (ef6c5c5a)
* no need to unmarshal empty lists (19232528)
* Added PullRequest List (72acbb01)
* Profile create new flags (d949e25a)
* Profiles (c3994497)
* Read configuration from .env if provided (1217d23b)
* Cobra Skeleton (22fae567)
* Initial Skeleton (b216ab1d)
* Initial commit (f0cea3fc)
