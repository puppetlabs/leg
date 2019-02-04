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

Specifically, we require commit messages to use the following format exactly:

```
<Tag>: <Short description> (<Op> <PI-NNN>)

<Long description>
```

`<Tag>` is one of the following:

| Name | Use |
|------|-----|
| `Fix` | A bug fix |
| `Update` | A backward-compatible enhancement |
| `New` | A new feature |
| `Breaking` | A backward-incompatible enhancement or feature |
| `Docs` | A change to internal (non-user-facing) documentation only |
| `Build` | A change to build processes only |
| `Upgrade` | A dependency upgrade |
| `Chore` | Refactoring, adding tests, and other non-user-facing changes |

`<Op> <PI-NNN>` references the Jira ticket `PI-NNN` if applicable. For example,
if the commit resolves PI-123, you would write `fixes PI-123`. If the commit
does not completely resolve the issue, you would write `refs PI-123` instead.

A full example of a commit message might be:

```
Fix: Stop consuming unbounded memory (fixes PI-123)

Sometimes users search for their own name. This caused us to run out of memory.
We no longer allow users to search for their own name.
```

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
