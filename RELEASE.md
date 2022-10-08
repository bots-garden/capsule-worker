# Manual release process

- Change the version number (`vN.N.N`) in:
  - `./commons/version.go`
  - `/README.md`
  - `install-capsule-worker.sh`
- 🖐 Check **every dependency** for every module
- Update and run `update-modules-for-release.sh`
- On GitHub: create a release + a tag (`vN.N.N`)
