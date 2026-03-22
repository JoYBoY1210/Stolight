# TODO

## Backend / Storage Improvements

- [x] Improve the HTTP endpoints  
  *Previously messy — now cleaned up and fixed.*

- [ ] Add presigned temporary URLs  
  *Avoid proxying uploads/downloads through the backend.  
  Frontend should directly make requests to the storage provider.*

- [x] Add bucket existence validation during project creation  
  *When an admin creates a new project, the backend must verify that the specified bucket actually exists.*

- [x] Allow admins to update bucket lists later  
  *Provide functionality to modify the list of allowed buckets after project creation.*

- [] Instead of sending the whole file over one http req, split them in numerous http req and then send the chunks in together. Helps to prevent restart due to network disruptions.

---

## Priority Requirement

- [x] Add a way to change the list of allowed buckets for a project  
- [x] EC in async. Keep the uploaded files in a temporary folder on server, return 200 OK to user. push this   job to Queue then after that do everything in the background.
- [x] FUCKKKKKK. right now i am not allowing files wiht the same name to exist in the same bucket.
- [x] complete the staging part(u need to add the worker in the end).
- [x] Change middleware to work with fileId instead of FileName.
- [x] change download handler.
- [x] update CLI to use fileID instead of name.
- [x] Update Delete File handler.
- [] Change the directories of everything to a permanent location(rn it is ./whatever).

---

## CLI Features

- [x] CLI for download support  

# extras
- [x] Hash the pieces and store hashes in DB, this prevents corruption. Check hashes when stitching the files.
- [] Garabage collector.
- [] compression before parity bits are formed.
