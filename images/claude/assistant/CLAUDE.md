- Think deeply about what you are being asked to do. Do not use more code to solve a problem if modifying existing code or making small tweaks can do the job.

- Solving a problem with a smaller amount of elegant and clean code is always preferred

- Try as hard as possible to prevent codebase complexity from increasing linearly with problem complexity

- The ideal codebase solves many problems with a small amount of code, cleanly, simply, and generally

- If a .llmcontext file exists it has been created by another LLM for your benefit. Use it to improve context and understand the project before making changes.

- Always import the @README.md file in the local directory on startup if it exists. When asked to work on a task import any @./\*/README.md files in any sub directories that you come across.

- Always update any documentation, including the @README.md file if it exists, when you have made changes that warrant it

- Always update the .llmcontext file, if it exists, after you have made changes that any future LLM might find valuable or if the changes you have made obsolete previous entries. Remember this file content is to be compressed for LLM usage, and does not need to be human readable.

- Always use `cargo` with the `--all-features` flag on rust projects, never use rustc directly
