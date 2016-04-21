
class Intent {
  constructor(activity){
    this.act = new activity();
  }

  onCreate(){
    this.act.onCreate();
  }
}

var module = module || {exports: {}};
module.exports.Intent = Intent;
