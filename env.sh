
if [ "$#" -ne 1 ]; then
    echo "You forgot to mention which env you want.  Prod?  Dev?"
else
  export SLACK_TOKEN=`sneaker $1 d /console/slack/token`
  export GITHUB_TOKEN=`sneaker $1 d /console/github/token`
  export APP_ROOT=./dist
  export GOREST_CONFIG=$APP_ROOT/config
fi
