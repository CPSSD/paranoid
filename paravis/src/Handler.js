// Connection server and port to local instance running the websocket server
const WebsocketServerUrl = "ws://127.0.0.1:10100";

class Handler {
  constructor(view){
    var self = this;

    console.log("Handler initialized");
    this.m_animationTime = 1000;

    // Action queue
    self.m_queue = new Queue();

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

    // Start handling
    this.handleMessage(self);

    // var popup = new Popup("Title", "new_node");
    // popup.show();

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
  handleMessage(self){
    this.m_ws.onmessage = function(e){
      console.log("Server: ", e.data);
      var msg = JSON.parse(e.data);
      self.m_queue.add(msg);

      // Handle Messages at specific interval
      var loop = function(){
        if(!self.m_queue.empty()){
          var m = self.m_queue.pop()
          switch(m.type){
            case 'event':
              self.handleEvent(m.data);
              break;
            case "state":
              self.handleStatus(m.data);
              break;
            case 'nodechange':
              self.handleNodeUpdate(m.data);
              break;
            default:
              console.error("unrecognished message type", m.type);
              break;
          }
        }

        setTimeout(loop, self.getProcessInterval());
      }
      loop()
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
      }).show());
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
        var newNode = new Node({
          name: n.commonName,
          uuid: n.uuid,
          address: n.addr,
          state: getStateFromString(n.state),
          parent: this.m_viewer,
          shown: true
        });
        newNode.show();

        this.addNode(newNode);
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
    target = this.nodes[data.target]
    source = this.nodes[data.source]
  }
}
module.exports.Handler = Handler;
