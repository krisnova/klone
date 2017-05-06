# klone

Checkout repositories like a pro

### Klone a project

```
klone kops
```

 - Will check my `.gitconfig` and discover that I am `@kris-nova`
 - Will look up that I have a repository called `kops`
 - Will discover that the repository was forked from `kubernetes/kops`
 - Will check if either `~` or `$GOPATH` has a `.gitmodules` file, and add the repo to `.gitmodules`
 - Will checkout the `kris-nova/kops` codebase to `$GOPATH/src/github.com/kris-nova/kops`
 - Will add the remote `kris-nova/kops` as `origin`
 - Will add the remote `kubernetes/kops` as `upstream`
 - Will `cd` to the new directory after `klone`
