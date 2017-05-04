# Change Log

## 0.6.0 - 2017-05-04
### Added
- `govm use` will spawn a UAC prompt to request elevation if necessary

## 0.5.0
### Fixed
- Reduced amount of file descriptors used while extracting downloaded Go archives
- Improved verbose mode implementation

## 0.4.0
### Added
- `govm install` command added to download and install new Go versions
- `govm uninstall` command added to remove unwanted Go versions

## 0.3.0
### Changed
- Project name changed from __gvm__ to __govm__
- CLI framework changed from [urfave/cli](https://github.com/urfave/cli) to [spf13/cobra](https://github.com/spf13/cobra)

## 0.2.0
### Added
- `gvm use` command implemented to switch between different Go versions

## 0.1.0
### Added
- Initial release
