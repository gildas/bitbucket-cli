# Contributing to Bitbucket-cli

First off, thank you for taking the time to contribute to the [Bitbucket cli](https://github.com/gildas/bitbucket-cli)!

We welcome contributions from the community to help make this the best command-line interface for the Bitbucket platform.

It’s folks like you that make `bitbucket-cli` a better tool for everyone.

---

## Getting Started

1. **Fork the repository**: To get started, please **fork the repository** to your own GitHub account.
2. **Clone your fork**: Clone the forked repository to your local machine.
3. **Create a feature branch**: Create a branch for your changes, ensuring it is based off the latest code.

---

## Pull Request Guidelines

To maintain code quality and a streamlined workflow, we enforce the following rules for all Pull Requests:

### 1. Reporting Issues

If you find a bug, please check if the issue you are addressing has already been reported. If not, please create a new [issue](https://github.com/gildas/bitbucket-cli/issues) with a clear description of the problem and link that issue in your Pull Request.

### 2. Target the `dev` Branch

All Pull Requests **must** be targeted at the `dev` branch. 
> [!IMPORTANT]
> PRs opened against the `master` branch will be closed or you will be asked to retarget them to `dev`.

### 3. Signed Commits

Integrity is key. **All commits in Pull Requests must be signed** (GPG, SSH, or X.509). 

* PRs containing unsigned commits will be closed or asked to be retargeted once the commits are signed.
* If you aren't sure how to do this, check out [GitHub's guide on signing commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/about-commit-signature-verification).

### 4. Command Structure

* **Resources and commands:**  
`bb` is built as a modern CLI using subcommands. Ensure new features follow this pattern (e.g., `bb <resource> <subresource...> <command>`).  
Resources should be nouns (e.g., `repository`, `pullrequest`), and commands should be verbs (e.g., `list`, `create`, `delete`). Resources should support the standard CRUD operations (Create -> `create`, Read -> `list` and `get`, Update -> `update`, Delete -> `delete`) where applicable. Additional commands are welcome.
* **Dry Run Support:**  
All commands that modify data on Bitbucket should support the --dry-run flag to allow users to preview changes.
* **Output Formats:**  
Ensure list and get commands remain compatible with various supported output formats (JSON, YAML, Table, etc.).

---

## Style & Standards

* **Formatting**:  
Ensure your code follows the standard Go language conventions (you can run `make fmt` in the project root).
* **Documentation**:  
If you are adding a feature, please update any relevant documentation or help text within the CLI and the [README.md](README.md) file.
* **Tests**:  
Verify your changes by running existing tests and adding new ones where applicable.  
If you add JSON paylods in the tests, make sure to add them in the `testdata` directory and reference them in your test code. You can find examples in the existing test files.  
Please ensure that the payloads are anonymized enough and do not contain any sensitive information.  
You can run all tests with `make test`.

---

## Code of Conduct

We ask that all contributors adhere to the [Code of Conduct](CODE_OF_CONDUCT.md) to maintain a welcoming and inclusive environment for everyone.

---

## License

By contributing to Bitbucket-cli, you agree that your contributions will be licensed under the project's current license.

You can find the license details in the [LICENSE](LICENSE) file.

---

## Thank You!

Thank you again for your interest in contributing to Bitbucket-cli! We look forward to your contributions and are excited to see how you can help improve the project. If you have any questions or need assistance, please don't hesitate to reach out by opening an issue or joining our discussions. Happy coding!
