var path = require("path");
var fs = require("fs");

function getUserHome() {
  return process.env.HOME || process.env.USERPROFILE;
}

function getFilesystems() {
  var fileSystemsDir = path.join(getUserHome(), ".pfs", "filesystems");
  var fileNames = fs.readdirSync(fileSystemsDir);
  var filesystems = [];
  for (var i=0; i<fileNames.length; i++) {
    var filesystem = {
      name: fileNames[i],
      path: path.join(fileSystemsDir, fileNames[i]),
      mounted: false,
      //attributes: require(path.join(fileSystemsDir, fileNames[i], "meta", "attributes"))
    };
    filesystems.push(filesystem);
  }

  return filesystems;
}

function drawFileSystem(i) {
  var fileSystem = fileSystems[i];
  var heading = '<h1>' + fileSystem.name + '</h1>';
  var status = '';
  if (fileSystem.mounted) {
    status += '<span class="label label-success">Mounted</span>';
  } else {
    status += '<span class="label label-default">Unmounted</span>';
  }

  var buttonGroupHeader = '<div id="buttons"><div class="btn-group">';
  var groupBodyMount = '<button type="button" class="btn btn-info btn-block" onclick="mountFs(' + i + ')">Mount</button>';
  var groupBodyUnmount = '<button type="button" class="btn btn-warning btn-block" onclick="unmountFs(' + i + ')">Unmount</button>';
  var groupBodyDelete = '<button type="button" class="btn btn-danger btn-block" onclick="deleteFs(' + i + ')">Delete</button></div></div>';

  $(".content").html(heading + status + buttonGroupHeader + groupBodyMount + groupBodyUnmount + groupBodyDelete);
}
/*
function drawFileSystem(fileS) {
  var items = [];
  $.each(fileSystems, function(i, item) {
    var panelIdentifier = '<div class="panel panel-';
    if (item.mounted) {
      panelIdentifier += 'success">';
    } else {
      panelIdentifier += 'primary">';
    }

    var panelheading = '<div class="panel-heading"><h3 class="panel-title">' + item.name + '</h3></div>';

    var panelBodyheader = '<div class="panel-body"><b>Path: </b>' + item.path + '<br>';
    var panelBodyMount = '<div class="row"><div class="col-md-3"><button type="button" class="btn btn-info btn-block" onclick="mountFs(' + i + ')">Mount</button></div>';
    var panelBodyUnmount = '<div class="col-md-3"><button type="button" class="btn btn-warning btn-block" onclick="unmountFs(' + i + ')">Unmount</button></div>';
    var panelBodyDelete = '<div class="col-md-3"><button type="button" class="btn btn-danger btn-block" onclick="deleteFs(' + i + ')">Delete</button></div></div>';

    var panelBody = panelBodyheader + panelBodyMount + panelBodyUnmount +
      panelBodyDelete + "</div>";

    items.push(panelIdentifier + panelheading + panelBody + "</div>");
  });

  $("#filist").append(items.join(' '));
}*/

function newfs(form) {
  console.log(form);
  var exec = require('child_process').exec;
  var cmd = "paranoid-cli init ";
  if (!form.secure.checked) {
    cmd += "-u ";
  }
  if (!form.network.checked) {
    cmd += "--networkoff ";
  }
  if (!form.encrypted.checked) {
    cmd += "--unencrypted ";
  }
  if (form.cert.value !== "") {
    cmd += "--cert " + form.cert.value + " ";
  }
  if (form.key.value !== "") {
    cmd += "--key " + form.key.value + " ";
  }
  if (form.pool.value !== "") {
    cmd += "--pool " + form.pool.value + " ";
  }
  cmd += form.name.value;

  exec(cmd, function(error, stdout, stderr) {
    console.log(error);
    fileSystems = getFilesystems();
    $("#filist").empty();
    rowClicked(-1);
  });
}

function deleteFs(i) {
  var exec = require('child_process').exec;
  var cmd = "paranoid-cli delete " + fileSystems[i].name;
  exec(cmd, function(error, stdout, stderr) {
    console.log(error);
    fileSystems = getFilesystems();
    $("#nav").empty();
    loadSideBar();
    rowClicked(i);
  });
}

function mountFs(i) {
  fileSystems[i].mounted = true;
  $("#nav").empty();
  loadSideBar();
  rowClicked(i);
}

function unmountFs(i) {
  var exec = require('child_process').exec;
  var cmd = "paranoid-cli unmount " + fileSystems[i].name;
  exec(cmd, function(error, stdout, stderr) {
    console.log(error);
    fileSystems[i].mounted = false;
    $("#nav").empty();
    loadSideBar();
    rowClicked(i);
  });
}
