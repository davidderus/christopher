# Changelog

## v1.0.0-alpha.5 - 2018/01/21

- Not using Travis for cross-platform deployment
- Fixing invalid yml indent in goreleaser config
- Using goreleaser default script and not relying on go for deploy

## v1.0.0-alpha.4 - 2018/01/21

- Adding execute bit to deploy script
- Removing osx from travis env config (the tests were always failing)

## v1.0.0-alpha.3 - 2018/01/21

- Fixing wrong shasum error for bootstrap css and js

## v1.0.0-alpha.2 - 2017/10/22

- Renaming config key `rpcURL` to `rpc_url` for consistency
- Introducing Teller, a text/json logger
- Releasing as binary for next releases
- Fixing licence copyright

## v1.0.0-alpha.1 - 2017/09/30

- Adding GoRelease config
- Adding Travis config
- Initial release
