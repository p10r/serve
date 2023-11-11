# Serve

Get the schedule your favourite volleyball leagues straight to your discord channel of choice.

The outcome looks as such:
![img.png](img.png)

## Deployment
```bash
brew install flyctl

flyctl auth login

fly launch

# set env variables
fly secrets set API_KEY=<flashscore-api-key>                                                                                                                        
fly secrets set DISCORD_URI=<full-discord-uri>
fly secrets list

# run only one container to execute CRON only once
fly scale count 1

# update running app
fly deploy
```