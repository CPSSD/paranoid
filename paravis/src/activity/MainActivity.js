class MainActivity extends Activity {
    constructor() {
        super();
        this.m_layout = "main";
        this.m_db = new KeyValueDB("Main", 1, PermanentDB);
    }

    name(){
      return "MainActivity";
    }

    onCreate(){
      loadLayout(this.m_layout);

      if(this.m_db.get("filesystems") == undefined){
        // var p = new Popup(A.strings.pool_new_pool, "new_node");
        // p.onOk(() => {
        //   // Get the data from the input
        // });
      }

      $(document).click(() => {
        var handler = new Handler($('#viewer'));
      });
    }
}

// Register the activity
activities[MainActivity.name] = MainActivity;


// getStateFromString converts a | separated strings corresponding to node
// state and returns its representation as acceptable state
// Warning: This function does not resolve issues with state race conditions
//          and it presumes a valid state
function getStateFromString(stateString){
  var state = 0; // return state

  var states = stateString.split("|");

  //TODO: Handle the state change
  for(var i in states){
    switch(states[i]){
      case 'Follower':
        state &= FollowerState;
      case 'Leader':
        state &= LeaderState;
      case 'Inactive':
        state &= InactiveState;
      case 'Current':
        state &= CurrentState;
    }
  }

  return state;
}
