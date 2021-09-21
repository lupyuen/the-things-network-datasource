::  Build the Grafana Data Source

::  Add NPM and Go to PATH
@PATH=%PATH%;%APPDATA%\npm;%GOPATH%\bin;

::  Build the project
yarn build
