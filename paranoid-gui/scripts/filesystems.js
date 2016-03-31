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
    var fsPath = path.join(fileSystemsDir, fileNames[i]);
    var filesystem = {
      name: fileNames[i],
      path: path.join(fileSystemsDir, fileNames[i]),
      mounted: false,
      attributes: loadJsonFile(path.join(fsPath, "meta", "attributes")),
      pool: readFile(path.join(fsPath, "meta", "pool")),
      uuid: readFile(path.join(fsPath, "meta", "uuid"))
    };
    filesystems.push(filesystem);
  }

  return filesystems;
}

function drawFileSystem(i) {
  var fileSystem = fileSystems[i];

  var str = readFile("html/filesystem.html");
  var html = $.parseHTML(str);
  $(".content").html(html);

  $(".content #fsName").html(fileSystem.name);
  $(".content #fsMountSection").hide();

  // attributes section
  if (!fileSystem.attributes.encrypted) {
    $(".content #fsAttributes #fsAttributeEncrypted #badge").html('<b>NO<b>');
    $(".content #fsAttributes #fsAttributeEncrypted #badge").addClass("label label-danger");
  }

  if (fileSystem.attributes.networkoff) {
    $(".content #fsAttributes #fsAttributeNetwork #badge").html('<b>NO<b>');
    $(".content #fsAttributes #fsAttributeNetwork #badge").addClass("label label-danger");
  }

  if (!fileSystem.attributes.keygenerated) {
    $(".content #fsAttributes #fsAttributeKeygen #badge").html('<b>NO<b>');
    $(".content #fsAttributes #fsAttributeKeygen #badge").addClass("label label-danger");
  }

  $(".content #fsAttributes #fsAttributePool #badge").html('<b>' + fileSystem.pool + '<b>');
  $(".content #fsAttributes #fsAttributeUuid #badge").html('<b>' + fileSystem.uuid + '<b>');

  // delete section
  $(".content #fsDeleteSection #deleteButton").attr("onClick", 'deleteClicked($(".content #fsDeleteSection #nameCheckText").val(), ' + i + ');');

}

function newfs(form) {
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
    fileSystems = getFilesystems();
    $("#filist").empty();
    rowClicked(-1);
  });
}

function deleteClicked(val, i) {
  var fileSystem = fileSystems[i];
  if (val != fileSystem.name) {
    alert("Incorrect FileSystem name");
  } else {
    deleteFs(i);
  }
}

function deleteFs(i) {
  var exec = require('child_process').exec;
  var cmd = "paranoid-cli delete " + fileSystems[i].name;
  exec(cmd, function(error, stdout, stderr) {
    if (error !== null) {
      alert(error);
    }
    fileSystems = getFilesystems();
    $("#nav").empty();
    loadSideBar();
    rowClicked(i);
  });
}

function pathExists(filePath) {
  var fs = require("fs");
  try {
    fs.accessSync(filePath, fs.F_OK);
    return true;
  } catch(e) {
    return false;
  }
}

function readFile(filePath) {
  var fs = require("fs");
  return fs.readFileSync(filePath, "utf8");
}

function loadJsonFile(filePath) {
  return JSON.parse(readFile(filePath));
}
