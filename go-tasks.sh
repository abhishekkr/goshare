#!/bin/bash

# managing go deps
#############################################################################
##### from github.com/abhishekkr/dotfiles/shell_profile/a.golang.sh
goenv_on_at(){
  if [ $# -eq 0 ]; then
    _GOPATH_VALUE="${PWD}/.goenv"
  else
    cd "$1" ; _GOPATH_VALUE="${1}/.goenv" ; cd -
  fi
  if [ ! -d $_GOPATH_VALUE ]; then
    mkdir -p "${_GOPATH_VALUE}/site"
  fi
  export _OLD_GOPATH=$GOPATH
  export _OLD_PATH=$PATH
  export GOPATH=$_GOPATH_VALUE/site
  export PATH=$PATH:$GOPATH/bin

  echo "your new GOPATH is at $GOPATH"
}
alias goenv_on="goenv_on_at \$PWD"
alias goenv_off="export GOPATH=$_OLD_GOPATH ; export PATH=$_OLD_PATH ; unset _OLD_PATH ; unset _OLD_GOPATH"

go_get_pkg_help(){
  echo "go_get_pkg handles your Golang Project dependencies."
  echo "* Create new dependency list or install from existing:"
  echo "  $ go_get_pkg"
  echo "* Install from existing with updated dependencies"
  echo "  $ GO_GET_UPDATE=true go_get_pkg"
  echo "* Install from existing with re-prepared binaries (required on new Golang update or local changed dependency code)"
  echo "  $ GO_GET_RENEW=true go_get_pkg"
  echo "* Install from existing with updated dependencies (re-prepared binaries even if no updates)"
  echo "  $ GO_GET_RENEW=true GO_GET_UPDATE=true go_get_pkg"
}
go_get_pkg_list_create(){
  if [ ! -f "$1" ]; then
    PKG_LISTS_DIR=$(dirname $PKG_LISTS)
    mkdir -p "$PKG_LISTS_DIR" && unset PKG_LISTS_DIR
    touch "${1}"
    echo "Created GoLang Package empty list ${PKG_LISTS}"
    echo "Start adding package paths as separate lines."
    return 0
  fi
  return 1
}
go_get_pkg_install(){
  for pkg_list in $PKG_LISTS; do
    cat $pkg_list | while read pkg_path; do
        echo "fetching golag package: go get ${pkg_path}";
        pkg_import_path=$(echo $pkg_path | awk '{print $NF}')
        if [[ ! -z $GO_GET_RENEW ]]; then
          rm -rf "${GOPATH}/pkg/${GOOS}_${GOARCH}/${pkg_import_path}"
          echo "cleaning old pkg for ${pkg_import_path}"
        fi
        if [[ -z $GO_GET_UPDATE ]]; then
          echo $pkg_path | xargs go get
        else
          echo $pkg_path | xargs go get -u
        fi
    done
  done

  unset GO_GET_UPDATE GO_GET_RENEW
}
go_get_pkg(){
  if [[ "$1" == "help" ]]; then
    go_get_pkg_help
    return 0
  fi

  if [[ $# -eq 0 ]]; then
    PKG_LISTS="$PWD/go-get-pkg.txt"
  else
    PKG_LISTS=($@)
    if [[ -d "$PKG_LISTS" ]]; then
      PKG_LISTS="${PKG_LISTS}/go-get-pkg.txt"
    fi
  fi
  go_get_pkg_list_create $PKG_LISTS
  if [[ $? -eq 0 ]]; then
    return 0
  fi

  if [[ -z $GO_GET_ENV ]]; then
    _GO_GET_ENV=$(dirname $PKG_LISTS)
    GO_GET_ENV=$(cd $_GO_GET_ENV ; pwd ; cd - >/dev/null)
  fi
  goenv_on_at $GO_GET_ENV

  go_get_pkg_install "$PKG_LISTS"

  unset _GO_GET_ENV GO_GET_ENV PKG_LISTS
}
##############################################################################


_OLD_PWD=$PWD
cd $(dirname $0)

if [[ "$1" == "deps" ]]; then
  go_get_pkg

elif [[ "$1" == "gorun" ]]; then
  go run -tags zmq_4_x "$2"

elif [[ "$1" == "test" ]]; then
  $0 build
  echo
  echo "~~~~~Test Pieces~~~~~"
  go test ./...
  echo
  echo "~~~~~Test Features~~~~~"
  for feature_test in `ls ./tests/go*_client.go`; do
    echo ">> Testing: "$feature_test
    ./bin/goshare_daemon -daemon=start -dbpath=/tmp/GOSHARE.TEST.DB
    go run $feature_test
    ./bin/goshare_daemon -daemon=stop
    rm -rf /tmp/GOSHARE.TEST.DB
  done

elif [[ "$1" == "wiki" ]]; then
  $0 bin
  echo
  echo "~~~~~Visit wiki at GoShare HTTP~~~~~"
  echo "~~~~~   http://0.0.0.0:9999    ~~~~~"
  ./bin/goshare_server

elif [[ "$1" == "build" ]]; then
  bash $0 deps
  goenv_on_at $PWD
  mkdir -p ./bin
  cd ./bin
  for go_code_to_build in `ls ../zxtra/goshare_*.go`; do
    echo "Building: "$go_code_to_build
    go build -tags zmq_4_x $go_code_to_build
  done

else
  echo "Use it wisely..."
  echo "Build usable binaries: '$0 build'"
  echo "Install tall Go lib dependencies: '$0 deps'"
  echo "Run all Tests: '$0 test'"

fi

cd $_OLD_PWD
