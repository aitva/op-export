# 1Password Export

I did not find an easy way to print all my passwords. So, I have built
this small app which query the 1Password server using their command line tool.

The app opens a webpage in your browser serving the passwords, so nothing
is written on your disk. But, there is an option to save the passwords
as CSV or HTML file.

## Instructions

1. [download and install](https://support.1password.com/command-line-getting-started/) the `op` command
2. [sign into your account](https://support.1password.com/command-line-getting-started/#get-started-with-the-command-line-tool) from the command line `eval $(op signin example)`
3. run `go get -u github.com/aitva/op-export`
4. run `op-export`
