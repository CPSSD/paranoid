class Activity {
  constructor(){
    console.log("Action initialized");
    this.layout = "";
  }

  layout(){
    return this.layout;
  }

  onCreate() {
    console.log("onCreate called");
  }
}

// All activies must be registered
// activities[Activity.name] =  Activity;
