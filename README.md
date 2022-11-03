# Food-track-be

## Description
This app is a simple food tracking app. It allows you to track your food intake and calories.

To use this app, first you have to install [grocery-be](https://github.com/nico-iaco/grocery-be) which provide food availability
and all the food data.

## Features

- [x] Add meals
- [x] Add food consumed
- [x] Calculate meal calories and price


## Requirements

- [PostgreSQL](https://www.postgresql.org/)
- [grocery-be](https://github.com/nico-iaco/grocery-be)

## Installation

### Cluster installation

To install this app in a cluster, first create grocery namespace, then modify the kustomization.yaml file in /k8s/overlays/qa 
changing the property to match your configuration and run the following command:

```bash
kubectl apply -k k8s/overlays/qa
```

### Local installation

You can run this app locally with docker. To do so, run the following command:

```bash
docker run -p 8080:8080 ghcr.io/nico-iaco/food-track-be:latest -e {ALL_ENV_VARIABLES}
```

## Environment variables

| Name                 | Description                                  | Default value |
|----------------------|----------------------------------------------|---------------|
| PORT                 | Port on which the app will listen            | 8080          |
| GIN_MODE             | Release type of app                          |               |
| DB_HOST              | Database host                                |               |
| DB_PORT              | Database port                                |               |
| DB_NAME              | Database name                                |               |
| DB_USER              | Database user                                |               |
| DB_PASSWORD          | Database password                            |               |
| GROCERY_BASE_URL     | Base url for grocery-be app                  |               |

## Database

To create the database, run the following command with the database user:

```sql
CREATE DATABASE food_track;
```

```sql
create table meal (
    id uuid primary key,
    user_id varchar(255) not null,
    name varchar(255) not null,
    description varchar(255),
    meal_type varchar(255) not null,
    date date not null
);
```

```sql
create table food_consumption (
    id uuid primary key,
    meal_id uuid not null,
    food_id uuid not null,
    transaction_id uuid not null,
    food_name varchar(255) not null,
    quantity_used float not null,
    quantity_used_std float not null,
    unit varchar(255) not null,
    kcal float not null,
    cost float not null,
    foreign key (meal_id) references meal(id)
);
```
