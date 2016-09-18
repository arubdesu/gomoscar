# gomoscar
Moscargo v2

# Build
`git clone git@github.com:arubdesu/gomoscar.git $GOPATH/src/github.com/arubdesu/gomoscar`
`cd $GOPATH/src/github.com/arubdesu/gomoscar`
`glide install`
`go build -i` or `go install`

## Installing and updating dependencies

Moscargo v2 uses [Glide](https://github.com/Masterminds/glide#install) to manage external dependencies

To install the latest version of required packages, use the `glide install` command, 
which will inspect the `glide.lock` file and download the correct versions into the `/vendor` folder.
