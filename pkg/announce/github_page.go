/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package announce

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"sigs.k8s.io/bom/pkg/spdx"
	"github.com/Priyankasaggu11929/release-sdk/git"
	"github.com/Priyankasaggu11929/release-sdk/github"
	"sigs.k8s.io/release-utils/hash"
	"sigs.k8s.io/release-utils/util"
)

const (
	sbomFileName      = "sbom.spdx"
	assetDownloadPath = "/releases/download/"
)

// ghPageBody is a generic template to build the GitHub
// rekease page.
const ghPageBody = `{{ if .Substitutions.logo }}
![Logo]({{ .Substitutions.logo }} "Logo")
{{ end }}
{{ .Substitutions.intro }}
{{ if .Substitutions.changelog }}
See [the CHANGELOG]({{ .Substitutions.changelog }}) for more details.
{{ end }}
{{ if .Substitutions.ReleaseNotes }}
### Release Notes

{{ .Substitutions.ReleaseNotes }}
{{ end }}

`

// GitHubPageOptions data for building the release page
type GitHubPageOptions struct {
	// ReleaseType indicates if we are dealing with an alpha,
	// beta, rc or official
	ReleaseType string

	// AssetFiles is a list of paths of files to be uploaded
	// as assets of this release
	AssetFiles []string

	// Tag is the release the github page will be edited
	Tag string

	// The release can have a name
	Name string

	// Owner GitHub organization which owns the repository
	Owner string

	// Name of the repository where we will publish the
	// release page. The specified tag has to exist there already
	Repo string

	// Run the whole process in non-mocked mode. Which means that it uses
	// production remote locations for storing artifacts and modifying git
	// repositories.
	NoMock bool

	// Create a draft release
	Draft bool

	// If the release exists, we do not overwrite the release page
	// unless specified so.
	UpdateIfReleaseExists bool

	// We can use a custom page template by spcifiying the path. The
	// file is a go template file that renders markdown.
	PageTemplate string

	// File to read the release notes from
	ReleaseNotesFile string

	// We automatizally calculate most values, but more substitutions for
	// the template can be supplied
	Substitutions map[string]string
}

type SBOMOptions struct {
	ReleaseName   string
	Repo          string
	RepoDirectory string
	Tag           string // Version Tag
	Assets        []Asset
}

type Asset struct {
	Path     string // Path where the artifact will be listed
	ReadFrom string // LocalPath to read the information
	Label    string // Label for the asset
}

// GenerateReleaseSBOM creates an SBOM describing the release
func GenerateReleaseSBOM(opts *SBOMOptions) (string, error) {
	// Create a temporary file to write the sbom
	dir, err := os.MkdirTemp("", "project-sbom-")
	if err != nil {
		return "", fmt.Errorf("creating temporary directory to write sbom: %w", err)
	}

	sbomFile := filepath.Join(dir, sbomFileName)
	logrus.Infof("SBOM will be temporarily written to %s", sbomFile)

	builder := spdx.NewDocBuilder()
	builderOpts := &spdx.DocGenerateOptions{
		ProcessGoModules: true,
		ScanLicenses:     true,
		Name:             opts.ReleaseName,
		Namespace:        github.GitHubURL + opts.Repo + "@" + opts.Tag,
		Directories:      []string{opts.RepoDirectory},
	}

	doc, err := builder.Generate(builderOpts)
	if err != nil {
		return "", fmt.Errorf("generating initial SBOM: %w", err)
	}

	// Add the downlad location and version to the first
	// SPDX package (which represents the repo)
	for t := range doc.Packages {
		doc.Packages[t].Version = opts.Tag
		doc.Packages[t].DownloadLocation = "git+" + github.GitHubURL + opts.Repo + "@" + opts.Tag
		break
	}

	// List all artifacts and add them
	spdxClient := spdx.NewSPDX()
	for _, f := range opts.Assets {
		logrus.Infof("Adding file %s to SBOM", f.Path)
		spdxFile, err := spdxClient.FileFromPath(f.ReadFrom)
		if err != nil {
			return "", fmt.Errorf("adding %s to SBOM: %w", f.ReadFrom, err)
		}
		spdxFile.Name = f.Path
		spdxFile.BuildID() // This is a boog in the spdx pkg, we have to call manually
		spdxFile.DownloadLocation = github.GitHubURL + filepath.Join(
			opts.Repo, assetDownloadPath, opts.Tag, f.Path,
		)
		if err := doc.AddFile(spdxFile); err != nil {
			return "", fmt.Errorf("adding %s as SPDX file to SBOM: %w", f.ReadFrom, err)
		}
	}

	if err := doc.Write(sbomFile); err != nil {
		return "", fmt.Errorf("writing sbom to disk: %w", err)
	}

	return sbomFile, nil
}

