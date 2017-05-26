# About

`klone` is a command line tool that makes it easy to fork and clone a repository locally.

# Installing

```bash
go get -u github.com/kris-nova/klone
```

# Example

```bash
klone kubernetes
```

1. Klone will look for a git server for the project.
2. In the case of `kubernetes` we will detect **github.com**.
3. Klone will then attempt to find the organization via the GitHub API, which will be the same value: `kubernetes`.
4. Klone will then attempt to find the repository via the GitHub API, which will be the same value: `kubernetes`.
5. After finding `github.com/kubernetes/kubernetes` klone check and see if you (the authenticated user) has forked the repository yet.
6. If needed, klone will use the GitHub API to fork the repo for you.
7. Klone will detect the `Go` programming language, and use the offical `Go` implementation for checking out the program.
8. Klone will then check out the repository with the following remotes:


After a `klone` you should have the following `git remote -v` configuration

| Remote        | URL                                         |
| ------------- | ------------------------------------------- |
| origin        | git@github.com:$you/$repo                 |
| upstream      | git@github.com:$them/$repo                |

# Kloners

`klone` is designed to offer opinionated `kloners` for many programming languages.

 - [C Kernel Modules](https://github.com/kris-nova/klone#c-kernel-module)
 - [Go](https://github.com/kris-nova/klone#go)
 - [Simple](https://github.com/kris-nova/klone#simple)

## C Kernel Module

 **Will** attempt to reason about your architecture and check out source accordingly.

| Architecture  | Path                                        |
| ------------- | ------------------------------------------- |
| Linux         | /usr/local/src                              |
| FreeBSD       | /usr/src                                    |

**Will** still create `origin` and `upstream` (even if you are the only owner!)

## Go

 - **Will** respect your `$GOPATH`'s (Yes, more than one)
 - **Will** check out the **parent** repository into your gopath (so it compiles), but you will **always be** ``origin``!


Also with custom language `kloner`'s and `.Klonefile`'s you could even have `klone` run custom logic **after** cloning (like checkout out dependencies!)

## Simple

 - `klone` will by default clone to your `$HOME` (`~`) directory
 - `klone` will still create `origin` and `upstream` (even if you are the only owner!)


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