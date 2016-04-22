
class Intent {
  constructor(activity, transition){
    this.transition = transition;
    if(activity.name == A.mainActivity.name()){
      console.error("This activity is already active");
      return;
    }
    this.act = new activity();
  }

  onCreate(){
    // Run the onCreate after the transition is done
    this.transition.Run(() => {
      this.act.onCreate();
    })

  }
}

var module = module || {exports: {}};
module.exports.Intent = Intent;
