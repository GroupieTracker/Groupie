# Groupie

* [The base](#the-base)
    * [Project Objective](#Project-Objective)
    * [Platform Features](#Platform-Features)
* [Game Descriptions](#Game-Descriptions)
    * [Guess the Song](#Guess-the-Song)
    * [Blind Test](#Blind-Test)
    * [Scattergories](#Scattergories)
    * [Instructions](#Instructions)
* [Required](#required)
    * [Golang](#golang)
* [Run Locally](#run-locally)

## The base

### Project Objective
The objective of this project is to design a web platform that provides access to three multiplayer music-related games. The platform will facilitate the following functionalities:

### Platform Features:
User Account Creation

Users can create an account on the platform. 

## Game Descriptions 

### Guess the Song
Guess the Song displays a snippet of lyrics from a song, and players must guess the original song title. Points are awarded based on the speed of each correct response, with points doubling each round. The round ends when all players have guessed correctly or when the time limit expires. At the end of the game, players are ranked based on their accumulated points.

### Blind Test
Blind Test plays a random song at the beginning of each round, and players must quickly identify the corresponding song. The song excerpts are randomly selected, and players can configure the minimum duration of the audio clips. Points are awarded based on the speed of each correct response, with options to play the audio clips in various modes such as slowed down, sped up, reversed, or normal. At the end of the game, players are ranked based on their accumulated points.

### Scattergories
Scattergories is an adaptation of the classic game where players must find words starting with a randomly assigned letter within different categories. The categories include Artist, Album, Music Group, Musical Instrument, and Featuring. The round ends when a player has provided a word for each category or when the time limit is reached. Points are awarded based on the validity and uniqueness of the answers. At the end of the game, players are ranked based on their accumulated points.

### Instructions
- Each game ends after a set number of rounds, determined at the start of the game.
- Players are ranked based on their total points at the end of each game.
- Customizable settings include round duration, audio clip options, and number of rounds per game.


## Required

### Browser

up-to-date browser

### Golang

Required version : `go1.21.9`

[Documentation of GOLANG](https://go.dev/doc/)

## Run Locally

**Clone the project**

```bash
git clone https://github.com/GroupieTracker/Groupie.git
```

**Go to the project directory**

```bash
cd Groupie
```

**Put your Ip in code for working**
Go on your IDE

For vscode ```code . ```

and replace all ```localhost```
by your IP or you're localhost if you want run locally


**Start game**

```bash
go run server.go
```


