package manifest

import "github.com/dtnitsch/manifestor/internal/build"

func DefaultManifestMeta() ManifestMeta {
    return ManifestMeta{
        Version: "0.3",
        Generator: GeneratorMeta{
            Name:      build.Name,
            Version:   build.Version,
            BuildTime: build.BuildTime,
            Commit:    build.CommitSHA,
        },
        Schema: SchemaMeta{
            Node:   "node.v1",
            Rollup: "rollup.v1",
        },
    }
}

