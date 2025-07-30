# Blog Aggregator

## Introduction
This blog aggregator was created as part of the boot.dev curriculum. It is a CLI tool that allows users to 
- Register with a username
- Add and follow RSS feeds from websites 
- Save the posts in a PostgreSQL database 
- Browse the saved posts.

## Setup
In order to use the blog aggregator, you will need to have **Go** and **PostgreSQL** installed. 

To install Go and PostgreSQL, please follow the instructions from the links below:

**Install Go**: https://go.dev/doc/install
**Install PostgreSQL**: https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql

Next, install the program by running `go install github.com/ehumba/blog-aggregator`. 

## Commands Overview
The blog aggregator is used by typing commands into the terminal. Most commands take one or more arguments. Here is an overview of the most important commands to get you started:

`blog-aggregator register username`
Registers a new user under the provided username and logs them in.

`blog-aggregator login username`
Logs the provided user in.

`blog-aggregator addfeed name url`
Adds a new feed and automatically follows it. The name argument is the name you want to give the new feed, the url argument is the blog's URL.

`blog-aggregator unfollow url`
Unfollows the feed from the given URL.

`blog-aggregator agg duration`
Scrapes all followed feeds and saves them to the database. The duration argument specifies the interval at which the program sends requests to the website (e.g. 10s, 3m, 1h, ...). Avoid intervals that are too short to prevent being rate-limited or blocked by websites.

`blog-aggregator browse limit`
Displays the recent posts from followed feeds. You can add an optional limit argument to limit the amount of posts displayed. If no argument is provided, the limit defaults to 2.