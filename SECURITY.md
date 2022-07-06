# Security Policy

Contact: [security@tawesoft.co.uk](mailto:security@tawesoft.co.uk)


## Announcements

It is our policy to publicly announce security issues and fixes through the 
GitHub "Security Advisories" feature for this repository.

You can subscribe to security announcements for this repository by
[configuring your "watching" settings](https://docs.github.com/en/account-and-profile/managing-subscriptions-and-notifications-on-github/setting-up-notifications/configuring-notifications#configuring-your-watch-settings-for-an-individual-repository)
and subscribing to security alerts.

On a case-by-case basis, we are prepared to pre-announce security issues and 
fixes to any downstream consumer of this repository who can provide evidence 
that any security issues would have a particularly high impact on their 
services, such as operators of a service that processes personal data or is 
used by a large volume of users. Email us with the subject
"join golib security-preannounce list" for details.


## Backporting fixes

Applicable security fixes will always be backported to legacy 
packages, to the previous two major versions of normal packages, and the 
previous one major version of candidate packages.

Security fixes may not always be backwards compatible, even between minor 
versions.


## Reporting a vulnerability

Please disclose responsibly so that we can notify the users of our software
with a fix and/or instructions, including a pre-announcement where appropriate.
Do not report security issues through the public issue tracker in the first
instance, unless it is being actively exploited in the wild.

Instead, please email information about vulnerabilities to
[security@tawesoft.co.uk](mailto:security@tawesoft.co.uk).

To help prioritise your report, format the subject line as follows:

`vulnerability: repository-url - description`

For example:

`vulnerability: github.com/example/repo - denial of service in foo/bar`

If you don't receive an acknowledgement within 48 hours, please contact us
through any contact method listed on the Tawesoft website.

If we have not fixed or disclosed a vulnerability after 90 days, then you may
reserve your right to disclose this publicly.