// UpdateGitHubPage updates a github page with data from the release
func UpdateGitHubPage(opts *GitHubPageOptions) (err error) {
	token := os.Getenv(github.TokenEnvKey)
	if token == "" {
		return errors.New("cannot update release page without a GitHub token")
	}

	gh := github.New()
	releaseVerb := "Posting"
	semver, err := util.TagStringToSemver(opts.Tag)
	if err != nil {
		return fmt.Errorf("parsing semver from tag: %w", err)
	}

	// Determine if this is a prerelase
	// // [[ "$FLAGS_type" == official ]] && prerelease="false"
	isPrerelease := false
	if len(semver.Pre) > 0 {
		isPrerelease = true
	}

	// Process the specified assets
	releaseAssets, err := processAssetFiles(opts.AssetFiles)
	if err != nil {
		return fmt.Errorf("processing the asset file list: %w", err)
	}

	// Substitution struct for the template
	subs := struct {
		Substitutions map[string]string
		Assets        []map[string]string
	}{
		Substitutions: opts.Substitutions,
		Assets:        releaseAssets,
	}

	// If we have a release notes file defined and set a substitution
	// entry for its contents
	if opts.ReleaseNotesFile != "" {
		rnData, err := os.ReadFile(opts.ReleaseNotesFile)
		if err != nil {
			return fmt.Errorf("reading release notes file: %w", err)
		}
		subs.Substitutions["ReleaseNotes"] = string(rnData)
	}

	// Open the template file (if a custom)
	templateText := ghPageBody
	if opts.PageTemplate != "" {
		logrus.Debugf("Using custom page template %s", opts.PageTemplate)
		templateText = opts.PageTemplate
	}
	// Parse the template we will use to build the release page
	tmpl, err := template.New("GitHubPage").Parse(templateText)
	if err != nil {
		return fmt.Errorf("parsing github page template: %w", err)
	}

	// Run the template to verify the output.
	output := new(bytes.Buffer)
	err = tmpl.Execute(output, subs)
	if err != nil {
		return fmt.Errorf("executing page template: %w", err)
	}

	// If we are in mock, we write it to stdout and exit. All checks
	// performed to the repo are skipped as the tag may not exist yet.
	if !opts.NoMock {
		logrus.Info("Mock mode, outputting the release page")
		_, err := os.Stdout.Write(output.Bytes())
		if err != nil {
			return fmt.Errorf("writing github page to stdout: %w", err)
		}
		return nil
	}

	// Check to see that a tag exists.
	// non-draft release posts to github create a tag.  We don't want to
	// create any tags on the repo this way. The tag should already exist
	// as a result of the release process.
	tagFound, err := gh.TagExists(opts.Owner, opts.Repo, opts.Tag)
	if err != nil {
		return fmt.Errorf("checking if the tag already exists in GitHub: %w", err)
	}
	if !tagFound {
		logrus.Warnf("The %s tag doesn't exist yet on GitHub.", opts.Tag)
		logrus.Warnf("That can't be good.")
		logrus.Warnf("We certainly cannot publish a release without a tag.")
		return errors.New("tag not found while trying to publish release page")
	}

	// Get the release we are looking for
	releases, err := gh.Releases(opts.Owner, opts.Repo, true)
	if err != nil {
		return fmt.Errorf("listing the repositories releases: %w", err)
	}

	// Does the release exist yet?
	var releaseID int64
	commitish := ""
	for _, release := range releases {
		if release.GetTagName() == opts.Tag {
			releaseID = release.GetID()
			commitish = release.GetTargetCommitish()
		}
	}

	if releaseID != 0 {
		logrus.Warnf("The %s is already published on github.", opts.Tag)
		if !opts.UpdateIfReleaseExists {
			return errors.New("release " + opts.Tag + " already exists. Left intact")
		}
		logrus.Infof("Using release id %d to update existing release.", releaseID)
		releaseVerb = "Updating"
	}

	// Post release data
	logrus.Infof("%s the %s release on github...", releaseVerb, opts.Tag)

	// Call GitHub to set the release page
	release, err := gh.UpdateReleasePage(
		opts.Owner, opts.Repo, releaseID,
		opts.Tag, commitish, opts.Name, output.String(),
		opts.Draft, isPrerelease,
	)
	if err != nil {
		return fmt.Errorf("updating the release on GitHub: %w", err)
	}

	// Releases often take a bit of time to show up in the API
	// after creating the page. If the release does not appear
	// in the API right away , sleep 3 secs and retry 3 times.
	for checkAttempts := 3; checkAttempts >= 0; checkAttempts-- {
		releaseFound := false
		releases, err = gh.Releases(opts.Owner, opts.Repo, true)
		if err != nil {
			return fmt.Errorf("listing releases in repository: %w", err)
		}
		// Check if the page shows up in the API
		for _, testRelease := range releases {
			if testRelease.GetID() == release.GetID() {
				releaseFound = true
				break
			}
		}
		if releaseFound {
			break
		}

		if checkAttempts == 0 {
			return errors.New("release not found, even when call to github was successful")
		}
		logrus.Info("Release page not yet returned by the GitHub API, sleeping and retrying")
		time.Sleep(3 * time.Second)
	}

	// Delete any assets reviously uploaded
	if err := deleteReleaseAssets(gh, opts.Owner, opts.Repo, release.GetID()); err != nil {
		return fmt.Errorf("deleting the existing release assets: %w", err)
	}

	// publish binary
	for _, assetData := range releaseAssets {
		logrus.Infof("Uploading %s as release asset", assetData["realpath"])
		asset, err := gh.UploadReleaseAsset(opts.Owner, opts.Repo, release.GetID(), assetData["rawpath"])
		if err != nil {
			return fmt.Errorf("uploading %s to the release: %w", assetData["realpath"], err)
		}
		logrus.Info("Successfully uploaded asset #", asset.GetID())
	}
	logrus.Infof("Release %s published on GitHub", opts.Tag)
	return nil
}

