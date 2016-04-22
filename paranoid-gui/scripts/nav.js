var selected = -1;

function loadSideBar() {
  $("#nav").empty();
  var items = [
  ];

  var starter = '<li><a href="#" ';
  if (selected == -1) {
    starter += 'class="selected" ';
  }
  starter += 'onclick="rowClicked(-1)"><img src="images/buildings.png"/>  Home</a></li>';
  items.push(starter);

  $.each(fileSystems, function(i, item) {
    var entry = '<li><a href="#" ';
    if (selected == i) {
      entry += 'class="selected" ';
    }
    entry += 'onclick="rowClicked(' + i + ')"> ';
    if (item.mounted) {
      entry += '<img src="images/interface_green.png"/>  ';
    } else {
      entry += '<img src="images/interface.png"/>  ';
    }
    entry += item.name + '</a></li>';
    items.push(entry);
  });

  var newFsRow = '<li><a href="#" ';
  if (selected == fileSystems.length) {
    newFsRow += 'class="selected" ';
  }
  newFsRow += 'onclick="rowClicked(fileSystems.length)"> <img src="images/circle.png"/>  New filesystem</a></li>';
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
