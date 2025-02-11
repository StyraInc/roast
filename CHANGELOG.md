# Changelog

## [0.8.0] - 2025-02-11

### Changed

- Repo now resides under the StyraInc org
- Don't print location for object rule with implied true value
- Copy concurrent map code from Regal to here
- Copy Set implementation from Regal to here
- Use ast.ValueName from OPA now that it's been upstreamed
- Add a few more common string terms for interning
- Add ToAST function to build entire Regal AST in Roast
- Better organization of interned values and terms

## [0.7.0] - 2025-01-27

### Changed

- Bump OPA dependency to 1.1.0
- Annotations scoped `rule` or `document` no longer serialized under the `package` node, but found under each
  respective rule only. Marginal performance improvement, but certainly more correct.

## [0.6.0] - 2025-01-14

### Changed

- Bump OPA dependency to 1.0.0, and update imports to v1
- Faster encoding of Term's by avoiding OPA's TypeName

## [0.5.0] - 2024-12-11

### Added

- New `ToOPAInputValue` function to prepare a map or slice for use as `rego.EvalParsedInput`. This is much
  faster than letting OPA do the conversion, but will only work for inputs created by this library.
- Add optimized `AnyToValue` implementation similar to `InterfaceToValue` provided by OPA, but tailored
  only for the use case of converting an AST `map[string]any` to a `ast.Value`. Highly optimized.
- New `encoding.JSONRoundTrip(from, to)` and `encoding.MustJSONRoundTrip(from, to)` convenience functions
- Bump OPA dependency to v0.70.0

## [0.4.2] - 2024-10-03

### Changed

- Fixed potential data race in package path serialization

## [0.4.1] - 2024-10-01

### Changed

- Update actual dependencies used (i.e. `go mod tidy`)

### Changed

## [0.4.0] - 2024-10-01

### Changed

- New location format
- Removed `name` attribute from rules in favor of using the rule's `ref` to infer name
- Updated OPA version from v0.68.0 to v0.69.0

## [0.3.0] - 2024-09-25

### Changed

- Removed `annotations` from module, in favor of annotations on `package` and `rules`

## [0.2.0] - 2024-09-09

### Changed

- OPA version updated from v0.67.1 to v0.68.0

## [0.1.1] - 2024-09-09

### Changed

- Fixed issue in annotations encoding, where multiple `custom` attributes wouldn't be encoded
  with a `,` separator.

## [0.1.0] - 2024-08-20

First release!
