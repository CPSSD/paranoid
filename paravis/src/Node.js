// Node states
const LeaderState     = 1<<1;
const FollowerState   = 1<<2;
const CandidateState  = 1<<3;
const CurrentState    = 1<<4;
const InactiveState   = 1<<5;

class Node {
  constructor(details){
    console.info("Node: Creating new node %s (%s)", details.name, details.uuid);
    this.m_name = details.name;
    this.m_uuid = details.uuid;
    this.m_position = details.position || {x: NaN, y: 0};
    this.m_state = details.state || FollowerState;
    this.m_parent = details.parent || {};
    this.m_shown = details.shown || false,
    this.m_address = details.address || "::"

    this.d_node = document.createElement("node");
    this.d_node.id = "_node"+this.m_name;

    if(this.m_parent != {}){
      this.appendTo(this.m_parent);
    }

    if( this.shown ){
      this.show();
    } else {
      this.hide();
    }

    if(this.m_position.x != NaN){
      this.setPosition(this.m_position);
    }

    switch(this.m_state){
      case LeaderState:
        this.changeToLeader();
        break;
      case FollowerState:
        this.changeToFollower();
        break;
      case CurrentState:
        this.changeToCurrent();
        break;
    }
  }

  getUuid(){
    return this.m_uuid;
  }

  getPosition(){
    return this.m_position;
  }

  setState(setState){

  }

  setPosition(position){
    this.m_position = position;
    $(this.d_node).css({
      top: position.x,
      left: position.y
    });
  }

  show(){
    $(this.d_node).show();
  }

  hide(){
    $(this.d_node).hide();
  }

  getNode(){
    return this.d_node;
  }

  appendTo(parent){
    parent.append(this.d_node);
  }

  changeToLeader(){
    this.d_node.setAttribute('class', 'leader');
  }

  changeToFollower(){
    this.d_node.setAttribute('class', '');
  }

  changeToCurrent(){
    this.d_node.setAttribute('class', 'current')
  }

  sendToPeer(peer, data) {

  }
}

module.exports = {
  Node: Node,
  LeaderState: LeaderState,
  CurrentState: CurrentState,
  FollowerState: FollowerState,
  CandidateState: CandidateState,
  InactiveState: InactiveState
}
