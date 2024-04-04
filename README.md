[![Go Report Card](https://goreportcard.com/badge/github.com/bitnami/go-version)](https://goreportcard.com/report/github.com/bitnami/go-version)
[![CI](https://github.com/bitnami/go-version/actions/workflows/go.yml/badge.svg)](https://github.com/bitnami/go-version/actions/workflows/go.yml)

# go-version

go-version is a library for parsing Bitnami package versions and version constraints and verifying versions against a set of constraints.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Usage](#usage)
  - [Version parsing and comparison](#version-parsing-and-comparison)
  - [Sorting](#sorting)
  - [Version constraints](#version-constraints)
    - [Version revision](#version-revision)
    - [Missing major/minor/patch versions](#missing-majorminorpatch-versions)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Usage

Versions used with the `version` package must follow [Semantic Versioning](https://semver.org/) with small adjustments: **pre-release versions are considered revisions** and, therefore, the bigger the revision number, the newer the version.

### Version parsing and comparison

When two versions are compared using functions such as Compare, LessThan, and others, it will follow the specification and always include revisions within the comparison.
It will provide an answer that is valid with the comparison section of [the spec](https://semver.org/#spec-item-11).

```go
v1, _ := version.Parse("1.2.0")
v2, _ := version.Parse("1.2.1")

// Comparison example. There is also GreaterThan, Equal, and just
// a simple Compare that returns an int allowing easy >=, <=, etc.
if v1.LessThan(v2) {
    fmt.Printf("%s is less than %s", v1, v2)
}
```

### Sorting

Collections of versions can be sorted using the `sort.Sort` function from the standard library.

```go
versionsRaw := []string{"1.1.0", "0.7.1", "1.4.0", "1.4.0-alpha", "1.4.1-beta", "1.4.0-alpha.2+20130313144700"}
versions := make(version.Collection, len(versionsRaw))
for i, raw := range versionsRaw {
    v, _ := version.Parse(raw)
    versions[i] = v
}

// After this, the versions are properly sorted
sort.Sort(versions)
```

### Version constraints

Comma-separated version constraints are considered an `AND`. For example, `>= 1.2.3, < 2.0.0` means the version needs to be greater than or equal to `1.2` and less than `3.0.0`.
In addition, they can be separated by `|| (OR)`. For example, `>= 1.2.3, < 2.0.0 || > 4.0.0` means the version needs to be greater than or equal to `1.2` and less than `3.0.0`, or greater than `4.0.0`.

```go
v, _ := version.Parse("2.1.0")
c, _ := version.NewConstraints(">= 1.0, < 1.4 || > 2.0")

if c.Check(v) {
    fmt.Printf("%s satisfies constraints '%s'", v, c)
}
```

Supported operators:

- `=` : you accept that exact version
- `!=` : not equal
- `>` : you accept any version higher than the one you specify
- `>=` : you accept any version equal to or higher than the one you specify
- `<` : you accept any version lower to the one you specify
- `<=` : you accept any version equal or lower to the one you specify
- `^` : it will only do updates that do not change the leftmost non-zero number.
  - e.g. `^1.2.3` := `>=1.2.3, <2.0.0`
- `~` : allows patch-level changes if a minor version is specified on the comparator. Allows minor-level changes if not.
  - e.g. `~1.2.3` := `>=1.2.3, <1.3.0`

#### Version revision

A revision may be denoted by appending a hyphen and a series of dot separated identifiers immediately following the patch version. Revision has a greater precedence than the associated normal version (e.g. `1.2.3-1 > 1.2.3`).

```go
v, _ := version.Parse("2.0.0")
c, _ := version.NewConstraints(">2.0.0-2")

c.Check(v) // false
```

Comparisons include revisions even with no revisions constraint:

```go
v, _ := version.Parse("2.0.0-1")
c, _ := version.NewConstraints(">2.0.0")

c.Check(v) // true
```

#### Missing major/minor/patch versions

If some of the major/minor/patch versions are not specified, it is treated as `*` by default. In short, `3.1.3` satisfies `= 3` because `= 3` is converted to `= 3.*.*`.

```go
v, _ := version.Parse("2.3.4")
c, _ := version.NewConstraints("=2")

c.Check(v) // true
```

Then, `2.2.3` doesn't satisfy `> 2` as `> 2` is treated as `> 2.*.*` = `>= 3.0.0`

```go
v, _ := version.Parse("2.2.3")
c, _ := version.NewConstraints(">2")

c.Check(v) // false
```

`3.3.9` satisifies `= 3.3`, and `5.1.2` doesn't satisfy `> 5.1` likewise.

If you want to treat them as 0, you can pass `version.WithZeroPadding(true)` as an argument of `version.NewConstraints`

```go
v, _ := version.Parse("2.3.4")
c, _ := version.NewConstraints("= 2", version.WithZeroPadding(true))

c.Check(v) // false
```

## License

Copyright &copy; 2024 Broadcom. The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.

You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.
