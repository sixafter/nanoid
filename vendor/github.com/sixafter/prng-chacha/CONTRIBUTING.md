# Contributing

Thank you for your interest in contributing to **PRNG-CHACHA**! Your contributions help make this project better for everyone. This guide outlines the process for contributing to the project, including reporting issues, submitting pull requests, and adhering to the project's code of conduct.

## 📜 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How to Contribute](#how-to-contribute)
    - [Reporting Bugs](#reporting-bugs)
    - [Requesting Features](#requesting-features)
    - [Submitting Changes](#submitting-changes)
- [Coding Guidelines](#coding-guidelines)
    - [Style Guidelines](#style-guidelines)
- [Security Considerations](#-security-considerations)
- [Pull Request Process](#pull-request-process)
---

## 🛡️ Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [osint@sixafter.com](mailto:osint@sixafter.com).

---

## Versioning

This software adheres to the [Semantic Versioning 2.0](https://semver.org/spec/v2.0.0.html) standard for version numbering as quoted here:

Given a version number MAJOR.MINOR.PATCH, increment the:

1.	MAJOR version when you make incompatible API changes
2.	MINOR version when you add functionality in a backward compatible manner
3.	PATCH version when you make backward compatible bug fixes

Additional labels for pre-release and build metadata are available as extensions to the MAJOR.MINOR.PATCH format.

---

## 🤝 How to Contribute

There are several ways you can contribute:

### 🐛 Reporting Bugs

If you encounter a bug or unexpected behavior:

1. **Search Existing Issues**: Check if the issue has already been reported [here](https://github.com/sixafter/prng-chacha/issues).
2. **Open a New Issue**: If not, open a new issue describing the bug in detail.
    - **Provide a Clear Title**: Summarize the problem.
    - **Describe the Behavior**: Explain what you expected to happen versus what actually happened.
    - **Steps to Reproduce**: Include code snippets or commands that can help reproduce the issue.
    - **Environment Details**: Mention your Go version, operating system, and any other relevant information.

### 🌟 Requesting Features

To suggest new features or improvements:

1. **Search Existing Features**: Ensure your idea hasn't been discussed [here](https://github.com/sixafter/prng-chacha/issues?q=is%3Aissue+is%3Aopen+label%3Afeature).
2. **Open a New Feature Request**: Provide a detailed description of the feature.
    - **Describe the Feature**: Explain what the feature is and why it's needed.
    - **Use Cases**: Provide examples of how the feature would be used.
    - **Potential Implementation**: If possible, suggest how it could be implemented.

### ✨ Submitting Changes

Contributions in the form of bug fixes, improvements, or new features are welcome!

#### 1. Fork the Repository

Fork the repository to your GitHub account by clicking the **Fork** button at the top right of the repository page.

#### 2. Clone Your Fork

Clone the forked repository to your local machine:

```bash
git clone https://github.com/sixafter/prng-chacha.git
cd prng-chacha
```

#### 3. Create a New Branch

Create a new branch for your changes:

```bash
git checkout -b feature/your-feature-name
```

#### 4. Make Your Changes

#### 5. Run Tests and Linters

Ensure all tests pass and the code adheres to the style guidelines:

```bash
make lint
make test
```

#### 6. Commit Your Changes

Commit your changes with a clear and descriptive message:

```bash
git add .
git commit -m "Add descriptive commit message"
```

#### 7. Push to Your Fork

Push your changes to your forked repository:

```bash
git push origin feature/your-feature-name
```

#### 8. Open a Pull Request

Navigate to the original repository and click New Pull Request. Provide a clear description of your changes and link any related issues.

## 🎨 Coding Guidelines

Adhering to consistent coding standards ensures the codebase remains clean, readable, and maintainable.

### 🖋️ Style Guidelines

* **Formatting**: Use `make fmt` to format your code. 
* **Linting**: Follow the linting rules defined in `.golangci.yaml`. Ensure that your code passes all linters before submitting. 
* **Documentation**: Document public functions, types, and methods using Go's standard documentation conventions. 
* **Error Handling**: Handle errors gracefully and consistently. Use the predefined error types where applicable.

---

## 🔒 Security Considerations

* **Randomness**: Ensure that all randomness sources use cryptographically secure methods (crypto/rand). 
* **Data Sanitization**: Avoid exposing sensitive data through IDs or logs.

---

## 🚀 Pull Request Process

Follow these steps to create a successful pull request (PR):

1. Ensure Your Branch is Up-to-Date
   * Before opening a PR, make sure your branch is based on the latest main branch.
2. Create a Pull Request 
3. Address Feedback 
   * **Respond Promptly**: Engage with reviewers by responding to comments and making necessary changes. 
   * **Update Your PR**: Push additional commits to your branch to address feedback.
4. Merge the PR 
   * Once approved and all checks pass, your PR will be merged by a maintainer.

---

## 📝 Additional Resources

* [Go Documentation](https://go.dev/doc/) 
* [GolangCI-Lint Documentation](https://golangci-lint.run) 
* [GitHub Flow](https://docs.github.com/en/get-started/using-github/github-flow) 
* [Contributor Covenant Code of Conduct](https://github.com/sixafter/prng-chacha/blob/main/CODE_OF_CONDUCT.md)

## 🙏 Thank You!

We appreciate your interest in contributing to `prng-chacha`! Your efforts help improve the project and support the community. If you have any questions or need assistance, feel free to reach out by opening an issue or contacting the maintainers.

Happy coding! 🎉