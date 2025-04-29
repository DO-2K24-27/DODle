# DODle
Like wordle but guess DO related people

## Languages

- For front we will use Next.js
- For backend we will use Go
- For database we will use PostgreSQL
- For deployment we will use Docker Container and Kubernetes
- For hosting we will use our hardware
- For CI/CD we will use Github Actions

## Features

With DODle you will be able to guess DO related people every day.
- You have infinite tries to guess the person
- Each guess will give you a hint on which parameters are correct or not
- When you find it you will be able to share it with your friends with leaderboard
- You will be able to see the leaderboard of the day
- You will be able to see the leaderboard of the week
- You will be able to see the leaderboard of the month

## Workflow

1. Everyday at 10am a new person will be picked from the database
2. The person will be stored in the database
3. Every client will be able to guess the person of the day
4. Each guess will be stored in the database
5. Each guess will give you a hint on which parameters are correct or not
6. When you find it your score will be sent in the database

Every day we will pick a new person from the database and store it in the database
Every day we will reset the database of guesses
Every month we will reset the database of leaderboard

## Participants

- Pierre-Louis will handle the backend
- Joris will handle the frontend