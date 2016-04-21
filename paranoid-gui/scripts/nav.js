var selected = -1;

function loadSideBar() {
  $("#nav").empty();
  var items = [
  ];

  var starter = '<li><a href="#" ';
  if (selected == -1) {
    starter += 'class="selected" ';
  }
  starter += 'onclick="rowClicked(-1)"><img src="images/home_32.png"/> Home</a></li>';
  items.push(starter);

  $.each(fileSystems, function(i, item) {
    var entry = '<li><a href="#" ';
    if (selected == i) {
      entry += 'class="selected" ';
    }
    entry += 'onclick="rowClicked(' + i + ')"> ';
    if (item.mounted) {
      entry += '<img src="images/Green_checkmark.png"/> ';
    } else {
      entry += '<img src="images/file_32.png"/> ';
    }
    entry += item.name + '</a></li>';
    items.push(entry);
  });

  var newFsRow = '<li><a href="#" ';
  if (selected == fileSystems.length) {
    newFsRow += 'class="selected" ';
  }
  newFsRow += 'onclick="rowClicked(fileSystems.length)"> <img src="images/plus_edited_32.png"/> New FileSystem</a></li>';
  items.push(newFsRow);

  $("#nav").append(items.join(' '));
}

function rowClicked(i) {
  if (i < -1 || i >= fileSystems.length +1) {
    i = -1;
  }
  selected = i;
  if (i == -1) {
    setupHome();
  } else if (i == fileSystems.length){
    $(".content").load("html/newFs.html");
  } else {
    drawFileSystem(i);
  }
  loadSideBar();
}
