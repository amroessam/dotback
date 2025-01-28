# User Stories for dotback

## 1. `dotback login`

**User Story**  
- *As a user, I want to authenticate with GitHub (via username and PAT), so that I can securely access my private repositories and manage my backups.*

---

## 2. `dotback scan`

**User Story**  
- *As a user, I want to quickly scan my filesystem for relevant configuration files and installed applications, so that I know what can be backed up or restored.*

---

## 3. `dotback scan --verbose`

**User Story**  
- *As a user, I want to see a detailed, verbose list of all discovered configuration files and applications, so that I can verify every item that will be included in the backup.*

---

## 4. `dotback backup`

**User Story**  
- *As a user, I want an interactive wizard to back up my configs and apps to a chosen private GitHub repository, so that I can preserve my environment for easy restoration.*

**Key Points**  
- If not logged in, the wizard should prompt me to login.  
- I can either select an existing private repo or create a new one.  
- I can select or create a “machine” entry (based on hostname).  
- The backup process copies dotfiles/configs and commits them to GitHub.

---

## 5. `dotback restore`

**User Story**  
- *As a user, I want an interactive wizard to restore configuration files and apps from my private GitHub repository, so that I can quickly replicate my setup on a new or existing machine.*

**Key Points**  
- If not logged in, prompt for GitHub login.  
- I can select or create a repository if needed.  
- I can select which machine to restore (or the one matching my hostname).  
- The tool symlinks restored files to their appropriate locations on the local system.

---

## 6. Machine Management

**User Story**  
- *As a user, I want each backup to be tied to a “machine” (identified by hostname), so that I can maintain multiple macOS devices with distinct sets of configs and applications in the same repository.*