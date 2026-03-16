# food-planner-3
Food planner 3 (codename FoodSmash) is a web app for creating recipes, sharing them, and planning meals and weekly shops. 

## Features 

- Create and edit versioned recipes 

## Tech stack 

Frontend
- React
- TypeScript 
- Vite 
- Tailwind 
- React Query 
- GraphQL Codegen

Backend
- Go 
- GraphQL 
- PostgreSQL

## Quick start 

Prerequisites
- Docker
- Go
- NPM (or PNPM)

1. Clone branch locally 
2. Run `api` service in `docker-compose.yml` (starts up Postgres DB, migrates it, and then starts the API)
3. Run `npm install` in `/web`
4. Run `npm run dev` (starts up the front end)