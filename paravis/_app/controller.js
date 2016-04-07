var win = nw.Window.get();
var gui = require('nw.gui')
var fs = require('fs');
var app = require('../App.json');

// List of available activities
var activities = [];


var A = {
  args: [], // CLI Argumens
  strings: {}, // All available strings
  name: gui.App.manifest.name, // Name of the application
  version: gui.App.manifest.version, // Version of the application
  debug: false, // Whether the debug mode is active
  lang: "en", // Language of the applications. Default: en
  mainActivity: {} // Object storing the main acitity
};


// Perform all actions when the window loads
$(window).load(function(){
  console.info("Initializing...");
  init();

  if(A.debug){
    gui.showDevTools();
  }

  console.info("Startup finished");
});


// Initialize everything
function init(){
  console.info("Starting main activity...");
  A.mainActivity = new activities[app.main]();
  console.info(app.main, "started.");

  // Load layout for the activity
  loadLayout(A.mainActivity.layout);

  // Call onCreate function of the activity
  A.mainActivity.onCreate();

  // Get the CLI arguments
  A.args = gui.App.argv;

  // Show DevTools if --debug argument is given
  if(A.args.indexOf("--debug") != -1 ){
    A.debug = true;
  }

  // Determine what language is to be used. Default is english
  var ix = 0;
  if((ix = A.args.indexOf("--lang")) != -1){
    A.lang = A.args[ix+1];
  }

  console.info("Loading strings...");
  // Load appropiate language text
  var langFile = fs.readFileSync("res/strings/strings_"+A.lang+".json");
  if(langFile == ""){
    langFile = fs.readFileSync("res/strings/strings_en.json");
  }
  if(langFile == ""){
    console.error("Unable to read strings file");
  }

  try {
    A.strings = JSON.parse(langFile);
  } catch(err){
    console.Error("Unable to parse strings JSON:", err);
  }
  console.info("Finished loading strings.");
}

// handleActions wraps around other handle functions
function handleActions(){
  handleIncludes();
  handleAppText();
  handleCloseButton();
}

// Handle whenever there is an <<include> in code
function handleIncludes(){
  var elem = {};
  $("include").each(function(){
    elem = $(this);
    if(elem.text().length != 0){
      return;
    }
    var layoutName = elem.attr("app:layout");
    elem.load("/res/layout/"+layoutName+".html", () => handleActions());
    console.info("Loaded", layoutName, "layout");
  });
}

// Handle app:text attribute
function handleAppText(){
  var elem = {};
  var textDescriptor = ["",""];
  var items = document.body.getElementsByTagName("*");
  for(var i = items.length; i--;){
    elem = items[i];
    if(!elem.hasAttribute("app:text")){
      continue; // Don't persist if the element does not have the attribute
    }
    textDescriptor = elem.getAttribute("app:text").split(":");
    switch(textDescriptor[0]){
      case "app":
        elem.innerHTML = A[textDescriptor[1]] + elem.innerHTML;
        break;
      case "string":
        elem.innerHTML = A.strings[textDescriptor[1]] + elem.innerHTML;
        break;
      default:
        break;
    }
  }
}

// Handle button to close the app
// TODO: Make it work...
function handleCloseButton(){
  $('#app_close_button').click(function(){
    gui.App.quit();
  });
}

// Loads the specific layout
function loadLayout(name){
  $(document.body).empty(); // Empty the body before writing to it
  $(document.body).load("/res/layout/"+name+".html layout", () => handleActions());
}
