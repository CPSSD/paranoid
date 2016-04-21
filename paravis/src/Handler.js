var node = require('../src/Node.js');
var ToRad = require('../src/utils/utils.js').ToRad;
var Queue = require('../src/utils/Queue.js').Queue;
var Popup = require('../_app/Popup.js').Popup;
var console = require('console');

// Connection server and port to local instance running the websocket server
const WebsocketServerUrl = "ws://127.0.0.1:10100";

class Handler {
  constructor(view){
    console.log("Handler initialized");
    this.m_animationTime = 1000;
    this.m_ws = {};
    try {
      this.m_ws = new WebSocket(WebsocketServerUrl);
    } catch( err ){
      console.log("Unable to connect to websocket: ", err );
    }
    // Get the viewer which would contain all elements
    this.m_viewer = view;

    // Store all nodes in an map
    this.nodes = {};

    // Action queue
    this.m_queue = new Queue();

    // Start handling
    this.handleMessage();

    var popup = new Popup("Title", "new_node", this.m_viewer);
    popup.show();

    console.log("popup:", popup);

    // ------------------------------------------------------------------------
    // TODO: REMOVE THIS
    var n1uuid ="1234-abcd-5678-efgh-9012"
    var n1 = new node.Node({
      name: "node1",
      uuid: n1uuid,
      shown: false,
      state: node.CurrentState,
      // position: {x: 100, y: 100},
      parent: this.m_viewer
    });

    var n2uuid = "4567-abcd-8901-efgh-2345"
    var n2 = new node.Node({
      name: 'node2',
      uuid: n2uuid,
      shown: false,
      state: node.LeaderState,
      // position: {x: 0, y: 200},
      parent: this.m_viewer
    });

    var n3uuid = "4567-abce-8901-efgh-2345"
    var n3 = new node.Node({
      name: 'node4',
      uuid: n3uuid,
      shown: false,
      state: node.FollowerState,
      parent: this.m_viewer
    });

    var n4uuid = "8567-abce-8901-efgh-2345"
    var n4 = new node.Node({
      name: 'node4',
      uuid: n4uuid,
      shown: false,
      state: node.InactiveState,
      parent: this.m_viewer
    });

    var n5uuid = "8568-abce-8901-efgh-2345"
    var n5 = new node.Node({
      name: 'node5',
      uuid: n5uuid,
      shown: false,
      state: node.CandidateState,
      parent: this.m_viewer
    });

    this.addNode(n1);
    n1.show();

    setTimeout( () => {
      this.addNode(n2);
      n2.show();

      this.addNode(n5);
      n5.show();

      setTimeout( () => {
        this.addNode(n3);
        n3.show();

        setTimeout( () => {
          this.addNode(n4);
          n4.show();
        }, 2000);
      }, 2000);
    }, 2000);
    // ------------------------------------------------------------------------

  }

  resetNodeLayout(){
    var center = {x: this.m_viewer.height()/2, y: this.m_viewer.width()/2};
    var nodesLen = Object.keys(this.nodes).length;
    var angle = ToRad(360 / nodesLen); console.log(angle);
    var distance = 100 + 10*nodesLen;

    console.log("Nodes:", Object.keys(this.nodes).length);

    var nodeSize = this.m_viewer.width()/10;

    var s, c, xpos, ypos;
    var currentAngle = ToRad(270);
    for(var i in this.nodes){
      s = Math.sin(currentAngle);
      c = Math.cos(currentAngle);
      xpos = center.x + (distance*c) - (distance*s) - nodeSize/2;
      ypos = center.y + (distance*s) + (distance*c) - nodeSize/2;
      console.log("New position: x:%d, y:%d (%d radians)", xpos, ypos, currentAngle);
      this.nodes[i].setPosition({x: xpos, y: ypos});
      currentAngle += angle;
    }
  }

  addNode(node){
    this.nodes[node.getUuid()] = node;
    this.resetNodeLayout();
  }

  changeAnimationTime(newTime){
    this.m_animationTime = newTime;
  }

  // Listen to whatever hanges in the raft log
  handleMessage(){
    this.m_ws.onmessage = function(e){
      var msg = JSON.parse(e);
      console.log("Server: ", msg);
      this.m_queue.add(msg);

      // Handle Messages at specific interval
      var loop = function(){
        switch(msg.type){
          case 'event':
            this.handleEvent(msg.data);
            break;
          case "status":
            this.handleStatus(msg.data);
            break;
          case 'nodechange':
            this.handleNodeUpdate(msg.data);
            break;
          default:
            console.error("unrecognished messsage type", msg.type);
            break;
        }

        setTimeout(loop, this.getProcessInterval());
      }

    }
  }

  getProcessInterval(){
    //TODO: Add a way to change this value
    return 500; //ms
  }

  // handleStatus takes care of the initial status which is send as first message
  // after connection is established
  handleStatus(data){
    var nodes = data.nodes || null;

    if(nodes == null){
      console.error("nodes cannot be null, cannot process status");
    }

    for(var n in nodes){
      this.addNode(new Node({
        name: nodes[n].commonName,
        uuid: nodes[n].uuid,
        state: getStateFromString(nodes[n].state),
        shown: true,
        parent: this.m_viewer
      }));
    }
  }

  handleNodeUpdate(data){
    var n = data.node;

    if(n === null){
      console.error("Node is null, unable to process nodechange");
      return;
    }

    switch(data.action){
      case 'add':
      case 'change':
        this.nodes[n.uuid] = new Node({
          name: n.commonName,
          uuid: n.uuid,
          address: n.addr,
          state: getStateFromString(n.state),
          parent: this.m_viewer,
          shown: true
        });
        break;
      case 'delete':
        delete(this.nodes[n.uuid]);
        break;
      default:
        console.error("Unknown action", data.action);
        return;
    }
  }

  handleEvent(data){
    //TODO: Implement event handling
  }
}
module.exports.Handler = Handler;
