class Transition{
  constructor(){
    console.info("Running transition...");
  }

  // Run the transtion
  Run( callback ){
    console.info("Transition finished");
    callback()
  }
}

class TransitionFill extends Transition {
  constructor(){
    super();
  }

  Run( callback ){
    //TODO: Implement the fill transition

    super.Run(callback);
  }
}

transitions = {
  None: new Transition(),
  Fill: new TransitionFill(),
}
