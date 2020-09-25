package modules

import (
	"fmt"
	"net/url"

	"github.com/joshcarp/gop/gop"

	"gopkg.in/yaml.v2"
)

/* LoadVersion returns the version from a version */
func LoadVersion(retriever gop.Retriever, cacher gop.Cacher, resolver gop.Resolver, cacheFile, resource string) (string, error) {
	var content []byte
	if cacheFile != "" {
		content, _, _ = retriever.Retrieve(cacheFile)
	}
	repo, path, ver := gop.ProcessRepo(resource)
	var repoVer = repo
	if ver != "" {
		repoVer += "@" + ver
	}

	mod := Modules{}
	if err := yaml.Unmarshal(content, &mod); err != nil {
		return "", err
	}
	if val, ok := mod.Imports[repoVer]; ok {
		return AddPath(val, path), nil
	}
	if cacher == nil {
		return resource, nil
	}
	hash, err := resolver(repoVer)
	if err != nil {
		return "", gop.GithubFetchError
	}
	resolve := gop.CreateResource(repo, "", hash)
	if mod.Imports == nil {
		mod.Imports = map[string]string{}
	}
	mod.Imports[repoVer] = resolve
	newfile, err := yaml.Marshal(mod)
	if err != nil {
		return "", err
	}
	if err := cacher.Cache(cacheFile, newfile); err != nil {
		return "", err
	}
	return AddPath(resolve, path), err
}
func AddPath(repoVer string, path string) string {
	a, _, c, _ := gop.ProcessRequest(repoVer)
	return gop.CreateResource(a, path, c)
}
func GetApiURL(resource string) string {
	requestedurl, _ := url.Parse("https://" + resource)
	switch requestedurl.Host {
	case "github.com":
		return "api.github.com"
	default:
		return fmt.Sprintf("%s/api/v3", requestedurl.Host)
	}
	return ""
}