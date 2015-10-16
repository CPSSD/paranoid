Writing Code for Paranoid
=========================

## Code Formatting ##

All `go` code checked into the repository must be formatted according to `go fmt`. To lint your
code for style errors, use [golint](https://github.com/golang/lint). All code must also be
checked by the `go vet` command, which will report valid but poor code. This is separate
from `golint`, which report style errors.

## Code Style ##

Go doesn't have an official style guide as such, but a list of common mistakes can be
found [here](https://github.com/golang/go/wiki/CodeReviewComments). Also worth a look is
[Effective Go](https://golang.org/doc/effective_go.html), a slightly more advanced guide
to writing Go.

## Commenting ##

Commenting will be done in accordance with the standard go style used for generating godocs. 
You should write comments to explain any code you write that does something in a way that may not be obvious.
An explaination of standard go commenting practice and godoc can be found [here](https://blog.golang.org/godoc-documenting-go-code).

## Branching ##

Branches should have a short prefix describing what type of changes are contained in it.
These prefixes should be one of the following:

* **feature/** -- for changes which add/affect a feature.
* **doc/** -- for changes to the documentation.
* **hotfix/** -- for quick bugfixes.

Branches must also contain the name of the contributor, i.e. `doc/terry/...`.

## Version Numbering ##

We will be using [semantic versioning](http://semver.org/). Every binary will have a separate
version number (server, client, etc.). This will begin at 0.1.0 and will be incremented strictly
according to the semantic versioning guidelines.

## Code Review ##

All code must be submitted via pull requests, *not* checked straight into the repo.
A pull request will be from a single feature branch, which will not be used for any other
features after merging. Ideally, it will be deleted after a merge.

All pull requests must be reviewed by at least one person who is not the submitter. You can
ask in Hipchat for a review, or use Github's assign feature to assign the pull request to a
specific person.
