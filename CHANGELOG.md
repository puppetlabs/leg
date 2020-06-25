# [2.13.0](https://github.com/puppetlabs/horsehead/compare/v2.12.0...v2.13.0) (2020-06-25)


### Update

* Adds Content-Security-Policy middleware builder ([f0b5949fb58660bc37980a9927e68eaec5e63553](https://github.com/puppetlabs/horsehead/commit/f0b5949fb58660bc37980a9927e68eaec5e63553))

# [2.12.0](https://github.com/puppetlabs/horsehead/compare/v2.11.1...v2.12.0) (2020-06-23)


### feat

* Use token-based auth with Intercom ([ff10479de6353a81a17cd05e5f1a3d607d8e1be7](https://github.com/puppetlabs/horsehead/commit/ff10479de6353a81a17cd05e5f1a3d607d8e1be7))

### New

* Use access token authentication with Intercom ([e27685d22582529a0ee71253b0eb9a9add3b3010](https://github.com/puppetlabs/horsehead/commit/e27685d22582529a0ee71253b0eb9a9add3b3010))

## [2.11.1](https://github.com/puppetlabs/horsehead/compare/v2.11.0...v2.11.1) (2020-06-23)


### Fix

* Respect control termination in activity reporting. ([aa4da115236c56e9a2f6dbd798ae3b0b4711eec8](https://github.com/puppetlabs/horsehead/commit/aa4da115236c56e9a2f6dbd798ae3b0b4711eec8))

# [2.11.0](https://github.com/puppetlabs/horsehead/compare/v2.10.0...v2.11.0) (2020-06-22)


### Build

* Remove vendor directory ([321b22757c559b049f71f5f4503f342c19ecf702](https://github.com/puppetlabs/horsehead/commit/321b22757c559b049f71f5f4503f342c19ecf702))

### New

* Add Intercom support ([bd7d922b95c6a92f7544b61ba31b6809fbb8beb3](https://github.com/puppetlabs/horsehead/commit/bd7d922b95c6a92f7544b61ba31b6809fbb8beb3))

# [2.10.0](https://github.com/puppetlabs/horsehead/compare/v2.9.0...v2.10.0) (2020-06-17)


### New

* Import utilities from other programs ([2086c81fc4226b2ba4c9ad8cbc52741b613f8cf2](https://github.com/puppetlabs/horsehead/commit/2086c81fc4226b2ba4c9ad8cbc52741b613f8cf2))

# [2.9.0](https://github.com/puppetlabs/horsehead/compare/v2.8.0...v2.9.0) (2020-06-11)


### Update

* Splits out preflight and middleware-enabled CORS handlers ([2ad1cceec574000a16ddade63be23b4a83bedf23](https://github.com/puppetlabs/horsehead/commit/2ad1cceec574000a16ddade63be23b4a83bedf23))

# [2.8.0](https://github.com/puppetlabs/horsehead/compare/v2.7.0...v2.8.0) (2020-06-09)


### Update

* http/api.CORSBuilder: adds AllowOrigins for setting origins that are allowed ([c10aefdd5d4d4c2f07da3ead34aa5754e7be26d6](https://github.com/puppetlabs/horsehead/commit/c10aefdd5d4d4c2f07da3ead34aa5754e7be26d6))

# [2.7.0](https://github.com/puppetlabs/horsehead/compare/v2.6.0...v2.7.0) (2020-02-25)


### Update

* Add support for a generic interface type that supports binary data in our transfer package ([6b2565e44cb75c160337fd7a077cba6aff97ae5d](https://github.com/puppetlabs/horsehead/commit/6b2565e44cb75c160337fd7a077cba6aff97ae5d))

# [2.6.0](https://github.com/puppetlabs/horsehead/compare/v2.5.0...v2.6.0) (2020-02-10)


### Update

* Add noop metrics delegate for testing ([28e8bb7ea31a4f0587a144c4a27026da0a8e2305](https://github.com/puppetlabs/horsehead/commit/28e8bb7ea31a4f0587a144c4a27026da0a8e2305))

# [2.5.0](https://github.com/puppetlabs/horsehead/compare/v2.4.0...v2.5.0) (2020-01-17)


### chore

* adds *.swp to .gitignore ([410cd51b8022e091c472a77fe75520039019f8fb](https://github.com/puppetlabs/horsehead/commit/410cd51b8022e091c472a77fe75520039019f8fb))

### scheduler

* adds backoff to RecoveryDescriptor ([61f9de47b617b08eadaaeb907cfd448d813354dc](https://github.com/puppetlabs/horsehead/commit/61f9de47b617b08eadaaeb907cfd448d813354dc))

### Update

* Scheduler: adds backoff and max retries to RecoveryDescriptor ([9f7d39211c9b69197bba3e09fad281b2e236ee52](https://github.com/puppetlabs/horsehead/commit/9f7d39211c9b69197bba3e09fad281b2e236ee52))

### vendor

* updates testify version ([01fad74b0b5826f3ef6a9ab750ad561b8bbe9236](https://github.com/puppetlabs/horsehead/commit/01fad74b0b5826f3ef6a9ab750ad561b8bbe9236))

# [2.4.0](https://github.com/puppetlabs/horsehead/compare/v2.3.0...v2.4.0) (2019-11-14)


### metrics

* adds label override to example metrics server ([00b489e83ab265b23f43d7c55a4d134a28cc4d9d](https://github.com/puppetlabs/horsehead/commit/00b489e83ab265b23f43d7c55a4d134a28cc4d9d))
* fixes panic when using WithLabels without passing in labels through ObserveDuration ([996564070c3b6eb5b86ee427e2cf4d43a9ec6a25](https://github.com/puppetlabs/horsehead/commit/996564070c3b6eb5b86ee427e2cf4d43a9ec6a25))
* removes delegate field from prometheus timer metric ([afda64785a3f2459aa1eb86d77ef87a871c16682](https://github.com/puppetlabs/horsehead/commit/afda64785a3f2459aa1eb86d77ef87a871c16682))

### Update

* Updates metrics library ([332a70b6b29857b1365479920f0c12810aa4b172](https://github.com/puppetlabs/horsehead/commit/332a70b6b29857b1365479920f0c12810aa4b172))

# [2.3.0](https://github.com/puppetlabs/horsehead/compare/v2.2.0...v2.3.0) (2019-11-05)


### Update

* Add support for specifying the sensitivity to use to encode errors in httputil/api ([0da358bb55d6306c14e0058fb729df6b38dd63dd](https://github.com/puppetlabs/horsehead/commit/0da358bb55d6306c14e0058fb729df6b38dd63dd))

# [2.2.0](https://github.com/puppetlabs/horsehead/compare/v2.1.2...v2.2.0) (2019-10-25)


### New

* Add data structures and graph libraries ([fc99ca1c6b0b9001d78d97cbbc9475b589654c23](https://github.com/puppetlabs/horsehead/commit/fc99ca1c6b0b9001d78d97cbbc9475b589654c23))

## [2.1.2](https://github.com/puppetlabs/horsehead/compare/v2.1.1...v2.1.2) (2019-10-23)


### Build

* Fix semantic-release vulnerabilities ([47908e9ae877b2c7a3fbd80092eaf038192fa39d](https://github.com/puppetlabs/horsehead/commit/47908e9ae877b2c7a3fbd80092eaf038192fa39d))
* Remove nonexistent-but-private package from go.sum ([42c802ca3f66a3081cc67d61bd2c1e74c5410176](https://github.com/puppetlabs/horsehead/commit/42c802ca3f66a3081cc67d61bd2c1e74c5410176))

### Chore

* Remove old CODEOWNERS ([c982592f3f21f37f4199bfa4ffb7708cedb470d5](https://github.com/puppetlabs/horsehead/commit/c982592f3f21f37f4199bfa4ffb7708cedb470d5))

### Fix

* Sync dependencies as a release ([5c247e63eaf48c283e17667bcaf055a1c91873e9](https://github.com/puppetlabs/horsehead/commit/5c247e63eaf48c283e17667bcaf055a1c91873e9))

## [2.1.1](https://github.com/puppetlabs/horsehead/compare/v2.1.0...v2.1.1) (2019-10-23)


### Fix

* Fix error domains; prepare for public release ([80a8dabc371637e1c438c0e669eca4539ace9c5d](https://github.com/puppetlabs/horsehead/commit/80a8dabc371637e1c438c0e669eca4539ace9c5d))

# [2.1.0](https://github.com/puppetlabs/horsehead/compare/v2.0.1...v2.1.0) (2019-10-22)


### New

* Add JSON support for encoding/transfer ([dc59f1999540d228e32b29b3018acd880f01bc1d](https://github.com/puppetlabs/horsehead/commit/dc59f1999540d228e32b29b3018acd880f01bc1d))

## [2.0.1](https://github.com/puppetlabs/horsehead/compare/v2.0.0...v2.0.1) (2019-09-11)


### Fix

* Comply with Go module versioning rules ([89917dea170c47e73aaff2b54dacf7a179a52787](https://github.com/puppetlabs/horsehead/commit/89917dea170c47e73aaff2b54dacf7a179a52787))

# [2.0.0](https://github.com/puppetlabs/horsehead/compare/v1.11.0...v2.0.0) (2019-09-11)


### Breaking

* Revise scheduler to correctly propagate error behavior ([db6c5f1c90160b25466dd3950a4a7ea5cb195ce6](https://github.com/puppetlabs/horsehead/commit/db6c5f1c90160b25466dd3950a4a7ea5cb195ce6))

# [1.11.0](https://github.com/puppetlabs/horsehead/compare/v1.10.0...v1.11.0) (2019-08-15)


### New

* BlobStore.Get range request ([8b8f381036b85761c48daf6d4cdaab4bbadb0688](https://github.com/puppetlabs/horsehead/commit/8b8f381036b85761c48daf6d4cdaab4bbadb0688))

# [1.10.0](https://github.com/puppetlabs/horsehead/compare/v1.9.0...v1.10.0) (2019-08-14)


### New

* Add encoding/transfer package ([66198273ffd1cd9cc9ddafbad9e3d6cb352789bc](https://github.com/puppetlabs/horsehead/commit/66198273ffd1cd9cc9ddafbad9e3d6cb352789bc))
* Add secrets/encoding package to standardize secret value encoding for storage ([d188a7b75c35dacec656c431fc335452e061ff5b](https://github.com/puppetlabs/horsehead/commit/d188a7b75c35dacec656c431fc335452e061ff5b))

# [1.9.0](https://github.com/puppetlabs/horsehead/compare/v1.8.0...v1.9.0) (2019-08-12)


### New

* storage/testutils makes it easy to generate a temp directory and use the file:// storage backend for testing. ([4d043182bab5075f5ee807f308845ed6e3c2d531](https://github.com/puppetlabs/horsehead/commit/4d043182bab5075f5ee807f308845ed6e3c2d531))

# [1.8.0](https://github.com/puppetlabs/horsehead/compare/v1.7.0...v1.8.0) (2019-08-08)


### Update

* adding a utility package for managing working directories ([50a74c4a5edcaa40313748cce3e0fa7894f6a736](https://github.com/puppetlabs/horsehead/commit/50a74c4a5edcaa40313748cce3e0fa7894f6a736))

# [1.7.0](https://github.com/puppetlabs/horsehead/compare/v1.6.1...v1.7.0) (2019-08-02)


### New

* storage API. ([b68b74287f7bc210ed8fd355ae79893ec40a14d9](https://github.com/puppetlabs/horsehead/commit/b68b74287f7bc210ed8fd355ae79893ec40a14d9))

## [1.6.1](https://github.com/puppetlabs/horsehead/compare/v1.6.0...v1.6.1) (2019-08-02)

# [1.6.0](https://github.com/puppetlabs/horsehead/compare/v1.5.0...v1.6.0) (2019-07-30)


### New

* Expose TrackingResponseWriter ([fc89e1fb201012cdefc3bd0ab6b712df69d98ec9](https://github.com/puppetlabs/horsehead/commit/fc89e1fb201012cdefc3bd0ab6b712df69d98ec9))

# [1.5.0](https://github.com/puppetlabs/horsehead/compare/v1.4.0...v1.5.0) (2019-07-26)


### New

* Add support for parsing Range header ([5e8297f9536abc6374b54141638e4b69082a4fec](https://github.com/puppetlabs/horsehead/commit/5e8297f9536abc6374b54141638e4b69082a4fec))

# [1.4.0](https://github.com/puppetlabs/horsehead/compare/v1.3.0...v1.4.0) (2019-06-18)


### New

* Import insights-logging ([594b93c7eeb169878416ccb0d8513e974d8d64f7](https://github.com/puppetlabs/horsehead/commit/594b93c7eeb169878416ccb0d8513e974d8d64f7))

# [1.3.0](https://github.com/puppetlabs/horsehead/compare/v1.2.0...v1.3.0) (2019-06-18)


### Chore

* Import from insights-stdlib ([5c3eea2dcc0b8a51b0bffba30b80909dac928d4b](https://github.com/puppetlabs/horsehead/commit/5c3eea2dcc0b8a51b0bffba30b80909dac928d4b))
* Release 1.0.0 ([c75017ad0897652de4aa9d33490be398c8a1a1a7](https://github.com/puppetlabs/horsehead/commit/c75017ad0897652de4aa9d33490be398c8a1a1a7))
* Release 1.1.0 ([5678ee08ab0548a9cf3592b02bb081c77231df27](https://github.com/puppetlabs/horsehead/commit/5678ee08ab0548a9cf3592b02bb081c77231df27))
* Reset CHANGELOG to fix tagging ([94cef6fe804cc441f6def6c408e434b0eb3a9a16](https://github.com/puppetlabs/horsehead/commit/94cef6fe804cc441f6def6c408e434b0eb3a9a16))

### New

* Import insights-instrumentation ([7d469842c310e447b7f70fb2b65eab9e10e0d686](https://github.com/puppetlabs/horsehead/commit/7d469842c310e447b7f70fb2b65eab9e10e0d686))

### Trivial

* align noop capturer with Sentry (again) ([25a3532380c155af6ba254886ac810d5fd73a473](https://github.com/puppetlabs/horsehead/commit/25a3532380c155af6ba254886ac810d5fd73a473))

# [1.2.0](https://github.com/puppetlabs/horsehead/compare/v1.1.1...v1.2.0) (2019-06-17)


### New

* Add sqlutil.WithTx() for handling nested SQL transactions ([b15b3871dbe1ab34add8cdaaa53b996aa3f7441e](https://github.com/puppetlabs/horsehead/commit/b15b3871dbe1ab34add8cdaaa53b996aa3f7441e))

## [1.1.1](https://github.com/puppetlabs/horsehead/compare/v1.1.0...v1.1.1) (2019-02-08)


### Fix

* Make api.SetContentDispositionHeader thread-safe ([4982868e2b1c9cfcb072ef649af35463495b424d](https://github.com/puppetlabs/horsehead/commit/4982868e2b1c9cfcb072ef649af35463495b424d))

# [1.1.0](https://github.com/puppetlabs/horsehead/compare/v1.0.1...v1.1.0) (2019-02-05)


### New

* Add Content-Disposition header filename sanitizer ([5936f7b8e44cdd727a26b55f483a7ae873e0459a](https://github.com/puppetlabs/horsehead/commit/5936f7b8e44cdd727a26b55f483a7ae873e0459a))

## [1.0.1](https://github.com/puppetlabs/horsehead/compare/v1.0.0...v1.0.1) (2019-02-04)


### Fix

* Update changelog format correctly ([53fd55c4254bb8ada045b284e6b59a966e2b7f4f](https://github.com/puppetlabs/horsehead/commit/53fd55c4254bb8ada045b284e6b59a966e2b7f4f))
