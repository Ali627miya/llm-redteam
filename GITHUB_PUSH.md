# Push this project to GitHub (Ali627miya)

I can’t sign in to your GitHub account from Cursor. Run these commands **on your Mac** after you log in once.

## 1. Log in to GitHub CLI (one time)

```bash
gh auth login
```

Choose GitHub.com → HTTPS or SSH → authenticate in the browser.

## 2. Create the repo and push (from this folder)

```bash
cd "/Users/alimiyagilani/LLM PROJECT"

git init
git add .
git commit -m "Initial commit: LLM red-team framework and demo site"

gh repo create llm-redteam --public --source=. --remote=origin --push
```

That creates **https://github.com/Ali627miya/llm-redteam** and uploads everything.

### If `llm-redteam` already exists or you prefer another name

```bash
gh repo create YOUR-REPO-NAME --public --source=. --remote=origin --push
```

Then update links in `redteam/website/index.html` to match.

## 3. Without `gh` (website only)

1. On GitHub: **New repository** → name `llm-redteam` → create **without** README.
2. Then:

```bash
cd "/Users/alimiyagilani/LLM PROJECT"
git init
git add .
git commit -m "Initial commit: LLM red-team framework"
git branch -M main
git remote add origin https://github.com/Ali627miya/llm-redteam.git
git push -u origin main
```

Use a [Personal Access Token](https://github.com/settings/tokens) as the password if prompted.

## 4. Deploy the buyer-facing site (Vercel)

- New project → import **Ali627miya/llm-redteam**
- **Root Directory:** `redteam/website`
- Deploy

See `redteam/DEMO.md` for more detail.
