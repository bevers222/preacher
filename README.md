# preacher

## spread your message

### Purpose
preacher is a command line tool written in Go. Its job is to communicate with the Github API, pull down a specified file from every repository in an organization. Once they are delivered, you can mass evaluate them and apply corrections where necessary.

### Installation

#### The Go Way
If you have go installed....

#### The Manual Way
grab the release ...

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