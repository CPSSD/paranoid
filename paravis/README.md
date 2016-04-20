Paravis
=======
Experimental Paranoid Visualizer

---

## Features ##
Paranoid Visualizer shows all actions occurring in
the cluster in a nice visual way

Features (both current and planned):
- [x] NWJS framework
  - [ ] Multiple Activity support
  - (only basic version is currently available, more features will be added)
- [ ] Downloading necessary binaries
- [ ] Filesystem operations
  - [ ] Creating new filesystems
  - [ ] Modifying existing ones
- [ ] Visualizing the data flow
  - [ ] Leader
  - [ ] Peers

---

## Contributing ##
If you wish to contribute please make sure to read the steps outlined in the top level `CONTRIBUTING.md` file. The following rules only apply to __paravis__

Please clone the repo into `$GOPATH/src/github.com/cpssd/paranoid` as the backend is written in Go
and won't run outside `$GOPATH`

After cloning make sure to run `npm install`

Here are the steps to set up paravis for development:
To enable Chrome Dev Tools, download SDK build of nwjs available [here](http://nwjs.io).


The following instructions are just a matter of personal preference for
installing the NWJS SDK, you can install it wherever you feel appropriate for
your system. Please read over the commands:

1. Create a new folder: `sudo mkdir -p /opt/nwjs-sdk`  
2. Navigate to the folder `cd /opt/nwjs-sdk`  
3. Extract the nwjs SDK: `sudo tar xvf ~/Download/$nwjs-sdk.tar.gz`  
4. Move everything to from extracted folder back: `sudo mv nwjs-sdk/* ./`  
5. Create new group: `sudo groupadd sdkusers`  
6. Add yourself to the group: `sudo usermod -aG sdkusers $USER`  
7. Change nwjs permissions: `sudo chown -R :sdkusers .`  
8. Copy the user permissions to group permissions using
`sudo chmod g+x $thing`, `sudo chmod g+r $thing` and similar.  
9. Consult `man` pages for `chmod`
if you need help. You will get errors until you finish this step  
10. Add SDK to PATH: `export PATH=$PATH:/opt/nwjs-sdk` (and remember to `source`)

---

## Building ##
If you wish to build __paravis__ into standalone application please follow the steps from official nwjs documentation available [here](http://docs.nwjs.io/en/latest/For%20Users/Package%20and%20Distribute/)

---

## Troubleshooting ##
You might run into problems while developing or running. Here is a compilation of the possible issues that were found. If you find more issues please make sure to report it with [Issues](https://github.com/cpssd/paranoid/issues)

#### I can't install start NWJS SDK ####
This is probably caused by bad permissions. Please make sure that `nw` can be executed, `.so`, `.pak`, `.bin` files read. Do not forget about the `locales` and `lib` folder. Just play around with it until it work ;)

#### I can't use nwbuild ####
This is most likely a issue with your permissions for npm. Not sure how to fix it :/

---

## Known issues ##
This list contains issues that are very hard to fix or we didn't really look at it yet
