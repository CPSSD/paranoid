compare() {
  if [ -z "$(diff -arq ~/test1 ~/test2)" ]; then
    echo "Files in both instances match"
  else
    echo "ERROR: Diff was found: $(diff -arq ~/test1 ~/test2)"
    exit 1
  fi
}

blackBox() {
  discovery-server 2> /dev/null &
  sleep 5s #Wait for discovery server to start
  echo "Starting up Paranoid system"
  paranoid-cli init test1 -u -p testPool 2> /dev/null
  paranoid-cli init test2 -u -p testPool 2> /dev/null
  mkdir ~/test1
  mkdir ~/test2
  yes | paranoid-cli mount test1 ~/test1 -d 127.0.0.1:10101 2> /dev/null
  echo "Waiting for PFSD 1 mount"
  sleep 5s #waiting for both to come to life
  yes | paranoid-cli mount test2 ~/test2 -d 127.0.0.1:10101 2> /dev/null
  echo "Waiting for PFSD 2 mount"
  sleep 15s #waiting for both to come to life. It seems to take some time before paranoid-cli is ready to start
  echo "CP Test"
  cp /bin/cp ~/test1/
  cp /bin/mv ~/test2/
  sleep 10s #Making sure PFSD has enough time to transfer

  compare

  if [ -z "$(cmp ~/test1/cp ~/test2/cp)" ]; then
    echo "Files Match"
  else
    echo "ERROR: Files Didnt Match"
    exit 1
  fi
}

torture() {
  echo "Torture"
  cp -Rf /bin/ ~/test1/
  sleep 10s #Making sure PFSD has enough time to transfer
  echo "Testing Large Random Files"
  head -c 50M /dev/urandom > ~/test1/newRand.txt #This is going to take a while
  head -c 250M /dev/urandom > ~/test2/random.txt
  sleep 5s
  compare
  if [ -z "$(cmp ~/test1/random.txt ~/test2/random.txt)" ]; then
    echo "Files Match"
  else
    echo "ERROR: Files Didnt Match"
    exit 1
  fi
}

cleanup() {
  echo "Cleaning Up Tests"
  paranoid-cli unmount test1
  sleep 5s
  paranoid-cli unmount test2
  sleep 5s
  paranoid-cli delete test1
  paranoid-cli delete test2
  pkill discovery-server
  rm -rf ~/test2 ~/test1
}

blackBox
torture
sleep 5s
cleanup
