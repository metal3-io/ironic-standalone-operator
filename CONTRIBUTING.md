# How to Contribute

Metal3 projects are [Apache 2.0 licensed](LICENSE) and accept contributions via
GitHub pull requests.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Certificate of Origin](#certificate-of-origin)
   - [Git commit Sign-off](#git-commit-sign-off)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](DCO) file for details.

### Git commit Sign-off

Commit message should contain signed off section with full name and email. For example:

 ```text
  Signed-off-by: John Doe <jdoe@example.com>
 ```

When making commits, include the `-s` flag and `Signed-off-by` section will be
automatically added to your commit message. If you want GPG signing too, add
the `-S` flag alongside `-s`.

```bash
  # Signing off commit
  git commit -s

  # Signing off commit and also additional signing with GPG
  git commit -s -S
```
