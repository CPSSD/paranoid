function setupHome() {
  $(".content").load("html/home.html");
  pfsdRunning(function(running) {
    if (running) {
      $(".content #runningLabel").attr("class", "label label-success");
      $(".content #runningLabel").html("running");
    } else {
      $(".content #runningLabel").attr("class", "label label-warning");
      $(".content #runningLabel").html("not running");
    }
  });
}

function pfsdRunning(callback) {
  var exec = require("child_process").exec;
  var cmd = "pidof pfsd";

  exec(cmd, function(error, stdout, stderr) {
    if (error !== null) {
      return callback(false);
    }
    callback(true);
  });
}
