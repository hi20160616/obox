# obox
Just an object collector for data files storage.

# Usage

1. `server` is used for a web site, port is rand generated by nanosecond.
2. `dbmanager` is used for turn on or turn of the database: folder db(configed at `configs/configs.json`)
    1. If you close the `server` don't forget run `dbmanager` to close the database for encrypt the data.
    2. If you open the `server` don't forget run `dbmanager` once to decrypt the data to open the database
3. Start Password configed at `configs/configs.json`

# TODO

- [x] upload attachments.
- [ ] optimized implements to solve load and walk too many times.
- [x] use walk instead of recurse in zipWriter
- [x] add open and close data to a single cmd exec
- [ ] innerLink invoke and use in render
- [ ] listAttachments and ListAttachments redundancy
- [ ] recurse dir but lost sub-dir, so the files in subfolder cannot be visited.
