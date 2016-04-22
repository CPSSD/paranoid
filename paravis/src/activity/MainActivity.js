class MainActivity extends Activity {
    constructor() {
        super();
        this.m_layout = "main";
    }

    onCreate(){
      loadLayout(this.m_layout);

      setTimeout(() => {
        var handler = new Handler($('#viewer'));
      }, 100);

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
