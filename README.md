# ssl-check

Simple utility that tells you expiration dates for all the domains specified in a zonefile. 


This tool is motivated by the problem many orgs often face with expiring certificates. Please note that a better solution to this problem is to use managed certificates like [Lets Encrypt]. However, many organizations continue to rely on manually provisioned ssl certs for various reasons, and this tool is meant to assist in maintaining them.

[Lets Encrypt]: https://letsencrypt.org/
