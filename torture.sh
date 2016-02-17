torture(){
  #discovery-server 2> /dev/null &
  paranoid-cli init test1 -u 2> /dev/null
  paranoid-cli init test2 -u 2> /dev/null
  mkdir ~/test1
  mkdir ~/test2
  paranoid-cli mount 127.0.0.1:10101 test1 ~/test1 -n 2> /dev/null
  echo "Waiting for PFSD 1 mount"
  sleep 5s #waiting for both to come to life
  paranoid-cli mount 127.0.0.1:10101 test2 ~/test2 -n 2> /dev/null
  echo "Waiting for PFSD 2 mount"
  sleep 15s #waiting for both to come to life. It seems to take some time before paranoid-cli is ready to start
  echo "CP Test"
  cp /bin/cp ~/test1/
  cp /bin/mv ~/test2/
  sleep 10s #Making sure PFSD has enough time to transfer
  if [ -z "$(diff -arq ~/test1 ~/test2)" ]; then
    echo "no Diff Found"
  else
    echo "ERROR: Diff was found"
    echo $(diff -arq ~/test1 ~/test2)
    exit 1
  fi

  if [ -z "$(cmp ~/test1/cp ~/test2/cp)" ]; then
    echo "Files Match"
  else
    echo "ERROR: Files Didnt Match"
    exit 1
  fi
}

cleanup(){
  echo "Cleaning Up Tests"
  paranoid-cli unmount test1
  sleep 5s
  paranoid-cli unmount test2
  sleep 5s
  paranoid-cli delete test1
  paranoid-cli delete test2
  rm -rf ~/test2 ~/test1
}

torture
sleep 5s
cleanup
