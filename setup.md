# golang template

This file describes how to setup this template.

## Features

* Adds binaries automatically when a release is created. The binaries will contain the tag of the release as version.
* Provide pull-request template with information for versioning

## Assumptions

* This template is tested only with a single module and a single binary. However, it could work with more.

## Usage

### Setup

1. Create a project based on this template
   * Select this template when creating a new project on GitHub
   * Go to [https://github.com/yawn77/tmplgolang/](https://github.com/yawn77/tmplgolang/) and click on *Use this template*
2. Adapt `README.md`. Enter the project name and adapt sections to your needs (further inspiration: [https://www.makeareadme.com/](https://www.makeareadme.com/))
3. Select a license (e.g., from [https://choosealicense.com/](https://choosealicense.com/))
   * Replace the `LICENSE` file
   * Edit the license link at the end of `README.md` file
4. Set *repository URL* and *go version* in `go.mod`
5. Rename directory `tmplgolang` in the `cmd` directory according to your needs
6. Adapt GitHub actions to your needs
   * Add go get to GitHub actions if required
   * `build.yml`/`test.yml`: select platforms
   * `publish.yml`: select platforms and update *project_path*
   * `docker.yml`: select platforms and set tags
7. Remove `setup.md`

## Maintenance

* Update date in `LICENSE` file on new year's day

## Recommendations

* Use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) to be compatible with future versions of this template
* Maintain the [changelog](https://keepachangelog.com/en/1.0.0/) with each commit. There will be features which will automate the maintenance of this file to some extent in the future.
* Branch protection for `main` and `v[0-9]*`
  * *Enable:* Require a pull request before merging
  * *Enable:* Require status checks to pass before merging
    * *Enable:* Require branches to be up to date before merging
    * Select status checks
  * *Enable:* Require conversation resolution before merging
  * *Enable:* Do not allow bypassing the above settings
