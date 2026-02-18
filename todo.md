# TODO

## Backend / Storage Improvements

- [x] Improve the HTTP endpoints  
  *Previously messy — now cleaned up and fixed.*

- [ ] Add presigned temporary URLs  
  *Avoid proxying uploads/downloads through the backend.  
  Frontend should directly make requests to the storage provider.*

- [ ] Add bucket existence validation during project creation  
  *When an admin creates a new project, the backend must verify that the specified bucket actually exists.*

- [ ] Allow admins to update bucket lists later  
  *Provide functionality to modify the list of allowed buckets after project creation.*

---

## Priority Requirement

- [x] Add a way to change the list of allowed buckets for a project  

---

## CLI Features

- [x] CLI for download support  
