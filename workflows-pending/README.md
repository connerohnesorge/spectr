# Pending Workflow Files

Due to GitHub App permissions restrictions, workflow files must be manually moved to `.github/workflows/`.

## To activate the CI workflow:

```bash
mkdir -p .github/workflows
mv workflows-pending/ci.yml .github/workflows/ci.yml
git add .github/workflows/ci.yml
git rm -r workflows-pending
git commit -m "Activate CI workflow"
git push
```

The CI workflow is ready to use once moved to the correct location.
