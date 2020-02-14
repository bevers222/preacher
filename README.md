# ![preacher](doc/logo.png)

[![GitHub stars](https://img.shields.io/github/stars/bevers222/preacher.svg?style=social&label=Star)](https://github.com/bevers222/preacher)
> spread your message

### Purpose
preacher is a command line tool written in Go. Its job is to communicate with the Github API, pull down a specified file from every repository in an organization. Once they are delivered, you can mass evaluate them and apply corrections where necessary.

**Note:** preacher is only built to work with Github. Support may come in the next major version for other SCMS, but none in the works as of yet!

### Installation

#### The Go Way
If you have your local Go environment set up, just run this command to add preacher to your GOBIN

`go get github.com/bevers222/preacher/cmd/preacher`

#### The Manual Way
Download the compiled binary from the releases page and place it where it will be on your PATH.

### Commands
The current commands and their purpose are listed below:

- `fetch [file name] (default: Jenkinsfile)`: fetch will build a list of repositories in the organization and search the default branch in the repository for the given file. If no file is specified, it will grab the `Jenkinsfile`. 

### Flags
The current flags and their purpose are listed below:

- `debug (default: false)`: debug will tell preacher to print our available debug information as it runs. This information includes delivery messages, error messages, and a list of skipped repositories.
- `directory, d (default: ./preacher)`: directory will tell preacher where to create the folders and store the downloaded files. Default is a new directory called `preacher` in your present working directory.
- `branch, b (default: repository default branch)`: branch will tell preacher what branch to pull the files from in the repository. Preacher will use the default branch in every repository unless this flag is set.
- `help, h`: help shows the help command.
- `token, t (default: read from config file)`: token will tell preacher what Github API token to use when making calls. This is setup during the first run of preacher when you configure it. This flag should only be used to override the value in the config file.
- `org, o (default: read from config file)`: org will tell preacher what Github Organization to poll for repositories. This is setup during the first run of preacher when you configure it. This flag should only be used to override the value in the config file.


### Future Plans
v1 was released with just the fetch command. Once v1 is deemed stabe and I have the time, I'll start work on v2!

v2 will include the update command. This will allow you to change however many files you want and it will open pull requests to each repository with a change.


### Inspiration
This project was inspired by [octopus](https://github.com/uptake/octopus) and by the fact that I just want to build cool Go stuff.

### Help!
I'm happy to consider and PR's, changes, tips, etc. that you have! Feel free to reach out!
