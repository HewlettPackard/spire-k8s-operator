# Contributing to the Spire Operator Project

The change management process for the Spire Operator Project is designed to be transparent, fair, and
efficient. Anyone may contribute to a project in the Spire Operator repository that they have read access to, provided they:

* Abide by the SPIFFE [code of conduct](https://github.com/spiffe/spiffe/blob/main/CODE-OF-CONDUCT.md)
* Can certify the clauses in the [Developer Certificate of Origin](https://github.com/spiffe/spiffe/blob/main/DCO)

## Getting started

* First, [README](/README.md) to become familiar with how the Spire Operator Project is managed.
* Make sure you're familiar with our [Coding Conventions](#coding-conventions).

## Sending a pull request

1. Fork the repo
2. Commit changes to your fork
3. Update the docs, if necessary
4. Ensure your branch is based on the latest commit in `main`
5. Ensure all tests pass (see project docs for more information)
6. Make sure your commit messages contain a `Signed-off-by: <your-email-address>` line (see `git-commit --signoff`) to certify the [DCO](https://github.com/spiffe/spiffe/blob/main/DCO)
7. Make sure your all your commits are [GPG-signed](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)
8. Open a [pull request](https://help.github.com/articles/creating-a-pull-request-from-a-fork/)
   against the upstream `main` branch

All changes to Spire Operator project must be code reviewed in a pull request (this goes for everyone, even those who have
merge rights).

## After your pull request is submitted

Pull requests are approved according to the process described in our [governance
policies](/GOVERNANCE.md). At least one maintainer must approve the pull request.

Once your pull request is submitted, it's your responsibility to:

* Respond to reviewer's feedback
* Keep it merge-ready at all times until it has been approved and actually merged

Following approval, the pull request will be merged by the submitter of the pull request.

## Coding Conventions <a name="conventions"></a>

Coding conventions will follow
the [SPIFFE project conventions](https://github.com/spiffe/spiffe/blob/main/CONTRIBUTING.md#coding-conventions-).


## Third-party code

When third-party code must be included, all licenses must be preserved. This includes modified
third-party code and excerpts, as well.

## Repositories and Licenses

All repositories under this project should include:

* A detailed `README.md` which includes a link back to this file
* A `LICENSE` file with the Apache 2.0 license
* A [CODEOWNERS](https://help.github.com/articles/about-codeowners/) file listing the maintainers

All code projects should use the [Apache License version 2.0](https://www.apache.org/licenses/LICENSE-2.0), and all
documentation repositories should use
the [Creative Commons License version 4.0](https://creativecommons.org/licenses/by/4.0/legalcode).