## github-archive

simple command line tool that allows one to archive an organizations [Github](https://github.com/) 
repositories. It assumes that where it is being run already has access to the repositories by way 
of the [git](https://git-scm.com/) command over [SSH](https://help.github.com/articles/connecting-to-github-with-ssh/).

You'll need to also generate a [Token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) to
access the Github API to find the list of repositories to archive.

The repositories will be cloned or updated (if they already exist) in the current working directory from where `github-archive`
is invoked.

### examples

Archive all repositories for an organization.

`github-archive --token TOKEN --org MyCoolDotCom`

Archive a single repository for an organization.

`github-archive --token TOKEN --org MyCoolDotCom --repo github-archive`

increase the number of archive processes from the default of four (4) to ten (10)

`github-archive --token TOKEN --org MyCoolDotCom --max 10`