// processAssetFiles reads the command line strings and returns
// a map holding the needed info from the asset files
func processAssetFiles(assetFiles []string) (releaseAssets []map[string]string, err error) {
	// Check all asset files and get their hashes
	for _, path := range assetFiles {
		assetData := map[string]string{
			"rawpath": path,
			"name":    "",
		}
		// Check if asset path has a label
		if strings.Contains(path, ":") {
			p := strings.SplitN(path, ":", 2)
			if len(p) == 2 {
				path = p[0]
				assetData["name"] = p[1]
			}
		}

		logrus.Debugf("Checking asset file %s", path)

		// Verify path exists
		if !util.Exists(path) {
			// TODO: psaggu
			//return nil, errors.New("unable to render release page, asset file does not exist")
			return nil, fmt.Errorf("unable to render release page, asset file does not exist. Path: %s", path)
		}

		assetData["realpath"] = path
		assetData["filename"] = filepath.Base(path)

		fileHashes, err := getFileHashes(path)
		if err != nil {
			return nil, fmt.Errorf("getting the hashes: %w", err)
		}

		assetData["sha512"] = fileHashes["512"]
		assetData["sha256"] = fileHashes["256"]

		releaseAssets = append(releaseAssets, assetData)
	}
	return releaseAssets, nil
}

func deleteReleaseAssets(gh *github.GitHub, owner, repo string, releaseID int64) error {
	// If the release already contains assets, delete them to match
	// the new uploads we are sending
	currentAssets, err := gh.ListReleaseAssets(owner, repo, releaseID)
	if err != nil {
		return fmt.Errorf("while checking if the release already has assets: %w", err)
	}
	if len(currentAssets) == 0 {
		logrus.Info("No assets found in release")
		return nil
	}

	logrus.Warnf("Deleting %d release assets to upload the latest files", len(currentAssets))
	for _, asset := range currentAssets {
		logrus.Infof("Deleting %s", asset.GetName())
		if err := gh.DeleteReleaseAsset(owner, repo, asset.GetID()); err != nil {
			return fmt.Errorf("deleting existing release assets: %w", err)
		}
	}
	return nil
}

// getFileHashes obtains a file's sha256 and 512
func getFileHashes(path string) (hashes map[string]string, err error) {
	sha256, err := hash.SHA256ForFile(path)
	if err != nil {
		return nil, fmt.Errorf("get sha256: %w", err)
	}

	sha512, err := hash.SHA512ForFile(path)
	if err != nil {
		return nil, fmt.Errorf("get sha512: %w", err)
	}

	return map[string]string{"256": sha256, "512": sha512}, nil
}

// Validate the GitHub page options to ensure they are correct
func (o *GitHubPageOptions) Validate() error {
	// TODO: Check that the tag is well formed
	if o.Tag == "" {
		return errors.New("cannot update github page without a tag")
	}
	if o.Repo == "" {
		return errors.New("cannot update github page, repository not defined")
	}
	if o.Owner == "" {
		return errors.New("cannot update github page, github organization not defined")
	}

	return nil
}

// ParseSubstitutions gets a slice of strings with the substitutions
// for the template and parses it as Substitutions in the options
func (o *GitHubPageOptions) ParseSubstitutions(subs []string) error {
	o.Substitutions = map[string]string{}
	for _, sString := range subs {
		p := strings.SplitN(sString, ":", 2)
		if len(p) != 2 || p[0] == "" {
			return errors.New("substitution value not well formed: " + sString)
		}
		o.Substitutions[p[0]] = p[1]
	}
	return nil
}

// SetRepository takes a repository slug in the form org/repo,
// paeses it and assigns the values to the options
func (o *GitHubPageOptions) SetRepository(repoSlug string) error {
	org, repo, err := git.ParseRepoSlug(repoSlug)
	if err != nil {
		return fmt.Errorf("parsing repository slug: %w", err)
	}
	o.Owner = org
	o.Repo = repo
	return nil
}

// ReadTemplate reads a custom template from a file and sets
// the PageTemplate option with its content
func (o *GitHubPageOptions) ReadTemplate(templatePath string) error {
	// If path is empty, no custom template will be used
	if templatePath == "" {
		o.PageTemplate = ""
		return nil
	}

	// Otherwise, read a custom template from a file
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("reading page template text: %w", err)
	}
	logrus.Infof("Using custom template from %s", templatePath)
	o.PageTemplate = string(templateData)
	return nil
}
