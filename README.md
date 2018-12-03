# Warships
[![Travis CI](https://api.travis-ci.org/jason-meredith/warships.svg?branch=master)](https://travis-ci.org/jason-meredith/warships)
[![Maintainability](https://api.codeclimate.com/v1/badges/fca64f446cb24582a339/maintainability)](https://codeclimate.com/github/jason-meredith/warships/maintainability)

An exciting new twist on Battleship, play Warships over the internet with countless friends and
enemies engaging in cut-throat warfare on the high seas.

Imagine Battleship with a few changes:
* **Teams** No one fights alone. Don't like your team? Switch. Don't like anyone? Mutiny and start your own.
* **Badass Command Line Interface** Ever want to feel like a submarine commander in an 80s Tom Clancy adaptation? It's your lucky day.
* **Essentially unlimited play area** The difference between fighting over a lake and domination of the Pacific is limited only by the server admins
imagination.
* **Control the ships** They're ships, after all. Need more? Spend deployment points to deploy more ships.

For more details - how to play, how to run, how to modify - check out the [wiki](https://github.com/jason-meredith/warships/wiki).

## Usage

Everything runs right out of the box. No external libraries to install (all handwritten myself). There are a few things you'll need
though...

### Prerequisites

So far only tested successfully on **Linux** - doesn't work on Windows *yet* - not sure about OSX.

```
$ wget -O warships https://github.com/jason-meredith/warships/releases/download/v0.1.0-alpha/warships-0.1.0-alpha
$ sudo chmod u+x warships
$ ./warships
```

Maybe you're a tinkerer... maybe you don't like how I made Warships... maybe you just want to make it your own

*Requires* ``` go version 1.11.2 ```

```
go get github.com/jason-meredith/warships/main
```


## Contributing

Feel free to make pull requests. I'd love all the help I can get.

## License

This project is licensed under the GNU General Public License - see the [LICENSE.md](LICENSE.md) file for details