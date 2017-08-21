REPOWEB
---

Repoweb is a simple way to manage a collection of static websites checked into Gitlab but served from another machine.

For example you may have a nginx or apache web server with a number of virtual host directives each pointing to the working directory 
of an individiual git repo; however, you might also have the the upstream version of these repos hosted on an external gitserver à la gitlab. Repoweb enables you to simply setup a git hook, that is called upon a push event that will update the local repositories to match origin/master without you having to log on to the 
webserver and manually checkout the changes.

## Installation and Configuration

 Download or build the binary and edit the `config.toml` to match your setup. For example:

 ```
[[website]]
name = "my-fancy-git-website"
worktree = "/srv/fancy"
gitdir = "/home/git/fancy-website.git"
token = "jJwixrBLuBfpZegsZYoNRKFvRd9HTzfDGQic3J8WO5XoNeXqyllT4CBKqyt5EWJe"
```

 * `name` is the unique identifier that is exposed by the webhook. 
 * `worktree` is the the webroot for a virtualhost. This location must be writable by the the user that repoweb is run under.
 * `gitdir` is the location of the bare git directory for the repository. It should also be writable.
 * `token` is checked against the `X-Gitlab-Token` sent in the hook request headers to verify that the request is legitimate.

You may add as many `[[website]]` blocks to `config.toml` as you want.

### Configuration Details 

The recomended setup is to have a working tree seprate from your gitdir stored to avoid leaking the contents of `.git` to the web. You can do this easily by cloning the bare repository with `git clone --bare http://<gitlabserver>/<user>/<my-site>.git <my-site>.git` this cloned path is what you will specify for `gitdir`. Please note that any changes you make in the work tree will be overwritten unless upon the next update so make sure that you do not edit the work tree directly.

You must also ensure that the user that runs repoweb has a ssh key that corresponds to a user on gitlab with at least read access to the remote repository. Instructions for setting up
an ssh key for a user can easily be found by googling.

Finally you must add the webhook to github. It will be of the form `http://<servername>:<port>/update/my-fancy-git-website`.

You may also want to add a system service to start repoweb as the desired user automatically on boot.

## Roadmap

  - [ ] `git pull` to a directory outside of the webroot and then use `mv` to atomically switch to the new version
  - [ ] Enable running post-pull scripts to run inside of a chroot or restricted shell
  - [ ] Support git lfs