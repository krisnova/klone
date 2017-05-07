# klone

`klone` helps you make quick and easy work out of contributing to a GitHub repository by handling `git clone` for you.

`klone` is written in [Go](https://golang.org/) but is designed to work with any programming language.

`klone` looks to see what **programming language** you are trying to `git clone` and puts it in a smart place accordingly.

With a few simple commands, you can have `klone` do exactly what you want, and you can always turn it off.

# You can (literally) do this:

```
$ go get github.com/kris-nova/klone
$ klone kubernetes/kubernetes
```

#  Here is what it does:

 - Authenticates you with your GitHub account
   - Sure, you can set `$VARIABLES` here, but just run the command and `klone` will prompt you and cache in `~/.klone/auth`
 - Looks up `kubernetes/kubernetes` at runtime
   - Then, discovers it's a GitHub repository
 - Checks your account to see if you have already forked it
   - Don't worry if you haven't we will take care of that
 - Reasons about what needs to be done to get things to a [desired state](https://github.com/kris-nova/klone#desired-state)
 - Looks up a `kloner` based on things like :
   - What **Programming Language** the repository is written in?
   - Does the repository explicitly call out a `.Klonefile`?
   - Don't worry - we have a simple `kloner` we always default to..
 - Checks your `git` configuration (yes we use the C library, so it's a real check)
 - Satisfies all concerns with your unique `git` configuration, and our [desired state](https://github.com/kris-nova/klone#desired-state)
 - Makes the [desired state](https://github.com/kris-nova/klone#desired-state) so (we actually `git clone` a repo, and `git checkout` out for you)
 - Drop you off (`cd`) in your new local workspace
 - You can now `git push origin master` to push to your fork
 - You can also `git rebase upstream/master` to rebase your fork


> GitHub is happy. You are happy. No conflicts. Just good clean contribution, the way you want it.

# Desired State

After a `klone` you should have the following `git remote -v` configuration

| Remote        | URL                                         |
| ------------- | ------------------------------------------- |
| origin        | git@github.com:$owner/$repo                 |
| upstream      | git@github.com:$parent/$repo                |

The goal here is also make it so your GitHub account is happy with this new configuration and you can

 - `push`
 - `fetch`
 - `pull`
 - `rebase`
 - `etc..`

without much trouble.

# Kloners

`klone` is designed to offer opinionated `kloners` for many programming languages.

 - [Go](https://github.com/kris-nova/klone#go)

## Go

 - `klone` **will** respect your `$GOPATH`'s (Yes, more than one)
 - `klone` **will** check out the **parent** repository into your gopath (so it compiles), but you will **always be** ``origin``!

Ex:

```
klone golang/dep
```

 - Will ensure you have a working fork of `golang/dep` at `$owner/dep` in GitHub
 - Will checkout `$owner/dep` to `golang/dep`

 (Hint: `$owner` is your GitHub login)

 | Local Path                               | Remote     | URL                          |
 | ---------------------------------------- | ---------- | ---------------------------- |
 | `$GOPATH`/src/github.com/$parent/$repo   | origin     | git@github.com:$owner/$repo  |
 | `$GOPATH`/src/github.com/$parent/$repo   | upstream   | git@github.com:$parent/$repo |


Also with custom language `kloner`'s and `.Klonefile`'s you could even have `klone` run custom logic **after** cloning (like checkout out dependencies!)


# GitHub Credentials

`klone` will access the GitHub API either by your username and password, or via an access token.

**Ideally you shouldn't do anything here.**
Just run `klone` and enter your username and password. (Don't worry if you use MFA, we will prompt you)
We store them in memory for the duration of the program's execution.
We never write them to disk.

`klone` will (and only with your account credentials) then create a unique Auth token, that we will use moving forward.
Delete it at any time, and yes we leave a very visible note on the token letting you know where it's from.

If you want it back, just run `klone` again.

We **will however** cache your auth token to disk in `~/.klone/auth`
Go ahead and delete it whenever you like, we will delete/create tokens as necessary and never leave a mess.

Also there are some friendly Environmental Variables you can use here:

```bash
export KLONE_GITHUBTOKEN="..abc123.."
```


```bash

# Both of these work fine:
export KLONE_GITHUBUSER="@kris-nova"
#export KLONE_GITHUBUSER="kris@nivenly.com"

# Not encrypted so be careful:
export KLONE_GITHUBPASS="password"

```
