# appUtils
The goal of the project is making tool allowing for testing and deploying applications using docker compose. Example usage of deployment pipeline:
`apputils deploy -f compose.yml -f compose.prod.yml cron -c '0 0 * * *'` will check if there are any updates of the code in git repo and if there are any it will
redeploy the app at given interval, in this case at midnight every day.

My next goal is making sure the app passes tests before redeploying.