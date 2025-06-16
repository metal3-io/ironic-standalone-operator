# ironic-standalone-operator releasing

This document details the steps to create a release for
`ironic-standalone-operator` aka IrSO.

**NOTE**: Always follow release documentation from the main branch. Release
documentation in release branches may be outdated.

## Before making a release

Things you should check before making a release:

- Check the
  [Metal3 release process](https://github.com/metal3-io/metal3-docs/blob/main/processes/releasing.md)
  for high-level process and possible follow-up actions
- Use the `./hack/verify-release.sh` script as helper to identify possible
  issues to be addressed before creating any release tags. It verifies issues
  like:
  - Verify any other direct or indirect dependency is uplifted to close any public
    vulnerabilities

To make sure the correct Ironic versions are supported and tested:

- Add any missing versions and their branch names to the [list of versions
  in the API][supported-versions].
- Extend the [functional tests suite][suite_test] with new tests if needed:
  - update [the list of known Ironic API versions][known-api-versions] using
    the [Ironic API versions listing][api-versions-list].
  - upgrade from the newest version to `latest` with and without HA
  - upgrade to the newest version from the one before it (with and without HA)
  Hint: you can usually copy existing tests, modifying only the `spec.version`
  field and the API versions on the `TestAssumptions` structure.

[supported-versions]: https://github.com/metal3-io/ironic-standalone-operator/blob/b630805cdd782a51845fc16086e5f64fa77e29af/api/v1alpha1/ironic_types.go#L23-L34
[suite_test]: https://github.com/metal3-io/ironic-standalone-operator/blob/main/test/suite_test.go
[known-api-versions]: https://github.com/metal3-io/ironic-standalone-operator/blob/b630805cdd782a51845fc16086e5f64fa77e29af/test/suite_test.go#L56-L67
[api-versions-list]: https://docs.openstack.org/ironic/latest/contributor/webapi-version-history.html

## Permissions

Creating the first release in a series requires repository `write` permissions
for branch creation.

These permissions are implicit for the org admins and repository admins.
Release team member gets his/her permissions via `metal3-release-team`
membership. This GitHub team has the required permissions in each repository
required to release IrSO. Adding person to the team gives him/her the necessary
rights in all relevant repositories in the organization. Individual persons
should not be given permissions directly.

Patch releases don't require extra permissions.

## Process

IrSO uses [semantic versioning](https://semver.org).

- Regular releases: `v0.x.y`
- Beta releases: `v0.x.y-beta.z`
- Release candidate releases: `v0.x.y-rc.z`

### Repository setup

Clone the repository: `git clone git@github.com:metal3-io/ironic-standalone-operator`

or if using existing repository, make sure origin is set to the fork and
upstream is set to `metal3-io`. Verify if your remote is set properly or not
by using following command `git remote -v`.

- Fetch the remote (`metal3-io`): `git fetch upstream`
This makes sure that all the tags are accessible.

### Preparing a branch

IrSO requires a branch to be created and updated before the automation runs.
If (and only if) you're creating a release `v0.x.0` (i.e. a minor release):

- Switch to the main branch: `git checkout main`

- Identify the commit you wish to create the branch from, and create a branch
  `release-0.x`: `git checkout <sha> -b release-0.x` and push it to remote:
  `git push upstream release-0.x` to create it. Replace `upstream` with
  the actual remote name for the upstream source (not your private fork).

- Setup the CI for the new branch in the prow configuration. [Prior
  art](https://github.com/metal3-io/project-infra/pull/1034).

Create a development branch (e.g. `prepare-0.x`) from the newly created branch:

- Change [the default Ironic version][default-version] to the most recent
  branch of the Ironic image.

- Change the branch of IrSO itself from `latest` to `release-0.x` in two
  places: [IMG in Makefile][img-makefile] and
  [Kustomize configuration][kustomize].

- Commit your changes, push the new branch and create a pull request:
  - The commit and PR title should be
    ✨ Switch to Ironic X.Y by default, prepare release-0.x:
    -`git commit -S -s -m ":sparkles: Switch to Ironic X.Y by default, prepare release-0.x"`
     where X.Y is the Ironic branch you used [above][default-version]
    -`git push -u origin prepare-0.x`
  - The pull request must target the new branch (`release-0.x`), not `main`!

Wait for the pull request to be merged before proceeding.

[default-version]: https://github.com/metal3-io/ironic-standalone-operator/blob/b630805cdd782a51845fc16086e5f64fa77e29af/pkg/ironic/version.go#L14-L15
[img-makefile]: https://github.com/metal3-io/ironic-standalone-operator/blob/b630805cdd782a51845fc16086e5f64fa77e29af/Makefile#L74-L75
[kustomize]: https://github.com/metal3-io/ironic-standalone-operator/blob/b630805cdd782a51845fc16086e5f64fa77e29af/config/manager/manager.yaml#L47

### Creating Release Notes

- Switch to the main branch: `git checkout main`

- Create a new branch for the release notes**:
  `git checkout -b release-notes-0.x.x origin/main`

- Generate the release notes: `RELEASE_TAG=v0.x.x make release-notes`
  - Replace `v0.x.x` with the new release tag you're creating.
  - This command generates the release notes here
    `releasenotes/<RELEASE_TAG>.md` .

- Next step is to clean up the release note manually.
  - If release is not a beta or release candidate, check for duplicates,
    reverts, and incorrect classifications of PRs, and whatever release
    creation tagged to be manually checked.
  - For any superseded PRs (like same dependency uplifted multiple times, or
    commit revertion) that provide no value to the release, move them to
    Superseded section. This way the changes are acknowledged to be part of the
    release, but not overwhelming the important changes contained by the
    release.

- Commit your changes, push the new branch and create a pull request:
  - The commit and PR title should be 🚀 Release v0.x.y:
     -`git commit -S -s -m ":rocket: Release v0.x.x"`
     -`git push -u origin release-notes-0.x.x`
  - Important! The commit should only contain the release notes file, nothing
    else, otherwise automation will not work.

- Ask maintainers and release team members to review your pull request.

Once PR is merged following GitHub actions are triggered:

- GitHub action `Create Release` runs following jobs:
  - GitHub job `push_release_tags` will create and push the tags. This action
    will also create release branch if its missing and release is `rc` or
    minor.
  - GitHub job `create draft release` creates draft release. Don't publish the
    release until release tag is visible in. Running actions are visible on the
    [Actions](https://github.com/metal3-io/ironic-standalone-operator/actions)
    page, and draft release will be visible on top of the
    [Releases](https://github.com/metal3-io/ironic-standalone-operator/releases).
    If the release you're making is not a new major release, new minor release,
    or a new patch release from the latest release branch, uncheck the box for
    latest release. If it is a release candidate (RC) or a beta release,
    tick pre-release box.
  - GitHub job `build_irso` build release images with the
    release tag, and push them to Quay. Make sure the release tags are visible in
    Quay tags pages:
    - [IrSO](https://quay.io/repository/metal3-io/ironic-standalone-operator?tab=tags)
    If the new release tag is not available for any of the images, check if the
    action has failed and retrigger as necessary.

### Release artifacts

We need to verify all release artifacts are correctly built or generated by the
release workflow. You can use `./hack/verify-release.sh` to check for existence
of release artifacts, which should include the following:

Git tags pushed:

- Primary release tag: `v0.x.y`
- Go module tags: `api/v0.x.y` and `test/v0.x.y`

Container images built and tagged at Quay registry:

- [ironic-standalone-operator:v0.x.y](https://quay.io/repository/metal3-io/ironic-standalone-operator?tab=tags)

You can also check the draft release and its tags in the Github UI.

### Make the release

After everything is checked out, hit the `Publish` button your GitHub draft
release!

## Post-release actions for new release branches

Some post-release actions are needed if new minor or major branch was created.

### Branch protection rules

Branch protection rules need to be applied to the new release branch. Copy the
settings after the previous release branch. Branch protection rules require
user to have `admin` permissions in the repository.

### Documentation

Update the [user
guide](https://github.com/metal3-io/metal3-docs/tree/main/docs/user-guide/src):

- Update [supported
  versions](https://github.com/metal3-io/metal3-docs/blob/main/docs/user-guide/src/irso/introduction.md#supported-versions)
  with the Ironic versions that the new release supports.

- Consider if Ironic versions in the
  [examples](https://github.com/metal3-io/metal3-docs/blob/main/docs/user-guide/src/irso/install-basics.md)
  need updating if they are too old or no longer supported.

### Update milestones

- Make sure the next two milestones exist. For example, after 0.4 is out, 0.5
  and 0.6 should exist in Github.
- Set the next milestone date based on the expected release date, which usually
  happens shortly after the next Ironic release.
- Remove milestone date for passed milestones.

Milestones must also be updated in the Prow configuration
([example](https://github.com/metal3-io/project-infra/pull/1035)).

## Additional actions outside this repository

Further additional actions are required in the Metal3 project after IrSO release.
For that, please continue following the instructions provided in
[Metal3 release process](https://github.com/metal3-io/metal3-docs/blob/main/processes/releasing.md)
