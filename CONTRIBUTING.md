# Puppet Insights Data Platform Contribution Guidelines

Open tickets against the
[Puppet Insights Jira project](https://tickets.puppetlabs.com/browse/PI). You
will want to create a ticket and confirm the scope of your work with the team
before starting to write code.

## Testing

```console
$ go test -v ./...
```

## Committing Changes

We use the
[ESLint commit message conventions](https://eslint.org/docs/developer-guide/contributing/pull-requests#step-2-make-your-changes).
By following this format, we can automatically generate a changelog and
determine the correct semantic version for the next release of this software.

We support using [Commitizen](http://commitizen.github.io/cz-cli/) to generate
messages that conform to this format. To use it for this project:

```console
$ npm install -g commitizen
$ npm install
$ git cz
```

## Contact

* team-insights-data@puppet.com
* #team-insights on Slack
