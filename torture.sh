torture(){
  discovery-server 2> /dev/null &
  paranoid-cli init test1 -u
  paranoid-cli init test2 -u
  mkdir ~/test1
  mkdir ~/test2
  paranoid-cli mount 127.0.0.1:10101 test1 ~/test1 -n 2> /dev/null
  sleep 5s
  paranoid-cli mount 127.0.0.1:10101 test2 ~/test2 -n 2> /dev/null
  # touch ~/test1/file.html
  # #cp -Rf /bin/ ~/test1
  # sleep 5s
  # if [ $(diff -arq ~/test1 ~/test2) ]; then
  #   #Diff was found
  #   exit 1
  # else
  #   echo "no Diff Found"
#  fi
}

cleanup(){
  paranoid-cli unmount test1
  paranoid-cli unmount test2
  paranoid-cli delete test1
  paranoid-cli delete test2
  #rm -rf ~/test2 ~/test1
}

#torture
cleanup
