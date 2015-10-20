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

write_fn (){
  TYPE=$1

  if [ "$TYPE" = "-f" ]
    then
      echo 1
      exit 0
	#TODO Need to change
  elif [ "$TYPE" = "-n" ]
    then
      echo 1
      exit 0
  else
    echo "No option provided?"
    exit 1
  fi
}

readdir_fn (){
  echo "my_file.sh"
}

if [ "$1" = "read" ]
  then
    read_fn $2 $3
elif [ "$1" = "readdir" ]
  then
    readdir_fn
elif [ "$1" = "write" ]
  then
    write_fn $2
else
  echo "Missing one or more args"
  exit 1
fi
