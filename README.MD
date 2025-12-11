### what is this?
gator is a command-line RSS reader written in Go. it fetches feeds, stores them in a Postgres database and provides commands to manage feeds and read posts.

it's a CLI tool that allows users to:
- add RSS feeds from across the internet to be collected
- store the collected posts in a PostgreSQL database
- follow and unfollow RSS feeds that other users have added
- view summaries of the aggregated posts in the terminal, with a link to the full post


### why?
i'm learning go because i need it for 95 [check it out](https://95ninefive.dev)

- learning how to integrate a Go application with a PostgreSQL database
- practicing SQL to query and migrate a database (using sqlc and goose)
- learning how to write a long-running service that continuously fetches new posts from RSS feeds and stores them in the database


## pre-requisites
- download and install Go
- verify the installation:
```sh
go version
```
- start Postgres and create a database (name it gator)
- install the Gator CLI via `go install` or clone this project and use `go build`


### notes (for me)
> general structure

```XML
<rss xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
<channel>
  <title>RSS Feed Example</title>
  <link>https://www.example.com</link>
  <description>This is an example RSS feed</description>
  <item>
    <title>First Article</title>
    <link>https://www.example.com/article1</link>
    <description>This is the content of the first article.</description>
    <pubDate>Mon, 06 Sep 2021 12:00:00 GMT</pubDate>
  </item>
  <item>
    <title>Second Article</title>
    <link>https://www.example.com/article2</link>
    <description>Here's the content of the second article.</description>
    <pubDate>Tue, 07 Sep 2021 14:30:00 GMT</pubDate>
  </item>
</channel>
</rss>
```

> running this app will create a hidden file named .gatorconfig.json in your home directory. this is used to track the username you provide (or change it) and the database URL.
