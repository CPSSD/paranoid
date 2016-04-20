class Queue {
  constructor(count){
    this.start = 0; // start of the queue
    this.end = 0; // end of the queue
    this.data = new Array(count); // data
  }

  add(item){
    this.data.append(item);
    this.end++;
  }

  pop(){
    this.start++;
    return this.data[this.start-1];
  }

}

module.exports.Queue = Queue;
