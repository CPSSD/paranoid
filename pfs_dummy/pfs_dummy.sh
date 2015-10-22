#/bin/sh

read_fn () {
  OPTION=$1
  OFFSET=$2

  if [ "$OFFSET" = "" ]
    then
      echo "Dummy Content from Fake file"
  else
    echo "Dummy Content"
  fi
}

readdir_fn (){
  echo "my_file.sh"
}

stat_fn(){
  echo "{\"length\":10,\"ctime\":\"2015-10-22T03:08:47.541016+01:00\"," \
      "\"mtime\":\"2015-10-22T03:08:47.541016+01:00\"," \
      "\"atime\":\"2015-10-22T03:07:33.877016+01:00\"}"
}

if [ "$1" = "read" ]
  then
    read_fn $2 $3
elif [ "$1" = "readdir" ]
  then
    readdir_fn
elif [ "$1" = "stat" ]
  then
    stat_fn
elif [ "$1" = "write" ]
  then
    exit 0
else
  echo "Missing one or more args"
  exit 0
fi
